package client

import (
	"context"
	"github.com/arvians-id/go-apriori-microservice/services/payment-service/config"
	"github.com/arvians-id/go-apriori-microservice/services/payment-service/pb"
	"google.golang.org/grpc"
	"log"
)

type UserOrderServiceClient struct {
	Client pb.UserOrderServiceClient
}

func NewUserOrderServiceClient(configuration *config.Config) UserOrderServiceClient {
	connection, err := grpc.Dial(configuration.UserOrderSvcUrl, grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Could not connect:", err)
	}

	return UserOrderServiceClient{
		Client: pb.NewUserOrderServiceClient(connection),
	}
}

func (client *UserOrderServiceClient) FindAllByPayloadId(ctx context.Context, payloadId int64) (*pb.ListUserOrderResponse, error) {
	request := &pb.GetUserOrderByPayloadIdRequest{
		PayloadId: payloadId,
	}

	response, err := client.Client.FindAllByPayloadId(ctx, request)
	if err != nil {
		log.Println("[UserOrderServiceClient][FindAllByPayloadId] problem calling find all by payload id on user order service, err: ", err.Error())
		return nil, err
	}

	return response, nil
}

func (client *UserOrderServiceClient) Create(ctx context.Context, req *pb.CreateUserOrderRequest) (*pb.GetUserOrderResponse, error) {
	response, err := client.Client.Create(ctx, req)
	if err != nil {
		log.Println("[UserOrderServiceClient][Create] problem calling create on user order service, err: ", err.Error())
		return nil, err
	}

	return response, nil
}
