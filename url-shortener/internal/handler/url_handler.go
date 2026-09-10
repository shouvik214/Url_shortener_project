package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"url-shortener/internal/repository"
	"url-shortener/internal/service"
)

type URLHandler struct {
	service *service.URLService
	baseURL string
}

func NewURLHandler(
	service *service.URLService,
	baseURL string,
) *URLHandler {
	return &URLHandler{
		service: service,
		baseURL: baseURL,
	}
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

func (h *URLHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()

	var request ShortenRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&request)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if decoder.Decode(&struct{}{}) != io.EOF {
		writeError(
			w,
			http.StatusBadRequest,
			"Request body must contain a single JSON object",
		)
		return
	}

	originalURL := strings.TrimSpace(request.URL)

	if originalURL == "" {
		writeError(w, http.StatusBadRequest, "URL is required")
		return
	}

	parsedURL, err := url.ParseRequestURI(originalURL)
	if err != nil ||
		parsedURL.Host == "" ||
		(parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {

		writeError(w, http.StatusBadRequest, "Invalid URL")
		return
	}

	ctx, cancel := context.WithTimeout(
		r.Context(),
		3*time.Second,
	)
	defer cancel()

	url, err := h.service.CreateShortURL(
		ctx,
		originalURL,
	)

	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	response := ShortenResponse{
		ShortURL: h.baseURL + "/" + url.ShortCode,
	}

	writeJSON(w, http.StatusCreated, response)
}

func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	shortCode := strings.TrimPrefix(r.URL.Path, "/")

	if shortCode == "" {
		writeError(w, http.StatusNotFound, "URL not found")
		return
	}
	ctx, cancel := context.WithTimeout(
		r.Context(),
		3*time.Second,
	)
	defer cancel()

	originalURL, err := h.service.GetOriginalURL(
		ctx,
		shortCode,
	)

	if err != nil {
		if errors.Is(err, repository.ErrURLNotFound) {
			writeError(w, http.StatusNotFound, "URL not found")
			return
		}

		writeError(
			w,
			http.StatusInternalServerError,
			"Internal server error",
		)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}

func (h *URLHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.Header().Set("Allow", http.MethodDelete)
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	shortCode := strings.TrimPrefix(r.URL.Path, "/")

	if shortCode == "" {
		writeError(w, http.StatusNotFound, "URL not found")
		return
	}
	ctx, cancel := context.WithTimeout(
		r.Context(),
		3*time.Second,
	)
	defer cancel()

	err := h.service.DeleteURL(
		ctx,
		shortCode,
	)

	if err != nil {
		if errors.Is(err, repository.ErrURLNotFound) {
			writeError(w, http.StatusNotFound, "URL not found")
			return
		}
		//logger.ErrorLogger.Printf(
		//	"failed to delete URL %s: %v",
		//	shortCode,
		//	err,
		//)
		writeError(
			w,
			http.StatusInternalServerError,
			"Internal server error",
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
