package service

import (
	"context"
	"database/sql"
	pbpayment "github.com/arvians-id/go-apriori-microservice/adapter/pkg/payment/pb"
	"github.com/arvians-id/go-apriori-microservice/adapter/pkg/user-order/pb"
	pbuser "github.com/arvians-id/go-apriori-microservice/adapter/pkg/user/pb"
	"github.com/arvians-id/go-apriori-microservice/services/user-order-service/repository"
	"github.com/arvians-id/go-apriori-microservice/util"
	"log"
)

type UserOrderService struct {
	UserOrderRepository repository.UserOrderRepository
	PaymentService      pbpayment.PaymentServiceClient
	UserService         pbuser.UserServiceClient
	DB                  *sql.DB
}

func NewUserOrderService(
	userOrderRepository *repository.UserOrderRepository,
	PpaymentService pbpayment.PaymentServiceClient,
	userService pbuser.UserServiceClient,
	db *sql.DB,
) pb.UserOrderServiceServer {
	return &UserOrderService{
		UserOrderRepository: *userOrderRepository,
		PaymentService:      PpaymentService,
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

	userOrders, err := service.UserOrderRepository.FindAllByPayloadId(ctx, tx, util.IntToStr(payloadId))
	if err != nil {
		log.Println("[UserOrderService][FindAllByPayloadId][FindAllByPayloadId] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return userOrders, nil
}

func (service *UserOrderService) FindAllByUserId(ctx context.Context, req *pb.GetUserOrderByUserIdRequest) (*pb.ListUserOrderResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[UserOrderService][FindAllByUserId] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	_, err = service.UserRepository.FindById(ctx, tx, userId)
	if err != nil {
		log.Println("[UserOrderService][FindAllByUserId][FindById] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	userOrders, err := service.UserOrderRepository.FindAllByUserId(ctx, tx, userId)
	if err != nil {
		log.Println("[UserOrderService][FindAllByUserId][FindAllByUserId] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return userOrders, nil
}

func (service *UserOrderService) FindById(ctx context.Context, req *pb.GetUserOrderByIdRequest) (*pb.GetUserOrderResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[UserOrderService][FindById] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	userOrder, err := service.UserOrderRepository.FindById(ctx, tx, id)
	if err != nil {
		log.Println("[UserOrderService][FindById][FindById] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return userOrder, nil
}

func (service *UserOrderService) Create(ctx context.Context, req *pb.CreateUserOrderRequest) (*pb.GetUserOrderResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[UserOrderService][Create] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	userOrder, err := service.UserOrderRepository.Create(ctx, tx, req)
	if err != nil {
		log.Println("[UserOrderService][Create][Create] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return userOrder, nil
}
