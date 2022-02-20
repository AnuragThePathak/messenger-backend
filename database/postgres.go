package database

import (
	"context"
	"errors"
	"fmt"
	"log"

	"messenger-backend/models"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	Postgres *pgxpool.Pool
}

var _ models.DB = (*DB)(nil)

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

func (db *DB) CreateAccount(ctx context.Context, user models.User) error {
	const query = `INSERT INTO users ("username", "name", "email", "hash")
	 VALUES ($1, $2, $3, $4);`

	switch _, err := db.Postgres.Exec(ctx, query, user.Username, user.Name,
		user.Email, user.Password); {
	case errors.Is(err, context.Canceled), errors.Is(err,
		context.DeadlineExceeded):
		return err
	case err != nil:
		if sqlErr := db.accountPgError(err); sqlErr != nil {
			return sqlErr
		}
		log.Printf("cannot create acoount on database: %v\n", err)
		return errors.New("cannot create account on database")
	default:
		return nil
	}
}

func (db *DB) IfEmailOrUsernameExists(ctx context.Context, credentialType string,
	credential string) (bool, error) {
	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM users 
		WHERE %s = $1)`, credentialType)
	var exists bool

	switch err := db.Postgres.QueryRow(ctx, query, credential).
		Scan(&exists); {
	case errors.Is(err, context.Canceled), errors.Is(err,
		context.DeadlineExceeded):
		return false, err
	case errors.Is(err, pgx.ErrNoRows):
		return false, errors.New("no response from database")
	case err != nil:
		log.Println(err)
		return false, errors.New("can't make query")
	default:
		return exists, nil
	}
}
