package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/arvians-id/go-apriori-microservice/adapter/config"
	"github.com/arvians-id/go-apriori-microservice/adapter/pb"
	"github.com/arvians-id/go-apriori-microservice/adapter/pkg/notification"
	"github.com/arvians-id/go-apriori-microservice/adapter/tests/setup"
	"github.com/arvians-id/go-apriori-microservice/adapter/third-party/redis"
	"github.com/arvians-id/go-apriori-microservice/adapter/util"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
)

var _ = Describe("Notification API", func() {
	var server *gin.Engine
	var tokenJWT string
	var cookie *http.Cookie
	var notification1 *pb.Notification
	configuration, err := config.LoadConfig("../../config/envs")
	if err != nil {
		log.Fatal(err)
	}

	BeforeEach(func() {
		// Setup Configuration
		router, _ := setup.ModuleSetup(configuration)
		server = router

		// User Authentication
		// Create user
		requestBody := strings.NewReader(`{"name": "Widdy","email": "widdy@gmail.com","address":"nganjok","phone":"082299","password": "Rahasia123"}`)
		request := httptest.NewRequest(http.MethodPost, "/api/auth/register", requestBody)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("X-API-KEY", configuration.XApiKey)

		writer := httptest.NewRecorder()
		server.ServeHTTP(writer, request)

		var responseBodyCreateUser map[string]interface{}
		_ = json.NewDecoder(writer.Result().Body).Decode(&responseBodyCreateUser)
		idUser := int64(responseBodyCreateUser["data"].(map[string]interface{})["id_user"].(float64))

		// Login
		requestBody = strings.NewReader(`{"email": "widdy@gmail.com","password":"Rahasia123"}`)
		request = httptest.NewRequest(http.MethodPost, "/api/auth/login", requestBody)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("X-API-KEY", configuration.XApiKey)

		writer = httptest.NewRecorder()
		server.ServeHTTP(writer, request)

		var responseBody map[string]interface{}
		_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

		tokenJWT = responseBody["data"].(map[string]interface{})["access_token"].(string)
		for _, c := range writer.Result().Cookies() {
			if c.Name == "token" {
				cookie = c
			}
		}

		notificationService := notification.NewNotificationServiceClient(configuration)

		notificationOne, _ := notificationService.Create(context.Background(), &pb.CreateNotificationRequest{
			UserId:      idUser,
			Title:       "First notification",
			Description: "This is first notification",
			URL:         "https://google.com",
			IsRead:      false,
		})

		_, _ = notificationService.Create(context.Background(), &pb.CreateNotificationRequest{
			UserId:      idUser,
			Title:       "Second notification",
			Description: "This is second notification",
			URL:         "https://facebook.com",
			IsRead:      false,
		})

		notification1 = notificationOne.Notification
	})

	AfterEach(func() {
		// Setup Configuration
		_, db := setup.ModuleSetup(configuration)
		defer db.Close()

		cacheService := redis.NewCacheService(configuration)
		_ = cacheService.FlushDB(context.Background())

		err := setup.TearDownTest(db)
		if err != nil {
			log.Fatal(err)
		}
	})

	Describe("Find All Notification /notifications", func() {
		When("the notification is exists", func() {
			It("should return successful find all notifications response", func() {
				// Find All Notification
				request := httptest.NewRequest(http.MethodGet, "/api/notifications", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))

				userOrderResponse := responseBody["data"].([]interface{})
				Expect(userOrderResponse[1].(map[string]interface{})["title"]).To(Equal("First Notification"))
				Expect(userOrderResponse[1].(map[string]interface{})["description"]).To(Equal("This is first notification"))
				Expect(userOrderResponse[1].(map[string]interface{})["url"]).To(Equal("https://google.com"))
				Expect(userOrderResponse[1].(map[string]interface{})["is_read"].(bool)).To(BeFalse())

				Expect(userOrderResponse[0].(map[string]interface{})["title"]).To(Equal("Second Notification"))
				Expect(userOrderResponse[0].(map[string]interface{})["description"]).To(Equal("This is second notification"))
				Expect(userOrderResponse[0].(map[string]interface{})["url"]).To(Equal("https://facebook.com"))
				Expect(userOrderResponse[0].(map[string]interface{})["is_read"].(bool)).To(BeFalse())
			})
		})
	})

	Describe("Find All Notification By User Id /notifications/user", func() {
		When("the notification is exists", func() {
			It("should return successful find all notifications response", func() {
				// Find All Notification By User Id
				request := httptest.NewRequest(http.MethodGet, "/api/notifications/user", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))

				userOrderResponse := responseBody["data"].([]interface{})
				Expect(userOrderResponse[1].(map[string]interface{})["title"]).To(Equal("First Notification"))
				Expect(userOrderResponse[1].(map[string]interface{})["description"]).To(Equal("This is first notification"))
				Expect(userOrderResponse[1].(map[string]interface{})["url"]).To(Equal("https://google.com"))
				Expect(userOrderResponse[1].(map[string]interface{})["is_read"].(bool)).To(BeFalse())

				Expect(userOrderResponse[0].(map[string]interface{})["title"]).To(Equal("Second Notification"))
				Expect(userOrderResponse[0].(map[string]interface{})["description"]).To(Equal("This is second notification"))
				Expect(userOrderResponse[0].(map[string]interface{})["url"]).To(Equal("https://facebook.com"))
				Expect(userOrderResponse[0].(map[string]interface{})["is_read"].(bool)).To(BeFalse())
			})
		})
	})

	Describe("Mark All Notification By User Id /notifications/mark", func() {
		When("the notification is exists", func() {
			It("should return successful find all notifications response", func() {
				// Mark All Notification By User Id
				request := httptest.NewRequest(http.MethodPatch, "/api/notifications/mark", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				// Find All Notification By User Id
				request = httptest.NewRequest(http.MethodGet, "/api/notifications/user", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))

				userOrderResponse := responseBody["data"].([]interface{})
				Expect(userOrderResponse[1].(map[string]interface{})["title"]).To(Equal("First Notification"))
				Expect(userOrderResponse[1].(map[string]interface{})["description"]).To(Equal("This is first notification"))
				Expect(userOrderResponse[1].(map[string]interface{})["url"]).To(Equal("https://google.com"))
				Expect(userOrderResponse[1].(map[string]interface{})["is_read"].(bool)).To(BeTrue())

				Expect(userOrderResponse[0].(map[string]interface{})["title"]).To(Equal("Second Notification"))
				Expect(userOrderResponse[0].(map[string]interface{})["description"]).To(Equal("This is second notification"))
				Expect(userOrderResponse[0].(map[string]interface{})["url"]).To(Equal("https://facebook.com"))
				Expect(userOrderResponse[0].(map[string]interface{})["is_read"].(bool)).To(BeTrue())
			})
		})
	})

	Describe("Mark One Notification Id /notifications/mark/:id", func() {
		When("the notification is exists", func() {
			It("should return successful find all notifications response with different is read status", func() {
				// Mark One Notification By Id
				request := httptest.NewRequest(http.MethodPatch, "/api/notifications/mark/"+util.IntToStr(int(notification1.IdNotification)), nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				// Find All Notification By User Id
				request = httptest.NewRequest(http.MethodGet, "/api/notifications/user", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))

				userOrderResponse := responseBody["data"].([]interface{})
				Expect(userOrderResponse[1].(map[string]interface{})["title"]).To(Equal("First Notification"))
				Expect(userOrderResponse[1].(map[string]interface{})["description"]).To(Equal("This is first notification"))
				Expect(userOrderResponse[1].(map[string]interface{})["url"]).To(Equal("https://google.com"))
				Expect(userOrderResponse[1].(map[string]interface{})["is_read"].(bool)).To(BeTrue())

				Expect(userOrderResponse[0].(map[string]interface{})["title"]).To(Equal("Second Notification"))
				Expect(userOrderResponse[0].(map[string]interface{})["description"]).To(Equal("This is second notification"))
				Expect(userOrderResponse[0].(map[string]interface{})["url"]).To(Equal("https://facebook.com"))
				Expect(userOrderResponse[0].(map[string]interface{})["is_read"].(bool)).To(BeFalse())
			})
		})
	})
})
