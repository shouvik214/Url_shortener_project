package router

import (
	"net/http"

	"url-shortener/internal/handler"
	"url-shortener/internal/middleware"

	"log/slog"
)

func NewRouter(
	urlHandler *handler.URLHandler,
	healthHandler *handler.HealthHandler,
	logger *slog.Logger,
	frontendURL string,
) http.Handler {

	mux := http.NewServeMux()

	// --------------------------------------------------
	// Routes
	// --------------------------------------------------

	// Create short URL
	mux.HandleFunc("/shorten", urlHandler.Shorten)

	// Health checks
	mux.HandleFunc("/health", healthHandler.Health)
	mux.HandleFunc("/ready", healthHandler.Ready)

	// Redirect / Delete
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodGet:
			urlHandler.Redirect(w, r)

		//case http.MethodDelete:
		//	urlHandler.Delete(w, r)

		default:
			http.Error(
				w,
				"Method not allowed",
				http.StatusMethodNotAllowed,
			)
		}
	})

	// --------------------------------------------------
	// Middleware
	// --------------------------------------------------

	appHandler := middleware.SecurityHeaders(mux)
	appHandler = middleware.CORS(frontendURL)(appHandler)
	appHandler = middleware.Logging(logger)(appHandler)
	appHandler = middleware.Recovery(appHandler)
	appHandler = middleware.RequestID(appHandler)

	return appHandler
}
