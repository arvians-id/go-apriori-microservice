package main

import (
	"database/sql"
	"fmt"
	"github.com/arvians-id/go-apriori-microservice/adapter/pkg/user-order/pb"
	"github.com/arvians-id/go-apriori-microservice/config"
	"github.com/arvians-id/go-apriori-microservice/services/user-order-service/client"
	"github.com/arvians-id/go-apriori-microservice/services/user-order-service/repository"
	"github.com/arvians-id/go-apriori-microservice/services/user-order-service/usecase"
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

func NewInitializedServices(configuration *config.Config) (pb.UserOrderServiceServer, error) {
	db, err := NewInitializedDatabase(configuration)
	if err != nil {
		return nil, err
	}

	userClient := client.NewUserServiceClient(configuration)

	userOrderRepository := repository.NewUserOrderRepository()
	userOrderService := usecase.NewUserOrderService(userOrderRepository, userClient, db)

	return userOrderService, nil
}

func main() {
	configuration, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	connection, err := net.Listen("tcp", configuration.UserOrderSvcUrl)
	if err != nil {
		log.Fatalln("Failed at listening", err)
	}

	fmt.Println("User Order service is running on port", configuration.UserOrderSvcUrl)

	services, err := NewInitializedServices(configuration)
	if err != nil {
		log.Fatalln("Failed at initializing services", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserOrderServiceServer(grpcServer, services)

	if err := grpcServer.Serve(connection); err != nil {
		log.Fatalln("Failed at serving", err)
	}
}
