package auth

import (
	"database/sql"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

func TestTokenStore_Load(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	store, err := NewTokenStore(db)
	if err != nil {
		t.Fatal(err)
	}
	tok := &oauth2.Token{AccessToken: "a", RefreshToken: "r", Expiry: time.Unix(100, 0)}
	if err := store.Save(tok); err != nil {
		t.Fatal(err)
	}
	var id int
	if err := db.QueryRow(`SELECT id FROM tokens`).Scan(&id); err != nil {
		t.Fatal(err)
	}
	got, err := store.Load(id)
	if err != nil {
		t.Fatal(err)
	}
	if got.AccessToken != "a" || got.RefreshToken != "r" || !got.Expiry.Equal(time.Unix(100, 0)) {
		t.Fatalf("unexpected token: %+v", got)
	}
}

func TestTokenStore_Delete(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	store, err := NewTokenStore(db)
	if err != nil {
		t.Fatal(err)
	}
	tok := &oauth2.Token{AccessToken: "a", RefreshToken: "r", Expiry: time.Now()}
	if err := store.Save(tok); err != nil {
		t.Fatal(err)
	}
	var id int
	if err := db.QueryRow(`SELECT id FROM tokens`).Scan(&id); err != nil {
		t.Fatal(err)
	}
	if err := store.Delete(id); err != nil {
		t.Fatal(err)
	}
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM tokens`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("expected 0 tokens, got %d", count)
	}
}
