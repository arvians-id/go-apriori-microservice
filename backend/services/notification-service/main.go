package main

import (
	"database/sql"
	"fmt"
	"github.com/arvians-id/go-apriori-microservice/services/notification-service/config"
	"github.com/arvians-id/go-apriori-microservice/services/notification-service/pb"
	"github.com/arvians-id/go-apriori-microservice/services/notification-service/repository"
	"github.com/arvians-id/go-apriori-microservice/services/notification-service/usecase"
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

func NewInitializedServices(configuration *config.Config) (pb.NotificationServiceServer, error) {
	db, err := NewInitializedDatabase(configuration)
	if err != nil {
		return nil, err
	}

	notificationRepository := repository.NewNotificationRepository()
	notificationService := usecase.NewNotificationService(notificationRepository, db)

	return notificationService, nil
}

func main() {
	configuration, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	port := ":" + strings.Split(configuration.NotificationSvcUrl, ":")[1]
	connection, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln("Failed at listening", err)
	}

	fmt.Println("Notification service is running on port", port)

	services, err := NewInitializedServices(configuration)
	if err != nil {
		log.Fatalln("Failed at initializing services", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterNotificationServiceServer(grpcServer, services)

	if err := grpcServer.Serve(connection); err != nil {
		log.Fatalln("Failed at serving", err)
	}
}
