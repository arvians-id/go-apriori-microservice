package notification

import (
	"errors"
	"github.com/arvians-id/go-apriori-microservice/adapter/middleware"
	"github.com/arvians-id/go-apriori-microservice/adapter/pb"
	"github.com/arvians-id/go-apriori-microservice/adapter/response"
	"github.com/arvians-id/go-apriori-microservice/config"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"log"
	"strconv"
)

type ServiceClient struct {
	NotificationService pb.NotificationServiceClient
}

func NewNotificationServiceClient(configuration *config.Config) pb.NotificationServiceClient {
	connection, err := grpc.Dial(configuration.NotificationSvcUrl, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}

	return pb.NewNotificationServiceClient(connection)
}

func RegisterRoutes(router *gin.Engine, configuration *config.Config) *ServiceClient {
	serviceClient := &ServiceClient{
		NotificationService: NewNotificationServiceClient(configuration),
	}

	authorized := router.Group("/api", middleware.AuthJwtMiddleware())
	{
		authorized.GET("/notifications", serviceClient.FindAll)
		authorized.GET("/notifications/user", serviceClient.FindAllByUserId)
		authorized.PATCH("/notifications/mark", serviceClient.MarkAll)
		authorized.PATCH("/notifications/mark/:id", serviceClient.Mark)
	}

	return serviceClient
}

func (client *ServiceClient) FindAll(c *gin.Context) {
	notifications, err := client.NotificationService.FindAll(c.Request.Context(), new(empty.Empty))
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", notifications)
}

func (client *ServiceClient) FindAllByUserId(c *gin.Context) {
	id, isExist := c.Get("id_user")
	if !isExist {
		response.ReturnErrorUnauthorized(c, errors.New("unauthorized"), nil)
		return
	}

	notifications, err := client.NotificationService.FindAllByUserId(c.Request.Context(), &pb.GetNotificationByUserIdRequest{
		UserId: int64(id.(float64)),
	})
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", notifications)
}

func (client *ServiceClient) MarkAll(c *gin.Context) {
	id, isExist := c.Get("id_user")
	if !isExist {
		response.ReturnErrorUnauthorized(c, errors.New("unauthorized"), nil)
		return
	}

	_, err := client.NotificationService.MarkAll(c.Request.Context(), &pb.GetNotificationByUserIdRequest{
		UserId: int64(id.(float64)),
	})
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", nil)
}

func (client *ServiceClient) Mark(c *gin.Context) {
	idParam, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	_, err = client.NotificationService.Mark(c.Request.Context(), &pb.GetNotificationByIdRequest{
		Id: idParam,
	})
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", nil)
}
