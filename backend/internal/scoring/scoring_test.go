package scoring

import (
	"testing"
	"time"

	"github.com/roof/dupes/internal/storage"
)

func TestExactFilePairIsHighRisk(t *testing.T) {
	cfg := Config{PhashThreshold: 10}
	newP := storage.Photo{SHA256: "abc", ContractorID: "c1", UploadDate: time.Now()}
	oldP := storage.Photo{SHA256: "abc", ContractorID: "c2", UploadDate: time.Now().Add(-time.Hour)}
	oNew := &storage.Object{Address: "Москва"}
	oOld := &storage.Object{Address: "Казань"}

	score, simType, dist, factors := cfg.ScorePhotoPair(newP, oldP, oNew, oOld)
	if simType != "exact" || dist != 0 {
		t.Fatalf("ожидался exact/0, получили %s/%d", simType, dist)
	}
	if score < 86 {
		t.Fatalf("точный дубль у другого подрядчика должен быть высоким риском, got %d", score)
	}
	if len(factors) == 0 {
		t.Fatal("должны быть факторы объяснения")
	}
}

func TestSameRoofIDObjectPair(t *testing.T) {
	cfg := Config{PhashThreshold: 10}
	a := storage.Object{RoofID: "84022", Name: "ТЦ Опалиха", ContractorID: "c1", Year: "2024"}
	b := storage.Object{RoofID: "84022", Name: "ТЦ Опалиха", ContractorID: "c1", Year: "2024"}
	score, factors := cfg.ScoreObjectPair(a, b, 3)
	if score < 61 {
		t.Fatalf("один и тот же объект ROOF должен давать высокий риск, got %d", score)
	}
	if len(factors) == 0 {
		t.Fatal("должны быть факторы")
	}
}

func TestBands(t *testing.T) {
	cases := map[int]string{10: "Низкий риск", 45: "Требуется проверка", 70: "Вероятный дубль", 95: "Почти точный дубль / высокий риск"}
	for score, want := range cases {
		if band, _ := Band(score); band != want {
			t.Errorf("score %d: ожидали %q, получили %q", score, want, band)
		}
	}
}
