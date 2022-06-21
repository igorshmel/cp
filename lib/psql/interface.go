package psql

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

// ServerPgxPool интерфейс для получения доступа к БД из gRPC перехватчиков
type ServerPgxPool interface {
	GetPgxPool() *pgxpool.Pool
}
