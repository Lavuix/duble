// Команда seed загружает тестовые данные в систему: метаданные объектов из xlsx
// и реальные фотографии из downloaded_photos. Идёт через публичный API бэкенда.
package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

var (
	backend   = env("BACKEND_URL", "http://backend:8080")
	photosDir = env("PHOTOS_DIR", "/seed/photos")
	xlsxPath  = env("XLSX_PATH", "/seed/table.xlsx")
)

// row{N}_{roofId}_{slug}
var folderRe = regexp.MustCompile(`^row(\d+)_(\d+)_(.*)$`)
var bracketID = regexp.MustCompile(`\[(\d+)\]`)
var urlRe = regexp.MustCompile(`https?://\S+`)

type meta struct {
	roofID, name, address, gps, customer, year, materials, system, region string
	contractorRaw                                                         string
	photoURLs                                                             []string
}

func main() {
	waitHealthy()

	// Идемпотентность: если данные уже загружены — выходим (one-shot повторный запуск).
	if st := stats(); st["photos"] > 0 {
		log.Printf("seed: данные уже загружены (photos=%d), пропускаю", st["photos"])
		return
	}

	rowsMeta := readXLSX()

	entries, err := os.ReadDir(photosDir)
	if err != nil {
		log.Fatalf("seed: не читается %s: %v", photosDir, err)
	}
	totalPhotos := 0
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		m := folderRe.FindStringSubmatch(e.Name())
		if m == nil {
			continue
		}
		rowNum, roofID, slug := m[1], m[2], m[3]
		objectID := "row" + rowNum
		// Метаданные берём по roof_id из папки: номера строк в xlsx могут сдвигаться
		// между версиями таблицы, а roof_id стабилен (папки именованы по V1-строкам).
		md := rowsMeta[roofID]
		if md.roofID == "" {
			md.roofID = roofID
		}
		if md.name == "" {
			md.name = strings.ReplaceAll(slug, "_", " ")
		}

		contractorID, contractorName := parseContractor(md.contractorRaw, objectID)
		upsertContractor(contractorID, contractorName)
		upsertObject(map[string]string{
			"object_id": objectID, "roof_id": md.roofID, "contractor_id": contractorID,
			"name": md.name, "address": md.address, "gps": md.gps, "customer": md.customer,
			"year": md.year, "materials": md.materials, "system": md.system, "region": md.region,
		})

		files := imageFiles(filepath.Join(photosDir, e.Name()))
		for i, fpath := range files {
			srcURL := ""
			if i < len(md.photoURLs) {
				srcURL = md.photoURLs[i]
			}
			ingestPhoto(objectID, contractorID, fpath, srcURL)
			totalPhotos++
		}
		log.Printf("seed: объект %s (roof=%s) подрядчик=%s, фото=%d", objectID, md.roofID, contractorID, len(files))
	}

	log.Printf("seed: загружено фото=%d, жду обработку очереди...", totalPhotos)
	waitQueueDrained()

	// Пересчёт дублей-данных (object_context) после того, как фото-совпадения посчитаны.
	post("/api/v1/ext/objects:recheck", nil)
	st := stats()
	log.Printf("seed: готово. objects=%d photos=%d matches=%d", st["objects"], st["photos"], st["matches"])
}

func readXLSX() map[string]meta {
	out := map[string]meta{}
	f, err := excelize.OpenFile(xlsxPath)
	if err != nil {
		log.Printf("seed: xlsx не открыт (%v) — продолжаю только по папкам", err)
		return out
	}
	defer f.Close()
	sheet := f.GetSheetList()[0]
	rows, err := f.GetRows(sheet)
	if err != nil {
		log.Printf("seed: не читаются строки xlsx: %v", err)
		return out
	}
	for idx, row := range rows {
		if idx == 0 {
			continue // заголовок
		}
		roofID := cell(row, 5) // колонка F — ID объекта ROOF
		if roofID == "" {
			continue
		}
		photos := urlRe.FindAllString(cell(row, 22), -1)
		out[roofID] = meta{ // ключ — roof_id (см. main: сопоставление папок по roof_id)
			name: cell(row, 0), roofID: roofID, address: cell(row, 6),
			year: cell(row, 9), customer: cell(row, 11), materials: cell(row, 12),
			gps: cell(row, 13), contractorRaw: cell(row, 16), region: cell(row, 17),
			system: cell(row, 18), photoURLs: photos,
		}
	}
	return out
}

func parseContractor(raw, objectID string) (id, name string) {
	raw = strings.TrimSpace(raw)
	if m := bracketID.FindStringSubmatch(raw); m != nil {
		return "c" + m[1], raw
	}
	if raw != "" {
		return "c_" + objectID, raw
	}
	return "c_" + objectID, "Подрядчик " + objectID
}

func imageFiles(dir string) []string {
	entries, _ := os.ReadDir(dir)
	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		n := strings.ToLower(e.Name())
		if strings.HasPrefix(e.Name(), ".") {
			continue
		}
		if strings.HasSuffix(n, ".jpg") || strings.HasSuffix(n, ".jpeg") || strings.HasSuffix(n, ".jfif") ||
			strings.HasSuffix(n, ".png") || strings.HasSuffix(n, ".webp") {
			files = append(files, filepath.Join(dir, e.Name()))
		}
	}
	sort.Strings(files)
	return files
}

// ---- HTTP-клиент к бэкенду ----

func ingestPhoto(objectID, contractorID, path, srcURL string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("seed: не читается %s: %v", path, err)
		return
	}
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	_ = mw.WriteField("object_id", objectID)
	_ = mw.WriteField("contractor_id", contractorID)
	_ = mw.WriteField("original_filename", filepath.Base(path))
	_ = mw.WriteField("source_url", srcURL)
	fw, _ := mw.CreateFormFile("file", filepath.Base(path))
	_, _ = fw.Write(data)
	mw.Close()

	req, _ := http.NewRequest(http.MethodPost, backend+"/api/v1/ext/photos:ingest", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("seed: ingest %s: %v", path, err)
		return
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
}

func upsertContractor(id, name string) {
	post("/api/v1/ext/contractors", map[string]string{"contractor_id": id, "name": name})
}

func upsertObject(o map[string]string) { post("/api/v1/ext/objects", o) }

func post(path string, body interface{}) {
	var r io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		r = bytes.NewReader(b)
	}
	req, _ := http.NewRequest(http.MethodPost, backend+path, r)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("seed: POST %s: %v", path, err)
		return
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
}

func stats() map[string]int {
	out := map[string]int{}
	resp, err := http.Get(backend + "/api/v1/ext/stats")
	if err != nil {
		return out
	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&out)
	return out
}

func waitHealthy() {
	for i := 0; i < 60; i++ {
		resp, err := http.Get(backend + "/api/v1/ext/health")
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()
			return
		}
		if resp != nil {
			resp.Body.Close()
		}
		log.Printf("seed: жду backend... (%d)", i+1)
		time.Sleep(2 * time.Second)
	}
	log.Fatal("seed: backend не поднялся")
}

func waitQueueDrained() {
	for i := 0; i < 120; i++ {
		if stats()["queue_depth"] == 0 {
			return
		}
		time.Sleep(time.Second)
	}
	log.Print("seed: таймаут ожидания очереди, продолжаю")
}

func cell(row []string, idx int) string {
	if idx < len(row) {
		return strings.TrimSpace(row[idx])
	}
	return ""
}

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
