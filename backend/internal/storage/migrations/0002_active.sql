-- 0002: флаг активности карточки объекта (модератор может деактивировать объект, ТЗ §12).
ALTER TABLE objects ADD COLUMN IF NOT EXISTS active BOOLEAN NOT NULL DEFAULT true;
