package repository

import (
	"context"
	"database/sql"
	"github.com/arvians-id/go-apriori-microservice/model"
)

type PaymentRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx) ([]*model.Payment, error)
	FindAllByUserId(ctx context.Context, tx *sql.Tx, userId int) ([]*model.Payment, error)
	FindByOrderId(ctx context.Context, tx *sql.Tx, orderId string) (*model.Payment, error)
	Create(ctx context.Context, tx *sql.Tx, payment *model.Payment) (*model.Payment, error)
	Update(ctx context.Context, tx *sql.Tx, payment *model.Payment) error
	UpdateReceiptNumber(ctx context.Context, tx *sql.Tx, payment *model.Payment) error
	Delete(ctx context.Context, tx *sql.Tx, orderId *string) error
}
