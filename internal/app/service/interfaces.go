package service

type (
	Shortener interface {
		ShortenURL(originalURL string) string
		ResolveURL(shortURL string) (string, error)
	}

	ShortenerRepo interface {
		Save(shortURL, originalURL string)
		Get(shortURL string) (string, bool)
	}
)
