package usecase

import (
	"context"
	"database/sql"
	"github.com/arvians-id/go-apriori-microservice/adapter/pb"
	"github.com/arvians-id/go-apriori-microservice/model"
	"github.com/arvians-id/go-apriori-microservice/services/product-service/client"
	"github.com/arvians-id/go-apriori-microservice/services/product-service/repository"
	"github.com/arvians-id/go-apriori-microservice/third-party/aws"
	"github.com/arvians-id/go-apriori-microservice/util"
	"github.com/golang/protobuf/ptypes/empty"
	"log"
	"strings"
	"time"
)

type ProductService struct {
	ProductRepository repository.ProductRepository
	AprioriService    client.AprioriServiceClient
	StorageS3         aws.StorageS3
	DB                *sql.DB
}

func NewProductService(
	productRepository repository.ProductRepository,
	aprioriService client.AprioriServiceClient,
	storageS3 *aws.StorageS3,
	db *sql.DB,
) pb.ProductServiceServer {
	return &ProductService{
		ProductRepository: productRepository,
		AprioriService:    aprioriService,
		StorageS3:         *storageS3,
		DB:                db,
	}
}

func (service *ProductService) FindAllByAdmin(ctx context.Context, empty *empty.Empty) (*pb.ListProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[ProductService][FindAllByAdmin] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	products, err := service.ProductRepository.FindAllByAdmin(ctx, tx)
	if err != nil {
		log.Println("[ProductService][FindAllByAdmin][FindAllByAdmin] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	var productListResponse []*pb.Product
	for _, product := range products {
		productListResponse = append(productListResponse, product.ToProtoBuff())
	}

	return &pb.ListProductResponse{
		Product: productListResponse,
	}, nil
}

func (service *ProductService) FindAll(ctx context.Context, req *pb.GetProductByFiltersRequest) (*pb.ListProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[ProductService][FindAll] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	products, err := service.ProductRepository.FindAll(ctx, tx, req.Search, req.Category)
	if err != nil {
		log.Println("[ProductService][FindAll][FindAll] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	var productListResponse []*pb.Product
	for _, product := range products {
		productListResponse = append(productListResponse, product.ToProtoBuff())
	}

	return &pb.ListProductResponse{
		Product: productListResponse,
	}, nil
}

func (service *ProductService) FindAllBySimilarCategory(ctx context.Context, req *pb.GetProductByProductCodeRequest) (*pb.ListProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[ProductService][FindAllBySimilarCategory] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	product, err := service.ProductRepository.FindByCode(ctx, tx, req.Code)
	if err != nil {
		log.Println("[ProductService][FindAllBySimilarCategory][FindByCode] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	categoryArray := strings.Split(product.Category, ", ")
	categoryString := strings.Join(categoryArray, "|")
	productCategories, err := service.ProductRepository.FindAllBySimilarCategory(ctx, tx, categoryString)
	if err != nil {
		log.Println("[ProductService][FindAllBySimilarCategory][FindAllBySimilarCategory] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	var productListResponse []*pb.Product
	for _, productCategory := range productCategories {
		if productCategory.Code != req.Code {
			productListResponse = append(productListResponse, productCategory.ToProtoBuff())
		}
	}

	return &pb.ListProductResponse{
		Product: productListResponse,
	}, nil
}

func (service *ProductService) FindAllRecommendation(ctx context.Context, req *pb.GetProductByProductCodeRequest) (*pb.ListProductRecommendationResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[ProductService][FindAllRecommendation] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	product, err := service.ProductRepository.FindByCode(ctx, tx, req.Code)
	if err != nil {
		log.Println("[ProductService][FindAllRecommendation][FindByCode] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	apriories, err := service.AprioriService.FindAllByActive(ctx)
	if err != nil {
		log.Println("[ProductService][FindAllRecommendation][FindAllByActive] problem in getting from repository, err: ", err.Error())
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

func (service *ProductService) FindByCode(ctx context.Context, req *pb.GetProductByProductCodeRequest) (*pb.GetProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[ProductService][FindByCode] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	productResponse, err := service.ProductRepository.FindByCode(ctx, tx, req.Code)
	if err != nil {
		log.Println("[ProductService][FindByCode][FindByCode] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return &pb.GetProductResponse{
		Product: productResponse.ToProtoBuff(),
	}, nil
}

func (service *ProductService) FindByName(ctx context.Context, req *pb.GetProductByProductNameRequest) (*pb.GetProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[ProductService][FindByCode] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	productResponse, err := service.ProductRepository.FindByName(ctx, tx, req.Name)
	if err != nil {
		log.Println("[ProductService][FindByCode][FindByCode] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return &pb.GetProductResponse{
		Product: productResponse.ToProtoBuff(),
	}, nil
}

func (service *ProductService) Create(ctx context.Context, req *pb.CreateProductRequest) (*pb.GetProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[ProductService][Create] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	timeNow, err := time.Parse(util.TimeFormat, time.Now().Format(util.TimeFormat))
	if err != nil {
		log.Println("[ProductService][Create] problem in parsing to time, err: ", err.Error())
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
		log.Println("[ProductService][Create][Create] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return &pb.GetProductResponse{
		Product: productResponse.ToProtoBuff(),
	}, nil
}

func (service *ProductService) Update(ctx context.Context, req *pb.UpdateProductRequest) (*pb.GetProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[ProductService][Update] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	product, err := service.ProductRepository.FindByCode(ctx, tx, req.Code)
	if err != nil {
		log.Println("[ProductService][Update][FindByCode] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	timeNow, err := time.Parse(util.TimeFormat, time.Now().Format(util.TimeFormat))
	if err != nil {
		log.Println("[ProductService][Update] problem in parsing to time, err: ", err.Error())
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
		log.Println("[ProductService][Update][Update] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return &pb.GetProductResponse{
		Product: productResponse.ToProtoBuff(),
	}, nil
}

func (service *ProductService) Delete(ctx context.Context, req *pb.GetProductByProductCodeRequest) (*empty.Empty, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[ProductService][Delete] problem in db transaction, err: ", err.Error())
		return new(empty.Empty), err
	}
	defer util.CommitOrRollback(tx)

	product, err := service.ProductRepository.FindByCode(ctx, tx, req.Code)
	if err != nil {
		log.Println("[ProductService][Delete][FindByCode] problem in getting from repository, err: ", err.Error())
		return new(empty.Empty), err
	}

	err = service.ProductRepository.Delete(ctx, tx, product.Code)
	if err != nil {
		log.Println("[ProductService][Delete][FindByCode] problem in getting from repository, err: ", err.Error())
		return new(empty.Empty), err
	}

	_ = service.StorageS3.DeleteFromAWS(product.Image)

	return new(empty.Empty), nil
}
