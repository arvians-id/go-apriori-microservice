package main

import (
	"database/sql"
	"fmt"
	"github.com/arvians-id/go-apriori-microservice/adapter/pb"
	"github.com/arvians-id/go-apriori-microservice/config"
	"github.com/arvians-id/go-apriori-microservice/services/payment-service/client"
	"github.com/arvians-id/go-apriori-microservice/services/payment-service/repository"
	"github.com/arvians-id/go-apriori-microservice/services/payment-service/usecase"
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

func NewInitializedServices(configuration *config.Config) (pb.PaymentServiceServer, error) {
	db, err := NewInitializedDatabase(configuration)
	if err != nil {
		return nil, err
	}

	userOrderClient := client.NewUserOrderServiceClient(configuration)
	transactionClient := client.NewTransactionServiceClient(configuration)

	paymentRepository := repository.NewPaymentRepository()
	paymentService := usecase.NewPaymentService(configuration, paymentRepository, userOrderClient, transactionClient, db)

	return paymentService, nil
}

func main() {
	configuration, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	connection, err := net.Listen("tcp", configuration.PaymentSvcUrl)
	if err != nil {
		log.Fatalln("Failed at listening", err)
	}

	fmt.Println("Payment service is running on port", configuration.PaymentSvcUrl)

	services, err := NewInitializedServices(configuration)
	if err != nil {
		log.Fatalln("Failed at initializing services", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPaymentServiceServer(grpcServer, services)

	if err := grpcServer.Serve(connection); err != nil {
		log.Fatalln("Failed at serving", err)
	}
}
