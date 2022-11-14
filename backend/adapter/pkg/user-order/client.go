package user_order

import (
	"errors"
	"github.com/arvians-id/go-apriori-microservice/adapter/middleware"
	"github.com/arvians-id/go-apriori-microservice/adapter/pb"
	"github.com/arvians-id/go-apriori-microservice/adapter/pkg/payment"
	"github.com/arvians-id/go-apriori-microservice/adapter/response"
	"github.com/arvians-id/go-apriori-microservice/config"
	"github.com/arvians-id/go-apriori-microservice/third-party/aws"
	"github.com/arvians-id/go-apriori-microservice/util"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"strconv"
)

type ServiceClient struct {
	UserOrderService pb.UserOrderServiceClient
	PaymentService   pb.PaymentServiceClient
	StorageS3        aws.StorageS3
}

func NewUserOrderServiceClient(configuration *config.Config) pb.UserOrderServiceClient {
	connection, err := grpc.Dial(configuration.UserOrderSvcUrl, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}

	return pb.NewUserOrderServiceClient(connection)
}

func RegisterRoutes(router *gin.Engine, configuration *config.Config) *ServiceClient {
	serviceClient := &ServiceClient{
		UserOrderService: NewUserOrderServiceClient(configuration),
		PaymentService:   payment.NewCommentServiceClient(configuration),
	}

	authorized := router.Group("/api", middleware.AuthJwtMiddleware(configuration))
	{
		authorized.GET("/user-order", serviceClient.FindAll)
		authorized.GET("/user-order/user", serviceClient.FindAllByUserId)
		authorized.GET("/user-order/:order_id", serviceClient.FindAllById)
		authorized.GET("/user-order/single/:id", serviceClient.FindById)
	}

	return serviceClient
}

func (client *ServiceClient) FindAll(c *gin.Context) {
	id, isExist := c.Get("id_user")
	if !isExist {
		response.ReturnErrorUnauthorized(c, errors.New("unauthorized"), nil)
		return
	}

	payments, err := client.PaymentService.FindAllByUserId(c.Request.Context(), &pb.GetPaymentByUserIdRequest{
		UserId: int64(id.(float64)),
	})
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", payments.GetPayment())
}

func (client *ServiceClient) FindAllByUserId(c *gin.Context) {
	id, isExist := c.Get("id_user")
	if !isExist {
		response.ReturnErrorUnauthorized(c, errors.New("unauthorized"), nil)
		return
	}

	userOrders, err := client.UserOrderService.FindAllByUserId(c.Request.Context(), &pb.GetUserOrderByUserIdRequest{
		UserId: int64(id.(float64)),
	})
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", userOrders.GetUserOrder())
}

func (client *ServiceClient) FindAllById(c *gin.Context) {
	orderIdParam := c.Param("order_id")
	paymentResponse, err := client.PaymentService.FindByOrderId(c.Request.Context(), &pb.GetPaymentByOrderIdRequest{
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

	userOrder, err := client.UserOrderService.FindAllByPayloadId(c.Request.Context(), &pb.GetUserOrderByPayloadIdRequest{
		PayloadId: paymentResponse.Payment.IdPayload,
	})
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", userOrder.GetUserOrder())
}

func (client *ServiceClient) FindById(c *gin.Context) {
	idParam, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	userOrder, err := client.UserOrderService.FindById(c.Request.Context(), &pb.GetUserOrderByIdRequest{
		Id: idParam,
	})
	if err != nil {
		if err.Error() == util.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", userOrder.GetUserOrder())
}
