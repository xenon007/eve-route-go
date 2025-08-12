package main

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/oauth2"

	"github.com/tkhamez/eve-route-go/internal/api"
	"github.com/tkhamez/eve-route-go/internal/auth"
	"github.com/tkhamez/eve-route-go/internal/capital"
	"github.com/tkhamez/eve-route-go/internal/config"
	"github.com/tkhamez/eve-route-go/internal/db"
	routepkg "github.com/tkhamez/eve-route-go/internal/route"
)

//go:embed frontend/dist
var frontendFS embed.FS

// mustEnv возвращает значение переменной окружения или завершает программу.
func mustEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("%s is not set", key)
	}
	return val
}

// initStore инициализирует хранилище данных на основе DATABASE_URL.
func initStore(ctx context.Context, urlStr string) db.Store {
	if urlStr == "" {
		log.Println("DATABASE_URL not set, using in-memory store")
		return db.NewMemory(nil, nil, capital.DefaultSystems())
	}
	u, err := url.Parse(urlStr)
	if err != nil {
		log.Printf("invalid DATABASE_URL %q: %v; using in-memory store", urlStr, err)
		return db.NewMemory(nil, nil, capital.DefaultSystems())
	}
	switch u.Scheme {
	case "postgres", "postgresql":
		conn, err := sql.Open("postgres", urlStr)
		if err != nil {
			log.Printf("postgres connection error: %v; using in-memory store", err)
			return db.NewMemory(nil, nil, capital.DefaultSystems())
		}
		return db.NewPostgres(conn)
	case "mongodb", "mongo":
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(urlStr))
		if err != nil {
			log.Printf("mongo connection error: %v; using in-memory store", err)
			return db.NewMemory(nil, nil, capital.DefaultSystems())
		}
		dbName := strings.TrimPrefix(u.Path, "/")
		if dbName == "" {
			dbName = "eve"
		}
		return db.NewMongo(client, dbName)
	case "sqlite":
		path := strings.TrimPrefix(u.Path, "/")
		if path == "" {
			log.Println("sqlite DATABASE_URL missing path, using in-memory store")
			return db.NewMemory(nil, nil, capital.DefaultSystems())
		}
		conn, err := sql.Open("sqlite", path)
		if err != nil {
			log.Printf("sqlite connection error: %v; using in-memory store", err)
			return db.NewMemory(nil, nil, capital.DefaultSystems())
		}
		return db.NewSQLite(conn)
	default:
		log.Printf("unsupported DATABASE_URL scheme %q, using in-memory store", u.Scheme)
		return db.NewMemory(nil, nil, capital.DefaultSystems())
	}
}

func main() {
	ctx := context.Background()
	cfg := config.FromEnv()

	store := initStore(ctx, cfg.DatabaseURL)

	tokenDB, err := sql.Open("sqlite", "tokens.db")
	if err != nil {
		log.Fatal(err)
	}
	defer tokenDB.Close()
	tokenStore, err := auth.NewTokenStore(tokenDB)
	if err != nil {
		log.Fatal(err)
	}

	oauthConf := &oauth2.Config{
		ClientID:     mustEnv("OAUTH_CLIENT_ID"),
		ClientSecret: mustEnv("OAUTH_CLIENT_SECRET"),
		RedirectURL:  mustEnv("REDIRECT_URL"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.eveonline.com/v2/oauth/authorize",
			TokenURL: "https://login.eveonline.com/v2/oauth/token",
		},
	}
	h := auth.NewHandler(oauthConf, tokenStore)

	r := mux.NewRouter()
	r.HandleFunc("/login", h.Login).Methods("GET")
	r.HandleFunc("/callback", h.Callback).Methods("GET")

	api.RegisterAnsiblexRoutes(r, mustEnv("API_SECRET"))

	mustEnv("SESSION_KEY")
	auth.NewManager()

	planner, err := capital.NewPlanner(store, 5)
	if err != nil {
		log.Fatalf("cannot create planner: %v", err)
	}
	r.HandleFunc("/api/capital", func(w http.ResponseWriter, r *http.Request) {
		start := r.URL.Query().Get("start")
		end := r.URL.Query().Get("end")
		if start == "" || end == "" {
			log.Printf("capital api: missing start or end (start=%q end=%q)", start, end)
			http.Error(w, "missing start or end", http.StatusBadRequest)
			return
		}
		path, err := planner.Plan(start, end)
		if err != nil {
			log.Printf("capital api: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"route": path})
	}).Methods("GET")

	rp, err := routepkg.NewRoute(store, nil, nil)
	if err != nil {
		log.Fatalf("cannot create route planner: %v", err)
	}
	r.HandleFunc("/api/route/{from}/{to}", api.NewRouteHandler(rp)).Methods("GET")

	r.PathPrefix("/").Handler(http.FileServer(http.FS(frontendFS)))

	csrfKey := mustEnv("CSRF_KEY")
	csrfMiddleware := csrf.Protect([]byte(csrfKey), csrf.Secure(false))

	addr := ":" + cfg.Port
	log.Printf("server started on %s", addr)
	log.Fatal(http.ListenAndServe(addr, csrfMiddleware(r)))
}
