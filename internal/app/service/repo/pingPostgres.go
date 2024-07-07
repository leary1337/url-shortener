package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/leary1337/url-shortener/internal/app/service"
)

var _ service.PingRepo = (*PingPostgres)(nil)

type PingPostgres struct {
	pg *pgxpool.Pool
}

func NewPingPostgres(pool *pgxpool.Pool) *PingPostgres {
	return &PingPostgres{
		pg: pool,
	}
}

func (p *PingPostgres) Ping(ctx context.Context) error {
	return p.pg.Ping(ctx)
}
