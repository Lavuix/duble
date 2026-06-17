package matcher

import "testing"

func TestSimilarityRatio(t *testing.T) {
	if r := SimilarityRatio("ТЦ Опалиха", "тц  опалиха"); r < 0.99 {
		t.Fatalf("нормализованные одинаковые строки должны давать ~1, got %f", r)
	}
	if r := SimilarityRatio("Москва", "Казань"); r > 0.5 {
		t.Fatalf("разные строки не должны быть похожи, got %f", r)
	}
}

func TestGPSDistance(t *testing.T) {
	if d := GPSDistanceMeters("55.855078,38.570047", "55.855078,38.570047"); d < 0 || d > 1 {
		t.Fatalf("одинаковые координаты ≈ 0 м, got %f", d)
	}
	if d := GPSDistanceMeters("bad", "55.0,38.0"); d != -1 {
		t.Fatalf("битые координаты → -1, got %f", d)
	}
}
