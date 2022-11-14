package usecase

import (
	"context"
	"database/sql"
	"github.com/arvians-id/go-apriori-microservice/adapter/pb"
	"github.com/arvians-id/go-apriori-microservice/model"
	"github.com/arvians-id/go-apriori-microservice/services/notification-service/repository"
	"github.com/arvians-id/go-apriori-microservice/util"
	"github.com/golang/protobuf/ptypes/empty"
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

func (service *NotificationService) FindAll(ctx context.Context, empty *empty.Empty) (*pb.ListNotificationResponse, error) {
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

func (service *NotificationService) MarkAll(ctx context.Context, req *pb.GetNotificationByUserIdRequest) (*empty.Empty, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[NotificationService][MarkAll] problem in db transaction, err: ", err.Error())
		return new(empty.Empty), err
	}
	defer util.CommitOrRollback(tx)

	err = service.NotificationRepository.MarkAll(ctx, tx, req.UserId)
	if err != nil {
		log.Println("[NotificationService][Create][MarkAll] problem in getting from repository, err: ", err.Error())
		return new(empty.Empty), err
	}

	return new(empty.Empty), nil
}

func (service *NotificationService) Mark(ctx context.Context, req *pb.GetNotificationByIdRequest) (*empty.Empty, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[NotificationService][Mark] problem in db transaction, err: ", err.Error())
		return new(empty.Empty), err
	}
	defer util.CommitOrRollback(tx)

	err = service.NotificationRepository.Mark(ctx, tx, req.Id)
	if err != nil {
		log.Println("[NotificationService][Create][Mark] problem in getting from repository, err: ", err.Error())
		return new(empty.Empty), err
	}

	return new(empty.Empty), nil
}
