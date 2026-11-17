package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/ppieniazek/chirpy/internal/database"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	db             *database.Queries
}

func main() {
	const (
		filepathRoot = "."
		port         = "8080"
	)

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Couldn't connect to database: %v", err)
	}
	apiCfg := apiConfig{db: database.New(dbConn)}

	rootFileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(rootFileServerHandler))

	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetMetricsHandler)

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusOK)))
	})
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)

	s := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(s.ListenAndServe())
}
