package storage

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
)

func NewID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// ---- Contractors ----

func (s *Store) UpsertContractor(c Contractor) error {
	_, err := s.db.Exec(`INSERT INTO contractors(contractor_id,name) VALUES($1,$2)
		ON CONFLICT(contractor_id) DO UPDATE SET name=EXCLUDED.name`, c.ContractorID, c.Name)
	return err
}

// ---- Objects ----

func (s *Store) UpsertObject(o Object) error {
	_, err := s.db.Exec(`INSERT INTO objects(object_id,roof_id,contractor_id,name,address,gps,customer,year,materials,system,region)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		ON CONFLICT(object_id) DO UPDATE SET roof_id=EXCLUDED.roof_id,contractor_id=EXCLUDED.contractor_id,
		  name=EXCLUDED.name,address=EXCLUDED.address,gps=EXCLUDED.gps,customer=EXCLUDED.customer,
		  year=EXCLUDED.year,materials=EXCLUDED.materials,system=EXCLUDED.system,region=EXCLUDED.region`,
		o.ObjectID, o.RoofID, nullStr(o.ContractorID), o.Name, o.Address, o.GPS, o.Customer, o.Year, o.Materials, o.System, o.Region)
	return err
}

func (s *Store) GetObject(id string) (*Object, error) {
	o := &Object{}
	var contractor sql.NullString
	err := s.db.QueryRow(`SELECT object_id,roof_id,COALESCE(contractor_id,''),name,address,gps,customer,year,materials,system,region,active
		FROM objects WHERE object_id=$1`, id).
		Scan(&o.ObjectID, &o.RoofID, &contractor, &o.Name, &o.Address, &o.GPS, &o.Customer, &o.Year, &o.Materials, &o.System, &o.Region, &o.Active)
	if err != nil {
		return nil, err
	}
	o.ContractorID = contractor.String
	return o, nil
}

