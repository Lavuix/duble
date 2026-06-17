package storage

import "time"

type Contractor struct {
	ContractorID string `json:"contractor_id"`
	Name         string `json:"name"`
}

type Object struct {
	ObjectID     string `json:"object_id"`
	RoofID       string `json:"roof_id"`
	ContractorID string `json:"contractor_id"`
	Name         string `json:"name"`
	Address      string `json:"address"`
	GPS          string `json:"gps"`
	Customer     string `json:"customer"`
	Year         string `json:"year"`
	Materials    string `json:"materials"`
	System       string `json:"system"`
	Region       string `json:"region"`
	Active       bool   `json:"active"`
}

type Photo struct {
	PhotoID          string    `json:"photo_id"`
	ContractorID     string    `json:"contractor_id"`
	ObjectID         string    `json:"object_id"`
	UploadDate       time.Time `json:"upload_date"`
	OriginalFilename string    `json:"original_filename"`
	SourceURL        string    `json:"source_url"`
	BlobPath         string    `json:"blob_path"`
	SHA256           string    `json:"file_hash_sha256"`
	PHash            int64     `json:"phash"`
	DHash            int64     `json:"dhash"`
	AHash            int64     `json:"ahash"`
	Width            int       `json:"width"`
	Height           int       `json:"height"`
	FileSize         int64     `json:"file_size"`
	ExifCamera       string    `json:"exif_camera"`
	ExifDateTime     string    `json:"exif_datetime"`
	ExifGPS          string    `json:"exif_gps"`
	ModerationStatus string    `json:"moderation_status"`
	DuplicateScore   int       `json:"duplicate_score"`
}

// Factor — одна причина срабатывания скоринга (объяснимость, ТЗ §11).
type Factor struct {
	Code   string `json:"code"`
	Label  string `json:"label"`
	Points int    `json:"points"`
}

type Match struct {
	MatchID              string    `json:"match_id"`
	PairKey              string    `json:"-"`
	PhotoID1             string    `json:"photo_id_1"`
	PhotoID2             string    `json:"photo_id_2"`
	ObjectID1            string    `json:"object_id_1"`
	ObjectID2            string    `json:"object_id_2"`
	ContractorID1        string    `json:"contractor_id_1"`
	ContractorID2        string    `json:"contractor_id_2"`
	SimilarityType       string    `json:"similarity_type"`
	HashDistance         int       `json:"hash_distance"`
	FirstUploadedPhotoID string    `json:"first_uploaded_photo_id"`
	RiskScore            int       `json:"risk_score"`
	RiskFactors          []Factor  `json:"risk_factors"`
	Status               string    `json:"status"`
	ModeratorComment     string    `json:"moderator_comment"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type Job struct {
	JobID    string
	PhotoID  string
	Status   string
	Attempts int
}
