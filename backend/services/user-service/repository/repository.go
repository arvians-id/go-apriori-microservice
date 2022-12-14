package repository

import (
	"context"
	"database/sql"
	"github.com/arvians-id/go-apriori-microservice/services/user-service/model"
)

type UserRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx) ([]*model.User, error)
	FindById(ctx context.Context, tx *sql.Tx, id int64) (*model.User, error)
	FindByEmail(ctx context.Context, tx *sql.Tx, email string) (*model.User, error)
	Create(ctx context.Context, tx *sql.Tx, user *model.User) (*model.User, error)
	Update(ctx context.Context, tx *sql.Tx, user *model.User) (*model.User, error)
	UpdatePassword(ctx context.Context, tx *sql.Tx, user *model.User) error
	Delete(ctx context.Context, tx *sql.Tx, id int64) error
}
