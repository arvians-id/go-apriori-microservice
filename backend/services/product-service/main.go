package main

import (
	"database/sql"
	"fmt"
	"github.com/arvians-id/go-apriori-microservice/adapter/pb"
	"github.com/arvians-id/go-apriori-microservice/config"
	"github.com/arvians-id/go-apriori-microservice/services/product-service/client"
	"github.com/arvians-id/go-apriori-microservice/services/product-service/repository"
	"github.com/arvians-id/go-apriori-microservice/services/product-service/usecase"
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

func NewInitializedServices(configuration *config.Config) (pb.ProductServiceServer, error) {
	db, err := NewInitializedDatabase(configuration)
	if err != nil {
		return nil, err
	}

	storageS3 := aws.NewStorageS3(configuration)

	aprioriClient := client.NewAprioriServiceClient(configuration)

	productRepository := repository.NewProductRepository()
	productService := usecase.NewProductService(productRepository, aprioriClient, storageS3, db)

	return productService, nil
}

func main() {
	configuration, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	connection, err := net.Listen("tcp", configuration.ProductSvcUrl)
	if err != nil {
		log.Fatalln("Failed at listening", err)
	}

	fmt.Println("Product service is running on port", configuration.ProductSvcUrl)

	services, err := NewInitializedServices(configuration)
	if err != nil {
		log.Fatalln("Failed at initializing services", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterProductServiceServer(grpcServer, services)

	if err := grpcServer.Serve(connection); err != nil {
		log.Fatalln("Failed at serving", err)
	}
}
