package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	redisLib "github.com/arvians-id/go-apriori-microservice/services/category-service/third-party/redis"
	"github.com/go-redis/redis/v8"

	"github.com/arvians-id/go-apriori-microservice/services/category-service/model"
	"github.com/arvians-id/go-apriori-microservice/services/category-service/pb"
	"github.com/arvians-id/go-apriori-microservice/services/category-service/repository"
	"github.com/arvians-id/go-apriori-microservice/services/category-service/util"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CategoryServiceCache struct {
	CategoryRepository repository.CategoryRepository
	Redis              redisLib.Redis
	DB                 *sql.DB
}

func NewCategoryServiceCache(categoryRepository repository.CategoryRepository, redis *redisLib.Redis, db *sql.DB) pb.CategoryServiceServer {
	return &CategoryServiceCache{
		CategoryRepository: categoryRepository,
		Redis:              *redis,
		DB:                 db,
	}
}

func (service *CategoryServiceCache) FindAll(ctx context.Context, empty *emptypb.Empty) (*pb.ListCategoryResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[CategoryServiceCache][FindAll] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	categoriesCache, err := service.Redis.Get(ctx, "categories")
	if err != redis.Nil {
		var categories []*pb.Category
		err = json.Unmarshal(categoriesCache, &categories)
		if err != nil {
			log.Println("[CategoryServiceCache][FindAll] unable to unmarshal json, err: ", err.Error())
			return nil, err
		}

		return &pb.ListCategoryResponse{
			Categories: categories,
		}, nil
	}

	categories, err := service.CategoryRepository.FindAll(ctx, tx)
	if err != nil {
		log.Println("[CategoryServiceCache][FindAll][FindAll] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	var categoriesResponse []*pb.Category
	for _, category := range categories {
		categoriesResponse = append(categoriesResponse, category.ToProtoBuff())
	}

	err = service.Redis.Set(ctx, "categories", categoriesResponse)
	if err != nil {
		log.Println("[CategoryServiceCache][FindById][Set] unable to set value to redis cache, err: ", err.Error())
		return nil, err
	}

	return &pb.ListCategoryResponse{
		Categories: categoriesResponse,
	}, nil
}

func (service *CategoryServiceCache) FindById(ctx context.Context, req *pb.GetCategoryByIdRequest) (*pb.GetCategoryResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[CategoryServiceCache][FindById] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	key := fmt.Sprintf("category:%d", req.Id)
	categoryCache, err := service.Redis.Get(ctx, key)
	if err != redis.Nil {
		var category model.Category
		err = json.Unmarshal(categoryCache, &category)
		if err != nil {
			log.Println("[CategoryServiceCache][FindById] unable to unmarshal json, err: ", err.Error())
			return nil, err
		}

		return &pb.GetCategoryResponse{
			Category: category.ToProtoBuff(),
		}, nil
	}

	category, err := service.CategoryRepository.FindById(ctx, tx, req.Id)
	if err != nil {
		log.Println("[CategoryServiceCache][FindById][FindById] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	err = service.Redis.Set(ctx, key, category)
	if err != nil {
		log.Println("[CategoryServiceCache][FindById][Set] unable to set value to redis cache, err: ", err.Error())
		return nil, err
	}

	return &pb.GetCategoryResponse{
		Category: category.ToProtoBuff(),
	}, nil
}

func (service *CategoryServiceCache) Create(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.GetCategoryResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[CategoryServiceCache][Create] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	timeNow, err := time.Parse(util.TimeFormat, time.Now().Format(util.TimeFormat))
	if err != nil {
		log.Println("[CategoryServiceCache][Create] problem in parsing to time, err: ", err.Error())
		return nil, err
	}

	categoryRequest := model.Category{
		Name:      util.UpperWords(req.Name),
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	category, err := service.CategoryRepository.Create(ctx, tx, &categoryRequest)
	if err != nil {
		log.Println("[CategoryServiceCache][Create][Create] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	err = service.Redis.Del(ctx, "categories")
	if err != nil {
		log.Println("[CategoryServiceCache][Create][Del] unable to deleting specific key cache, err: ", err.Error())
	}

	return &pb.GetCategoryResponse{
		Category: category.ToProtoBuff(),
	}, nil

}

func (service *CategoryServiceCache) Update(ctx context.Context, req *pb.UpdateCategoryRequest) (*pb.GetCategoryResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[CategoryServiceCache][Update] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	category, err := service.CategoryRepository.FindById(ctx, tx, req.IdCategory)
	if err != nil {
		log.Println("[CategoryServiceCache][Update][FindById] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	timeNow, err := time.Parse(util.TimeFormat, time.Now().Format(util.TimeFormat))
	if err != nil {
		log.Println("[CategoryServiceCache][Update] problem in parsing to time, err: ", err.Error())
		return nil, err
	}
	category.Name = util.UpperWords(req.Name)
	category.UpdatedAt = timeNow

	_, err = service.CategoryRepository.Update(ctx, tx, category)
	if err != nil {
		log.Println("[CategoryServiceCache][Update][Update] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	key := fmt.Sprintf("category:%d", category.IdCategory)
	err = service.Redis.Del(ctx, "categories", key)
	if err != nil {
		log.Println("[CategoryServiceCache][Update][Del] unable to deleting specific key cache, err: ", err.Error())
	}

	return &pb.GetCategoryResponse{
		Category: category.ToProtoBuff(),
	}, nil
}

func (service *CategoryServiceCache) Delete(ctx context.Context, req *pb.GetCategoryByIdRequest) (*emptypb.Empty, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[CategoryServiceCache][Delete] problem in db transaction, err: ", err.Error())
		return new(emptypb.Empty), err
	}
	defer util.CommitOrRollback(tx)

	category, err := service.CategoryRepository.FindById(ctx, tx, req.Id)
	if err != nil {
		log.Println("[CategoryServiceCache][Delete][FindById] problem in getting from repository, err: ", err.Error())
		return new(emptypb.Empty), err
	}

	err = service.CategoryRepository.Delete(ctx, tx, category.IdCategory)
	if err != nil {
		log.Println("[CategoryServiceCache][Delete][Delete] problem in getting from repository, err: ", err.Error())
		return new(emptypb.Empty), err
	}

	key := fmt.Sprintf("category:%d", category.IdCategory)
	err = service.Redis.Del(ctx, "categories", key)
	if err != nil {
		log.Println("[CategoryServiceCache][Delete][Del] unable to deleting specific key cache, err: ", err.Error())
	}

	return new(emptypb.Empty), nil
}
