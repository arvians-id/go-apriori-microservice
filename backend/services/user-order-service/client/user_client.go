package client

import (
	"context"
	"github.com/arvians-id/go-apriori-microservice/adapter/pb"
	"github.com/arvians-id/go-apriori-microservice/config"
	"google.golang.org/grpc"
	"log"
)

type UserServiceClient struct {
	Client pb.UserServiceClient
}

func NewUserServiceClient(configuration *config.Config) UserServiceClient {
	connection, err := grpc.Dial(configuration.UserSvcUrl, grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Could not connect:", err)
	}

	return UserServiceClient{
		Client: pb.NewUserServiceClient(connection),
	}
}

func (client *UserServiceClient) FindById(ctx context.Context, id int64) (*pb.GetUserResponse, error) {
	request := &pb.GetUserByIdRequest{
		Id: id,
	}

	response, err := client.Client.FindById(ctx, request)
	if err != nil {
		log.Println("[TransactionServiceClient][FindAllItemSet] problem calling transaction service, err: ", err.Error())
		return nil, err
	}

	return response, nil
}
