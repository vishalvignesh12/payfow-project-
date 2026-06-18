package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/vishalvignesh12/payflow/internal/config"
	"github.com/vishalvignesh12/payflow/internal/db"
)

func main() {

	cfg := config.Load()

	postgres, err := db.NewPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to postgres")
	}

	defer postgres.Close()

	redis, err := db.Newredis(cfg.RedisURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to redis")
	}

	defer redis.Close()

	log.Info().Msg("conncted to postgres and redis")

	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := postgres.PingContext(context.Background()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, `{"status": "unhealty","db": "error"}`)
		}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status": "healty","db": "ok", "redis": "ok"}`)
	})

	addr := ":" + cfg.Port
	log.Info().Str("adr", addr).Msg("Server is starting...")
	http.ListenAndServe(addr, r)
}
