package ext

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"

	"github.com/roof/dupes/internal/httpmw"
)

// Router собирает все роуты с middleware (request-id, recover, метрики по шаблону роута).
func Router(h *Handlers, log zerolog.Logger) http.Handler {
	r := mux.NewRouter()

	// Метрики отдаются на корне (observability §2).
	r.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)

	api := r.PathPrefix("/api/v1/ext").Subrouter()
	route := func(method, path, name string, fn http.HandlerFunc) {
		api.Handle(path, httpmw.Metrics(name, fn)).Methods(method)
	}

	route(http.MethodGet, "/health", "/health", h.Health)
	route(http.MethodGet, "/stats", "/stats", h.Stats)
	route(http.MethodPost, "/contractors", "/contractors", h.UpsertContractor)
	route(http.MethodGet, "/objects", "/objects", h.ListObjects)
	route(http.MethodPost, "/objects", "/objects", h.UpsertObject)
	route(http.MethodPost, "/objects:recheck", "/objects:recheck", h.RecheckObjects)
	route(http.MethodPost, "/objects:setactive", "/objects:setactive", h.SetObjectsActive)
	route(http.MethodPost, "/import/xlsx", "/import/xlsx", h.ImportXLSX)
	route(http.MethodPost, "/data:reset", "/data:reset", h.ResetData)
	route(http.MethodPost, "/reindex", "/reindex", h.Reindex)
	route(http.MethodPost, "/photos:ingest", "/photos:ingest", h.IngestPhoto)
	route(http.MethodGet, "/photos/{id}/content", "/photos/{id}/content", h.PhotoContent)
	route(http.MethodGet, "/matches", "/matches", h.ListMatches)
	route(http.MethodGet, "/matches/{id}", "/matches/{id}", h.MatchDetail)
	route(http.MethodPost, "/matches/{id}/decision", "/matches/{id}/decision", h.DecideMatch)

	// Общие middleware на всё дерево.
	return httpmw.RequestID(httpmw.Recover(log)(r))
}
