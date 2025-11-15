package notes

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"context"

	"github.com/gorilla/mux"
)

type Handler struct {
	repo Repository
	logger *log.Logger
}

func NewHandler(r Repository, l *log.Logger) *Handler {
	return &Handler{repo: r, logger: l}
}

func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var input struct { Title string `json:"title"`; Content string `json:"content"` }
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if input.Title == "" || input.Content == "" {
		writeError(w, http.StatusBadRequest, "title and content required")
		return
	}
	n := &Note{Title: input.Title, Content: input.Content}
	if err := h.repo.Create(r.Context(), n); err != nil {
		writeError(w, http.StatusInternalServerError, "db error")
		return
	}
	writeJSON(w, http.StatusCreated, n)
}

func (h *Handler) GetNote(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	n, err := h.repo.Get(context.Background(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "db error")
		return
	}
	if n == nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	writeJSON(w, http.StatusOK, n)
}

func (h *Handler) ListNotes(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "db error")
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}

func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}
