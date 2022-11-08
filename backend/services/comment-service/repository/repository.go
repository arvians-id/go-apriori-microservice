package repository

import (
	"context"
	"database/sql"
	"github.com/arvians-id/go-apriori-microservice/model"
)

type CommentRepository interface {
	FindAllRatingByProductCode(ctx context.Context, tx *sql.Tx, productCode string) ([]*model.RatingFromComment, error)
	FindAllByProductCode(ctx context.Context, tx *sql.Tx, productCode string, rating string, tags string) ([]*model.Comment, error)
	FindById(ctx context.Context, tx *sql.Tx, id int64) (*model.Comment, error)
	FindByUserOrderId(ctx context.Context, tx *sql.Tx, userOrderId int64) (*model.Comment, error)
	Create(ctx context.Context, tx *sql.Tx, comment *model.Comment) (*model.Comment, error)
}
