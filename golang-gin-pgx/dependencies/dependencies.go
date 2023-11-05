package dependencies

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"

	"example-server/database"
)

type Dependencies struct {
	Validator *validator.Validate
	DBPool    database.PgxPoolIface
}

func NewDependencies(validator *validator.Validate, pgxPool database.PgxPoolIface) *Dependencies {
	return &Dependencies{
		Validator: validator,
		DBPool:    pgxPool,
	}
}

func (deps *Dependencies) CleanupDependencies() {
	deps.DBPool.Close()
}

type PgxIface interface {
	Begin(context.Context) (pgx.Tx, error)
	Close(context.Context) error
}
