package usecase

import (
	"context"
	"database/sql"
	"github.com/arvians-id/go-apriori-microservice/services/notification-service/model"
	"github.com/arvians-id/go-apriori-microservice/services/notification-service/pb"
	"github.com/arvians-id/go-apriori-microservice/services/notification-service/repository"
	"github.com/arvians-id/go-apriori-microservice/services/notification-service/util"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"strings"
	"time"
)

type NotificationService struct {
	NotificationRepository repository.NotificationRepository
	DB                     *sql.DB
}

func NewNotificationService(
	notificationRepository repository.NotificationRepository,
	db *sql.DB,
) pb.NotificationServiceServer {
	return &NotificationService{
		NotificationRepository: notificationRepository,
		DB:                     db,
	}
}

func (service *NotificationService) FindAll(ctx context.Context, empty *emptypb.Empty) (*pb.ListNotificationResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[NotificationService][FindAll] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	notifications, err := service.NotificationRepository.FindAll(ctx, tx)
	if err != nil {
		log.Println("[NotificationService][FindAll][FindAll] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	var notificationListResponse []*pb.Notification
	for _, notification := range notifications {
		notificationListResponse = append(notificationListResponse, notification.ToProtoBuff())
	}

	return &pb.ListNotificationResponse{
		Notification: notificationListResponse,
	}, nil
}

func (service *NotificationService) FindAllByUserId(ctx context.Context, req *pb.GetNotificationByUserIdRequest) (*pb.ListNotificationResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[NotificationService][FindAllByUserId] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	notifications, err := service.NotificationRepository.FindAllByUserId(ctx, tx, req.UserId)
	if err != nil {
		log.Println("[NotificationService][FindAllByUserId][FindAllByUserId] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	var notificationListResponse []*pb.Notification
	for _, notification := range notifications {
		notificationListResponse = append(notificationListResponse, notification.ToProtoBuff())
	}

	return &pb.ListNotificationResponse{
		Notification: notificationListResponse,
	}, nil
}

func (service *NotificationService) Create(ctx context.Context, req *pb.CreateNotificationRequest) (*pb.GetNotificationResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[NotificationService][Create] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	timeNow, err := time.Parse(util.TimeFormat, time.Now().Format(util.TimeFormat))
	if err != nil {
		log.Println("[NotificationService][Create] problem in parsing to time, err: ", err.Error())
		return nil, err
	}

	notificationResponse, err := service.NotificationRepository.Create(ctx, tx, &model.Notification{
		UserId:      req.UserId,
		Title:       strings.Title(req.Title),
		Description: &req.Description,
		URL:         &req.URL,
		IsRead:      false,
		CreatedAt:   timeNow,
	})
	if err != nil {
		log.Println("[NotificationService][Create][Create] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return &pb.GetNotificationResponse{
		Notification: notificationResponse.ToProtoBuff(),
	}, nil
}

func (service *NotificationService) MarkAll(ctx context.Context, req *pb.GetNotificationByUserIdRequest) (*emptypb.Empty, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[NotificationService][MarkAll] problem in db transaction, err: ", err.Error())
		return new(emptypb.Empty), err
	}
	defer util.CommitOrRollback(tx)

	err = service.NotificationRepository.MarkAll(ctx, tx, req.UserId)
	if err != nil {
		log.Println("[NotificationService][MarkAll][MarkAll] problem in getting from repository, err: ", err.Error())
		return new(emptypb.Empty), err
	}

	return new(emptypb.Empty), nil
}

func (service *NotificationService) Mark(ctx context.Context, req *pb.GetNotificationByIdRequest) (*emptypb.Empty, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[NotificationService][Mark] problem in db transaction, err: ", err.Error())
		return new(emptypb.Empty), err
	}
	defer util.CommitOrRollback(tx)

	err = service.NotificationRepository.Mark(ctx, tx, req.Id)
	if err != nil {
		log.Println("[NotificationService][Create][Mark] problem in getting from repository, err: ", err.Error())
		return new(emptypb.Empty), err
	}

	return new(emptypb.Empty), nil
}
