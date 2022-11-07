package auth

import (
	"errors"
	"fmt"
	"github.com/arvians-id/go-apriori-microservice/adapter/middleware"
	"github.com/arvians-id/go-apriori-microservice/adapter/pkg/auth/pb"
	"github.com/arvians-id/go-apriori-microservice/adapter/pkg/user"
	pbuser "github.com/arvians-id/go-apriori-microservice/adapter/pkg/user/pb"
	"github.com/arvians-id/go-apriori-microservice/adapter/response"
	"github.com/arvians-id/go-apriori-microservice/config"
	"github.com/arvians-id/go-apriori-microservice/model"
	"github.com/arvians-id/go-apriori-microservice/third-party/jwt"
	messaging "github.com/arvians-id/go-apriori-microservice/third-party/message-queue"
	"github.com/arvians-id/go-apriori-microservice/util"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type ServiceClient struct {
	PasswordResetService pb.PasswordResetServiceClient
	UserService          pbuser.UserServiceClient
	Jwt                  *jwt.JsonWebToken
	Producer             *messaging.Producer
}

func NewAuthServiceClient(configuration *config.Config) pb.PasswordResetServiceClient {
	connection, err := grpc.Dial(configuration.AuthSvcUrl, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}

	return pb.NewPasswordResetServiceClient(connection)
}

func RegisterRoutes(router *gin.Engine, configuration *config.Config, jwt *jwt.JsonWebToken, producer *messaging.Producer) *ServiceClient {
	serviceClient := &ServiceClient{
		PasswordResetService: NewAuthServiceClient(configuration),
		UserService:          user.NewUserServiceClient(configuration),
		Jwt:                  jwt,
		Producer:             producer,
	}

	authorized := router.Group("/api/jwt", middleware.AuthJwtMiddleware())
	{
		authorized.GET("/token", serviceClient.Token)
	}

	unauthorized := router.Group("/api/jwt")
	{
		unauthorized.POST("/login", serviceClient.Login)
		unauthorized.POST("/refresh", serviceClient.Refresh)
		unauthorized.POST("/forgot-password", serviceClient.ForgotPassword)
		unauthorized.POST("/verify", serviceClient.VerifyResetPassword)
		unauthorized.POST("/register", serviceClient.Register)
		unauthorized.DELETE("/logout", serviceClient.Logout)
	}

	return serviceClient
}

func (client *ServiceClient) Login(c *gin.Context) {
	var requestCredential GetUserCredentialRequest
	err := c.ShouldBindJSON(&requestCredential)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	userResponse, err := client.UserService.FindByEmail(c.Request.Context(), &pbuser.FindByEmailRequest{
		Email:    requestCredential.Email,
		Password: requestCredential.Password,
	})
	if err != nil {
		if err.Error() == util.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		if err.Error() == util.WrongPassword {
			response.ReturnErrorBadRequest(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	expiredTimeAccess, err := strconv.Atoi(os.Getenv("JWT_ACCESS_EXPIRED_TIME"))
	if err != nil {
		log.Println("[AuthController][Login] problem in conversion string to integer, err: ", err.Error())
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	expirationTime := time.Now().Add(time.Duration(expiredTimeAccess) * 24 * time.Hour)
	token, err := client.Jwt.GenerateToken(userResponse.User.IdUser, userResponse.User.Role, expirationTime)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    url.QueryEscape(token.AccessToken),
		Expires:  expirationTime,
		Path:     "/api",
		HttpOnly: true,
	})

	response.ReturnSuccessOK(c, "OK", gin.H{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	})
}

func (client *ServiceClient) Refresh(c *gin.Context) {
	var requestToken GetRefreshTokenRequest
	err := c.ShouldBindJSON(&requestToken)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	token, err := client.Jwt.RefreshToken(requestToken.RefreshToken)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	expiredTimeAccess, err := strconv.Atoi(os.Getenv("JWT_ACCESS_EXPIRED_TIME"))
	if err != nil {
		log.Println("[AuthController][Login] problem in conversion string to integer, err: ", err.Error())
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	expirationTime := time.Now().Add(time.Duration(expiredTimeAccess) * 24 * time.Hour)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    url.QueryEscape(token.AccessToken),
		Expires:  expirationTime,
		Path:     "/api",
		HttpOnly: true,
	})

	response.ReturnSuccessOK(c, "OK", gin.H{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	})
}

func (client *ServiceClient) Register(c *gin.Context) {
	var requestCreate CreateUserRequest
	err := c.ShouldBindJSON(&requestCreate)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	userResponse, err := client.UserService.Create(c.Request.Context(), &pbuser.CreateRequest{
		Name:     requestCreate.Name,
		Email:    requestCreate.Email,
		Password: requestCreate.Password,
		Address:  requestCreate.Address,
		Phone:    requestCreate.Phone,
	})
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "created", userResponse)
}

func (client *ServiceClient) ForgotPassword(c *gin.Context) {
	var requestCreate CreatePasswordResetRequest
	err := c.ShouldBindJSON(&requestCreate)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	result, err := client.PasswordResetService.CreateOrUpdateByEmail(c.Request.Context(), &pb.GetPasswordResetByEmailRequest{
		Email: requestCreate.Email,
	})
	if err != nil {
		if err.Error() == util.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	message := fmt.Sprintf("%s/jwt/reset-password?signature=%v", os.Getenv("APP_URL_FE"), result.PasswordReset.Token)
	emailService := model.EmailService{
		ToEmail: result.PasswordReset.Email,
		Subject: "Forgot Password",
		Message: message,
	}
	err = client.Producer.Publish("mail_topic", emailService)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "mail sent successfully", gin.H{
		"signature": result.PasswordReset.Token,
	})
}

func (client *ServiceClient) VerifyResetPassword(c *gin.Context) {
	var requestUpdate UpdateResetPasswordUserRequest
	err := c.ShouldBindJSON(&requestUpdate)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	_, err = client.PasswordResetService.Verify(c.Request.Context(), &pb.PasswordReset{
		Email: requestUpdate.Email,
		Token: c.Query("signature"),
	})
	if err != nil {
		if err.Error() == util.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		if err.Error() == util.VerificationExpired {
			response.ReturnErrorBadRequest(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "updated", nil)
}

func (client *ServiceClient) Logout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/api",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})

	response.ReturnSuccessOK(c, "OK", nil)
}

func (client *ServiceClient) Token(c *gin.Context) {
	_, isExist := c.Get("id_user")
	if !isExist {
		response.ReturnErrorUnauthorized(c, errors.New("unauthorized"), nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", nil)
}
