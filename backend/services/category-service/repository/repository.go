package repository

import (
	"context"
	"database/sql"
	"github.com/arvians-id/go-apriori-microservice/model"
)

type CategoryRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx) ([]*model.Category, error)
	FindById(ctx context.Context, tx *sql.Tx, id int64) (*model.Category, error)
	Create(ctx context.Context, tx *sql.Tx, category *model.Category) (*model.Category, error)
	Update(ctx context.Context, tx *sql.Tx, category *model.Category) (*model.Category, error)
	Delete(ctx context.Context, tx *sql.Tx, id int64) error
}
