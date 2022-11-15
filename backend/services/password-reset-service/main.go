package main

import (
	"database/sql"
	"fmt"
	"github.com/arvians-id/go-apriori-microservice/services/password-reset-service/client"
	"github.com/arvians-id/go-apriori-microservice/services/password-reset-service/config"
	"github.com/arvians-id/go-apriori-microservice/services/password-reset-service/pb"
	"github.com/arvians-id/go-apriori-microservice/services/password-reset-service/repository"
	"github.com/arvians-id/go-apriori-microservice/services/password-reset-service/usecase"
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

func NewInitializedServices(configuration *config.Config) (pb.PasswordResetServiceServer, error) {
	db, err := NewInitializedDatabase(configuration)
	if err != nil {
		return nil, err
	}

	userClient := client.NewUserServiceClient(configuration)

	passwordResetRepository := repository.NewPasswordResetRepository()
	passwordResetService := usecase.NewPasswordResetService(passwordResetRepository, userClient, db)

	return passwordResetService, nil
}

func main() {
	configuration, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	connection, err := net.Listen("tcp", configuration.PasswordResetSvcUrl)
	if err != nil {
		log.Fatalln("Failed at listening", err)
	}

	fmt.Println("Password Reset service is running on port", configuration.PasswordResetSvcUrl)

	services, err := NewInitializedServices(configuration)
	if err != nil {
		log.Fatalln("Failed at initializing services", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPasswordResetServiceServer(grpcServer, services)

	if err := grpcServer.Serve(connection); err != nil {
		log.Fatalln("Failed at serving", err)
	}
}
