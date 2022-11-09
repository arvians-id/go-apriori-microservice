package main

import (
	"database/sql"
	"fmt"
	"github.com/arvians-id/go-apriori-microservice/adapter/pkg/comment/pb"
	"github.com/arvians-id/go-apriori-microservice/config"
	"github.com/arvians-id/go-apriori-microservice/services/comment-service/client"
	"github.com/arvians-id/go-apriori-microservice/services/comment-service/repository"
	"github.com/arvians-id/go-apriori-microservice/services/comment-service/usecase"
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

func NewInitializedServices(configuration *config.Config) (pb.CommentServiceServer, error) {
	db, err := NewInitializedDatabase(configuration)
	if err != nil {
		return nil, err
	}

	productClient := client.NewProductServiceClient(configuration)

	commentRepository := repository.NewCommentRepository()
	commentService := usecase.NewCommentService(commentRepository, productClient, db)

	return commentService, nil
}

func main() {
	configuration, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	connection, err := net.Listen("tcp", configuration.CommentSvcUrl)
	if err != nil {
		log.Fatalln("Failed at listening", err)
	}

	fmt.Println("Comment service is running on port", configuration.CommentSvcUrl)

	services, err := NewInitializedServices(configuration)
	if err != nil {
		log.Fatalln("Failed at initializing services", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCommentServiceServer(grpcServer, services)

	if err := grpcServer.Serve(connection); err != nil {
		log.Fatalln("Failed at serving", err)
	}
}
