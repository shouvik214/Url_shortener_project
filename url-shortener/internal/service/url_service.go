package service

import (
	"context"
	"log"
	"time"

	"url-shortener/internal/model"
	"url-shortener/internal/repository"
)

const cacheTTL = 1 * time.Hour

type URLService struct {
	repo  *repository.URLRepository
	cache *repository.RedisRepository
}

func NewURLService(
	repo *repository.URLRepository,
	cache *repository.RedisRepository,
) *URLService {
	return &URLService{
		repo:  repo,
		cache: cache,
	}
}

func (s *URLService) CreateShortURL(
	ctx context.Context,
	originalURL string,
) (*model.URL, error) {

	// 1. Check whether URL already exists
	existingURL, err := s.repo.FindByOriginalURL(
		ctx,
		originalURL,
	)

	if err != nil {
		return nil, err
	}

	// 2. Return existing short URL
	if existingURL != nil {
		return existingURL, nil
	}

	// 3. Get unique PostgreSQL ID
	id, err := s.repo.GetNextID(ctx)
	if err != nil {
		return nil, err
	}

	// 4. Convert ID to Base62
	shortCode := encodeBase62(uint64(id))

	// 5. Insert URL
	url, err := s.repo.Create(
		ctx,
		id,
		shortCode,
		originalURL,
	)

	if err != nil {
		return nil, err
	}

	// 6. Cache the URL
	if err := s.cache.Set(
		ctx,
		shortCode,
		originalURL,
		cacheTTL,
	); err != nil {
		// Cache failure should not fail URL creation.
	}

	return url, nil
}
func (s *URLService) GetOriginalURL(
	ctx context.Context,
	shortCode string,
) (string, error) {

	// 1. Check Redis
	cachedURL, err := s.cache.Get(ctx, shortCode)

	if err == nil && cachedURL != "" {
		return cachedURL, nil
	}

	if err != nil {
		log.Printf(
			"Redis GET failed for key %s: %v",
			shortCode,
			err,
		)
	}

	// 2. Cache miss → PostgreSQL
	url, err := s.repo.GetByShortCode(ctx, shortCode)

	if err != nil {
		return "", err
	}

	// 3. Store in Redis
	err = s.cache.Set(
		ctx,
		shortCode,
		url.OriginalURL,
		cacheTTL,
	)

	if err != nil {
		log.Printf(
			"Redis SET failed for key %s: %v",
			shortCode,
			err,
		)
	}

	return url.OriginalURL, nil
}

const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func encodeBase62(n uint64) string {
	if n == 0 {
		return "0"
	}

	var result []byte

	for n > 0 {
		remainder := n % 62
		result = append(result, base62Chars[remainder])
		n /= 62
	}

	// Reverse the result
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}

//func generateShortCode(length int) (string, error) {
//
//	bytes := make([]byte, length)
//
//	if _, err := rand.Read(bytes); err != nil {
//		return "", err
//	}
//
//	return base64.RawURLEncoding.EncodeToString(bytes)[:length], nil
//}

func (s *URLService) InvalidateCache(
	ctx context.Context,
	shortCode string,
) error {
	return s.cache.Delete(ctx, shortCode)
}

func (s *URLService) DeleteURL(
	ctx context.Context,
	shortCode string,
) error {

	// 1. Delete from NeonDB
	err := s.repo.DeleteByShortCode(ctx, shortCode)
	if err != nil {
		return err
	}

	// 2. Invalidate Redis cache
	// Redis failure should not make the delete operation fail.
	err = s.cache.Delete(ctx, shortCode)

	if err != nil {
		log.Printf(
			"Redis DELETE failed for key %s: %v",
			shortCode,
			err,
		)
	}

	return nil
}

//func isDuplicateKeyError(err error) bool {
//
//	// PostgreSQL duplicate-key errors will be handled
//	// more precisely in the next step.
//	return errors.Is(err, repository.ErrDuplicateShortCode)
//}
//
