package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestRouter() (*Store, http.Handler) {
	s := NewStore()
	return s, setupRouter(s)
}

func TestHealthCheck(t *testing.T) {
	_, r := newTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestCreateContact_ValidPayload(t *testing.T) {
	_, r := newTestRouter()
	body := `{"name":"Alice","email":"alice@example.com","phone":"11999999999"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/contacts", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	var got Contact
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if got.ID == "" {
		t.Error("expected non-empty ID")
	}
	if got.Name != "Alice" {
		t.Errorf("expected Name=Alice, got %q", got.Name)
	}
	if got.Email != "alice@example.com" {
		t.Errorf("expected Email=alice@example.com, got %q", got.Email)
	}
}

func TestCreateContact_MissingName(t *testing.T) {
	_, r := newTestRouter()
	body := `{"email":"alice@example.com"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/contacts", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreateContact_InvalidEmail(t *testing.T) {
	_, r := newTestRouter()
	body := `{"name":"Alice","email":"not-an-email"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/contacts", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestListContacts_Empty(t *testing.T) {
	_, r := newTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/contacts", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var got []Contact
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if got == nil {
		t.Error("expected empty slice, not null")
	}
	if len(got) != 0 {
		t.Errorf("expected 0 contacts, got %d", len(got))
	}
}

func TestListContacts_AfterCreate(t *testing.T) {
	_, r := newTestRouter()
	for _, name := range []string{"Alice", "Bob"} {
		body := `{"name":"` + name + `","email":"` + name + `@example.com"}`
		req, _ := http.NewRequest(http.MethodPost, "/contacts", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(httptest.NewRecorder(), req)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/contacts", nil)
	r.ServeHTTP(w, req)
	var got []Contact
	json.NewDecoder(w.Body).Decode(&got)
	if len(got) != 2 {
		t.Errorf("expected 2 contacts, got %d", len(got))
	}
}

func TestGetContact_Found(t *testing.T) {
	_, r := newTestRouter()
	body := `{"name":"Alice","email":"alice@example.com"}`
	req, _ := http.NewRequest(http.MethodPost, "/contacts", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	wc := httptest.NewRecorder()
	r.ServeHTTP(wc, req)

	var created Contact
	json.NewDecoder(wc.Body).Decode(&created)

	w := httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/contacts/"+created.ID, nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var got Contact
	json.NewDecoder(w.Body).Decode(&got)
	if got.ID != created.ID {
		t.Errorf("ID mismatch: want %q, got %q", created.ID, got.ID)
	}
}

func TestGetContact_NotFound(t *testing.T) {
	_, r := newTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/contacts/nonexistent", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestUpdateContact_Found(t *testing.T) {
	_, r := newTestRouter()
	body := `{"name":"Alice","email":"alice@example.com"}`
	req, _ := http.NewRequest(http.MethodPost, "/contacts", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	wc := httptest.NewRecorder()
	r.ServeHTTP(wc, req)

	var created Contact
	json.NewDecoder(wc.Body).Decode(&created)

	updated := `{"name":"Alice Updated","email":"alice2@example.com"}`
	w := httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPut, "/contacts/"+created.ID, bytes.NewBufferString(updated))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var got Contact
	json.NewDecoder(w.Body).Decode(&got)
	if got.Name != "Alice Updated" {
		t.Errorf("expected Name=Alice Updated, got %q", got.Name)
	}
	if got.ID != created.ID {
		t.Errorf("ID must not change after update")
	}
}

func TestUpdateContact_NotFound(t *testing.T) {
	_, r := newTestRouter()
	body := `{"name":"Alice","email":"alice@example.com"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/contacts/nonexistent", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestDeleteContact_Found(t *testing.T) {
	_, r := newTestRouter()
	body := `{"name":"Alice","email":"alice@example.com"}`
	req, _ := http.NewRequest(http.MethodPost, "/contacts", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	wc := httptest.NewRecorder()
	r.ServeHTTP(wc, req)

	var created Contact
	json.NewDecoder(wc.Body).Decode(&created)

	w := httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodDelete, "/contacts/"+created.ID, nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func TestDeleteContact_NotFound(t *testing.T) {
	_, r := newTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/contacts/nonexistent", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
