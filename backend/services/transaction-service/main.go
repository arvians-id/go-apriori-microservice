package main

import (
	"database/sql"
	"fmt"
	"github.com/arvians-id/go-apriori-microservice/adapter/pkg/transaction/pb"
	"github.com/arvians-id/go-apriori-microservice/config"
	"github.com/arvians-id/go-apriori-microservice/services/transaction-service/repository"
	"github.com/arvians-id/go-apriori-microservice/services/transaction-service/usecase"
	"google.golang.org/grpc"
	"log"
	"net"
)

func NewInitializedDatabase(configuration *config.Config) (*sql.DB, error) {
	db, err := config.NewPostgresSQL(configuration)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewInitializedServices(configuration *config.Config) (pb.TransactionServiceServer, error) {
	db, err := NewInitializedDatabase(configuration)
	if err != nil {
		return nil, err
	}

	transactionRepository := repository.NewTransactionRepository()
	transactionService := usecase.NewTransactionService(transactionRepository, db)

	return transactionService, nil
}

func main() {
	configuration, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	connection, err := net.Listen("tcp", configuration.TransactionSvcUrl)
	if err != nil {
		log.Fatalln("Failed at listening", err)
	}

	fmt.Println("Transaction service is running on port", configuration.TransactionSvcUrl)

	services, err := NewInitializedServices(configuration)
	if err != nil {
		log.Fatalln("Failed at initializing services", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTransactionServiceServer(grpcServer, services)

	if err := grpcServer.Serve(connection); err != nil {
		log.Fatalln("Failed at serving", err)
	}
}
