package payment

import (
	"encoding/json"
	"github.com/arvians-id/go-apriori-microservice/adapter/middleware"
	"github.com/arvians-id/go-apriori-microservice/adapter/pb"
	"github.com/arvians-id/go-apriori-microservice/adapter/pkg/notification"
	"github.com/arvians-id/go-apriori-microservice/adapter/pkg/user"
	"github.com/arvians-id/go-apriori-microservice/adapter/response"
	"github.com/arvians-id/go-apriori-microservice/config"
	"github.com/arvians-id/go-apriori-microservice/model"
	messaging "github.com/arvians-id/go-apriori-microservice/third-party/message-queue"
	"github.com/arvians-id/go-apriori-microservice/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/veritrans/go-midtrans"
	"google.golang.org/grpc"
	"log"
)

type ServiceClient struct {
	PaymentService      pb.PaymentServiceClient
	UserService         pb.UserServiceClient
	NotificationService pb.NotificationServiceClient
	Producer            *messaging.Producer
}

func NewCommentServiceClient(configuration *config.Config) pb.PaymentServiceClient {
	connection, err := grpc.Dial(configuration.PaymentSvcUrl, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}

	return pb.NewPaymentServiceClient(connection)
}

func RegisterRoutes(router *gin.Engine, configuration *config.Config, producer *messaging.Producer) *ServiceClient {
	serviceClient := &ServiceClient{
		PaymentService:      NewCommentServiceClient(configuration),
		UserService:         user.NewUserServiceClient(configuration),
		NotificationService: notification.NewNotificationServiceClient(configuration),
		Producer:            producer,
	}

	authorized := router.Group("/api", middleware.AuthJwtMiddleware())
	{
		authorized.GET("/payments", middleware.SetupXApiKeyMiddleware(), serviceClient.FindAll)
		authorized.GET("/payments/:order_id", middleware.SetupXApiKeyMiddleware(), serviceClient.FindByOrderId)
		authorized.PATCH("/payments/:order_id", middleware.SetupXApiKeyMiddleware(), serviceClient.UpdateReceiptNumber)
	}

	unauthorized := router.Group("/api")
	{
		unauthorized.POST("/payments/pay", middleware.SetupXApiKeyMiddleware(), serviceClient.Pay)
		unauthorized.POST("/payments/notification", serviceClient.Notification)
		unauthorized.DELETE("/payments/:order_id", middleware.SetupXApiKeyMiddleware(), serviceClient.Delete)
	}

	return serviceClient
}

func (client *ServiceClient) FindAll(c *gin.Context) {
	payments, err := client.PaymentService.FindAll(c.Request.Context(), new(empty.Empty))
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", payments)
}

func (client *ServiceClient) FindByOrderId(c *gin.Context) {
	orderIdParam := c.Param("order_id")
	payment, err := client.PaymentService.FindByOrderId(c.Request.Context(), &pb.GetPaymentByOrderIdRequest{
		OrderId: orderIdParam,
	})
	if err != nil {
		if err.Error() == util.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", payment)
}

func (client *ServiceClient) UpdateReceiptNumber(c *gin.Context) {
	var requestPayment AddReceiptNumberRequest
	err := c.ShouldBindJSON(&requestPayment)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	requestPayment.OrderId = c.Param("order_id")
	payment, err := client.PaymentService.UpdateReceiptNumber(c.Request.Context(), &pb.UpdateReceiptNumberRequest{
		OrderId:       requestPayment.OrderId,
		ReceiptNumber: requestPayment.ReceiptNumber,
	})
	if err != nil {
		if err.Error() == util.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	// Send Notification
	notificationResponse, err := client.NotificationService.Create(c.Request.Context(), &pb.CreateNotificationRequest{
		UserId:      payment.Payment.UserId,
		Title:       "Receipt number arrived",
		Description: "Your receipt number had been entered by admin",
		URL:         "product",
	})
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	userResponse, err := client.UserService.FindById(c.Request.Context(), &pb.GetUserByIdRequest{
		Id: payment.Payment.UserId,
	})
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	// Send Email
	emailService := model.EmailService{
		ToEmail: userResponse.User.Email,
		Subject: notificationResponse.Notification.Title,
		Message: notificationResponse.Notification.Description,
	}
	err = client.Producer.Publish("mail_topic", emailService)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", nil)
}

func (client *ServiceClient) Pay(c *gin.Context) {
	var requestToken GetPaymentTokenRequest
	err := c.ShouldBind(&requestToken)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	data, err := client.PaymentService.GetToken(c.Request.Context(), &pb.GetPaymentTokenRequest{
		GrossAmount:    requestToken.GrossAmount,
		Items:          requestToken.Items,
		UserId:         requestToken.UserId,
		CustomerName:   requestToken.CustomerName,
		Address:        requestToken.Address,
		Courier:        requestToken.Courier,
		CourierService: requestToken.CourierService,
		ShippingCost:   requestToken.ShippingCost,
	})
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", data)
}

func (client *ServiceClient) Notification(c *gin.Context) {
	var payload midtrans.ChargeReqWithMap
	err := c.BindJSON(&payload)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	encode, err := json.Marshal(payload)
	if err != nil {
		log.Println("[PaymentController][Notification] unable to marshal json, err: ", err.Error())
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	resArray := make(map[string]interface{})
	err = json.Unmarshal(encode, &resArray)
	if err != nil {
		log.Println("[PaymentController][Notification] unable to unmarshal json, err: ", err.Error())
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	isSettlement, err := client.PaymentService.CreateOrUpdate(c.Request.Context(), &pb.CreatePaymentRequest{
		Payment: encode,
	})
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	if isSettlement.IsSuccess {
		idUser := util.StrToInt(resArray["custom_field1"].(string))
		// Send Notification
		notificationResponse, err := client.NotificationService.Create(c.Request.Context(), &pb.CreateNotificationRequest{
			UserId:      int64(idUser),
			Title:       "Transaction Successfully",
			Description: "You have successfully made a payment. Thank you for shopping at Ryzy Shop",
			URL:         "product",
		})
		if err != nil {
			response.ReturnErrorInternalServerError(c, err, nil)
			return
		}

		userResponse, err := client.UserService.FindById(c.Request.Context(), &pb.GetUserByIdRequest{
			Id: int64(idUser),
		})
		if err != nil {
			response.ReturnErrorInternalServerError(c, err, nil)
			return
		}

		// Send Email
		emailService := model.EmailService{
			ToEmail: userResponse.User.Email,
			Subject: notificationResponse.Notification.Title,
			Message: notificationResponse.Notification.Description,
		}
		err = client.Producer.Publish("mail_topic", emailService)
		if err != nil {
			response.ReturnErrorInternalServerError(c, err, nil)
			return
		}
	}

	response.ReturnSuccessOK(c, "OK", nil)
}

func (client *ServiceClient) Delete(c *gin.Context) {
	orderIdParam := c.Param("order_id")
	_, err := client.PaymentService.Delete(c.Request.Context(), &pb.GetPaymentByOrderIdRequest{
		OrderId: orderIdParam,
	})
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", nil)
}
