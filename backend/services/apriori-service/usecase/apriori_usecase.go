package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/arvians-id/go-apriori-microservice/adapter/pkg/apriori/pb"
	pbproduct "github.com/arvians-id/go-apriori-microservice/adapter/pkg/product/pb"
	pbtransaction "github.com/arvians-id/go-apriori-microservice/adapter/pkg/transaction/pb"
	"github.com/arvians-id/go-apriori-microservice/model"
	"github.com/arvians-id/go-apriori-microservice/services/apriori-service/repository"
	"github.com/arvians-id/go-apriori-microservice/third-party/aws"
	"github.com/arvians-id/go-apriori-microservice/util"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"math"
	"os"
	"strings"
	"time"
)

type AprioriService struct {
	AprioriRepository  repository.AprioriRepository
	ProductService     pbproduct.ProductServiceClient
	TransactionService pbtransaction.TransactionServiceClient
	StorageS3          aws.StorageS3
	DB                 *sql.DB
}

func NewAprioriService(
	aprioriRepository *repository.AprioriRepository,
	storageS3 *aws.StorageS3,
	db *sql.DB,
	productService pbproduct.ProductServiceClient,
	transactionService pbtransaction.TransactionServiceClient,
) pb.AprioriServiceServer {
	return &AprioriService{
		AprioriRepository:  *aprioriRepository,
		StorageS3:          *storageS3,
		DB:                 db,
		ProductService:     productService,
		TransactionService: transactionService,
	}
}

