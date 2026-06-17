package ext

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"path"
	"regexp"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"

	"github.com/roof/dupes/internal/storage"
)

var (
	bracketID = regexp.MustCompile(`\[(\d+)\]`)
	urlStart  = regexp.MustCompile(`https?://`)
)

// splitURLs разбивает значение ячейки «Фото» на отдельные ссылки. В выгрузке ROOF
// URL внутри одной ячейки идут подряд без разделителя, поэтому режем по началу каждой "http(s)://".
func splitURLs(s string) []string {
	idx := urlStart.FindAllStringIndex(s, -1)
	var out []string
	for i, loc := range idx {
		end := len(s)
		if i+1 < len(idx) {
			end = idx[i+1][0]
		}
		if u := trim(s[loc[0]:end]); u != "" {
			out = append(out, u)
		}
	}
	return out
}

// ImportXLSX загружает объекты/подрядчиков/фото из выгрузки Excel (ТЗ §2: столбец W «Фото» с URL).
// Фото скачиваются воркером по source_url; объекты и данные доступны сразу.
func (h *Handlers) ImportXLSX(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeErr(w, http.StatusBadRequest, "ожидается multipart-форма с файлом xlsx")
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "нужен файл в поле file")
		return
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		writeErr(w, http.StatusBadRequest, "не удалось открыть xlsx: "+err.Error())
		return
	}
	defer f.Close()
	rows, err := f.GetRows(f.GetSheetList()[0])
	if err != nil {
		writeErr(w, http.StatusBadRequest, "не удалось прочитать строки: "+err.Error())
		return
	}

	seen := map[string]int{}
	objects, photos := 0, 0
	for idx, row := range rows {
		if idx == 0 {
			continue // заголовок
		}
		roofID := cell(row, 5) // F — ID объекта ROOF
		if roofID == "" {
			continue
		}
		seen[roofID]++
		objectID := "imp_" + roofID
		if seen[roofID] > 1 { // один roof_id может встречаться несколько раз — это сигнал дубля
			objectID += "_" + strconv.Itoa(seen[roofID])
		}
		contractorID, contractorName := parseContractor(cell(row, 16), objectID)
		_ = h.Store.UpsertContractor(storage.Contractor{ContractorID: contractorID, Name: contractorName})
		if err := h.Store.UpsertObject(storage.Object{
			ObjectID: objectID, RoofID: roofID, ContractorID: contractorID,
			Name: cell(row, 0), Address: cell(row, 6), GPS: cell(row, 13),
			Customer: cell(row, 11), Year: cell(row, 9), Materials: cell(row, 12),
			System: cell(row, 18), Region: cell(row, 17),
		}); err != nil {
			h.Log.Error().Err(err).Str("object", objectID).Msg("import: upsert object")
			continue
		}
		objects++
		for _, u := range splitURLs(cell(row, 22)) {
			pid := storage.NewID()
			if err := h.Store.InsertPhoto(storage.Photo{
				PhotoID: pid, ContractorID: contractorID, ObjectID: objectID,
				UploadDate: time.Now(), OriginalFilename: path.Base(u), SourceURL: u,
				ModerationStatus: "new",
			}); err != nil {
				continue
			}
			_, _ = h.Store.EnqueueJob(pid)
			photos++
		}
	}

	// Пересчёт дублей-данных после того, как фото-очередь обработается.
	go h.recheckWhenIdle()
	writeJSON(w, http.StatusOK, map[string]int{"objects": objects, "photos": photos})
}

// recheckWhenIdle ждёт опустошения очереди и пересчитывает object_context-совпадения.
func (h *Handlers) recheckWhenIdle() {
	for i := 0; i < 300; i++ {
		if d, err := h.Store.QueueDepth(); err == nil && d == 0 {
			break
		}
		time.Sleep(time.Second)
	}
	if _, err := h.Pipe.RecheckObjects(); err != nil {
		h.Log.Error().Err(err).Msg("import: recheck objects")
	}
}

// SetObjectsActive — массовая (де)активация выбранных объектов (чекбоксы в UI, ТЗ §12).
func (h *Handlers) SetObjectsActive(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ObjectIDs []string `json:"object_ids"`
		Active    bool     `json:"active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || len(body.ObjectIDs) == 0 {
		writeErr(w, http.StatusBadRequest, "нужен непустой object_ids")
		return
	}
	if err := h.Store.SetObjectsActive(body.ObjectIDs, body.Active); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"updated": len(body.ObjectIDs), "active": body.Active})
}

// ResetData полностью очищает базу (кнопка «Очистить базу» в UI).
func (h *Handlers) ResetData(w http.ResponseWriter, r *http.Request) {
	if err := h.Store.Reset(); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "cleared"})
}

func parseContractor(raw, objectID string) (id, name string) {
	if m := bracketID.FindStringSubmatch(raw); m != nil {
		return "c" + m[1], raw
	}
	if raw != "" {
		return "c_" + objectID, raw
	}
	return "c_" + objectID, "Подрядчик " + objectID
}

func cell(row []string, idx int) string {
	if idx < len(row) {
		return trim(row[idx])
	}
	return ""
}

func trim(s string) string {
	// локальная замена strings.TrimSpace, чтобы не тянуть импорт ради одной функции
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\t' || s[0] == '\n' || s[0] == '\r') {
		s = s[1:]
	}
	for len(s) > 0 {
		c := s[len(s)-1]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			s = s[:len(s)-1]
		} else {
			break
		}
	}
	return s
}
