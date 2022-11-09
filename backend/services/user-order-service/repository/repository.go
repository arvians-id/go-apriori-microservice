package repository

import (
	"context"
	"database/sql"
	"github.com/arvians-id/go-apriori-microservice/model"
)

type UserOrderRepository interface {
	FindAllByPayloadId(ctx context.Context, tx *sql.Tx, payloadId int64) ([]*model.UserOrder, error)
	FindAllByUserId(ctx context.Context, tx *sql.Tx, userId int64) ([]*model.UserOrder, error)
	FindById(ctx context.Context, tx *sql.Tx, id int64) (*model.UserOrder, error)
	Create(ctx context.Context, tx *sql.Tx, userOrder *model.UserOrder) (*model.UserOrder, error)
}
