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

func (client *TransactionServiceClient) FindAllItemSet(ctx context.Context, startDate string, endDate string) (*pb.ListTransactionsResponse, error) {
	request := &pb.GetAllItemSetTransactionRequest{
		StartDate: startDate,
		EndDate:   endDate,
	}

	response, err := client.Client.FindAllItemSet(ctx, request)
	if err != nil {
		log.Println("[TransactionServiceClient][FindAllItemSet] problem calling transaction service, err: ", err.Error())
		return nil, err
	}

	return response, nil
}