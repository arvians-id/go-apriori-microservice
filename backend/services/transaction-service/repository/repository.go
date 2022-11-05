package repository

import (
	"context"
	"database/sql"
	"github.com/arvians-id/go-apriori-microservice/model"
)

type TransactionRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx) ([]*model.Transaction, error)
	FindAllItemSet(ctx context.Context, tx *sql.Tx, startDate string, endDate string) ([]*model.Transaction, error)
	FindByNoTransaction(ctx context.Context, tx *sql.Tx, noTransaction string) (*model.Transaction, error)
	CreateByCsv(ctx context.Context, tx *sql.Tx, transaction []*model.Transaction) error
	Create(ctx context.Context, tx *sql.Tx, transaction *model.Transaction) (*model.Transaction, error)
	Update(ctx context.Context, tx *sql.Tx, transaction *model.Transaction) (*model.Transaction, error)
	Delete(ctx context.Context, tx *sql.Tx, noTransaction string) error
	Truncate(ctx context.Context, tx *sql.Tx) error
}
