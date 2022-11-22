package client

import (
	"context"
	"github.com/arvians-id/go-apriori-microservice/services/product-service/config"
	"github.com/arvians-id/go-apriori-microservice/services/product-service/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

type AprioriServiceClient struct {
	Client pb.AprioriServiceClient
}

func NewAprioriServiceClient(configuration *config.Config) AprioriServiceClient {
	connection, err := grpc.Dial(configuration.AprioriSvcUrl, grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Could not connect:", err)
	}

	return AprioriServiceClient{
		Client: pb.NewAprioriServiceClient(connection),
	}
}

func (client *AprioriServiceClient) FindAllByActive(ctx context.Context) (*pb.ListAprioriResponse, error) {
	response, err := client.Client.FindAllByActive(ctx, new(emptypb.Empty))
	if err != nil {
		log.Println("[AprioriServiceClient][FindAllByActive] problem calling find all by active on apriori service, err: ", err.Error())
		return nil, err
	}

	return response, nil
}
