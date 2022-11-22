package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/arvians-id/go-apriori-microservice/services/product-service/client"
	"github.com/arvians-id/go-apriori-microservice/services/product-service/model"
	"github.com/arvians-id/go-apriori-microservice/services/product-service/pb"
	"github.com/arvians-id/go-apriori-microservice/services/product-service/repository"
	"github.com/arvians-id/go-apriori-microservice/services/product-service/third-party/aws"
	redisLib "github.com/arvians-id/go-apriori-microservice/services/product-service/third-party/redis"
	"github.com/arvians-id/go-apriori-microservice/services/product-service/util"
)

type ProductServiceCache struct {
	ProductRepository repository.ProductRepository
	AprioriService    client.AprioriServiceClient
	StorageS3         aws.StorageS3
	Redis             redisLib.Redis
	DB                *sql.DB
}

func NewProductServiceCache(
	productRepository repository.ProductRepository,
	aprioriService client.AprioriServiceClient,
	storageS3 *aws.StorageS3,
	redis *redisLib.Redis,
	db *sql.DB,
) pb.ProductServiceServer {
	return &ProductServiceCache{
		ProductRepository: productRepository,
		AprioriService:    aprioriService,
		StorageS3:         *storageS3,
		Redis:             *redis,
		DB:                db,
	}
}

