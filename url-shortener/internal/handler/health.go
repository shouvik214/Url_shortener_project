package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"url-shortener/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct {
	db    *pgxpool.Pool
	redis *repository.RedisRepository
}

func NewHealthHandler(
	db *pgxpool.Pool,
	redis *repository.RedisRepository,
) *HealthHandler {
	return &HealthHandler{
		db:    db,
		redis: redis,
	}
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	// Check PostgreSQL
	if err := h.db.Ping(ctx); err != nil {
		http.Error(
			w,
			`{"status":"not ready","database":"down"}`,
			http.StatusServiceUnavailable,
		)
		return
	}

	// Check Redis
	if err := h.redis.Ping(ctx); err != nil {
		http.Error(
			w,
			`{"status":"not ready","redis":"down"}`,
			http.StatusServiceUnavailable,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"status": "ready",
	})
}
