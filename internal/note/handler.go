package note

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /healthz", h.handleHealth)
	mux.HandleFunc("GET /notes", h.handleListNotes)
	mux.HandleFunc("POST /notes", h.handleCreateNote)
	mux.HandleFunc("GET /notes/{id}", h.handleGetNote)
	mux.HandleFunc("PUT /notes/{id}", h.handleUpdateNote)
	mux.HandleFunc("DELETE /notes/{id}", h.handleDeleteNote)
}

func (h *Handler) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) handleCreateNote(w http.ResponseWriter, r *http.Request) {
	var req CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if err := validateTitleAndContent(req.Title, req.Content); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.repo.Create(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "create note failed")
		return
	}

	n, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "load note failed")
		return
	}
	writeJSON(w, http.StatusCreated, n)
}

func (h *Handler) handleListNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := h.repo.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "list notes failed")
		return
	}
	writeJSON(w, http.StatusOK, notes)
}

func (h *Handler) handleGetNote(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	n, err := h.repo.GetByID(r.Context(), id)
	if errors.Is(err, ErrNoteNotFound) {
		writeError(w, http.StatusNotFound, "note not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "get note failed")
		return
	}
	writeJSON(w, http.StatusOK, n)
}

func (h *Handler) handleUpdateNote(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	var req UpdateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if err := validateTitleAndContent(req.Title, req.Content); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	err := h.repo.Update(r.Context(), id, req)
	if errors.Is(err, ErrNoteNotFound) {
		writeError(w, http.StatusNotFound, "note not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "update note failed")
		return
	}

	n, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "load note failed")
		return
	}
	writeJSON(w, http.StatusOK, n)
}

func (h *Handler) handleDeleteNote(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	err := h.repo.Delete(r.Context(), id)
	if errors.Is(err, ErrNoteNotFound) {
		writeError(w, http.StatusNotFound, "note not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "delete note failed")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func parseID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	idRaw := r.PathValue("id")
	id, err := strconv.ParseInt(idRaw, 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid note id")
		return 0, false
	}
	return id, true
}

func validateTitleAndContent(title, content string) error {
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)

	if title == "" {
		return errors.New("title is required")
	}
	if len(title) > 120 {
		return errors.New("title is too long (max: 120)")
	}
	if content == "" {
		return errors.New("content is required")
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
