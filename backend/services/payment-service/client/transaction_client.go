package client

import (
	"context"
	"github.com/arvians-id/go-apriori-microservice/adapter/pkg/transaction/pb"
	"github.com/arvians-id/go-apriori-microservice/config"
	"google.golang.org/grpc"
	"log"
)

type TransactionServiceClient struct {
	Client pb.TransactionServiceClient
}

func NewTransactionServiceClient(configuration *config.Config) pb.TransactionServiceClient {
	connection, err := grpc.Dial(configuration.TransactionSvcUrl, grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Could not connect:", err)
	}

	return pb.NewTransactionServiceClient(connection)
}

func (client *TransactionServiceClient) Create(ctx context.Context, productName string, customerName string, noTransaction *string) (*pb.GetTransactionResponse, error) {
	request := &pb.CreateTransactionRequest{
		ProductName:   productName,
		CustomerName:  customerName,
		NoTransaction: noTransaction,
	}

	response, err := client.Client.Create(ctx, request)
	if err != nil {
		log.Println("[TransactionServiceClient][FindAllItemSet] problem calling transaction service, err: ", err.Error())
		return nil, err
	}

	return response, nil
}
