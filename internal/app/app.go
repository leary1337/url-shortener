package app

import (
	"fmt"
	"io"
	"net/http"
)

const ShortURLLength = 8

type ServerHandler struct {
	serverAddr string
	urlMap     map[string]string
}

func NewServerHandler(serverAddr string) *ServerHandler {
	return &ServerHandler{
		serverAddr: serverAddr,
		urlMap:     map[string]string{},
	}
}

func (a *ServerHandler) MainPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil || len(body) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		shortURL, err := GenerateRandomString(ShortURLLength)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		a.urlMap[shortURL] = string(body)

		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(fmt.Sprintf("http://%s/%s", a.serverAddr, shortURL)))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	if r.Method == http.MethodGet {
		shortURL := r.URL.Path[1:]
		originalURL, ok := a.urlMap[shortURL]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Add("Location", originalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}
