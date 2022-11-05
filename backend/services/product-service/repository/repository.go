package repository

import (
	"context"
	"database/sql"
	"github.com/arvians-id/go-apriori-microservice/model"
)

type ProductRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx, search string, category string) ([]*model.Product, error)
	FindAllByAdmin(ctx context.Context, tx *sql.Tx) ([]*model.Product, error)
	FindAllBySimilarCategory(ctx context.Context, tx *sql.Tx, category string) ([]*model.Product, error)
	FindById(ctx context.Context, tx *sql.Tx, id int) (*model.Product, error)
	FindByCode(ctx context.Context, tx *sql.Tx, code string) (*model.Product, error)
	FindByName(ctx context.Context, tx *sql.Tx, name string) (*model.Product, error)
	Create(ctx context.Context, tx *sql.Tx, product *model.Product) (*model.Product, error)
	Update(ctx context.Context, tx *sql.Tx, product *model.Product) (*model.Product, error)
	Delete(ctx context.Context, tx *sql.Tx, code string) error
}
