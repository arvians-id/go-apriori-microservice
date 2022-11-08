package client

import (
	"context"
	"github.com/arvians-id/go-apriori-microservice/adapter/pkg/product/pb"
	"github.com/arvians-id/go-apriori-microservice/config"
	"google.golang.org/grpc"
	"log"
)

type ProductServiceClient struct {
	Client pb.ProductServiceClient
}

func NewProductServiceClient(configuration *config.Config) pb.ProductServiceClient {
	connection, err := grpc.Dial(configuration.ProductSvcUrl, grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Could not connect:", err)
	}

	return pb.NewProductServiceClient(connection)
}

func (client *ProductServiceClient) FindByCode(ctx context.Context, code string) (*pb.GetProductResponse, error) {
	request := &pb.GetProductByProductCodeRequest{
		Code: code,
	}

	response, err := client.Client.FindByCode(ctx, request)
	if err != nil {
		log.Println("[TransactionServiceClient][FindAllItemSet] problem calling transaction service, err: ", err.Error())
		return nil, err
	}

	return response, nil
}
