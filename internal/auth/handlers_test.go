package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/oauth2"
)

func TestLogin(t *testing.T) {
	conf := &oauth2.Config{Endpoint: oauth2.Endpoint{AuthURL: "http://example.com/auth"}}
	h := NewHandler(conf, &TokenStore{})
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	rr := httptest.NewRecorder()
	h.Login(rr, req)
	if rr.Code != http.StatusFound {
		t.Fatalf("expected %d, got %d", http.StatusFound, rr.Code)
	}
	loc := rr.Header().Get("Location")
	if loc == "" {
		t.Fatal("missing Location header")
	}
}

func TestCallback(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token":  "abc",
			"token_type":    "Bearer",
			"refresh_token": "ref",
			"expires_in":    3600,
		})
	}))
	defer ts.Close()

	conf := &oauth2.Config{
		ClientID:     "id",
		ClientSecret: "secret",
		RedirectURL:  "http://localhost/callback",
		Endpoint: oauth2.Endpoint{
			AuthURL:  ts.URL,
			TokenURL: ts.URL,
		},
	}

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	store, err := NewTokenStore(db)
	if err != nil {
		t.Fatal(err)
	}

	h := NewHandler(conf, store)
	req := httptest.NewRequest(http.MethodGet, "/callback?code=1", nil)
	rr := httptest.NewRecorder()
	h.Callback(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rr.Code)
	}
	var access, refresh string
	err = db.QueryRow("SELECT access_token, refresh_token FROM tokens").Scan(&access, &refresh)
	if err != nil {
		t.Fatal(err)
	}
	if access != "abc" || refresh != "ref" {
		t.Fatalf("unexpected tokens: %s %s", access, refresh)
	}
}
