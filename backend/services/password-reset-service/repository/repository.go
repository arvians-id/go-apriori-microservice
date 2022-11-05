package repository

import (
	"context"
	"database/sql"
	"github.com/arvians-id/go-apriori-microservice/model"
)

type PasswordResetRepository interface {
	FindByEmail(ctx context.Context, tx *sql.Tx, email string) (*model.PasswordReset, error)
	FindByEmailAndToken(ctx context.Context, tx *sql.Tx, passwordReset *model.PasswordReset) (*model.PasswordReset, error)
	Create(ctx context.Context, tx *sql.Tx, passwordReset *model.PasswordReset) (*model.PasswordReset, error)
	Update(ctx context.Context, tx *sql.Tx, passwordReset *model.PasswordReset) (*model.PasswordReset, error)
	Delete(ctx context.Context, tx *sql.Tx, email string) error
	Truncate(ctx context.Context, tx *sql.Tx) error
}
