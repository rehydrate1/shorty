package storage

import (
	"context"
	"errors"
)

var ErrURLNotFound = errors.New("url not found")

type URLSaver interface {
	SaveURL(ctx context.Context, alias, urlToSave string) error
	GetURL(ctx context.Context, alias string) (string, error)
}
