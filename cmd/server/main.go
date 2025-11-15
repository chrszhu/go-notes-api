package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
    
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/yourusername/resume-app/internal/notes"
)

func main() {
	logger := log.New(os.Stdout, "api ", log.LstdFlags|log.Lshortfile)

	cfg := loadConfig()

	db, err := sql.Open("postgres", cfg.ConnString())
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		logger.Fatal("database unreachable: ", err)
	}

	repo := notes.NewRepository(db)
	h := notes.NewHandler(repo, logger)

	r := mux.NewRouter()
	r.HandleFunc("/health", healthHandler).Methods(http.MethodGet)
	r.HandleFunc("/notes", h.CreateNote).Methods(http.MethodPost)
	r.HandleFunc("/notes", h.ListNotes).Methods(http.MethodGet)
	r.HandleFunc("/notes/{id}", h.GetNote).Methods(http.MethodGet)

	server := &http.Server{ 
		Addr: ":8080",
		Handler: r,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 60 * time.Second,
	}

	logger.Println("starting server on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal(err)
	}
}

type config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func loadConfig() config {
	return config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "resumeapp"),
	}
}

func (c config) ConnString() string {
	return "host=" + c.DBHost + " port=" + c.DBPort + " user=" + c.DBUser + " password=" + c.DBPassword + " dbname=" + c.DBName + " sslmode=disable"
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" { return v }
	return def
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
