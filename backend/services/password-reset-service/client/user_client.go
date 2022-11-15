package client

import (
	"context"
	"github.com/arvians-id/go-apriori-microservice/services/password-reset-service/config"
	"github.com/arvians-id/go-apriori-microservice/services/password-reset-service/pb"
	"github.com/golang/protobuf/ptypes/empty"
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

func (client *UserServiceClient) FindByEmail(ctx context.Context, email string) (*pb.GetUserResponse, error) {
	request := &pb.GetUserByEmailRequest{
		Email: email,
	}

	response, err := client.Client.FindByEmail(ctx, request)
	if err != nil {
		log.Println("[TransactionServiceClient][FindAllItemSet] problem calling transaction service, err: ", err.Error())
		return nil, err
	}

	return response, nil
}

func (client *UserServiceClient) UpdatePassword(ctx context.Context, req *pb.UpdateUserPasswordRequest) (*empty.Empty, error) {
	response, err := client.Client.UpdatePassword(ctx, req)
	if err != nil {
		log.Println("[TransactionServiceClient][FindAllItemSet] problem calling transaction service, err: ", err.Error())
		return nil, err
	}

	return response, nil
}
