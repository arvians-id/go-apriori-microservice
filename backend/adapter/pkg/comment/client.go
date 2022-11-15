package comment

import (
	"github.com/arvians-id/go-apriori-microservice/adapter/config"
	"github.com/arvians-id/go-apriori-microservice/adapter/middleware"
	"github.com/arvians-id/go-apriori-microservice/adapter/pb"
	"github.com/arvians-id/go-apriori-microservice/adapter/response"
	"github.com/arvians-id/go-apriori-microservice/adapter/util"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"strconv"
)

type ServiceClient struct {
	CommentService pb.CommentServiceClient
}

func NewCommentServiceClient(configuration *config.Config) pb.CommentServiceClient {
	connection, err := grpc.Dial(configuration.CommentSvcUrl, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}

	return pb.NewCommentServiceClient(connection)
}

func RegisterRoutes(router *gin.Engine, configuration *config.Config) *ServiceClient {
	serviceClient := &ServiceClient{
		CommentService: NewCommentServiceClient(configuration),
	}

	authorized := router.Group("/api", middleware.AuthJwtMiddleware(configuration))
	{
		authorized.POST("/comments", serviceClient.Create)
	}

	unauthorized := router.Group("/api")
	{
		unauthorized.GET("/comments/:id", serviceClient.FindById)
		unauthorized.GET("/comments/rating/:product_code", serviceClient.FindAllRatingByProductCode)
		unauthorized.GET("/comments/product/:product_code", serviceClient.FindAllByProductCode)
		unauthorized.GET("/comments/user-order/:user_order_id", serviceClient.FindByUserOrderId)
	}

	return serviceClient
}

func (client *ServiceClient) FindAllByProductCode(c *gin.Context) {
	productCodeParam := c.Param("product_code")
	tagsQuery := c.Query("tags")
	ratingQuery := c.Query("rating")
	comments, err := client.CommentService.FindAllByProductCode(c.Request.Context(), &pb.GetCommentByFiltersRequest{
		ProductCode: productCodeParam,
		Tag:         tagsQuery,
		Rating:      ratingQuery,
	})
	if err != nil {
		if err.Error() == util.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", comments.GetComment())
}

func (client *ServiceClient) FindById(c *gin.Context) {
	idParam, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	comment, err := client.CommentService.FindById(c.Request.Context(), &pb.GetCommentByIdRequest{
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

	response.ReturnSuccessOK(c, "OK", comment.GetComment())
}

func (client *ServiceClient) FindAllRatingByProductCode(c *gin.Context) {
	productCodeParam := c.Param("product_code")
	comments, err := client.CommentService.FindAllRatingByProductCode(c.Request.Context(), &pb.GetCommentByProductCodeRequest{
		ProductCode: productCodeParam,
	})
	if err != nil {
		if err.Error() == util.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", comments.GetRatingFromComments())
}

func (client *ServiceClient) FindByUserOrderId(c *gin.Context) {
	userOrderIdParam, err := strconv.ParseInt(c.Param("user_order_id"), 10, 64)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	comment, err := client.CommentService.FindByUserOrderId(c.Request.Context(), &pb.GetCommentByUserOrderIdRequest{
		UserOrderId: userOrderIdParam,
	})
	if err != nil {
		if err.Error() == util.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", comment.GetComment())
}

func (client *ServiceClient) Create(c *gin.Context) {
	var requestCreate CreateCommentRequest
	if err := c.ShouldBindJSON(&requestCreate); err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	comment, err := client.CommentService.Create(c.Request.Context(), &pb.CreateCommentRequest{
		UserOrderId: requestCreate.UserOrderId,
		ProductCode: requestCreate.ProductCode,
		Description: requestCreate.Description,
		Tag:         requestCreate.Tag,
		Rating:      requestCreate.Rating,
	})
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", comment.GetComment())
}
