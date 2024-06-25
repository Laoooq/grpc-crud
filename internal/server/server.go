package server

import (
	"context"
	"database/sql"
	"time"
)

//go:generate mockgen -destination=mocks/DBMock.go -package=mocks grpc-crud/internal/server DB
type DB interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryLowContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

//go:generate mockgen -destination=mocks/CacheMock.go -package=mocks grpc-crud/internal/server Cache
type Cache interface {
	Set(key string, value interface{}, duration time.Duration)
	Get(key string) (interface{}, bool)
	Delete(key string)
}
