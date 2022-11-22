package main

import (
	"database/sql"
	"fmt"
	"github.com/arvians-id/go-apriori-microservice/services/product-service/client"
	"github.com/arvians-id/go-apriori-microservice/services/product-service/config"
	"github.com/arvians-id/go-apriori-microservice/services/product-service/pb"
	"github.com/arvians-id/go-apriori-microservice/services/product-service/repository"
	"github.com/arvians-id/go-apriori-microservice/services/product-service/third-party/aws"
	"github.com/arvians-id/go-apriori-microservice/services/product-service/third-party/redis"
	"github.com/arvians-id/go-apriori-microservice/services/product-service/usecase"
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

func NewInitializedServices(configuration *config.Config) (pb.ProductServiceServer, error) {
	db, err := NewInitializedDatabase(configuration)
	if err != nil {
		return nil, err
	}

	storageS3 := aws.NewStorageS3(configuration)
	redisLib := redis.NewCacheService(configuration)

	aprioriClient := client.NewAprioriServiceClient(configuration)

	productRepository := repository.NewProductRepository()
	productService := usecase.NewProductServiceCache(productRepository, aprioriClient, storageS3, redisLib, db)

	return productService, nil
}

func main() {
	configuration, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	port := ":" + strings.Split(configuration.ProductSvcUrl, ":")[1]
	connection, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln("Failed at listening", err)
	}

	fmt.Println("Product service is running on port", port)

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
