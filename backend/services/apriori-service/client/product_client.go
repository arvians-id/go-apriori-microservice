package client

import (
	"context"
	"github.com/arvians-id/go-apriori-microservice/services/apriori-service/config"
	"github.com/arvians-id/go-apriori-microservice/services/apriori-service/pb"
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

func (client *ProductServiceClient) FindByName(ctx context.Context, name string) (*pb.GetProductResponse, error) {
	request := &pb.GetProductByProductNameRequest{
		Name: name,
	}

	response, err := client.Client.FindByName(ctx, request)
	if err != nil {
		log.Println("[ProductServiceClient][FindByName] problem calling find by name on product service, err: ", err.Error())
		return nil, err
	}

	return response, nil
}
