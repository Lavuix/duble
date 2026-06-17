-- 0001_init: базовая схема системы поиска дублей (forward-only, без CREATE EXTENSION).
-- Модель данных по ТЗ §7. object_id — внутренний ключ карточки портфолио (напр. "row34"),
-- roof_id — id объекта в ROOF, который может повторяться (это сигнал дубля).

CREATE TABLE IF NOT EXISTS contractors (
    contractor_id TEXT PRIMARY KEY,
    name          TEXT NOT NULL DEFAULT '',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS objects (
    object_id     TEXT PRIMARY KEY,            -- внутренний ключ карточки (row_ref)
    roof_id       TEXT NOT NULL DEFAULT '',    -- id объекта в ROOF (может повторяться)
    contractor_id TEXT REFERENCES contractors(contractor_id),
    name          TEXT NOT NULL DEFAULT '',
    address       TEXT NOT NULL DEFAULT '',
    gps           TEXT NOT NULL DEFAULT '',
    customer      TEXT NOT NULL DEFAULT '',
    year          TEXT NOT NULL DEFAULT '',
    materials     TEXT NOT NULL DEFAULT '',
    system        TEXT NOT NULL DEFAULT '',
    region        TEXT NOT NULL DEFAULT '',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_objects_roof_id ON objects(roof_id);

CREATE TABLE IF NOT EXISTS photos (
    photo_id          TEXT PRIMARY KEY,
    contractor_id     TEXT,
    object_id         TEXT REFERENCES objects(object_id),
    upload_date       TIMESTAMPTZ NOT NULL DEFAULT now(),
    original_filename TEXT NOT NULL DEFAULT '',
    source_url        TEXT NOT NULL DEFAULT '',
    blob_path         TEXT NOT NULL DEFAULT '',
    file_hash_sha256  TEXT NOT NULL DEFAULT '',
    phash             BIGINT NOT NULL DEFAULT 0,
    dhash             BIGINT NOT NULL DEFAULT 0,
    ahash             BIGINT NOT NULL DEFAULT 0,
    width             INT NOT NULL DEFAULT 0,
    height            INT NOT NULL DEFAULT 0,
    file_size         BIGINT NOT NULL DEFAULT 0,
    exif_camera       TEXT NOT NULL DEFAULT '',
    exif_datetime     TEXT NOT NULL DEFAULT '',
    exif_gps          TEXT NOT NULL DEFAULT '',
    moderation_status TEXT NOT NULL DEFAULT 'new',
    duplicate_score   INT NOT NULL DEFAULT 0,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_photos_sha256 ON photos(file_hash_sha256);
CREATE INDEX IF NOT EXISTS idx_photos_object ON photos(object_id);

CREATE TABLE IF NOT EXISTS matches (
    match_id               TEXT PRIMARY KEY,
    pair_key               TEXT NOT NULL UNIQUE,   -- дедуп пары (photo/object)
    photo_id_1             TEXT,
    photo_id_2             TEXT,
    object_id_1            TEXT,
    object_id_2            TEXT,
    contractor_id_1        TEXT,
    contractor_id_2        TEXT,
    similarity_type        TEXT NOT NULL,          -- exact | perceptual | object_context
    hash_distance          INT NOT NULL DEFAULT 0,
    embedding_similarity   DOUBLE PRECISION,       -- зарезервировано (без ИИ = NULL)
    first_uploaded_photo_id TEXT,
    risk_score             INT NOT NULL DEFAULT 0,
    risk_factors           JSONB NOT NULL DEFAULT '[]'::jsonb,
    status                 TEXT NOT NULL DEFAULT 'pending',
    moderator_comment      TEXT NOT NULL DEFAULT '',
    created_at             TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_matches_status ON matches(status);
CREATE INDEX IF NOT EXISTS idx_matches_score ON matches(risk_score);
CREATE INDEX IF NOT EXISTS idx_matches_type ON matches(similarity_type);

CREATE TABLE IF NOT EXISTS ingest_jobs (
    job_id     TEXT PRIMARY KEY,
    photo_id   TEXT NOT NULL,
    status     TEXT NOT NULL DEFAULT 'pending',   -- pending | processing | done | failed
    attempts   INT NOT NULL DEFAULT 0,
    error      TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    started_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_jobs_status ON ingest_jobs(status);

CREATE TABLE IF NOT EXISTS settings (
    key   TEXT PRIMARY KEY,
    value TEXT NOT NULL
);
