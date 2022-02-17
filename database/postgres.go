package database

import (
	"context"
	"errors"
	"log"

	"messenger-backend/models"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	Postgres *pgxpool.Pool
}

var _ models.DB = (*DB)(nil)

func (db *DB) CreateAccount(ctx context.Context, user models.User) error {
	const sql = `INSERT INTO users ("username", "name", "email", "hash")
	 VALUES ($1, $2, $3, $4);`

	switch _, err := db.Postgres.Exec(ctx, sql, user.Username, user.Name,
		user.Email, user.Password); {
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return err
	case err != nil:
		if sqlErr := db.accountPgError(err); sqlErr != nil {
			return sqlErr
		}
		log.Printf("cannot create acoount on database: %v\n", err)
		return errors.New("cannot create account on database")
	}
	return nil
}

func (db *DB) accountPgError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return nil
	}
	if pgErr.Code == pgerrcode.UniqueViolation {
		return errors.New("account already exists")
	}
	return nil
}
