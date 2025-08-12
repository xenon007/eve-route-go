package auth

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

// Manager handles session storage using a cookie store.
type Manager struct {
	store *sessions.CookieStore
}

// NewManager creates a new session manager with secret from SESSION_KEY env.
func NewManager() *Manager {
	key := os.Getenv("SESSION_KEY")
	if key == "" {
		log.Println("SESSION_KEY is not set")
	}
	store := sessions.NewCookieStore([]byte(key))
	return &Manager{store: store}
}

// Get returns a session for the given request and name.
func (m *Manager) Get(r *http.Request, name string) (*sessions.Session, error) {
	return m.store.Get(r, name)
}
