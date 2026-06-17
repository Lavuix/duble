// Package hashing считает технические признаки фото без ИИ:
// точный хеш SHA-256, перцептивные pHash/dHash/aHash и EXIF (ТЗ §4, §6.1).
package hashing

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	"math/bits"

	// регистрируем декодеры форматов из ТЗ §2 (.jpg/.jpeg/.jfif → jpeg, .png, .webp)
	_ "image/jpeg"
	_ "image/png"

	"github.com/corona10/goimagehash"
	"github.com/rwcarlsen/goexif/exif"
	_ "golang.org/x/image/webp"
)

type Features struct {
	SHA256       string
	PHash        int64
	DHash        int64
	AHash        int64
	Width        int
	Height       int
	Size         int64
	ExifCamera   string
	ExifDateTime string
	ExifGPS      string
}

// Compute декодирует изображение и считает все признаки. Ошибка = «не картинка / битый файл»
// (ТЗ §11: это не дубль, а ошибка обработки — логируется отдельно вызывающим кодом).
func Compute(data []byte) (*Features, error) {
	f := &Features{Size: int64(len(data))}
	sum := sha256.Sum256(data)
	f.SHA256 = hex.EncodeToString(sum[:])

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}
	b := img.Bounds()
	f.Width, f.Height = b.Dx(), b.Dy()

	if p, e := goimagehash.PerceptionHash(img); e == nil {
		f.PHash = int64(p.GetHash())
	}
	if d, e := goimagehash.DifferenceHash(img); e == nil {
		f.DHash = int64(d.GetHash())
	}
	if a, e := goimagehash.AverageHash(img); e == nil {
		f.AHash = int64(a.GetHash())
	}

	// EXIF: best-effort, обычно есть только у jpeg.
	if x, e := exif.Decode(bytes.NewReader(data)); e == nil {
		if t, e := x.Get(exif.Model); e == nil {
			f.ExifCamera, _ = t.StringVal()
		}
		if t, e := x.Get(exif.DateTimeOriginal); e == nil {
			f.ExifDateTime, _ = t.StringVal()
		}
		if lat, lon, e := x.LatLong(); e == nil {
			f.ExifGPS = fmt.Sprintf("%.6f,%.6f", lat, lon)
		}
	}
	return f, nil
}

// Hamming — расстояние между двумя перцептивными хешами (число различных бит).
func Hamming(a, b int64) int {
	return bits.OnesCount64(uint64(a) ^ uint64(b))
}
