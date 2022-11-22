package main

import (
	"database/sql"
	"fmt"
	"github.com/arvians-id/go-apriori-microservice/services/category-service/config"
	"github.com/arvians-id/go-apriori-microservice/services/category-service/pb"
	"github.com/arvians-id/go-apriori-microservice/services/category-service/repository"
	"github.com/arvians-id/go-apriori-microservice/services/category-service/third-party/redis"
	"github.com/arvians-id/go-apriori-microservice/services/category-service/usecase"
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

func NewInitializedServices(configuration *config.Config) (pb.CategoryServiceServer, error) {
	db, err := NewInitializedDatabase(configuration)
	if err != nil {
		return nil, err
	}

	// Third party
	redisLib := redis.NewCacheService(configuration)

	// Main App
	categoryRepository := repository.NewCategoryRepository()
	categoryService := usecase.NewCategoryServiceCache(categoryRepository, redisLib, db)

	return categoryService, nil
}

func main() {
	configuration, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	port := ":" + strings.Split(configuration.CategorySvcUrl, ":")[1]
	connection, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln("Failed at listening", err)
	}

	fmt.Println("Category service is running on port", port)

	services, err := NewInitializedServices(configuration)
	if err != nil {
		log.Fatalln("Failed at initializing services", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCategoryServiceServer(grpcServer, services)

	if err := grpcServer.Serve(connection); err != nil {
		log.Fatalln("Failed at serving", err)
	}
}
