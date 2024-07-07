package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/leary1337/url-shortener/internal/app/service"
	"github.com/leary1337/url-shortener/pkg/logger"
)

type PingHandler struct {
	l       logger.Interface
	service service.Ping
}

func NewPingHandler(l logger.Interface, service service.Ping) *PingHandler {
	return &PingHandler{
		l:       l,
		service: service,
	}
}
func (p *PingHandler) RegisterRoutes(r chi.Router) {
	r.Get("/ping", p.PingDB)
}

func (p *PingHandler) PingDB(w http.ResponseWriter, r *http.Request) {
	err := p.service.PingDB(r.Context())
	if err != nil {
		p.l.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