func (s *Store) ListObjects() ([]Object, error) {
	rows, err := s.db.Query(`SELECT object_id,roof_id,COALESCE(contractor_id,''),name,address,gps,customer,year,materials,system,region,active FROM objects ORDER BY object_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Object
	for rows.Next() {
		var o Object
		if err := rows.Scan(&o.ObjectID, &o.RoofID, &o.ContractorID, &o.Name, &o.Address, &o.GPS, &o.Customer, &o.Year, &o.Materials, &o.System, &o.Region, &o.Active); err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, rows.Err()
}

// SetObjectsActive массово включает/выключает активность карточек (ТЗ §12).
func (s *Store) SetObjectsActive(ids []string, active bool) error {
	_, err := s.db.Exec(`UPDATE objects SET active=$1 WHERE object_id = ANY($2)`, active, pq.Array(ids))
	return err
}

// Reset полностью очищает данные (объекты, фото, совпадения, очередь). Подрядчики тоже.
func (s *Store) Reset() error {
	_, err := s.db.Exec(`TRUNCATE matches, ingest_jobs, photos, objects, contractors RESTART IDENTITY CASCADE`)
	return err
}

// ---- Photos ----

func (s *Store) InsertPhoto(p Photo) error {
	_, err := s.db.Exec(`INSERT INTO photos(photo_id,contractor_id,object_id,upload_date,original_filename,source_url,blob_path,
		file_hash_sha256,phash,dhash,ahash,width,height,file_size,exif_camera,exif_datetime,exif_gps,moderation_status,duplicate_score)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19)`,
		p.PhotoID, nullStr(p.ContractorID), nullStr(p.ObjectID), p.UploadDate, p.OriginalFilename, p.SourceURL, p.BlobPath,
		p.SHA256, p.PHash, p.DHash, p.AHash, p.Width, p.Height, p.FileSize, p.ExifCamera, p.ExifDateTime, p.ExifGPS,
		p.ModerationStatus, p.DuplicateScore)
	return err
}

func (s *Store) UpdatePhotoFeatures(p Photo) error {
	_, err := s.db.Exec(`UPDATE photos SET file_hash_sha256=$2,phash=$3,dhash=$4,ahash=$5,width=$6,height=$7,file_size=$8,
		exif_camera=$9,exif_datetime=$10,exif_gps=$11,duplicate_score=$12 WHERE photo_id=$1`,
		p.PhotoID, p.SHA256, p.PHash, p.DHash, p.AHash, p.Width, p.Height, p.FileSize,
		p.ExifCamera, p.ExifDateTime, p.ExifGPS, p.DuplicateScore)
	return err
}

func (s *Store) SetPhotoBlobPath(id, path string) error {
	_, err := s.db.Exec(`UPDATE photos SET blob_path=$2 WHERE photo_id=$1`, id, path)
	return err
}

func (s *Store) GetPhoto(id string) (*Photo, error) {
	return s.scanPhoto(s.db.QueryRow(photoCols+` WHERE photo_id=$1`, id))
}

// ListOtherPhotos — все фото, кроме указанного (для линейного поиска кандидатов на MVP-масштабе).
func (s *Store) ListOtherPhotos(exclude string) ([]Photo, error) {
	rows, err := s.db.Query(photoCols+` WHERE photo_id<>$1 AND file_hash_sha256<>''`, exclude)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return s.scanPhotos(rows)
}

func (s *Store) PhotosByObject(objectID string) ([]Photo, error) {
	rows, err := s.db.Query(photoCols+` WHERE object_id=$1`, objectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return s.scanPhotos(rows)
}

const photoCols = `SELECT photo_id,COALESCE(contractor_id,''),COALESCE(object_id,''),upload_date,original_filename,source_url,blob_path,
	file_hash_sha256,phash,dhash,ahash,width,height,file_size,exif_camera,exif_datetime,exif_gps,moderation_status,duplicate_score FROM photos`

func (s *Store) scanPhoto(row *sql.Row) (*Photo, error) {
	p := &Photo{}
	err := row.Scan(&p.PhotoID, &p.ContractorID, &p.ObjectID, &p.UploadDate, &p.OriginalFilename, &p.SourceURL, &p.BlobPath,
		&p.SHA256, &p.PHash, &p.DHash, &p.AHash, &p.Width, &p.Height, &p.FileSize, &p.ExifCamera, &p.ExifDateTime, &p.ExifGPS,
		&p.ModerationStatus, &p.DuplicateScore)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Store) scanPhotos(rows *sql.Rows) ([]Photo, error) {
	var out []Photo
	for rows.Next() {
		var p Photo
		if err := rows.Scan(&p.PhotoID, &p.ContractorID, &p.ObjectID, &p.UploadDate, &p.OriginalFilename, &p.SourceURL, &p.BlobPath,
			&p.SHA256, &p.PHash, &p.DHash, &p.AHash, &p.Width, &p.Height, &p.FileSize, &p.ExifCamera, &p.ExifDateTime, &p.ExifGPS,
			&p.ModerationStatus, &p.DuplicateScore); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// ---- Matches ----

// UpsertMatch сохраняет совпадение по pair_key. Решение модератора (status/comment) НЕ перезаписывается
// при пересчёте — false_positive продолжает подавлять пару (ТЗ §13).
func (s *Store) UpsertMatch(m Match) error {
	factors, _ := json.Marshal(m.RiskFactors)
	if m.MatchID == "" {
		m.MatchID = NewID()
	}
	_, err := s.db.Exec(`INSERT INTO matches(match_id,pair_key,photo_id_1,photo_id_2,object_id_1,object_id_2,
		contractor_id_1,contractor_id_2,similarity_type,hash_distance,first_uploaded_photo_id,risk_score,risk_factors,status)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13::jsonb,'pending')
		ON CONFLICT(pair_key) DO UPDATE SET similarity_type=EXCLUDED.similarity_type,hash_distance=EXCLUDED.hash_distance,
		  first_uploaded_photo_id=EXCLUDED.first_uploaded_photo_id,risk_score=EXCLUDED.risk_score,
		  risk_factors=EXCLUDED.risk_factors,updated_at=now()`,
		m.MatchID, m.PairKey, nullStr(m.PhotoID1), nullStr(m.PhotoID2), nullStr(m.ObjectID1), nullStr(m.ObjectID2),
		nullStr(m.ContractorID1), nullStr(m.ContractorID2), m.SimilarityType, m.HashDistance,
		nullStr(m.FirstUploadedPhotoID), m.RiskScore, string(factors))
	return err
}

type MatchFilter struct {
	MinRisk        int
	MaxRisk        int
	SimilarityType string
	Status         string
	ContractorID   string
	ObjectID       string
	Source         string // internal | external (на MVP всё internal)
	Limit          int
}

func (s *Store) ListMatches(f MatchFilter) ([]Match, error) {
	q := `SELECT match_id,pair_key,COALESCE(photo_id_1,''),COALESCE(photo_id_2,''),COALESCE(object_id_1,''),COALESCE(object_id_2,''),
		COALESCE(contractor_id_1,''),COALESCE(contractor_id_2,''),similarity_type,hash_distance,COALESCE(first_uploaded_photo_id,''),
		risk_score,risk_factors,status,moderator_comment,created_at,updated_at FROM matches WHERE 1=1`
	var args []interface{}
	add := func(cond string, val interface{}) {
		args = append(args, val)
		q += fmt.Sprintf(" AND %s$%d", cond, len(args))
	}
	if f.MinRisk > 0 {
		add("risk_score>=", f.MinRisk)
	}
	if f.MaxRisk > 0 {
		add("risk_score<=", f.MaxRisk)
	}
	if f.SimilarityType != "" {
		add("similarity_type=", f.SimilarityType)
	}
	if f.Status != "" {
		add("status=", f.Status)
	}
	if f.ContractorID != "" {
		args = append(args, f.ContractorID, f.ContractorID)
		q += fmt.Sprintf(" AND (contractor_id_1=$%d OR contractor_id_2=$%d)", len(args)-1, len(args))
	}
	if f.ObjectID != "" {
		args = append(args, f.ObjectID, f.ObjectID)
		q += fmt.Sprintf(" AND (object_id_1=$%d OR object_id_2=$%d)", len(args)-1, len(args))
	}
	q += " ORDER BY risk_score DESC, updated_at DESC"
	if f.Limit > 0 {
		q += fmt.Sprintf(" LIMIT %d", f.Limit)
	}
	rows, err := s.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMatches(rows)
}

func (s *Store) GetMatch(id string) (*Match, error) {
	row := s.db.QueryRow(`SELECT match_id,pair_key,COALESCE(photo_id_1,''),COALESCE(photo_id_2,''),COALESCE(object_id_1,''),COALESCE(object_id_2,''),
		COALESCE(contractor_id_1,''),COALESCE(contractor_id_2,''),similarity_type,hash_distance,COALESCE(first_uploaded_photo_id,''),
		risk_score,risk_factors,status,moderator_comment,created_at,updated_at FROM matches WHERE match_id=$1`, id)
	var m Match
	var factors []byte
	err := row.Scan(&m.MatchID, &m.PairKey, &m.PhotoID1, &m.PhotoID2, &m.ObjectID1, &m.ObjectID2,
		&m.ContractorID1, &m.ContractorID2, &m.SimilarityType, &m.HashDistance, &m.FirstUploadedPhotoID,
		&m.RiskScore, &factors, &m.Status, &m.ModeratorComment, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(factors, &m.RiskFactors)
	return &m, nil
}

func (s *Store) SetMatchDecision(id, status, comment string) error {
	_, err := s.db.Exec(`UPDATE matches SET status=$2,moderator_comment=$3,updated_at=now() WHERE match_id=$1`, id, status, comment)
	return err
}

func scanMatches(rows *sql.Rows) ([]Match, error) {
	var out []Match
	for rows.Next() {
		var m Match
		var factors []byte
		if err := rows.Scan(&m.MatchID, &m.PairKey, &m.PhotoID1, &m.PhotoID2, &m.ObjectID1, &m.ObjectID2,
			&m.ContractorID1, &m.ContractorID2, &m.SimilarityType, &m.HashDistance, &m.FirstUploadedPhotoID,
			&m.RiskScore, &factors, &m.Status, &m.ModeratorComment, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(factors, &m.RiskFactors)
		out = append(out, m)
	}
	return out, rows.Err()
}

// ---- Ingest jobs (очередь обработки) ----

func (s *Store) EnqueueJob(photoID string) (string, error) {
	id := NewID()
	_, err := s.db.Exec(`INSERT INTO ingest_jobs(job_id,photo_id,status) VALUES($1,$2,'pending')`, id, photoID)
	return id, err
}

// ClaimJob атомарно берёт одну pending-задачу в обработку (FOR UPDATE SKIP LOCKED).
func (s *Store) ClaimJob() (*Job, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	var j Job
	err = tx.QueryRow(`SELECT job_id,photo_id,status,attempts FROM ingest_jobs
		WHERE status='pending' ORDER BY created_at LIMIT 1 FOR UPDATE SKIP LOCKED`).
		Scan(&j.JobID, &j.PhotoID, &j.Status, &j.Attempts)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if _, err := tx.Exec(`UPDATE ingest_jobs SET status='processing',attempts=attempts+1,started_at=now(),updated_at=now() WHERE job_id=$1`, j.JobID); err != nil {
		return nil, err
	}
	return &j, tx.Commit()
}

func (s *Store) FinishJob(jobID, status, errMsg string) error {
	_, err := s.db.Exec(`UPDATE ingest_jobs SET status=$2,error=$3,updated_at=now() WHERE job_id=$1`, jobID, status, errMsg)
	return err
}

// RequeueStuck переотправляет «processing»-задачи старше cutoff (watchdog/restart, lessons §А3).
func (s *Store) RequeueStuck(olderThan time.Duration) (int64, error) {
	res, err := s.db.Exec(`UPDATE ingest_jobs SET status='pending',updated_at=now()
		WHERE status='processing' AND started_at < now() - $1::interval`, fmt.Sprintf("%d seconds", int(olderThan.Seconds())))
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// Counts — агрегаты для /stats и для идемпотентности сида.
func (s *Store) Counts() (objects, photos, matches int, err error) {
	if err = s.db.QueryRow(`SELECT count(*) FROM objects`).Scan(&objects); err != nil {
		return
	}
	if err = s.db.QueryRow(`SELECT count(*) FROM photos`).Scan(&photos); err != nil {
		return
	}
	err = s.db.QueryRow(`SELECT count(*) FROM matches`).Scan(&matches)
	return
}

func (s *Store) QueueDepth() (int, error) {
	var n int
	err := s.db.QueryRow(`SELECT count(*) FROM ingest_jobs WHERE status IN ('pending','processing')`).Scan(&n)
	return n, err
}

// AllPhotoCounts — число фото по каждому объекту (для обзора портфолио).
func (s *Store) AllPhotoCounts() (map[string]int, error) {
	rows, err := s.db.Query(`SELECT object_id, count(*) FROM photos WHERE object_id IS NOT NULL GROUP BY object_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	m := map[string]int{}
	for rows.Next() {
		var id string
		var c int
		if err := rows.Scan(&id, &c); err != nil {
			return nil, err
		}
		m[id] = c
	}
	return m, rows.Err()
}

// CountPhotoMatchesBetween — сколько фото совпало между двумя объектами (для object_context-скоринга).
func (s *Store) CountPhotoMatchesBetween(objA, objB string) (int, error) {
	var n int
	err := s.db.QueryRow(`SELECT count(*) FROM matches
		WHERE similarity_type IN ('exact','perceptual')
		  AND ((object_id_1=$1 AND object_id_2=$2) OR (object_id_1=$2 AND object_id_2=$1))`, objA, objB).Scan(&n)
	return n, err
}

func nullStr(s string) interface{} {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}
