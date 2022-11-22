package usecase

import (
	"context"
	"database/sql"
	"github.com/arvians-id/go-apriori-microservice/services/transaction-service/model"
	"github.com/arvians-id/go-apriori-microservice/services/transaction-service/pb"
	"github.com/arvians-id/go-apriori-microservice/services/transaction-service/repository"
	"github.com/arvians-id/go-apriori-microservice/services/transaction-service/util"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"strings"
	"time"
)

type TransactionService struct {
	TransactionRepository repository.TransactionRepository
	DB                    *sql.DB
}

func NewTransactionService(
	transactionRepository repository.TransactionRepository,
	db *sql.DB,
) pb.TransactionServiceServer {
	return &TransactionService{
		TransactionRepository: transactionRepository,
		DB:                    db,
	}
}

func (service *TransactionService) FindAll(ctx context.Context, empty *emptypb.Empty) (*pb.ListTransactionsResponse, error) {
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

	var transactionListResponse []*pb.Transaction
	for _, transaction := range transactions {
		transactionListResponse = append(transactionListResponse, transaction.ToProtoBuff())
	}

	return &pb.ListTransactionsResponse{
		Transaction: transactionListResponse,
	}, nil
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

	var transactionListResponse []*pb.Transaction
	for _, transaction := range transactions {
		transactionListResponse = append(transactionListResponse, transaction.ToProtoBuff())
	}

	return &pb.ListTransactionsResponse{
		Transaction: transactionListResponse,
	}, nil
}

func (service *TransactionService) FindByNoTransaction(ctx context.Context, req *pb.GetTransactionByNoTransactionRequest) (*pb.GetTransactionResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[TransactionService][FindByNoTransaction] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	transaction, err := service.TransactionRepository.FindByNoTransaction(ctx, tx, req.NoTransaction)
	if err != nil {
		log.Println("[TransactionService][FindByNoTransaction][FindByNoTransaction] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return &pb.GetTransactionResponse{
		Transaction: transaction.ToProtoBuff(),
	}, nil
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
		ProductName:   strings.ToLower(req.ProductName),
		CustomerName:  req.CustomerName,
		NoTransaction: noTransaction,
		CreatedAt:     timeNow,
		UpdatedAt:     timeNow,
	}

	transaction, err := service.TransactionRepository.Create(ctx, tx, &transactionRequest)
	if err != nil {
		log.Println("[TransactionService][Create][Create] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return &pb.GetTransactionResponse{
		Transaction: transaction.ToProtoBuff(),
	}, nil
}

func (service *TransactionService) CreateByCSV(ctx context.Context, req *pb.CreateTransactionByCSVRequest) (*emptypb.Empty, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[TransactionService][CreateByCsv] problem in db transaction, err: ", err.Error())
		return new(emptypb.Empty), err
	}
	defer util.CommitOrRollback(tx)

	data, err := util.OpenCsvFile(req.FilePath)
	if err != nil {
		log.Println("[TransactionService][CreateByCsv] problem in db transaction, err: ", err.Error())
		return new(emptypb.Empty), err
	}

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
		return new(emptypb.Empty), err
	}

	return new(emptypb.Empty), nil
}

func (service *TransactionService) Update(ctx context.Context, req *pb.UpdateTransactionRequest) (*pb.GetTransactionResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[TransactionService][Update] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	// Find Transaction by number transaction
	transaction, err := service.TransactionRepository.FindByNoTransaction(ctx, tx, req.NoTransaction)
	if err != nil {
		log.Println("[TransactionService][Update][FindByNoTransaction] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	timeNow, err := time.Parse(util.TimeFormat, time.Now().Format(util.TimeFormat))
	if err != nil {
		log.Println("[TransactionService][Update] problem in parsing to time, err: ", err.Error())
		return nil, err
	}

	transaction.ProductName = strings.ToLower(req.ProductName)
	transaction.CustomerName = req.CustomerName
	transaction.NoTransaction = req.NoTransaction
	transaction.UpdatedAt = timeNow

	_, err = service.TransactionRepository.Update(ctx, tx, transaction)
	if err != nil {
		log.Println("[TransactionService][Update][Update] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return &pb.GetTransactionResponse{
		Transaction: transaction.ToProtoBuff(),
	}, nil
}

func (service *TransactionService) Delete(ctx context.Context, req *pb.GetTransactionByNoTransactionRequest) (*emptypb.Empty, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[TransactionService][Delete] problem in db transaction, err: ", err.Error())
		return new(emptypb.Empty), err
	}
	defer util.CommitOrRollback(tx)

	transaction, err := service.TransactionRepository.FindByNoTransaction(ctx, tx, req.NoTransaction)
	if err != nil {
		log.Println("[TransactionService][Delete][FindByNoTransaction] problem in getting from repository, err: ", err.Error())
		return new(emptypb.Empty), err
	}

	err = service.TransactionRepository.Delete(ctx, tx, transaction.NoTransaction)
	if err != nil {
		log.Println("[TransactionService][Delete][Delete] problem in getting from repository, err: ", err.Error())
		return new(emptypb.Empty), err
	}

	return new(emptypb.Empty), nil
}

func (service *TransactionService) Truncate(ctx context.Context, emptys *emptypb.Empty) (*emptypb.Empty, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[TransactionService][Truncate] problem in db transaction, err: ", err.Error())
		return new(emptypb.Empty), err
	}
	defer util.CommitOrRollback(tx)

	err = service.TransactionRepository.Truncate(ctx, tx)
	if err != nil {
		log.Println("[TransactionService][Truncate][Truncate] problem in getting from repository, err: ", err.Error())
		return new(emptypb.Empty), err
	}

	return new(emptypb.Empty), nil
}
