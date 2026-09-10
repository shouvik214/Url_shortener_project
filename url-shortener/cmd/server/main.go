package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"url-shortener/internal/config"
	"url-shortener/internal/database"
	"url-shortener/internal/handler"
	"url-shortener/internal/repository"
	"url-shortener/internal/router"
	"url-shortener/internal/service"

	"github.com/joho/godotenv"
)

func main() {

	// --------------------------------------------------
	// Load environment
	// --------------------------------------------------

	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	// --------------------------------------------------
	// Configuration
	// --------------------------------------------------

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// --------------------------------------------------
	// Database
	// --------------------------------------------------

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Println("Connected to PostgreSQL successfully")

	// --------------------------------------------------
	// Repositories
	// --------------------------------------------------

	redisRepository := repository.NewRedisRepository(cfg.RedisAddr)
	defer redisRepository.Close()

	urlRepository := repository.NewURLRepository(db)

	// --------------------------------------------------
	// Service
	// --------------------------------------------------

	urlService := service.NewURLService(
		urlRepository,
		redisRepository,
	)

	// --------------------------------------------------
	// Handlers
	// --------------------------------------------------

	urlHandler := handler.NewURLHandler(
		urlService,
		cfg.BaseURL,
	)

	healthHandler := handler.NewHealthHandler(
		db,
		redisRepository,
	)

	// --------------------------------------------------
	// Logger
	// --------------------------------------------------

	logger := slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		),
	)

	// --------------------------------------------------
	// Router
	// --------------------------------------------------

	appHandler := router.NewRouter(
		urlHandler,
		healthHandler,
		logger,
		cfg.FrontendURL,
	)

	// --------------------------------------------------
	// HTTP Server
	// --------------------------------------------------

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: appHandler,

		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MB
	}

	// --------------------------------------------------
	// Start Server
	// --------------------------------------------------

	go func() {

		log.Printf(
			"Server running on :%s",
			cfg.Port,
		)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {

			log.Fatalf(
				"Server failed: %v",
				err,
			)
		}
	}()

	// --------------------------------------------------
	// Graceful Shutdown
	// --------------------------------------------------

	stop := make(chan os.Signal, 1)

	signal.Notify(
		stop,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-stop

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf(
			"Server shutdown error: %v",
			err,
		)
	}

	log.Println("Server stopped")
}
