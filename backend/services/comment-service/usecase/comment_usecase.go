package usecase

import (
	"context"
	"database/sql"
	"github.com/arvians-id/go-apriori-microservice/services/comment-service/client"
	"github.com/arvians-id/go-apriori-microservice/services/comment-service/model"
	"github.com/arvians-id/go-apriori-microservice/services/comment-service/pb"
	"github.com/arvians-id/go-apriori-microservice/services/comment-service/repository"
	"github.com/arvians-id/go-apriori-microservice/services/comment-service/util"
	"log"
	"strings"
	"time"
)

type CommentService struct {
	CommentRepository repository.CommentRepository
	ProductService    client.ProductServiceClient
	DB                *sql.DB
}

func NewCommentService(
	commentRepository repository.CommentRepository,
	productService client.ProductServiceClient,
	db *sql.DB,
) pb.CommentServiceServer {
	return &CommentService{
		CommentRepository: commentRepository,
		ProductService:    productService,
		DB:                db,
	}
}

func (service *CommentService) FindAllRatingByProductCode(ctx context.Context, req *pb.GetCommentByProductCodeRequest) (*pb.ListRatingFromCommentResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[CommentService][FindAllRatingByProductCode] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	product, err := service.ProductService.FindByCode(ctx, req.ProductCode)
	if err != nil {
		log.Println("[CategoryService][FindAllRatingByProductCode][FindByCode] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	ratings, err := service.CommentRepository.FindAllRatingByProductCode(ctx, tx, product.Product.Code)
	if err != nil {
		log.Println("[CategoryService][FindAllRatingByProductCode][FindAllRatingByProductCode] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	var ratingListResponse []*pb.RatingFromComment
	for _, rating := range ratings {
		ratingListResponse = append(ratingListResponse, rating.ToProtoBuff())
	}

	return &pb.ListRatingFromCommentResponse{
		RatingFromComments: ratingListResponse,
	}, nil
}

func (service *CommentService) FindAllByProductCode(ctx context.Context, req *pb.GetCommentByFiltersRequest) (*pb.ListCommentResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[CommentService][FindAllByProductCode] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	product, err := service.ProductService.FindByCode(ctx, req.ProductCode)
	if err != nil {
		log.Println("[CategoryService][FindAllByProductCode][FindByCode] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	tagArray := strings.Split(req.Tag, ",")
	tag := strings.Join(tagArray, "|")
	comments, err := service.CommentRepository.FindAllByProductCode(ctx, tx, product.Product.Code, req.Rating, tag)
	if err != nil {
		log.Println("[CategoryService][FindAllByProductCode][FindAllByProductCode] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	var commentListResponse []*pb.Comment
	for _, comment := range comments {
		commentListResponse = append(commentListResponse, comment.ToProtoBuff())
	}

	return &pb.ListCommentResponse{
		Comment: commentListResponse,
	}, nil
}

func (service *CommentService) FindById(ctx context.Context, req *pb.GetCommentByIdRequest) (*pb.GetCommentResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[CommentService][FindById] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	comment, err := service.CommentRepository.FindById(ctx, tx, req.Id)
	if err != nil {
		log.Println("[CategoryService][FindById][FindById] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return &pb.GetCommentResponse{
		Comment: comment.ToProtoBuff(),
	}, nil
}

func (service *CommentService) FindByUserOrderId(ctx context.Context, req *pb.GetCommentByUserOrderIdRequest) (*pb.GetCommentResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[CommentService][FindByUserOrderId] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	comment, err := service.CommentRepository.FindByUserOrderId(ctx, tx, req.UserOrderId)
	if err != nil {
		log.Println("[CategoryService][FindByUserOrderId][FindByUserOrderId] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return &pb.GetCommentResponse{
		Comment: comment.ToProtoBuff(),
	}, nil
}

func (service *CommentService) Create(ctx context.Context, req *pb.CreateCommentRequest) (*pb.GetCommentResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[CommentService][Create] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	timeNow, err := time.Parse(util.TimeFormat, time.Now().Format(util.TimeFormat))
	if err != nil {
		log.Println("[CommentService][Create] problem in parsing to time, err: ", err.Error())
		return nil, err
	}

	comment, err := service.CommentRepository.Create(ctx, tx, &model.Comment{
		UserOrderId: req.UserOrderId,
		ProductCode: req.ProductCode,
		Description: &req.Description,
		Rating:      req.Rating,
		Tag:         &req.Tag,
		CreatedAt:   timeNow,
	})
	if err != nil {
		log.Println("[CategoryService][Create][Create] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return &pb.GetCommentResponse{
		Comment: comment.ToProtoBuffWithoutRelation(),
	}, nil
}
