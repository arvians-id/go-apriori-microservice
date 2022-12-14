package client

import (
	"context"
	"github.com/arvians-id/go-apriori-microservice/services/payment-service/config"
	"github.com/arvians-id/go-apriori-microservice/services/payment-service/pb"
	"google.golang.org/grpc"
	"log"
)

type TransactionServiceClient struct {
	Client pb.TransactionServiceClient
}

func NewTransactionServiceClient(configuration *config.Config) TransactionServiceClient {
	connection, err := grpc.Dial(configuration.TransactionSvcUrl, grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Could not connect:", err)
	}

	return TransactionServiceClient{
		Client: pb.NewTransactionServiceClient(connection),
	}
}

func (client *TransactionServiceClient) Create(ctx context.Context, req *pb.CreateTransactionRequest) (*pb.GetTransactionResponse, error) {
	response, err := client.Client.Create(ctx, req)
	if err != nil {
		log.Println("[TransactionServiceClient][FindAllItemSet] problem calling create on transaction service, err: ", err.Error())
		return nil, err
	}

	return response, nil
}
