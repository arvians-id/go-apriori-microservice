package usecase

import (
	"context"
	"database/sql"
	pbapriori "github.com/arvians-id/go-apriori-microservice/adapter/pkg/apriori/pb"
	"github.com/arvians-id/go-apriori-microservice/adapter/pkg/product/pb"
	"github.com/arvians-id/go-apriori-microservice/model"
	"github.com/arvians-id/go-apriori-microservice/services/product-service/repository"
	"github.com/arvians-id/go-apriori-microservice/third-party/aws"
	"github.com/arvians-id/go-apriori-microservice/util"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/golang/protobuf/ptypes/empty"
	"log"
	"strings"
	"time"
)

type ProductService struct {
	ProductRepository repository.ProductRepository
	AprioriService    pbapriori.AprioriServiceClient
	StorageS3         aws.StorageS3
	DB                *sql.DB
}

func NewProductService(
	productRepository *repository.ProductRepository,
	aprioriService pbapriori.AprioriServiceClient,
	storageS3 *aws.StorageS3,
	db *sql.DB,
) pb.ProductServiceServer {
	return &ProductService{
		ProductRepository: *productRepository,
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

	return products, nil
}

func (service *ProductService) FindAll(ctx context.Context, req *pb.GetProductByFiltersRequest) (*pb.ListProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[ProductService][FindAll] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	products, err := service.ProductRepository.FindAll(ctx, tx, search, category)
	if err != nil {
		log.Println("[ProductService][FindAll][FindAll] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return products, nil
}

func (service *ProductService) FindAllBySimilarCategory(ctx context.Context, req *pb.GetProductByProductCodeRequest) (*pb.ListProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[ProductService][FindAllBySimilarCategory] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	product, err := service.ProductRepository.FindByCode(ctx, tx, code)
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

	var productResponses []*model.Product
	for _, productCategory := range productCategories {
		if productCategory.Code != code {
			productResponses = append(productResponses, productCategory)
		}
	}

	return productResponses, nil
}

func (service *ProductService) FindAllRecommendation(ctx context.Context, req *pb.GetProductByProductCodeRequest) (*pb.ListProductRecommendationResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[ProductService][FindAllRecommendation] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	product, err := service.ProductRepository.FindByCode(ctx, tx, code)
	if err != nil {
		log.Println("[ProductService][FindAllRecommendation][FindByCode] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	apriories, err := service.AprioriRepository.FindAllByActive(ctx, tx)
	if err != nil {
		log.Println("[ProductService][FindAllRecommendation][FindAllByActive] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	var productResponses []*model.ProductRecommendation
	for _, apriori := range apriories {
		productNames := strings.Split(apriori.Item, ",")
		var exists bool
		for _, productName := range productNames {
			if strings.ToLower(product.Name) == strings.TrimSpace(productName) {
				exists = true
			}
		}

		var totalPrice int
		if exists {
			for _, productName := range productNames {
				productByName, _ := service.ProductRepository.FindByName(ctx, tx, util.UpperWords(productName))
				totalPrice += productByName.Price
			}

			productResponses = append(productResponses, &model.ProductRecommendation{
				AprioriId:         apriori.IdApriori,
				AprioriCode:       apriori.Code,
				AprioriItem:       apriori.Item,
				AprioriDiscount:   apriori.Discount,
				ProductTotalPrice: totalPrice,
				PriceDiscount:     totalPrice - (totalPrice * int(apriori.Discount) / 100),
				AprioriImage:      apriori.Image,
			})
		}
	}

	return productResponses, nil
}

func (service *ProductService) FindByCode(ctx context.Context, req *pb.GetProductByProductCodeRequest) (*pb.GetProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[ProductService][FindByCode] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	productResponse, err := service.ProductRepository.FindByCode(ctx, tx, code)
	if err != nil {
		log.Println("[ProductService][FindByCode][FindByCode] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return productResponse, nil
}

func (service *ProductService) FindByName(ctx context.Context, req *pb.GetProductByProductNameRequest) (*pb.GetProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[ProductService][FindByCode] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	productResponse, err := service.ProductRepository.FindByCode(ctx, tx, code)
	if err != nil {
		log.Println("[ProductService][FindByCode][FindByCode] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return productResponse, nil
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
	if request.Image == "" {
		request.Image = "no-image.png"
	}

	productRequest := model.Product{
		Code:        request.Code,
		Name:        util.UpperWords(request.Name),
		Description: &request.Description,
		Price:       request.Price,
		Image:       &request.Image,
		Category:    util.UpperWords(request.Category),
		IsEmpty:     false,
		Mass:        request.Mass,
		CreatedAt:   timeNow,
		UpdatedAt:   timeNow,
	}

	productResponse, err := service.ProductRepository.Create(ctx, tx, &productRequest)
	if err != nil {
		log.Println("[ProductService][Create][Create] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return productResponse, nil
}

func (service *ProductService) Update(ctx context.Context, req *pb.UpdateProductRequest) (*pb.GetProductResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[ProductService][Update] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	product, err := service.ProductRepository.FindByCode(ctx, tx, request.Code)
	if err != nil {
		log.Println("[ProductService][Update][FindByCode] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	timeNow, err := time.Parse(util.TimeFormat, time.Now().Format(util.TimeFormat))
	if err != nil {
		log.Println("[ProductService][Update] problem in parsing to time, err: ", err.Error())
		return nil, err
	}

	product.Name = util.UpperWords(request.Name)
	product.Description = &request.Description
	product.Price = request.Price
	product.Category = util.UpperWords(request.Category)
	product.IsEmpty = request.IsEmpty
	product.Mass = request.Mass
	product.UpdatedAt = timeNow
	if request.Image != "" {
		_ = service.StorageS3.DeleteFromAWS(*product.Image)
		product.Image = &request.Image
	}

	productResponse, err := service.ProductRepository.Update(ctx, tx, product)
	if err != nil {
		log.Println("[ProductService][Update][Update] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return productResponse, nil
}

func (service *ProductService) Delete(ctx context.Context, req *pb.GetProductByProductCodeRequest) (*empty.Empty, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[ProductService][Delete] problem in db transaction, err: ", err.Error())
		return err
	}
	defer util.CommitOrRollback(tx)

	product, err := service.ProductRepository.FindByCode(ctx, tx, code)
	if err != nil {
		log.Println("[ProductService][Delete][FindByCode] problem in getting from repository, err: ", err.Error())
		return err
	}

	err = service.ProductRepository.Delete(ctx, tx, product.Code)
	if err != nil {
		log.Println("[ProductService][Delete][FindByCode] problem in getting from repository, err: ", err.Error())
		return err
	}

	_ = service.StorageS3.DeleteFromAWS(*product.Image)

	return nil
}
