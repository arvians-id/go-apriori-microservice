package usecase

import (
	"context"
	"database/sql"
	"github.com/arvians-id/apriori/cmd/library/aws"
	"github.com/arvians-id/apriori/internal/http/presenter/request"
	"github.com/arvians-id/apriori/internal/model"
	"github.com/arvians-id/apriori/internal/repository"
	"github.com/arvians-id/apriori/util"
	"log"
	"strings"
	"time"
)

type ProductServiceImpl struct {
	ProductRepository repository.ProductRepository
	AprioriRepository repository.AprioriRepository
	StorageS3         aws.StorageS3
	DB                *sql.DB
}

func NewProductService(
	productRepository *repository.ProductRepository,
	aprioriRepository *repository.AprioriRepository,
	storageS3 *aws.StorageS3,
	db *sql.DB,
) ProductService {
	return &ProductServiceImpl{
		ProductRepository: *productRepository,
		AprioriRepository: *aprioriRepository,
		StorageS3:         *storageS3,
		DB:                db,
	}
}

func (service *ProductServiceImpl) FindAllByAdmin(ctx context.Context) ([]*model.Product, error) {
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

func (service *ProductServiceImpl) FindAll(ctx context.Context, search string, category string) ([]*model.Product, error) {
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

func (service *ProductServiceImpl) FindAllBySimilarCategory(ctx context.Context, code string) ([]*model.Product, error) {
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

func (service *ProductServiceImpl) FindAllRecommendation(ctx context.Context, code string) ([]*model.ProductRecommendation, error) {
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

func (service *ProductServiceImpl) FindByCode(ctx context.Context, code string) (*model.Product, error) {
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

func (service *ProductServiceImpl) Create(ctx context.Context, request *request.CreateProductRequest) (*model.Product, error) {
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

func (service *ProductServiceImpl) Update(ctx context.Context, request *request.UpdateProductRequest) (*model.Product, error) {
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

func (service *ProductServiceImpl) Delete(ctx context.Context, code string) error {
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