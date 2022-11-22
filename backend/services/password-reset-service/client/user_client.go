package client

import (
	"context"
	"github.com/arvians-id/go-apriori-microservice/services/password-reset-service/config"
	"github.com/arvians-id/go-apriori-microservice/services/password-reset-service/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
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
		log.Println("[UserServiceClient][FindByEmail] problem calling find by email on product service, err: ", err.Error())
		return nil, err
	}

	return response, nil
}

func (client *UserServiceClient) UpdatePassword(ctx context.Context, req *pb.UpdateUserPasswordRequest) (*emptypb.Empty, error) {
	response, err := client.Client.UpdatePassword(ctx, req)
	if err != nil {
		log.Println("[UserServiceClient][UpdatePassword] problem calling update password on product service, err: ", err.Error())
		return new(emptypb.Empty), err
	}

	return response, nil
}
