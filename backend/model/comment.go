package model

import (
	"github.com/arvians-id/go-apriori-microservice/adapter/pkg/comment/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Comment struct {
	IdComment   int64     `json:"id_comment"`
	UserOrderId int64     `json:"user_order_id"`
	ProductCode string    `json:"product_code"`
	Description *string   `json:"description"`
	Tag         *string   `json:"tag"`
	Rating      int32     `json:"rating"`
	CreatedAt   time.Time `json:"created_at"`
}

func (comment *Comment) ToProtoBuff() *pb.Comment {
	return &pb.Comment{
		IdComment:   comment.IdComment,
		UserOrderId: comment.UserOrderId,
		ProductCode: comment.ProductCode,
		Description: comment.Description,
		Tag:         comment.Tag,
		Rating:      comment.Rating,
		CreatedAt:   timestamppb.New(comment.CreatedAt),
	}
}

type RatingFromComment struct {
	Rating        int32 `json:"rating"`
	ResultRating  int32 `json:"result_rating"`
	ResultComment int32 `json:"result_comment"`
}

func (ratingFromComment *RatingFromComment) ToProtoBuff() *pb.ListRatingFromCommentResponse_RatingFromComment {
	return &pb.ListRatingFromCommentResponse_RatingFromComment{
		Rating:        ratingFromComment.Rating,
		ResultRating:  ratingFromComment.ResultRating,
		ResultComment: ratingFromComment.ResultComment,
	}
}
