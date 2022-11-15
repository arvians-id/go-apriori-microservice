package repository

import (
	"context"
	"database/sql"
	"github.com/arvians-id/go-apriori-microservice/services/apriori-service/model"
)

type AprioriRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx) ([]*model.Apriori, error)
	FindAllByActive(ctx context.Context, tx *sql.Tx) ([]*model.Apriori, error)
	FindAllByCode(ctx context.Context, tx *sql.Tx, code string) ([]*model.Apriori, error)
	FindByCodeAndId(ctx context.Context, tx *sql.Tx, code string, id int64) (*model.Apriori, error)
	Create(ctx context.Context, tx *sql.Tx, apriories []*model.Apriori) error
	Update(ctx context.Context, tx *sql.Tx, apriori *model.Apriori) (*model.Apriori, error)
	Delete(ctx context.Context, tx *sql.Tx, code string) error
	UpdateAllStatus(ctx context.Context, tx *sql.Tx, status bool) error
	UpdateStatusByCode(ctx context.Context, tx *sql.Tx, code string, status bool) error
}
