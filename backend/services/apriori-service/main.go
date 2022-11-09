package main

import (
	"database/sql"
	"fmt"
	"github.com/arvians-id/go-apriori-microservice/adapter/pb"
	"github.com/arvians-id/go-apriori-microservice/config"
	"github.com/arvians-id/go-apriori-microservice/services/apriori-service/client"
	"github.com/arvians-id/go-apriori-microservice/services/apriori-service/repository"
	"github.com/arvians-id/go-apriori-microservice/services/apriori-service/usecase"
	"github.com/arvians-id/go-apriori-microservice/third-party/aws"
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

func NewInitializedServices(configuration *config.Config) (pb.AprioriServiceServer, error) {
	db, err := NewInitializedDatabase(configuration)
	if err != nil {
		return nil, err
	}

	storageS3 := aws.NewStorageS3(configuration)

	productClient := client.NewProductServiceClient(configuration)
	transactionClient := client.NewTransactionServiceClient(configuration)

	aprioriRepository := repository.NewAprioriRepository()
	aprioriService := usecase.NewAprioriService(aprioriRepository, storageS3, db, productClient, transactionClient)

	return aprioriService, nil
}

func main() {
	configuration, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	connection, err := net.Listen("tcp", configuration.AprioriSvcUrl)
	if err != nil {
		log.Fatalln("Failed at listening", err)
	}

	fmt.Println("Apriori service is running on port", configuration.AprioriSvcUrl)

	services, err := NewInitializedServices(configuration)
	if err != nil {
		log.Fatalln("Failed at initializing services", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAprioriServiceServer(grpcServer, services)

	if err := grpcServer.Serve(connection); err != nil {
		log.Fatalln("Failed at serving", err)
	}
}
