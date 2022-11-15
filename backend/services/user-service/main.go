package main

import (
	"database/sql"
	"fmt"
	"github.com/arvians-id/go-apriori-microservice/services/user-service/config"
	"github.com/arvians-id/go-apriori-microservice/services/user-service/pb"
	"github.com/arvians-id/go-apriori-microservice/services/user-service/repository"
	"github.com/arvians-id/go-apriori-microservice/services/user-service/usecase"
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

func NewInitializedServices(configuration *config.Config) (pb.UserServiceServer, error) {
	db, err := NewInitializedDatabase(configuration)
	if err != nil {
		return nil, err
	}

	userRepository := repository.NewUserRepository()
	userService := usecase.NewUserService(userRepository, db)

	return userService, nil
}

func main() {
	configuration, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	connection, err := net.Listen("tcp", configuration.UserSvcUrl)
	if err != nil {
		log.Fatalln("Failed at listening", err)
	}

	fmt.Println("User service is running on port", configuration.UserSvcUrl)

	services, err := NewInitializedServices(configuration)
	if err != nil {
		log.Fatalln("Failed at initializing services", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, services)

	if err := grpcServer.Serve(connection); err != nil {
		log.Fatalln("Failed at serving", err)
	}
}
