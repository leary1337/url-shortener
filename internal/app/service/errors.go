package service

import "errors"

var (
	ErrURLNotFound        = errors.New("URL not found")
	ErrPostgresPingFailed = errors.New("ping failed")
)
