package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suhrobdomoiZ/document-search-service/internal/es"
	"github.com/suhrobdomoiZ/document-search-service/internal/service"
	"github.com/suhrobdomoiZ/document-search-service/internal/utils"
)

type documentsHandler struct {
	service *service.DocumentsService
	logger  *slog.Logger
}

func newDocumentsHandler(
	logger *slog.Logger,
	pool *pgxpool.Pool,
	esClient *es.Client,
) *documentsHandler {
	service := service.NewDocumentsService(logger, pool, esClient)

	return &documentsHandler{service: service, logger: logger}
}

func (h *documentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleSearch(w, r)
	case http.MethodDelete:
		h.handleDelete(w, r)
	default:
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (h *documentsHandler) handleSearch(w http.ResponseWriter, r *http.Request) {
	searchText := r.URL.Query().Get("text")
	if searchText == "" {
		http.Error(w, `{"error": "query param 'text' is required"}`, http.StatusBadRequest)

		return
	}

	docs, err := h.service.Search(r.Context(), searchText)
	if err != nil {
		h.logger.Error("handler: search failed", "error", err)
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(docs)
	if err != nil {
		h.logger.Error("handler: failed to encode response", "error", err)
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
	}
}

func (h *documentsHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, `{"error": "document id is required"}`, http.StatusBadRequest)

		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "id must be an integer"}`, http.StatusBadRequest)

		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			http.Error(w, `{"error": "document not found"}`, http.StatusNotFound)

			return
		}

		h.logger.Error("handler: delete failed", "error", err)
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