func (service *ProductServiceCache) FindAllByAdmin(ctx context.Context, empty *emptypb.Empty) (*pb.ListProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[ProductServiceCache][FindAllByAdmin] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	productsCache, err := service.Redis.Get(ctx, "products:admin")
	if err != redis.Nil {
		var products []*pb.Product
		err = json.Unmarshal(productsCache, &products)
		if err != nil {
			log.Println("[ProductServiceCache][FindAllByAdmin] unable to unmarshal json, err: ", err.Error())
			return nil, err
		}

		return &pb.ListProductResponse{
			Product: products,
		}, nil
	}

	products, err := service.ProductRepository.FindAllByAdmin(ctx, tx)
	if err != nil {
		log.Println("[ProductServiceCache][FindAllByAdmin][FindAllByAdmin] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	var productListResponse []*pb.Product
	for _, product := range products {
		productListResponse = append(productListResponse, product.ToProtoBuff())
	}

	err = service.Redis.Set(ctx, "products:admin", productListResponse)
	if err != nil {
		log.Println("[ProductServiceCache][FindAllByAdmin][Set] unable to set value to redis cache, err: ", err.Error())
		return nil, err
	}

	return &pb.ListProductResponse{
		Product: productListResponse,
	}, nil
}

func (service *ProductServiceCache) FindAll(ctx context.Context, req *pb.GetProductByFiltersRequest) (*pb.ListProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[ProductServiceCache][FindAll] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	productsCache, err := service.Redis.Get(ctx, "products:all")
	if err != redis.Nil {
		var products []*pb.Product
		err = json.Unmarshal(productsCache, &products)
		if err != nil {
			log.Println("[ProductServiceCache][FindAll] unable to unmarshal json, err: ", err.Error())
			return nil, err
		}

		return &pb.ListProductResponse{
			Product: products,
		}, nil
	}

	products, err := service.ProductRepository.FindAll(ctx, tx, req.Search, req.Category)
	if err != nil {
		log.Println("[ProductServiceCache][FindAll][FindAll] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	var productListResponse []*pb.Product
	for _, product := range products {
		productListResponse = append(productListResponse, product.ToProtoBuff())
	}

	err = service.Redis.Set(ctx, "products:all", productListResponse)
	if err != nil {
		log.Println("[ProductServiceCache][FindAllByAdmin][Set] unable to set value to redis cache, err: ", err.Error())
		return nil, err
	}

	return &pb.ListProductResponse{
		Product: productListResponse,
	}, nil
}

func (service *ProductServiceCache) FindAllBySimilarCategory(ctx context.Context, req *pb.GetProductByProductCodeRequest) (*pb.ListProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[ProductServiceCache][FindAllBySimilarCategory] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	productsCache, err := service.Redis.Get(ctx, "products:similar")
	if err != redis.Nil {
		var products []*pb.Product
		err = json.Unmarshal(productsCache, &products)
		if err != nil {
			log.Println("[ProductServiceCache][FindAllBySimilarCategory] unable to unmarshal json, err: ", err.Error())
			return nil, err
		}

		return &pb.ListProductResponse{
			Product: products,
		}, nil
	}

	product, err := service.ProductRepository.FindByCode(ctx, tx, req.Code)
	if err != nil {
		log.Println("[ProductServiceCache][FindAllBySimilarCategory][FindByCode] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	categoryArray := strings.Split(product.Category, ", ")
	categoryString := strings.Join(categoryArray, "|")
	productCategories, err := service.ProductRepository.FindAllBySimilarCategory(ctx, tx, categoryString)
	if err != nil {
		log.Println("[ProductServiceCache][FindAllBySimilarCategory][FindAllBySimilarCategory] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	var productListResponse []*pb.Product
	for _, productCategory := range productCategories {
		if productCategory.Code != req.Code {
			productListResponse = append(productListResponse, productCategory.ToProtoBuff())
		}
	}

	err = service.Redis.Set(ctx, "products:similar", productListResponse)
	if err != nil {
		log.Println("[ProductServiceCache][FindAllBySimilarCategory][Set] unable to set value to redis cache, err: ", err.Error())
		return nil, err
	}

	return &pb.ListProductResponse{
		Product: productListResponse,
	}, nil
}

func (service *ProductServiceCache) FindAllRecommendation(ctx context.Context, req *pb.GetProductByProductCodeRequest) (*pb.ListProductRecommendationResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[ProductServiceCache][FindAllRecommendation] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	product, err := service.ProductRepository.FindByCode(ctx, tx, req.Code)
	if err != nil {
		log.Println("[ProductServiceCache][FindAllRecommendation][FindByCode] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	apriories, err := service.AprioriService.FindAllByActive(ctx)
	if err != nil {
		log.Println("[ProductServiceCache][FindAllRecommendation][FindAllByActive] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	var productResponses []*pb.ProductRecommendation
	for _, apriori := range apriories.Apriori {
		productNames := strings.Split(apriori.Item, ",")
		var exists bool
		for _, productName := range productNames {
			if strings.ToLower(product.Name) == strings.TrimSpace(productName) {
				exists = true
			}
		}

		var totalPrice int32
		if exists {
			for _, productName := range productNames {
				productByName, _ := service.ProductRepository.FindByName(ctx, tx, util.UpperWords(productName))
				totalPrice += productByName.Price
			}

			productResponses = append(productResponses, &pb.ProductRecommendation{
				AprioriId:         apriori.IdApriori,
				AprioriCode:       apriori.Code,
				AprioriItem:       apriori.Item,
				AprioriDiscount:   apriori.Discount,
				ProductTotalPrice: totalPrice,
				PriceDiscount:     totalPrice - (totalPrice * int32(apriori.Discount) / 100),
				AprioriImage:      apriori.Image,
			})
		}
	}

	return &pb.ListProductRecommendationResponse{
		ProductRecommendation: productResponses,
	}, nil
}

func (service *ProductServiceCache) FindByCode(ctx context.Context, req *pb.GetProductByProductCodeRequest) (*pb.GetProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[ProductServiceCache][FindByCode] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	key := fmt.Sprintf("product:%s", req.Code)
	productCache, err := service.Redis.Get(ctx, key)
	if err != redis.Nil {
		var product model.Product
		err = json.Unmarshal(productCache, &product)
		if err != nil {
			log.Println("[ProductServiceCache][FindByCode] unable to unmarshal json, err: ", err.Error())
			return nil, err
		}

		return &pb.GetProductResponse{
			Product: product.ToProtoBuff(),
		}, nil
	}

	productResponse, err := service.ProductRepository.FindByCode(ctx, tx, req.Code)
	if err != nil {
		log.Println("[ProductServiceCache][FindByCode][FindByCode] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	err = service.Redis.Set(ctx, key, productResponse)
	if err != nil {
		log.Println("[ProductServiceCache][FindByCode][Set] unable to set value to redis cache, err: ", err.Error())
		return nil, err
	}

	return &pb.GetProductResponse{
		Product: productResponse.ToProtoBuff(),
	}, nil
}

func (service *ProductServiceCache) FindByName(ctx context.Context, req *pb.GetProductByProductNameRequest) (*pb.GetProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[ProductServiceCache][FindByCode] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	productResponse, err := service.ProductRepository.FindByName(ctx, tx, req.Name)
	if err != nil {
		log.Println("[ProductServiceCache][FindByCode][FindByCode] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return &pb.GetProductResponse{
		Product: productResponse.ToProtoBuff(),
	}, nil
}

func (service *ProductServiceCache) Create(ctx context.Context, req *pb.CreateProductRequest) (*pb.GetProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[ProductServiceCache][Create] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	timeNow, err := time.Parse(util.TimeFormat, time.Now().Format(util.TimeFormat))
	if err != nil {
		log.Println("[ProductServiceCache][Create] problem in parsing to time, err: ", err.Error())
		return nil, err
	}
	if req.Image == "" {
		req.Image = "no-image.png"
	}

	productRequest := model.Product{
		Code:        req.Code,
		Name:        util.UpperWords(req.Name),
		Description: &req.Description,
		Price:       req.Price,
		Image:       &req.Image,
		Category:    util.UpperWords(req.Category),
		IsEmpty:     false,
		Mass:        req.Mass,
		CreatedAt:   timeNow,
		UpdatedAt:   timeNow,
	}

	productResponse, err := service.ProductRepository.Create(ctx, tx, &productRequest)
	if err != nil {
		log.Println("[ProductServiceCache][Create][Create] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	err = service.Redis.Del(ctx, "products:all", "products:admin", "products:similar")
	if err != nil {
		log.Println("[ProductServiceCache][Create][Del] unable to deleting specific key cache, err: ", err.Error())
	}

	return &pb.GetProductResponse{
		Product: productResponse.ToProtoBuff(),
	}, nil
}

func (service *ProductServiceCache) Update(ctx context.Context, req *pb.UpdateProductRequest) (*pb.GetProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[ProductServiceCache][Update] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	product, err := service.ProductRepository.FindByCode(ctx, tx, req.Code)
	if err != nil {
		log.Println("[ProductServiceCache][Update][FindByCode] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	timeNow, err := time.Parse(util.TimeFormat, time.Now().Format(util.TimeFormat))
	if err != nil {
		log.Println("[ProductServiceCache][Update] problem in parsing to time, err: ", err.Error())
		return nil, err
	}

	product.Name = util.UpperWords(req.Name)
	product.Description = &req.Description
	product.Price = req.Price
	product.Category = util.UpperWords(req.Category)
	product.IsEmpty = req.IsEmpty
	product.Mass = req.Mass
	product.UpdatedAt = timeNow
	if req.Image != "" {
		_ = service.StorageS3.DeleteFromAWS(product.Image)
		product.Image = &req.Image
	}

	productResponse, err := service.ProductRepository.Update(ctx, tx, product)
	if err != nil {
		log.Println("[ProductServiceCache][Update][Update] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	err = service.Redis.Del(ctx, "products:all", "products:admin", "products:similar", fmt.Sprintf("product:%s", req.Code))
	if err != nil {
		log.Println("[ProductServiceCache][Update][Del] unable to deleting specific key cache, err: ", err.Error())
	}

	return &pb.GetProductResponse{
		Product: productResponse.ToProtoBuff(),
	}, nil
}

func (service *ProductServiceCache) Delete(ctx context.Context, req *pb.GetProductByProductCodeRequest) (*emptypb.Empty, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[ProductServiceCache][Delete] problem in db transaction, err: ", err.Error())
		return new(emptypb.Empty), err
	}
	defer util.CommitOrRollback(tx)

	product, err := service.ProductRepository.FindByCode(ctx, tx, req.Code)
	if err != nil {
		log.Println("[ProductServiceCache][Delete][FindByCode] problem in getting from repository, err: ", err.Error())
		return new(emptypb.Empty), err
	}

	err = service.ProductRepository.Delete(ctx, tx, product.Code)
	if err != nil {
		log.Println("[ProductServiceCache][Delete][FindByCode] problem in getting from repository, err: ", err.Error())
		return new(emptypb.Empty), err
	}

	_ = service.StorageS3.DeleteFromAWS(product.Image)

	err = service.Redis.Del(ctx, "products:all", "products:admin", "products:similar", fmt.Sprintf("product:%s", req.Code))
	if err != nil {
		log.Println("[ProductServiceCache][Delete][Del] unable to deleting specific key cache, err: ", err.Error())
	}

	return new(emptypb.Empty), nil
}
