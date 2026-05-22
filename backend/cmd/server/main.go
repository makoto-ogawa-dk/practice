package main

import (
	"log"
	"net/http"
	"os"

	"github.com/makoto-ogawa-dk/practice/backend/internal/api"
	"github.com/makoto-ogawa-dk/practice/backend/internal/store"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://practice:practice@localhost:5432/practice?sslmode=disable"
	}

	db, err := store.NewDB(dsn)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer db.Close()

	h := api.NewHandler(db)

	addr := ":8080"
	log.Printf("backend listening on %s", addr)
	if err := http.ListenAndServe(addr, h.Router()); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
