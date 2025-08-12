package auth

import (
	"database/sql"
	"log"
	"time"

	"golang.org/x/oauth2"

	_ "modernc.org/sqlite"
)

// TokenStore persists OAuth2 tokens.
type TokenStore struct {
	db *sql.DB
}

// NewTokenStore returns a TokenStore and ensures schema.
func NewTokenStore(db *sql.DB) (*TokenStore, error) {
	s := &TokenStore{db: db}
	if err := s.init(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *TokenStore) init() error {
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS tokens (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        access_token TEXT,
        refresh_token TEXT,
        expiry INTEGER
    )`)
	if err != nil {
		log.Println("create table:", err)
	}
	return err
}

// Save stores the token in the database.
func (s *TokenStore) Save(t *oauth2.Token) error {
	_, err := s.db.Exec(`INSERT INTO tokens(access_token, refresh_token, expiry) VALUES(?,?,?)`,
		t.AccessToken, t.RefreshToken, t.Expiry.Unix())
	return err
}

// Load retrieves a token by its id.
func (s *TokenStore) Load(id int) (*oauth2.Token, error) {
	row := s.db.QueryRow(`SELECT access_token, refresh_token, expiry FROM tokens WHERE id=?`, id)
	var access, refresh string
	var expiry int64
	if err := row.Scan(&access, &refresh, &expiry); err != nil {
		return nil, err
	}
	tok := &oauth2.Token{
		AccessToken:  access,
		RefreshToken: refresh,
		Expiry:       time.Unix(expiry, 0),
	}
	return tok, nil
}

// Delete removes a token by its id.
func (s *TokenStore) Delete(id int) error {
	_, err := s.db.Exec(`DELETE FROM tokens WHERE id=?`, id)
	return err
}
