package usecase

import (
	"context"
	"database/sql"
	"github.com/arvians-id/go-apriori-microservice/services/user-order-service/client"
	"github.com/arvians-id/go-apriori-microservice/services/user-order-service/model"
	"github.com/arvians-id/go-apriori-microservice/services/user-order-service/pb"
	"github.com/arvians-id/go-apriori-microservice/services/user-order-service/repository"
	"github.com/arvians-id/go-apriori-microservice/services/user-order-service/util"
	"log"
)

type UserOrderService struct {
	UserOrderRepository repository.UserOrderRepository
	UserService         client.UserServiceClient
	DB                  *sql.DB
}

func NewUserOrderService(
	userOrderRepository repository.UserOrderRepository,
	userService client.UserServiceClient,
	db *sql.DB,
) pb.UserOrderServiceServer {
	return &UserOrderService{
		UserOrderRepository: userOrderRepository,
		UserService:         userService,
		DB:                  db,
	}
}

func (service *UserOrderService) FindAllByPayloadId(ctx context.Context, req *pb.GetUserOrderByPayloadIdRequest) (*pb.ListUserOrderResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[UserOrderService][FindAllByPayloadId] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	userOrders, err := service.UserOrderRepository.FindAllByPayloadId(ctx, tx, req.PayloadId)
	if err != nil {
		log.Println("[UserOrderService][FindAllByPayloadId][FindAllByPayloadId] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	var userOrderListResponse []*pb.UserOrder
	for _, userOrder := range userOrders {
		userOrderListResponse = append(userOrderListResponse, userOrder.ToProtoBuff())
	}

	return &pb.ListUserOrderResponse{
		UserOrder: userOrderListResponse,
	}, nil
}

func (service *UserOrderService) FindAllByUserId(ctx context.Context, req *pb.GetUserOrderByUserIdRequest) (*pb.ListUserOrderResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[UserOrderService][FindAllByUserId] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	user, err := service.UserService.FindById(ctx, req.UserId)
	if err != nil {
		log.Println("[UserOrderService][FindAllByUserId][FindById] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	userOrders, err := service.UserOrderRepository.FindAllByUserId(ctx, tx, user.User.IdUser)
	if err != nil {
		log.Println("[UserOrderService][FindAllByUserId][FindAllByUserId] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	var userOrderListResponse []*pb.UserOrder
	for _, userOrder := range userOrders {
		userOrderListResponse = append(userOrderListResponse, userOrder.ToProtoBuff())
	}

	return &pb.ListUserOrderResponse{
		UserOrder: userOrderListResponse,
	}, nil
}

func (service *UserOrderService) FindById(ctx context.Context, req *pb.GetUserOrderByIdRequest) (*pb.GetUserOrderResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[UserOrderService][FindById] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	userOrder, err := service.UserOrderRepository.FindById(ctx, tx, req.Id)
	if err != nil {
		log.Println("[UserOrderService][FindById][FindById] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return &pb.GetUserOrderResponse{
		UserOrder: userOrder.ToProtoBuff(),
	}, nil
}

func (service *UserOrderService) Create(ctx context.Context, req *pb.CreateUserOrderRequest) (*pb.GetUserOrderResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[UserOrderService][Create] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	userOrder, err := service.UserOrderRepository.Create(ctx, tx, &model.UserOrder{
		PayloadId:      req.PayloadId,
		Code:           req.Code,
		Name:           req.Name,
		Price:          req.Price,
		Image:          req.Image,
		Quantity:       req.Quantity,
		TotalPriceItem: req.TotalPriceItem,
	})
	if err != nil {
		log.Println("[UserOrderService][Create][Create] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return &pb.GetUserOrderResponse{
		UserOrder: userOrder.ToProtoBuff(),
	}, nil
}
