package usecase

import (
	"context"
	"database/sql"
	"github.com/arvians-id/go-apriori-microservice/adapter/pb"
	"github.com/arvians-id/go-apriori-microservice/model"
	"github.com/arvians-id/go-apriori-microservice/services/category-service/repository"
	"github.com/arvians-id/go-apriori-microservice/util"
	"github.com/golang/protobuf/ptypes/empty"
	"log"
	"time"
)

type CategoryService struct {
	CategoryRepository repository.CategoryRepository
	DB                 *sql.DB
}

func NewCategoryService(categoryRepository repository.CategoryRepository, db *sql.DB) pb.CategoryServiceServer {
	return &CategoryService{
		CategoryRepository: categoryRepository,
		DB:                 db,
	}
}

func (service *CategoryService) FindAll(ctx context.Context, empty *empty.Empty) (*pb.ListCategoryResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[CategoryService][FindAll] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	categories, err := service.CategoryRepository.FindAll(ctx, tx)
	if err != nil {
		log.Println("[CategoryService][FindAll][FindAll] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	var categoriesResponse []*pb.Category
	for _, category := range categories {
		categoriesResponse = append(categoriesResponse, category.ToProtoBuff())
	}

	return &pb.ListCategoryResponse{
		Categories: categoriesResponse,
	}, nil
}

func (service *CategoryService) FindById(ctx context.Context, req *pb.GetCategoryByIdRequest) (*pb.GetCategoryResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[CategoryService][FindById] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	category, err := service.CategoryRepository.FindById(ctx, tx, req.Id)
	if err != nil {
		log.Println("[CategoryService][FindById][FindById] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return &pb.GetCategoryResponse{
		Category: category.ToProtoBuff(),
	}, nil
}

func (service *CategoryService) Create(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.GetCategoryResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[CategoryService][Create] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	timeNow, err := time.Parse(util.TimeFormat, time.Now().Format(util.TimeFormat))
	if err != nil {
		log.Println("[CategoryService][Create] problem in parsing to time, err: ", err.Error())
		return nil, err
	}

	categoryRequest := model.Category{
		Name:      util.UpperWords(req.Name),
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	category, err := service.CategoryRepository.Create(ctx, tx, &categoryRequest)
	if err != nil {
		log.Println("[CategoryService][Create][Create] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return &pb.GetCategoryResponse{
		Category: category.ToProtoBuff(),
	}, nil

}

func (service *CategoryService) Update(ctx context.Context, req *pb.UpdateCategoryRequest) (*pb.GetCategoryResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[CategoryService][Update] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	category, err := service.CategoryRepository.FindById(ctx, tx, req.IdCategory)
	if err != nil {
		log.Println("[CategoryService][Update][FindById] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	timeNow, err := time.Parse(util.TimeFormat, time.Now().Format(util.TimeFormat))
	if err != nil {
		log.Println("[CategoryService][Update] problem in parsing to time, err: ", err.Error())
		return nil, err
	}
	category.Name = util.UpperWords(req.Name)
	category.UpdatedAt = timeNow

	_, err = service.CategoryRepository.Update(ctx, tx, category)
	if err != nil {
		log.Println("[CategoryService][Update][Update] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return &pb.GetCategoryResponse{
		Category: category.ToProtoBuff(),
	}, nil
}

func (service *CategoryService) Delete(ctx context.Context, req *pb.GetCategoryByIdRequest) (*empty.Empty, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[CategoryService][Delete] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	category, err := service.CategoryRepository.FindById(ctx, tx, req.Id)
	if err != nil {
		log.Println("[CategoryService][Delete][FindById] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	err = service.CategoryRepository.Delete(ctx, tx, category.IdCategory)
	if err != nil {
		log.Println("[CategoryService][Delete][Delete] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return nil, nil
}
