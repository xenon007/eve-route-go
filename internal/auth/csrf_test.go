package auth

import (
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/csrf"
)

func TestCSRFMissingToken(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-CSRF-Token", csrf.Token(r))
	})
	mux.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	csrfKey := []byte("01234567890123456789012345678901")
	h := csrf.Protect(csrfKey, csrf.Secure(false))(mux)
	srv := httptest.NewServer(h)
	defer srv.Close()

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	resp, err := client.Get(srv.URL + "/token")
	if err != nil {
		t.Fatalf("GET token: %v", err)
	}
	token := resp.Header.Get("X-CSRF-Token")
	resp.Body.Close()
	if token == "" {
		t.Fatalf("empty token")
	}

	req, _ := http.NewRequest("POST", srv.URL+"/submit", nil)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("POST without token: %v", err)
	}
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", resp.StatusCode)
	}
	resp.Body.Close()
}
