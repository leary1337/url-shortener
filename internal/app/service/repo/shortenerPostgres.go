package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/leary1337/url-shortener/internal/app/entity"
	"github.com/leary1337/url-shortener/internal/app/service"
)

var _ service.ShortenerRepo = (*ShortenerPostgres)(nil)

type ShortenerPostgres struct {
	pg *pgxpool.Pool
}

func NewShortenerPostgres(pool *pgxpool.Pool) *ShortenerPostgres {
	return &ShortenerPostgres{
		pg: pool,
	}
}

func (s *ShortenerPostgres) Init(ctx context.Context) error {
	_, err := s.pg.Exec(
		ctx,
		`CREATE TABLE IF NOT EXISTS "shorturl"
			(
				"Id" uuid NOT NULL,
				"ShortURL" text NOT NULL,
				"OriginalURL" text NOT NULL,
				PRIMARY KEY ("Id")
			);`,
	)
	return err
}

func (s *ShortenerPostgres) Save(ctx context.Context, shortURL *entity.ShortURL) error {
	_, err := s.pg.Exec(
		ctx,
		`INSERT INTO "shorturl" VALUES ($1, $2, $3)`,
		shortURL.UUID, shortURL.ShortURL, shortURL.OriginalURL,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *ShortenerPostgres) SaveBatch(ctx context.Context, shortURLs []entity.ShortURL) error {
	batch := &pgx.Batch{}
	query := `INSERT INTO "shorturl" ("Id", "ShortURL", "OriginalURL") VALUES (@id, @short_url, @original_url)`
	for _, url := range shortURLs {
		batch.Queue(
			query,
			pgx.NamedArgs{
				"id":           url.UUID,
				"short_url":    url.ShortURL,
				"original_url": url.OriginalURL,
			},
		)
	}
	r := s.pg.SendBatch(ctx, batch)
	defer func() {
		_ = r.Close()
	}()

	for _, url := range shortURLs {
		_, err := r.Exec()
		if err != nil {
			return fmt.Errorf("unable to insert row (%v): %w", url, err)
		}
	}
	return nil
}

func (s *ShortenerPostgres) GetByShortURL(ctx context.Context, shortURL string) (*entity.ShortURL, error) {
	row := s.pg.QueryRow(
		ctx,
		`SELECT "Id", "ShortURL", "OriginalURL" FROM "shorturl" WHERE "ShortURL" = $1::text`,
		shortURL,
	)
	var url entity.ShortURL
	err := row.Scan(&url.UUID, &url.ShortURL, &url.OriginalURL)
	if err != nil {
		return nil, err
	}
	return &url, nil
}
