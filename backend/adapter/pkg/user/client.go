package user

import (
	"errors"
	"github.com/arvians-id/go-apriori-microservice/adapter/config"
	"github.com/arvians-id/go-apriori-microservice/adapter/middleware"
	"github.com/arvians-id/go-apriori-microservice/adapter/pb"
	"github.com/arvians-id/go-apriori-microservice/adapter/response"
	"github.com/arvians-id/go-apriori-microservice/adapter/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"log"
	"strconv"
)

type ServiceClient struct {
	UserService pb.UserServiceClient
}

func NewUserServiceClient(configuration *config.Config) pb.UserServiceClient {
	connection, err := grpc.Dial(configuration.UserSvcUrl, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}

	return pb.NewUserServiceClient(connection)
}

func RegisterRoutes(router *gin.Engine, configuration *config.Config) *ServiceClient {
	serviceClient := &ServiceClient{
		UserService: NewUserServiceClient(configuration),
	}

	authorized := router.Group("/api", middleware.AuthJwtMiddleware(configuration))
	{
		authorized.GET("/users", serviceClient.FindAll)
		authorized.GET("/users/:id", serviceClient.FindById)
		authorized.POST("/users", serviceClient.Create)
		authorized.PATCH("/users/:id", serviceClient.Update)
		authorized.DELETE("/users/:id", serviceClient.Delete)
		authorized.GET("/profile", serviceClient.Profile)
		authorized.PATCH("/profile/update", serviceClient.UpdateProfile)
	}

	return serviceClient
}

func (client *ServiceClient) Profile(c *gin.Context) {
	id, isExist := c.Get("id_user")
	if !isExist {
		response.ReturnErrorUnauthorized(c, errors.New("unauthorized"), nil)
		return
	}

	user, err := client.UserService.FindById(c.Request.Context(), &pb.GetUserByIdRequest{
		Id: int64(id.(float64)),
	})
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", user.GetUser())
}

func (client *ServiceClient) UpdateProfile(c *gin.Context) {
	var requestUpdate UpdateUserRequest
	err := c.ShouldBindJSON(&requestUpdate)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	id, isExist := c.Get("id_user")
	if !isExist {
		response.ReturnErrorUnauthorized(c, errors.New("unauthorized"), nil)
		return
	}

	user, err := client.UserService.Update(c.Request.Context(), &pb.UpdateUserRequest{
		IdUser:   int64(id.(float64)),
		Role:     requestUpdate.Role,
		Name:     requestUpdate.Name,
		Email:    requestUpdate.Email,
		Address:  requestUpdate.Address,
		Phone:    requestUpdate.Phone,
		Password: requestUpdate.Password,
	})
	if err != nil {
		if err.Error() == util.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "updated", user.GetUser())
}

func (client *ServiceClient) FindAll(c *gin.Context) {
	users, err := client.UserService.FindAll(c.Request.Context(), new(empty.Empty))
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", users.GetUser())
}

func (client *ServiceClient) FindById(c *gin.Context) {
	idParam, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	user, err := client.UserService.FindById(c.Request.Context(), &pb.GetUserByIdRequest{
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

	response.ReturnSuccessOK(c, "OK", user.GetUser())
}

func (client *ServiceClient) Create(c *gin.Context) {
	var requestCreate CreateUserRequest
	err := c.ShouldBindJSON(&requestCreate)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	user, err := client.UserService.Create(c.Request.Context(), &pb.CreateUserRequest{
		Name:     requestCreate.Name,
		Email:    requestCreate.Email,
		Address:  requestCreate.Address,
		Phone:    requestCreate.Phone,
		Password: requestCreate.Password,
	})
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "created", user.GetUser())
}

func (client *ServiceClient) Update(c *gin.Context) {
	var requestUpdate UpdateUserRequest
	err := c.ShouldBindJSON(&requestUpdate)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	idParam, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	user, err := client.UserService.Update(c.Request.Context(), &pb.UpdateUserRequest{
		IdUser:   idParam,
		Role:     requestUpdate.Role,
		Name:     requestUpdate.Name,
		Email:    requestUpdate.Email,
		Address:  requestUpdate.Address,
		Phone:    requestUpdate.Phone,
		Password: requestUpdate.Password,
	})
	if err != nil {
		if err.Error() == util.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "updated", user.GetUser())
}

func (client *ServiceClient) Delete(c *gin.Context) {
	idParam, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	_, err = client.UserService.Delete(c.Request.Context(), &pb.GetUserByIdRequest{
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

	response.ReturnSuccessOK(c, "deleted", nil)
}
