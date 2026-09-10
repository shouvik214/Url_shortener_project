package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"url-shortener/internal/model"
)

var ErrURLNotFound = errors.New("URL not found")
var ErrDuplicateShortCode = errors.New("short code already exists")

type URLRepository struct {
	db *pgxpool.Pool
}

func NewURLRepository(db *pgxpool.Pool) *URLRepository {
	return &URLRepository{
		db: db,
	}
}

func (r *URLRepository) GetNextID(ctx context.Context) (int64, error) {
	var id int64

	err := r.db.QueryRow(
		ctx,
		"SELECT nextval('urls_id_seq')",
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *URLRepository) Create(
	ctx context.Context,
	id int64,
	shortCode string,
	originalURL string,
) (*model.URL, error) {

	var url model.URL

	query := `
		INSERT INTO urls (
			id,
			short_code,
			original_url
		)
		VALUES ($1, $2, $3)
		RETURNING id, short_code, original_url, created_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		id,
		shortCode,
		originalURL,
	).Scan(
		&url.ID,
		&url.ShortCode,
		&url.OriginalURL,
		&url.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &url, nil
}

func (r *URLRepository) GetByShortCode(
	ctx context.Context,
	shortCode string,
) (*model.URL, error) {

	var url model.URL

	query := `
		SELECT id, short_code, original_url, created_at
		FROM urls
		WHERE short_code = $1
	`

	err := r.db.QueryRow(
		ctx,
		query,
		shortCode,
	).Scan(
		&url.ID,
		&url.ShortCode,
		&url.OriginalURL,
		&url.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrURLNotFound
		}

		return nil, err
	}

	return &url, nil
}

func (r *URLRepository) DeleteByShortCode(
	ctx context.Context,
	shortCode string,
) error {

	result, err := r.db.Exec(
		ctx,
		"DELETE FROM urls WHERE short_code = $1",
		shortCode,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrURLNotFound
	}

	return nil
}

func (r *URLRepository) FindByOriginalURL(
	ctx context.Context,
	originalURL string,
) (*model.URL, error) {

	var url model.URL

	query := `
		SELECT id, short_code, original_url, created_at
		FROM urls
		WHERE original_url = $1
		LIMIT 1
	`

	err := r.db.QueryRow(
		ctx,
		query,
		originalURL,
	).Scan(
		&url.ID,
		&url.ShortCode,
		&url.OriginalURL,
		&url.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &url, nil
}
