package repository

import (
	"context"
	"database/sql"
	"github.com/arvians-id/go-apriori-microservice/model"
)

type NotificationRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx) ([]*model.Notification, error)
	FindAllByUserId(ctx context.Context, tx *sql.Tx, userId int) ([]*model.Notification, error)
	Create(ctx context.Context, tx *sql.Tx, notification *model.Notification) (*model.Notification, error)
	Mark(ctx context.Context, tx *sql.Tx, id int) error
	MarkAll(ctx context.Context, tx *sql.Tx, userId int) error
}
