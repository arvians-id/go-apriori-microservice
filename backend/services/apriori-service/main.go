package main

import (
	"database/sql"
	"fmt"
	"github.com/arvians-id/go-apriori-microservice/services/apriori-service/client"
	"github.com/arvians-id/go-apriori-microservice/services/apriori-service/config"
	"github.com/arvians-id/go-apriori-microservice/services/apriori-service/pb"
	"github.com/arvians-id/go-apriori-microservice/services/apriori-service/repository"
	"github.com/arvians-id/go-apriori-microservice/services/apriori-service/third-party/aws"
	"github.com/arvians-id/go-apriori-microservice/services/apriori-service/usecase"
	"google.golang.org/grpc"
	"log"
	"net"
	"strings"
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
	aprioriService := usecase.NewAprioriService(aprioriRepository, storageS3, db, productClient, transactionClient, configuration)

	return aprioriService, nil
}

func main() {
	configuration, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	port := ":" + strings.Split(configuration.AprioriSvcUrl, ":")[1]
	connection, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln("Failed at listening", err)
	}

	fmt.Println("Apriori service is running on port", port)

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
