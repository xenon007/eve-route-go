package auth

import (
	"log"
	"net/http"
	"strconv"

	"golang.org/x/oauth2"
)

// Handler provides OAuth2 login and callback handlers.
type Handler struct {
	Config *oauth2.Config
	Store  *TokenStore
}

// NewHandler creates a new Handler.
func NewHandler(config *oauth2.Config, store *TokenStore) *Handler {
	return &Handler{Config: config, Store: store}
}

// Login redirects the user to the OAuth2 provider.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	url := h.Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}

// Callback exchanges the code for a token and stores it.
func (h *Handler) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "missing code", http.StatusBadRequest)
		return
	}
	token, err := h.Config.Exchange(r.Context(), code)
	if err != nil {
		log.Println("token exchange failed:", err)
		http.Error(w, "exchange failed", http.StatusInternalServerError)
		return
	}
	if err := h.Store.Save(token); err != nil {
		log.Println("token save failed:", err)
		http.Error(w, "save failed", http.StatusInternalServerError)
		return
	}
	if _, err := h.Store.Load(1); err != nil {
		log.Println("token load failed:", err)
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

// Logout removes the token for the given user id.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.Store.Delete(id); err != nil {
		log.Println("token delete failed:", err)
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("logged out"))
}
