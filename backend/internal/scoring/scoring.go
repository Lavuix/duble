// Package scoring считает Duplicate Risk Score (0..100) по факторам ТЗ §6.2.
// Решение объяснимо: каждый фактор возвращается с подписью и вкладом в балл (ТЗ §11).
package scoring

import (
	"fmt"

	"github.com/roof/dupes/internal/matcher"
	"github.com/roof/dupes/internal/storage"
)

// Config — конфигурируемые пороги/веса (ТЗ §11: настройка без релиза кода).
type Config struct {
	PhashThreshold int // макс. Hamming для «почти-дубля»
}

func clamp(v int) int {
	if v > 100 {
		return 100
	}
	if v < 0 {
		return 0
	}
	return v
}

// Band возвращает классификацию и действие по диапазонам ТЗ §6.2.
func Band(score int) (band, action string) {
	switch {
	case score >= 86:
		return "Почти точный дубль / высокий риск", "Скрыть до решения модератора"
	case score >= 61:
		return "Вероятный дубль", "Приоритетная очередь модерации"
	case score >= 31:
		return "Требуется проверка", "Показать модератору карточку сравнения"
	default:
		return "Низкий риск", "Стандартный поток модерации"
	}
}

// ScorePhotoPair оценивает пару фотографий. new — только что загруженное фото, other — найденный кандидат.
func (c Config) ScorePhotoPair(new, other storage.Photo, oNew, oOther *storage.Object) (score int, simType string, distance int, factors []storage.Factor) {
	add := func(code, label string, pts int) {
		factors = append(factors, storage.Factor{Code: code, Label: label, Points: pts})
		score += pts
	}

	exact := new.SHA256 != "" && new.SHA256 == other.SHA256
	if exact {
		simType, distance = "exact", 0
		add("EXACT_FILE", "Точное совпадение файла (SHA-256)", 100)
	} else {
		simType = "perceptual"
		distance = hammingMin(new, other)
		switch {
		case distance <= 2:
			add("PERCEPTUAL_HASH", fmt.Sprintf("Очень близкий перцептивный хеш (расстояние %d)", distance), 80)
		case distance <= 6:
			add("PERCEPTUAL_HASH", fmt.Sprintf("Близкий перцептивный хеш (расстояние %d)", distance), 60)
		default:
			add("PERCEPTUAL_HASH", fmt.Sprintf("Похожий перцептивный хеш (расстояние %d)", distance), 40)
		}
	}

	// Контекстные сигналы поверх визуального совпадения.
	if new.ContractorID != "" && other.ContractorID != "" && new.ContractorID != other.ContractorID {
		add("CROSS_CONTRACTOR", "Совпавшее фото у другого подрядчика", 20)
	}
	if !new.UploadDate.IsZero() && !other.UploadDate.IsZero() &&
		other.UploadDate.Before(new.UploadDate) && new.ContractorID != other.ContractorID {
		add("UPLOADED_EARLIER_ELSEWHERE", "Это фото раньше загрузил другой подрядчик", 20)
	}
	if oNew != nil && oOther != nil && oNew.Address != "" && oOther.Address != "" &&
		matcher.SimilarityRatio(oNew.Address, oOther.Address) < 0.6 {
		add("DIFFERENT_ADDRESS", "Разные адреса объектов при совпавших фото", 20)
	}
	if new.ExifDateTime == "" && new.ExifGPS == "" {
		add("NO_EXIF", "У загруженного фото отсутствует EXIF", 8)
	}

	return clamp(score), simType, distance, factors
}

// ScoreObjectPair оценивает пару карточек объектов (дубль данных без обязательного фото-совпадения).
// sharedPhotos — число фото, совпавших между объектами по хешам.
func (c Config) ScoreObjectPair(a, b storage.Object, sharedPhotos int) (score int, factors []storage.Factor) {
	add := func(code, label string, pts int) {
		factors = append(factors, storage.Factor{Code: code, Label: label, Points: pts})
		score += pts
	}

	if a.RoofID != "" && a.RoofID == b.RoofID {
		add("SAME_ROOF_ID", fmt.Sprintf("Один и тот же объект ROOF (id=%s)", a.RoofID), 60)
	}
	if r := matcher.SimilarityRatio(a.Name, b.Name); r >= 0.9 {
		add("SAME_NAME", "Совпадает название объекта", 30)
	} else if r >= 0.7 {
		add("SIMILAR_NAME", "Похожее название объекта", 20)
	}
	if a.Address != "" && b.Address != "" && matcher.SimilarityRatio(a.Address, b.Address) >= 0.9 {
		add("SAME_ADDRESS", "Совпадает адрес объекта", 25)
	}
	if a.ContractorID != "" && a.ContractorID == b.ContractorID {
		add("SAME_CONTRACTOR", "Один и тот же подрядчик (возможен повторный ввод карточки)", 10)
	} else if sharedPhotos > 0 {
		add("CROSS_CONTRACTOR_OBJECT", "Объекты разных подрядчиков делят одни фото", 30)
	}
	if a.Year != "" && a.Year == b.Year {
		add("SAME_YEAR", "Совпадает год выполнения", 5)
	}
	if d := matcher.GPSDistanceMeters(a.GPS, b.GPS); d >= 0 && d < 200 {
		add("SAME_GPS", fmt.Sprintf("Координаты совпадают (≈%.0f м)", d), 15)
	}
	switch {
	case sharedPhotos >= 5:
		add("SHARED_PHOTOS", fmt.Sprintf("Совпало %d фото между объектами", sharedPhotos), 50)
	case sharedPhotos >= 2:
		add("SHARED_PHOTOS", fmt.Sprintf("Совпало %d фото между объектами", sharedPhotos), 35)
	case sharedPhotos == 1:
		add("SHARED_PHOTOS", "Совпало 1 фото между объектами", 20)
	}

	return clamp(score), factors
}

func hammingMin(a, b storage.Photo) int {
	d := ham(a.PHash, b.PHash)
	if x := ham(a.DHash, b.DHash); x < d {
		d = x
	}
	return d
}

func ham(a, b int64) int {
	x := uint64(a) ^ uint64(b)
	c := 0
	for x != 0 {
		c++
		x &= x - 1
	}
	return c
}
