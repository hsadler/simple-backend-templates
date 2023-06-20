package dependencies

import (
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Dependencies struct {
	Validator *validator.Validate
	DBPool    *pgxpool.Pool
}
