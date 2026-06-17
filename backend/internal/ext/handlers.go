// Package ext — внешний HTTP-слой /api/v1/ext/* (единственный публичный контракт, tech-stack §2).
package ext

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"

	"github.com/roof/dupes/internal/pipeline"
	"github.com/roof/dupes/internal/scoring"
	"github.com/roof/dupes/internal/storage"
)

type Handlers struct {
	Store   *storage.Store
	Pipe    *pipeline.Pipeline
	BlobDir string
	Version string
	Log     zerolog.Logger
}

func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	depth, _ := h.Store.QueueDepth()
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":               "ok",
		"version":              h.Version,
		"queue_depth":          depth,
		"deprecation_warnings": []string{},
	})
}

func (h *Handlers) Stats(w http.ResponseWriter, r *http.Request) {
	objects, photos, matches, err := h.Store.Counts()
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	depth, _ := h.Store.QueueDepth()
	writeJSON(w, http.StatusOK, map[string]int{
		"objects": objects, "photos": photos, "matches": matches, "queue_depth": depth,
	})
}

func (h *Handlers) UpsertContractor(w http.ResponseWriter, r *http.Request) {
	var c storage.Contractor
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil || c.ContractorID == "" {
		writeErr(w, http.StatusBadRequest, "нужен contractor_id")
		return
	}
	if err := h.Store.UpsertContractor(c); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, c)
}

type objectRow struct {
	storage.Object
	PhotoCount int `json:"photo_count"`
}

// ListObjects отдаёт все карточки портфолио (в т.ч. без проблем) с числом фото.
func (h *Handlers) ListObjects(w http.ResponseWriter, r *http.Request) {
	objs, err := h.Store.ListObjects()
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	counts, _ := h.Store.AllPhotoCounts()
	out := make([]objectRow, 0, len(objs))
	for _, o := range objs {
		out = append(out, objectRow{Object: o, PhotoCount: counts[o.ObjectID]})
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"items": out, "count": len(out)})
}

func (h *Handlers) UpsertObject(w http.ResponseWriter, r *http.Request) {
	var o storage.Object
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil || o.ObjectID == "" {
		writeErr(w, http.StatusBadRequest, "нужен object_id")
		return
	}
	if err := h.Store.UpsertObject(o); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, o)
}

// IngestPhoto принимает multipart-форму (file + метаданные), сохраняет байты в blob-хранилище,
// создаёт запись фото и ставит задачу в очередь обработки (ТЗ §6.1).
func (h *Handlers) IngestPhoto(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(64 << 20); err != nil {
		writeErr(w, http.StatusBadRequest, "ожидается multipart-форма")
		return
	}
	file, hdr, err := r.FormFile("file")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "нужен файл в поле file")
		return
	}
	defer file.Close()

	photoID := storage.NewID()
	if err := os.MkdirAll(h.BlobDir, 0o755); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	blobPath := filepath.Join(h.BlobDir, photoID)
	dst, err := os.Create(blobPath)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	if _, err := io.Copy(dst, file); err != nil {
		dst.Close()
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	dst.Close()

	origName := r.FormValue("original_filename")
	if origName == "" {
		origName = hdr.Filename
	}
	p := storage.Photo{
		PhotoID:          photoID,
		ContractorID:     r.FormValue("contractor_id"),
		ObjectID:         r.FormValue("object_id"),
		UploadDate:       time.Now(),
		OriginalFilename: origName,
		SourceURL:        r.FormValue("source_url"),
		BlobPath:         blobPath,
		ModerationStatus: "new",
	}
	if err := h.Store.InsertPhoto(p); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	jobID, err := h.Store.EnqueueJob(photoID)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]string{"photo_id": photoID, "job_id": jobID})
}

func (h *Handlers) PhotoContent(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	p, err := h.Store.GetPhoto(id)
	if err != nil {
		writeErr(w, http.StatusNotFound, "фото не найдено")
		return
	}
	http.ServeFile(w, r, p.BlobPath)
}

func (h *Handlers) ListMatches(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	f := storage.MatchFilter{
		MinRisk:        atoi(q.Get("min_risk")),
		MaxRisk:        atoi(q.Get("max_risk")),
		SimilarityType: q.Get("type"),
		Status:         q.Get("status"),
		ContractorID:   q.Get("contractor_id"),
		ObjectID:       q.Get("object_id"),
		Source:         q.Get("source"),
		Limit:          atoi(q.Get("limit")),
	}
	if f.Limit == 0 {
		f.Limit = 200
	}
	matches, err := h.Store.ListMatches(f)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	out := make([]matchView, 0, len(matches))
	for _, m := range matches {
		band, action := scoring.Band(m.RiskScore)
		out = append(out, matchView{Match: m, RiskBand: band, RiskAction: action})
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"items": out, "count": len(out)})
}

type matchView struct {
	storage.Match
	RiskBand   string `json:"risk_band"`
	RiskAction string `json:"risk_action"`
}

// MatchDetail отдаёт совпадение с обоими фото, обеими карточками объектов и расшифровкой риска (ТЗ §8).
func (h *Handlers) MatchDetail(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	m, err := h.Store.GetMatch(id)
	if err != nil {
		writeErr(w, http.StatusNotFound, "совпадение не найдено")
		return
	}
	band, action := scoring.Band(m.RiskScore)
	resp := map[string]interface{}{
		"match":       m,
		"risk_band":   band,
		"risk_action": action,
		"photo_1":     h.photoOrNil(m.PhotoID1),
		"photo_2":     h.photoOrNil(m.PhotoID2),
		"object_1":    h.objectOrNil(m.ObjectID1),
		"object_2":    h.objectOrNil(m.ObjectID2),
	}
	writeJSON(w, http.StatusOK, resp)
}

var allowedDecisions = map[string]bool{
	"confirmed": true, "false_positive": true, "similar_not_duplicate": true,
	"manufacturer_photo": true, "external_source": true, "pending": true,
}

func (h *Handlers) DecideMatch(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var body struct {
		Status  string `json:"status"`
		Comment string `json:"comment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || !allowedDecisions[body.Status] {
		writeErr(w, http.StatusBadRequest, "недопустимый статус решения")
		return
	}
	if err := h.Store.SetMatchDecision(id, body.Status, body.Comment); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": body.Status})
}

func (h *Handlers) RecheckObjects(w http.ResponseWriter, r *http.Request) {
	n, err := h.Pipe.RecheckObjects()
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"object_matches": n})
}

func (h *Handlers) Reindex(w http.ResponseWriter, r *http.Request) {
	n, err := h.Pipe.Reindex()
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]int{"queued_photos": n})
}

func (h *Handlers) photoOrNil(id string) interface{} {
	if id == "" {
		return nil
	}
	p, err := h.Store.GetPhoto(id)
	if err != nil {
		return nil
	}
	return p
}

func (h *Handlers) objectOrNil(id string) interface{} {
	if id == "" {
		return nil
	}
	o, err := h.Store.GetObject(id)
	if err != nil {
		return nil
	}
	return o
}

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
