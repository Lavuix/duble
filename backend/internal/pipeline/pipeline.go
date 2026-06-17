// Package pipeline — асинхронная обработка фото и пересчёт совпадений (ТЗ §6.1).
// Устойчивость по lessons §А3: семафор-пул воркеров, таймаут на задачу,
// watchdog возвращает зависшие задачи, RequeueUnfinished при старте.
package pipeline

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/rs/zerolog"

	"github.com/roof/dupes/internal/hashing"
	"github.com/roof/dupes/internal/metrics"
	"github.com/roof/dupes/internal/scoring"
	"github.com/roof/dupes/internal/storage"
)

type Pipeline struct {
	store   *storage.Store
	score   scoring.Config
	log     zerolog.Logger
	workers int
	persist int    // MatchPersistThreshold
	blobDir string // куда складывать скачанные по URL изображения
}

func New(store *storage.Store, score scoring.Config, persist, workers int, blobDir string, log zerolog.Logger) *Pipeline {
	return &Pipeline{store: store, score: score, log: log, workers: workers, persist: persist, blobDir: blobDir}
}

// Start запускает пул воркеров и watchdog до отмены ctx.
func (p *Pipeline) Start(ctx context.Context) {
	if n, err := p.store.RequeueStuck(0); err == nil && n > 0 {
		p.log.Info().Int64("requeued", n).Msg("RequeueUnfinished: незавершённые задачи возвращены в очередь")
	}
	for i := 0; i < p.workers; i++ {
		go p.worker(ctx, i)
	}
	go p.watchdog(ctx)
}

func (p *Pipeline) worker(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		job, err := p.store.ClaimJob()
		if err != nil {
			p.log.Error().Err(err).Msg("claim job")
			time.Sleep(time.Second)
			continue
		}
		if job == nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		p.processJob(ctx, job)
	}
}

func (p *Pipeline) watchdog(ctx context.Context) {
	t := time.NewTicker(time.Minute)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			if n, err := p.store.RequeueStuck(2 * time.Minute); err == nil && n > 0 {
				p.log.Warn().Int64("requeued", n).Msg("watchdog: зависшие задачи возвращены в очередь")
			}
			if d, err := p.store.QueueDepth(); err == nil {
				metrics.IngestQueueDepth.Set(float64(d))
			}
		}
	}
}

func (p *Pipeline) processJob(ctx context.Context, job *storage.Job) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := p.processPhoto(ctx, job.PhotoID); err != nil {
		// Ошибка обработки/скачивания — это НЕ дубль, логируем отдельно (ТЗ §11).
		p.log.Error().Err(err).Str("photo_id", job.PhotoID).Msg("ошибка обработки фото")
		_ = p.store.FinishJob(job.JobID, "failed", err.Error())
		metrics.IngestJobs.WithLabelValues("failed").Inc()
		return
	}
	_ = p.store.FinishJob(job.JobID, "done", "")
	metrics.IngestJobs.WithLabelValues("done").Inc()
	metrics.IngestLastSuccess.SetToCurrentTime()
}

