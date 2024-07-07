package service

import "context"

var _ Ping = (*PingService)(nil)

type PingService struct {
	repo PingRepo
}

func NewPingService(repo PingRepo) *PingService {
	return &PingService{repo: repo}
}

func (p *PingService) PingDB(ctx context.Context) error {
	err := p.repo.Ping(ctx)
	if err != nil {
		return ErrPostgresPingFailed
	}
	return nil
}
