package usecase

import (
	"context"
	"database/sql"
	pbproduct "github.com/arvians-id/go-apriori-microservice/adapter/pkg/product/pb"
	"github.com/arvians-id/go-apriori-microservice/adapter/pkg/transaction/pb"
	"github.com/arvians-id/go-apriori-microservice/model"
	"github.com/arvians-id/go-apriori-microservice/services/transaction-service/repository"
	"github.com/arvians-id/go-apriori-microservice/util"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"strings"
	"time"
)

type TransactionService struct {
	TransactionRepository repository.TransactionRepository
	ProductService        pbproduct.ProductServiceClient
	DB                    *sql.DB
}

func NewTransactionService(
	transactionRepository *repository.TransactionRepository,
	productService pbproduct.ProductServiceClient,
	db *sql.DB,
) pb.TransactionServiceClient {
	return &TransactionService{
		TransactionRepository: *transactionRepository,
		ProductService:        productService,
		DB:                    db,
	}
}

func (service *TransactionService) FindAll(ctx context.Context, empty *empty.Empty) (*pb.ListTransactionsResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[TransactionService][FindAll] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	transactions, err := service.TransactionRepository.FindAll(ctx, tx)
	if err != nil {
		log.Println("[TransactionService][FindAll][FindAll] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return transactions, nil
}

func (service *TransactionService) FindAllItemSet(ctx context.Context, req *pb.GetAllItemSetTransactionRequest) (*pb.ListTransactionsResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[TransactionService][FindAll] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	transactions, err := service.TransactionRepository.FindAll(ctx, tx)
	if err != nil {
		log.Println("[TransactionService][FindAll][FindAll] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return transactions, nil
}

func (service *TransactionService) FindByNoTransaction(ctx context.Context, req *pb.GetTransactionByNoTransactionRequest) (*pb.GetTransactionResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[TransactionService][FindByNoTransaction] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	transaction, err := service.TransactionRepository.FindByNoTransaction(ctx, tx, noTransaction)
	if err != nil {
		log.Println("[TransactionService][FindByNoTransaction][FindByNoTransaction] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return transaction, nil
}

func (service *TransactionService) Create(ctx context.Context, req *pb.CreateTransactionRequest) (*pb.GetTransactionResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[TransactionService][Create] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	timeNow, err := time.Parse(util.TimeFormat, time.Now().Format(util.TimeFormat))
	if err != nil {
		log.Println("[TransactionService][Create] problem in parsing to time, err: ", err.Error())
		return nil, err
	}

	noTransaction := util.CreateTransaction()
	if req.NoTransaction != nil {
		noTransaction = *req.NoTransaction
	}

	transactionRequest := model.Transaction{
		ProductName:   strings.ToLower(request.ProductName),
		CustomerName:  request.CustomerName,
		NoTransaction: noTransaction,
		CreatedAt:     timeNow,
		UpdatedAt:     timeNow,
	}

	transaction, err := service.TransactionRepository.Create(ctx, tx, &transactionRequest)
	if err != nil {
		log.Println("[TransactionService][Create][Create] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return transaction, nil
}

func (service *TransactionService) CreateByCSV(ctx context.Context, req *pb.CreateTransactionByCSVRequest) (*empty.Empty, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[TransactionService][CreateByCsv] problem in db transaction, err: ", err.Error())
		return err
	}
	defer util.CommitOrRollback(tx)

	var transactions []*model.Transaction
	for _, transaction := range data {
		createdAt, _ := time.Parse(util.TimeFormat, transaction[3]+" 00:00:00")
		transactions = append(transactions, &model.Transaction{
			ProductName:   transaction[0],
			CustomerName:  transaction[1],
			NoTransaction: transaction[2],
			CreatedAt:     createdAt,
			UpdatedAt:     createdAt,
		})
	}

	err = service.TransactionRepository.CreateByCsv(ctx, tx, transactions)
	if err != nil {
		log.Println("[TransactionService][CreateByCsv][CreateByCsv] problem in getting from repository, err: ", err.Error())
		return err
	}

	return nil
}

func (service *TransactionService) Update(ctx context.Context, req *pb.UpdateTransactionRequest) (*pb.GetTransactionResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[TransactionService][Update] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	// Find Transaction by number transaction
	transaction, err := service.TransactionRepository.FindByNoTransaction(ctx, tx, request.NoTransaction)
	if err != nil {
		log.Println("[TransactionService][Update][FindByNoTransaction] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	timeNow, err := time.Parse(util.TimeFormat, time.Now().Format(util.TimeFormat))
	if err != nil {
		log.Println("[TransactionService][Update] problem in parsing to time, err: ", err.Error())
		return nil, err
	}

	transaction.ProductName = strings.ToLower(request.ProductName)
	transaction.CustomerName = request.CustomerName
	transaction.NoTransaction = request.NoTransaction
	transaction.UpdatedAt = timeNow

	_, err = service.TransactionRepository.Update(ctx, tx, transaction)
	if err != nil {
		log.Println("[TransactionService][Update][Update] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return transaction, nil
}

func (service *TransactionService) Delete(ctx context.Context, req *pb.GetTransactionByNoTransactionRequest) (*emptypb.Empty, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[TransactionService][Delete] problem in db transaction, err: ", err.Error())
		return err
	}
	defer util.CommitOrRollback(tx)

	transaction, err := service.TransactionRepository.FindByNoTransaction(ctx, tx, noTransaction)
	if err != nil {
		log.Println("[TransactionService][Delete][FindByNoTransaction] problem in getting from repository, err: ", err.Error())
		return err
	}

	err = service.TransactionRepository.Delete(ctx, tx, transaction.NoTransaction)
	if err != nil {
		log.Println("[TransactionService][Delete][Delete] problem in getting from repository, err: ", err.Error())
		return err
	}

	return nil
}

func (service *TransactionService) Truncate(ctx context.Context, empty *empty.Empty) (*emptypb.Empty, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[TransactionService][Truncate] problem in db transaction, err: ", err.Error())
		return err
	}
	defer util.CommitOrRollback(tx)

	err = service.TransactionRepository.Truncate(ctx, tx)
	if err != nil {
		log.Println("[TransactionService][Truncate][Truncate] problem in getting from repository, err: ", err.Error())
		return err
	}

	return nil
}
