package notes_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	notes "github.com/yourusername/resume-app/internal/notes"
	"github.com/yourusername/resume-app/internal/testutil"
)

func TestCreateAndGetNote(t *testing.T) {
	stub := &testutil.StubRepo{}
	logger := log.New(bytes.NewBuffer(nil), "test", 0)
	h := notes.NewHandler(stub, logger)

	r := mux.NewRouter()
	r.HandleFunc("/notes", h.CreateNote).Methods(http.MethodPost)
	r.HandleFunc("/notes/{id}", h.GetNote).Methods(http.MethodGet)

	payload := map[string]string{"title": "Hello", "content": "World"}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewReader(b))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated { t.Fatalf("expected 201 got %d", w.Code) }

	var created notes.Note
	json.Unmarshal(w.Body.Bytes(), &created)
	if created.ID == 0 { t.Fatalf("expected ID > 0") }

	getReq := httptest.NewRequest(http.MethodGet, "/notes/1", nil)
	getW := httptest.NewRecorder()
	r.ServeHTTP(getW, getReq)
	if getW.Code != http.StatusOK { t.Fatalf("expected 200 got %d", getW.Code) }
}