func (p *Pipeline) processPhoto(ctx context.Context, photoID string) error {
	ph, err := p.store.GetPhoto(photoID)
	if err != nil {
		return fmt.Errorf("get photo: %w", err)
	}
	// Если локального файла нет, но есть source_url — скачиваем (ТЗ §6.1 п.1-2, §11:
	// таймаут на внешний вызов; недоступная ссылка логируется отдельно, это не дубль).
	if ph.BlobPath == "" && ph.SourceURL != "" {
		if err := p.downloadToBlob(ctx, ph); err != nil {
			return fmt.Errorf("скачивание по ссылке не удалось: %w", err)
		}
	}
	data, err := os.ReadFile(ph.BlobPath)
	if err != nil {
		return fmt.Errorf("read blob: %w", err)
	}
	feats, err := hashing.Compute(data)
	if err != nil {
		return err // битый/нечитаемый файл
	}
	ph.SHA256 = feats.SHA256
	ph.PHash, ph.DHash, ph.AHash = feats.PHash, feats.DHash, feats.AHash
	ph.Width, ph.Height, ph.FileSize = feats.Width, feats.Height, feats.Size
	ph.ExifCamera, ph.ExifDateTime, ph.ExifGPS = feats.ExifCamera, feats.ExifDateTime, feats.ExifGPS

	others, err := p.store.ListOtherPhotos(photoID)
	if err != nil {
		return fmt.Errorf("list photos: %w", err)
	}
	oNew, _ := p.store.GetObject(ph.ObjectID)

	maxScore := 0
	for _, other := range others {
		if other.ObjectID == ph.ObjectID {
			continue // фото внутри одной карточки естественно похожи — не считаем дублем
		}
		oOther, _ := p.store.GetObject(other.ObjectID)
		score, simType, dist, factors := p.score.ScorePhotoPair(*ph, other, oNew, oOther)

		isVisual := simType == "exact" || (simType == "perceptual" && dist <= p.score.PhashThreshold)
		if !isVisual {
			continue
		}
		if simType != "exact" && score < p.persist {
			continue
		}
		first := other.PhotoID
		if ph.UploadDate.Before(other.UploadDate) {
			first = ph.PhotoID
		}
		m := storage.Match{
			PairKey:              pairKey("photo", ph.PhotoID, other.PhotoID),
			PhotoID1:             ph.PhotoID,
			PhotoID2:             other.PhotoID,
			ObjectID1:            ph.ObjectID,
			ObjectID2:            other.ObjectID,
			ContractorID1:        ph.ContractorID,
			ContractorID2:        other.ContractorID,
			SimilarityType:       simType,
			HashDistance:         dist,
			FirstUploadedPhotoID: first,
			RiskScore:            score,
			RiskFactors:          factors,
		}
		if err := p.store.UpsertMatch(m); err != nil {
			p.log.Error().Err(err).Msg("upsert photo match")
			continue
		}
		metrics.MatchesCreated.WithLabelValues(simType).Inc()
		if score > maxScore {
			maxScore = score
		}
	}
	ph.DuplicateScore = maxScore
	return p.store.UpdatePhotoFeatures(*ph)
}

// downloadToBlob скачивает изображение по source_url во временное хранилище и проставляет blob_path.
func (p *Pipeline) downloadToBlob(ctx context.Context, ph *storage.Photo) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ph.SourceURL, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	if err := os.MkdirAll(p.blobDir, 0o755); err != nil {
		return err
	}
	path := filepath.Join(p.blobDir, ph.PhotoID)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := io.Copy(f, io.LimitReader(resp.Body, 64<<20)); err != nil {
		return err
	}
	ph.BlobPath = path
	return p.store.SetPhotoBlobPath(ph.PhotoID, path)
}

// RecheckObjects пересчитывает дубли-данных (object_context) по всем парам объектов (ТЗ §6.1 п.5-6).
// Лёгкая операция на масштабе MVP; вызывается сидом, /objects:recheck и /reindex.
func (p *Pipeline) RecheckObjects() (int, error) {
	objs, err := p.store.ListObjects()
	if err != nil {
		return 0, err
	}
	created := 0
	for i := 0; i < len(objs); i++ {
		for j := i + 1; j < len(objs); j++ {
			a, b := objs[i], objs[j]
			shared, _ := p.store.CountPhotoMatchesBetween(a.ObjectID, b.ObjectID)
			score, factors := p.score.ScoreObjectPair(a, b, shared)
			if score < p.persist {
				continue
			}
			m := storage.Match{
				PairKey:        pairKey("obj", a.ObjectID, b.ObjectID),
				ObjectID1:      a.ObjectID,
				ObjectID2:      b.ObjectID,
				ContractorID1:  a.ContractorID,
				ContractorID2:  b.ContractorID,
				SimilarityType: "object_context",
				RiskScore:      score,
				RiskFactors:    factors,
			}
			if err := p.store.UpsertMatch(m); err != nil {
				p.log.Error().Err(err).Msg("upsert object match")
				continue
			}
			metrics.MatchesCreated.WithLabelValues("object_context").Inc()
			created++
		}
	}
	return created, nil
}

// Reindex переотправляет все фото в очередь (фоновая переиндексация базы, ТЗ §11).
func (p *Pipeline) Reindex() (int, error) {
	objs, err := p.store.ListObjects()
	if err != nil {
		return 0, err
	}
	n := 0
	for _, o := range objs {
		photos, err := p.store.PhotosByObject(o.ObjectID)
		if err != nil {
			return n, err
		}
		for _, ph := range photos {
			if _, err := p.store.EnqueueJob(ph.PhotoID); err != nil {
				return n, err
			}
			n++
		}
	}
	return n, nil
}

func pairKey(kind, a, b string) string {
	ids := []string{a, b}
	sort.Strings(ids)
	return kind + ":" + ids[0] + "|" + ids[1]
}
