// Package matcher — сравнение карточек объектов без ИИ: нормализация строк,
// расстояние Левенштейна, близость GPS (ТЗ §4 «один объект с другого ракурса», §6.1 п.5).
package matcher

import (
	"math"
	"strconv"
	"strings"
)

func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(s), 64)
}

// Normalize приводит строку к сравнимому виду: нижний регистр, схлопнутые пробелы.
func Normalize(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	return strings.Join(strings.Fields(s), " ")
}

// SimilarityRatio возвращает близость двух строк в диапазоне [0..1] на базе Левенштейна.
func SimilarityRatio(a, b string) float64 {
	a, b = Normalize(a), Normalize(b)
	if a == "" && b == "" {
		return 1
	}
	if a == "" || b == "" {
		return 0
	}
	d := levenshtein([]rune(a), []rune(b))
	maxLen := len(a)
	if len(b) > maxLen {
		maxLen = len(b)
	}
	maxRunes := math.Max(float64(len([]rune(a))), float64(len([]rune(b))))
	return 1 - float64(d)/maxRunes
}

func levenshtein(a, b []rune) int {
	prev := make([]int, len(b)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(a); i++ {
		cur := make([]int, len(b)+1)
		cur[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			cur[j] = min3(cur[j-1]+1, prev[j]+1, prev[j-1]+cost)
		}
		prev = cur
	}
	return prev[len(b)]
}

func min3(a, b, c int) int {
	m := a
	if b < m {
		m = b
	}
	if c < m {
		m = c
	}
	return m
}

// GPSDistanceMeters парсит "lat,lon" и считает расстояние по формуле гаверсинуса.
// Возвращает -1, если хотя бы одна координата не распознана.
func GPSDistanceMeters(a, b string) float64 {
	la, loa, ok1 := parseLatLon(a)
	lb, lob, ok2 := parseLatLon(b)
	if !ok1 || !ok2 {
		return -1
	}
	const R = 6371000.0
	dLat := rad(lb - la)
	dLon := rad(lob - loa)
	h := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(rad(la))*math.Cos(rad(lb))*math.Sin(dLon/2)*math.Sin(dLon/2)
	return R * 2 * math.Atan2(math.Sqrt(h), math.Sqrt(1-h))
}

func parseLatLon(s string) (lat, lon float64, ok bool) {
	parts := strings.Split(strings.TrimSpace(s), ",")
	if len(parts) != 2 {
		return 0, 0, false
	}
	la, err1 := parseFloat(parts[0])
	lo, err2 := parseFloat(parts[1])
	if err1 != nil || err2 != nil {
		return 0, 0, false
	}
	return la, lo, true
}

func rad(d float64) float64 { return d * math.Pi / 180 }
