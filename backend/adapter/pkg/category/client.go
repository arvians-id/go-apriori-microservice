package category

import (
	"github.com/arvians-id/go-apriori-microservice/adapter/config"
	"github.com/arvians-id/go-apriori-microservice/adapter/middleware"
	"github.com/arvians-id/go-apriori-microservice/adapter/pb"
	"github.com/arvians-id/go-apriori-microservice/adapter/response"
	"github.com/arvians-id/go-apriori-microservice/adapter/util"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"strconv"
)

type ServiceClient struct {
	CategoryService pb.CategoryServiceClient
}

func NewCategoryServiceClient(configuration *config.Config) pb.CategoryServiceClient {
	connection, err := grpc.Dial(configuration.CategorySvcUrl, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}

	return pb.NewCategoryServiceClient(connection)
}

func RegisterRoutes(router *gin.Engine, configuration *config.Config) *ServiceClient {
	serviceClient := &ServiceClient{
		CategoryService: NewCategoryServiceClient(configuration),
	}

	authorized := router.Group("/api", middleware.AuthJwtMiddleware(configuration))
	{
		authorized.GET("/categories/:id", serviceClient.FindById)
		authorized.POST("/categories", serviceClient.Create)
		authorized.PATCH("/categories/:id", serviceClient.Update)
		authorized.DELETE("/categories/:id", serviceClient.Delete)
	}

	unauthorized := router.Group("/api")
	{
		unauthorized.GET("/categories", serviceClient.FindAll)
	}

	return serviceClient
}

func (client *ServiceClient) FindAll(c *gin.Context) {
	categories, err := client.CategoryService.FindAll(c.Request.Context(), new(emptypb.Empty))
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", categories.GetCategories())
}

func (client *ServiceClient) FindById(c *gin.Context) {
	idParam, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	category, err := client.CategoryService.FindById(c.Request.Context(), &pb.GetCategoryByIdRequest{
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

	response.ReturnSuccessOK(c, "OK", category.GetCategory())
}

func (client *ServiceClient) Create(c *gin.Context) {
	var requestCreate CreateCategoryRequest
	if err := c.ShouldBindJSON(&requestCreate); err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	category, err := client.CategoryService.Create(c.Request.Context(), &pb.CreateCategoryRequest{
		Name: requestCreate.Name,
	})
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", category.GetCategory())
}

func (client *ServiceClient) Update(c *gin.Context) {
	var requestUpdate UpdateCategoryRequest
	if err := c.ShouldBindJSON(&requestUpdate); err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	idParam, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	category, err := client.CategoryService.Update(c.Request.Context(), &pb.UpdateCategoryRequest{
		IdCategory: idParam,
		Name:       requestUpdate.Name,
	})
	if err != nil {
		if err.Error() == util.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", category.GetCategory())
}

func (client *ServiceClient) Delete(c *gin.Context) {
	idParam, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	_, err = client.CategoryService.Delete(c.Request.Context(), &pb.GetCategoryByIdRequest{
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

	response.ReturnSuccessOK(c, "OK", nil)
}
