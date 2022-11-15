package model

import (
	"github.com/arvians-id/go-apriori-microservice/services/category-service/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Notification struct {
	IdNotification int64     `json:"id_notification"`
	UserId         int64     `json:"user_id"`
	Title          string    `json:"title"`
	Description    *string   `json:"description"`
	URL            *string   `json:"url"`
	IsRead         bool      `json:"is_read"`
	CreatedAt      time.Time `json:"created_at"`
	User           *User     `json:"user"`
}

func (notification *Notification) ToProtoBuff() *pb.Notification {
	return &pb.Notification{
		IdNotification: notification.IdNotification,
		UserId:         notification.UserId,
		Title:          notification.Title,
		Description:    notification.Description,
		URL:            notification.URL,
		IsRead:         notification.IsRead,
		CreatedAt:      timestamppb.New(notification.CreatedAt),
		//User:           notification.User.ToProtoBuff(),
	}
}
