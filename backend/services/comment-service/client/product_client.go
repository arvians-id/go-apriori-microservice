package client

import (
	"context"
	"github.com/arvians-id/go-apriori-microservice/services/comment-service/config"
	"github.com/arvians-id/go-apriori-microservice/services/comment-service/pb"
	"google.golang.org/grpc"
	"log"
)

type ProductServiceClient struct {
	Client pb.ProductServiceClient
}

func NewProductServiceClient(configuration *config.Config) ProductServiceClient {
	connection, err := grpc.Dial(configuration.ProductSvcUrl, grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Could not connect:", err)
	}

	return ProductServiceClient{
		Client: pb.NewProductServiceClient(connection),
	}
}

func (client *ProductServiceClient) FindByCode(ctx context.Context, code string) (*pb.GetProductResponse, error) {
	request := &pb.GetProductByProductCodeRequest{
		Code: code,
	}

	response, err := client.Client.FindByCode(ctx, request)
	if err != nil {
		log.Println("[ProductServiceClient][FindAllItemSet] problem calling find by code on product service, err: ", err.Error())
		return nil, err
	}

	return response, nil
}