func (service *AprioriService) FindAll(ctx context.Context, empty *emptypb.Empty) (*pb.ListAprioriResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[AprioriService][FindAll] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	apriories, err := service.AprioriRepository.FindAll(ctx, tx)
	if err != nil {
		log.Println("[AprioriService][FindAll] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	var aprioriResponse []*pb.Apriori
	for _, apriori := range apriories {
		aprioriResponse = append(aprioriResponse, &pb.Apriori{
			IdApriori:   apriori.IdApriori,
			Code:        apriori.Code,
			Item:        apriori.Item,
			Discount:    apriori.Discount,
			Support:     apriori.Support,
			Confidence:  apriori.Confidence,
			RangeDate:   apriori.RangeDate,
			IsActive:    apriori.IsActive,
			Description: apriori.Description,
			Mass:        apriori.Mass,
			Image:       apriori.Image,
			CreatedAt:   timestamppb.New(apriori.CreatedAt),
		})
	}

	return &pb.ListAprioriResponse{
		Apriori: aprioriResponse,
	}, nil
}

func (service *AprioriService) FindAllByActive(ctx context.Context, empty *emptypb.Empty) (*pb.ListAprioriResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[AprioriService][FindAllByActive] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	apriories, err := service.AprioriRepository.FindAllByActive(ctx, tx)
	if err != nil {
		log.Println("[AprioriService][FindAllByActive][FindAllByActive] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	var aprioriResponse []*pb.Apriori
	for _, apriori := range apriories {
		aprioriResponse = append(aprioriResponse, &pb.Apriori{
			IdApriori:   apriori.IdApriori,
			Code:        apriori.Code,
			Item:        apriori.Item,
			Discount:    apriori.Discount,
			Support:     apriori.Support,
			Confidence:  apriori.Confidence,
			RangeDate:   apriori.RangeDate,
			IsActive:    apriori.IsActive,
			Description: apriori.Description,
			Mass:        apriori.Mass,
			Image:       apriori.Image,
			CreatedAt:   timestamppb.New(apriori.CreatedAt),
		})
	}

	return &pb.ListAprioriResponse{
		Apriori: aprioriResponse,
	}, nil
}

func (service *AprioriService) FindAllByCode(ctx context.Context, req *pb.GetAprioriByCodeRequest) (*pb.ListAprioriResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[AprioriService][FindAllByCode] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	apriories, err := service.AprioriRepository.FindAllByCode(ctx, tx, req.Code)
	if err != nil {
		log.Println("[AprioriService][FindAllByCode][FindAllByCode] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	var aprioriResponse []*pb.Apriori
	for _, apriori := range apriories {
		aprioriResponse = append(aprioriResponse, &pb.Apriori{
			IdApriori:   apriori.IdApriori,
			Code:        apriori.Code,
			Item:        apriori.Item,
			Discount:    apriori.Discount,
			Support:     apriori.Support,
			Confidence:  apriori.Confidence,
			RangeDate:   apriori.RangeDate,
			IsActive:    apriori.IsActive,
			Description: apriori.Description,
			Mass:        apriori.Mass,
			Image:       apriori.Image,
			CreatedAt:   timestamppb.New(apriori.CreatedAt),
		})
	}

	return &pb.ListAprioriResponse{
		Apriori: aprioriResponse,
	}, nil
}

func (service *AprioriService) FindByCodeAndId(ctx context.Context, req *pb.GetAprioriByCodeAndIdRequest) (*pb.GetAprioriByCodeAndIdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[AprioriService][FindByCodeAndId] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	apriori, err := service.AprioriRepository.FindByCodeAndId(ctx, tx, req.Code, int(req.Id))
	if err != nil {
		log.Println("[AprioriService][FindByCodeAndId][FindByCodeAndId] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	var totalPrice, mass int32
	productNames := strings.Split(apriori.Item, ",")
	for _, productName := range productNames {
		product, _ := service.ProductService.FindByName(ctx, &pbproduct.GetProductByProductNameRequest{
			Name: util.UpperWords(productName),
		})
		totalPrice += product.Product.Price
		mass += product.Product.Mass
	}

	return &pb.GetAprioriByCodeAndIdResponse{
		ProductRecommendation: &pbproduct.ProductRecommendation{
			AprioriId:          apriori.IdApriori,
			AprioriCode:        apriori.Code,
			AprioriItem:        apriori.Item,
			AprioriDiscount:    apriori.Discount,
			ProductTotalPrice:  totalPrice,
			PriceDiscount:      totalPrice - (totalPrice * int32(apriori.Discount) / 100),
			AprioriImage:       apriori.Image,
			Mass:               mass,
			AprioriDescription: apriori.Description,
		},
	}, nil
}

func (service *AprioriService) Create(ctx context.Context, req *pb.CreateAprioriRequest) (*emptypb.Empty, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[AprioriService][Create] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	timeNow, err := time.Parse(util.TimeFormat, time.Now().Format(util.TimeFormat))
	if err != nil {
		log.Println("[AprioriService][Create] problem in parsing to time, err: ", err.Error())
		return nil, err
	}

	var aprioriRequests []*model.Apriori
	code := util.RandomString(10)
	for _, requestItem := range req.CreateAprioriRequest {
		image := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/assets/%s", os.Getenv("AWS_BUCKET"), os.Getenv("AWS_REGION"), "no-image.png")
		aprioriRequests = append(aprioriRequests, &model.Apriori{
			Code:       code,
			Item:       requestItem.Item,
			Discount:   requestItem.Discount,
			Support:    requestItem.Support,
			Confidence: requestItem.Confidence,
			RangeDate:  requestItem.RangeDate,
			IsActive:   false,
			Image:      &image,
			CreatedAt:  timeNow,
		})
	}

	err = service.AprioriRepository.Create(ctx, tx, aprioriRequests)
	if err != nil {
		log.Println("[AprioriService][Create][Create] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return nil, nil
}

func (service *AprioriService) Update(ctx context.Context, req *pb.UpdateAprioriRequest) (*pb.GetAprioriResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[AprioriService][Update] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	apriori, err := service.AprioriRepository.FindByCodeAndId(ctx, tx, req.Code, int(req.IdApriori))
	if err != nil {
		log.Println("[AprioriService][Update][FindByCodeAndId] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	image := apriori.Image
	if req.Image != "" {
		_ = service.StorageS3.DeleteFromAWS(*apriori.Image)
		image = &req.Image
	}

	apriori.Image = image
	apriori.Description = &req.Description

	_, err = service.AprioriRepository.Update(ctx, tx, apriori)
	if err != nil {
		log.Println("[AprioriService][Update][Update] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return &pb.GetAprioriResponse{
		Apriori: apriori.ToProtoBuf(),
	}, nil
}

func (service *AprioriService) UpdateStatus(ctx context.Context, req *pb.GetAprioriByCodeRequest) (*emptypb.Empty, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[AprioriService][UpdateStatus] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	apriories, err := service.AprioriRepository.FindAllByCode(ctx, tx, req.Code)
	if err != nil {
		log.Println("[AprioriService][UpdateStatus][FindAllByCode] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	err = service.AprioriRepository.UpdateAllStatus(ctx, tx, false)
	if err != nil {
		log.Println("[AprioriService][UpdateStatus][UpdateAllStatus] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	status := true
	if apriories[0].IsActive {
		status = false
	}

	err = service.AprioriRepository.UpdateStatusByCode(ctx, tx, apriories[0].Code, status)
	if err != nil {
		log.Println("[AprioriService][UpdateStatus][UpdateStatusByCode] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return nil, nil
}

func (service *AprioriService) Delete(ctx context.Context, req *pb.GetAprioriByCodeRequest) (*emptypb.Empty, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[AprioriService][Delete] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	apriories, err := service.AprioriRepository.FindAllByCode(ctx, tx, req.Code)
	if err != nil {
		log.Println("[AprioriService][Delete][FindAllByCode] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	err = service.AprioriRepository.Delete(ctx, tx, apriories[0].Code)
	if err != nil {
		log.Println("[AprioriService][Delete][Delete] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	for _, apriori := range apriories {
		_ = service.StorageS3.DeleteFromAWS(*apriori.Image)
	}

	return nil, nil
}

func (service *AprioriService) Generate(ctx context.Context, req *pb.GenerateAprioriRequest) (*pb.GetGenerateAprioriResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[AprioriService][Generate] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	var apriori []*model.GenerateApriori
	// Get all transaction from database
	transactionsSet, err := service.TransactionService.FindAllItemSet(ctx, &pbtransaction.GetAllItemSetTransactionRequest{
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	})
	if err != nil {
		log.Println("[AprioriService][Generate][FindAllItemSet] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	var transactionsModel []*model.Transaction
	for _, transaction := range transactionsSet.Transaction {
		transactionsModel = append(transactionsModel, &model.Transaction{
			IdTransaction: transaction.IdTransaction,
			ProductName:   transaction.ProductName,
			CustomerName:  transaction.CustomerName,
			NoTransaction: transaction.NoTransaction,
			CreatedAt:     transaction.CreatedAt.AsTime(),
			UpdatedAt:     transaction.UpdatedAt.AsTime(),
		})
	}

	// Find first item set
	transactions, productName, propertyProduct := util.FindFirstItemSet(transactionsModel, float64(req.MinimumSupport))

	// Handle random maps problem
	oneSet, support, totalTransaction, isEligible, cleanSet := util.HandleMapsProblem(propertyProduct, float64(req.MinimumSupport))

	// Get one item set
	for i := 0; i < len(oneSet); i++ {
		apriori = append(apriori, &model.GenerateApriori{
			ItemSet:     []string{oneSet[i]},
			Support:     float32(support[i]),
			Iterate:     1,
			Transaction: int32(totalTransaction[i]),
			Description: isEligible[i],
			RangeDate:   req.StartDate + " - " + req.EndDate,
		})
	}

	oneSet = cleanSet
	// Looking for more than one item set
	var iterate int
	var dataTemp [][]string
	for {
		for i := 0; i < len(oneSet)-iterate; i++ {
			for j := i + 1; j < len(oneSet); j++ {
				var iterateCandidate []string

				iterateCandidate = append(iterateCandidate, oneSet[i])
				for z := 1; z <= iterate; z++ {
					iterateCandidate = append(iterateCandidate, oneSet[i+z])
				}
				iterateCandidate = append(iterateCandidate, oneSet[j])

				dataTemp = append(dataTemp, iterateCandidate)
			}
		}
		// Filter when the slice has duplicate values
		var cleanValues [][]string
		for i := 0; i < len(dataTemp); i++ {
			if isDuplicate := util.IsDuplicate(dataTemp[i]); !isDuplicate {
				cleanValues = append(cleanValues, dataTemp[i])
			}
		}
		dataTemp = cleanValues
		// Filter candidates by comparing slice to slice
		dataTemp = util.FilterCandidateInSlice(dataTemp)

		// Find item set by minimum support
		for i := 0; i < len(dataTemp); i++ {
			countCandidates := util.FindCandidate(dataTemp[i], transactions)
			result := float64(countCandidates) / float64(len(transactionsModel)) * 100
			if result >= float64(req.MinimumSupport) {
				apriori = append(apriori, &model.GenerateApriori{
					ItemSet:     dataTemp[i],
					Support:     float32(math.Round(result*100) / 100),
					Iterate:     int32(iterate + 2),
					Transaction: int32(countCandidates),
					Description: "Eligible",
					RangeDate:   req.StartDate + " - " + req.EndDate,
				})
			} else {
				apriori = append(apriori, &model.GenerateApriori{
					ItemSet:     dataTemp[i],
					Support:     float32(math.Round(result*100) / 100),
					Iterate:     int32(iterate + 2),
					Transaction: int32(countCandidates),
					Description: "Not Eligible",
					RangeDate:   req.StartDate + " - " + req.EndDate,
				})
			}
		}

		// Convert dataTemp slice of slice to one slice
		var test []string
		for i := 0; i < len(dataTemp); i++ {
			test = append(test, dataTemp[i]...)
		}
		oneSet = test

		// After finish operating, then clean the array
		dataTemp = [][]string{}

		var checkClean bool
		for _, value := range apriori {
			if value.Iterate == int32(iterate+2) && value.Description == "Eligible" {
				checkClean = true
				break
			}
		}

		countIterate := 0
		for i := 0; i < len(apriori); i++ {
			if apriori[i].Iterate == int32(iterate+2) {
				countIterate++
			}
		}

		if checkClean == false {
			for i := 0; i < len(apriori); i++ {
				if apriori[i].Iterate == int32(iterate+2) {
					apriori = append(apriori[:i], apriori[i+countIterate:]...)
				}
			}
			break
		}

		// if nothing else is sent to the candidate, then break it
		if int32(iterate+2) > apriori[len(apriori)-1].Iterate {
			break
		}

		// Add increment, if any candidate is submitted
		iterate++
	}

	// Find Association rules
	// Set confidence
	confidence := util.FindConfidence(apriori, productName, float64(req.MinimumSupport), float64(req.MinimumConfidence))

	// Set discount
	discount := util.FindDiscount(confidence, float64(req.MinimumDiscount), float64(req.MaximumDiscount))

	//// Remove last element in apriori as many association rules
	//for i := 0; i < len(discount); i++ {
	//	apriori = apriori[:len(apriori)-1]
	//}

	// Replace the last item set and add discount and confidence
	for i := 0; i < len(discount); i++ {
		if discount[i].Confidence >= req.MinimumConfidence {
			apriori = append(apriori, &model.GenerateApriori{
				ItemSet:     discount[i].ItemSet,
				Support:     float32(math.Round(float64(discount[i].Support*100)) / 100),
				Iterate:     discount[i].Iterate + 1,
				Transaction: discount[i].Transaction,
				Confidence:  float32(math.Round(float64(discount[i].Confidence*100)) / 100),
				Discount:    discount[i].Discount,
				Description: "Rules",
				RangeDate:   req.StartDate + " - " + req.EndDate,
			})
		}
	}

	var aprioriGenerateResponse []*pb.AprioriGenerate
	for _, value := range apriori {
		aprioriGenerateResponse = append(aprioriGenerateResponse, &pb.AprioriGenerate{
			ItemSet:     value.ItemSet,
			Support:     value.Support,
			Iterate:     value.Iterate,
			Transaction: value.Transaction,
			Confidence:  value.Confidence,
			Discount:    value.Discount,
			Description: value.Description,
			RangeDate:   value.RangeDate,
		})
	}

	return &pb.GetGenerateAprioriResponse{
		AprioriGenerate: aprioriGenerateResponse,
	}, nil
}