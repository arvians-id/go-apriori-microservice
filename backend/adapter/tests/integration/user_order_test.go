package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/arvians-id/go-apriori-microservice/adapter/config"
	"github.com/arvians-id/go-apriori-microservice/adapter/pb"
	"github.com/arvians-id/go-apriori-microservice/adapter/pkg/payment"
	user_order "github.com/arvians-id/go-apriori-microservice/adapter/pkg/user-order"
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

/*
	Error :
		- Find All User Order /user-order/user
*/
var _ = Describe("User Order API", func() {
	var server *gin.Engine
	var tokenJWT string
	var cookie *http.Cookie
	var paymentPayload *pb.Payment
	var userOrder *pb.UserOrder
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

		// Create Product
		requestBody = strings.NewReader(`{"code": "Lfanp","name": "Bantal Biasa","description": "Test Bang","category": "Bantal, Kasur","mass": 1000,"price": 7000}`)
		request = httptest.NewRequest(http.MethodPost, "/api/products", requestBody)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("X-API-KEY", configuration.XApiKey)
		request.AddCookie(cookie)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

		writer = httptest.NewRecorder()
		server.ServeHTTP(writer, request)

		// Create payment
		paymentService := payment.NewCommentServiceClient(configuration)
		paymentResponse, _ := paymentService.OnlyCreate(context.Background(), &pb.OnlyCreatePaymentRequest{
			UserId:         idUser,
			Address:        "Jl Merdeka",
			Courier:        "Widdy",
			CourierService: "TokopediaCourier",
		})

		// Create User Order
		userOrderService := user_order.NewUserOrderServiceClient(configuration)

		code := "aXksCj2"
		name := "Bantal Biasa"
		price := int64(20000)
		image := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/assets/%s", configuration.AwsBucket, configuration.AwsRegion, "no-image.png")
		var quantity int32 = 1
		totalPriceItem := int64(20000)
		userOrderOneResponse, _ := userOrderService.Create(context.Background(), &pb.CreateUserOrderRequest{
			PayloadId:      paymentResponse.Payment.IdPayload,
			Code:           &code,
			Name:           &name,
			Price:          &price,
			Image:          &image,
			Quantity:       &quantity,
			TotalPriceItem: &totalPriceItem,
		})

		name = "Guling"
		price = int64(10000)
		quantity = 2
		totalPriceItem = int64(20000)
		userOrderService.Create(context.Background(), &pb.CreateUserOrderRequest{
			PayloadId:      paymentResponse.Payment.IdPayload,
			Code:           &code,
			Name:           &name,
			Price:          &price,
			Image:          &image,
			Quantity:       &quantity,
			TotalPriceItem: &totalPriceItem,
		})

		userOrder = userOrderOneResponse.UserOrder
		paymentPayload = paymentResponse.Payment
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

	Describe("Find All Payment On User Order /user-order", func() {
		When("user not logged in yet", func() {
			It("should return error unauthorized/invalid token", func() {
				// Find All Payment on User Order
				request := httptest.NewRequest(http.MethodGet, "/api/user-order", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusUnauthorized))
				Expect(responseBody["status"]).To(Equal("invalid token"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("the user order is exists", func() {
			It("should return successful find all payment on user order response", func() {
				// Find All Payment on User Order
				request := httptest.NewRequest(http.MethodGet, "/api/user-order", nil)
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
				Expect(userOrderResponse[0].(map[string]interface{})["address"]).To(Equal("Jl Merdeka"))
				Expect(userOrderResponse[0].(map[string]interface{})["courier"]).To(Equal("Widdy"))
				Expect(userOrderResponse[0].(map[string]interface{})["courier_service"]).To(Equal("TokopediaCourier"))
			})
		})
	})

	//Describe("Find All User Order /user-order/user", func() {
	//	When("the user order is exists", func() {
	//		It("should return successful find all user order response", func() {
	//			// Find All User Order
	//			request := httptest.NewRequest(http.MethodGet, "/api/user-order/user", nil)
	//			request.Header.Add("Content-Type", "application/json")
	//			request.Header.Add("X-API-KEY", configuration.XApiKey)
	//			request.AddCookie(cookie)
	//			request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))
	//
	//			writer := httptest.NewRecorder()
	//			server.ServeHTTP(writer, request)
	//
	//			var responseBody map[string]interface{}
	//			_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)
	//
	//			Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
	//			Expect(responseBody["status"]).To(Equal("OK"))
	//
	//			userOrderResponse := responseBody["data"].([]interface{})
	//
	//			Expect(userOrderResponse[0].(map[string]interface{})["code"]).To(Equal(order1.Code))
	//			Expect(userOrderResponse[0].(map[string]interface{})["name"]).To(Equal(order1.Name))
	//			Expect(int(userOrderResponse[0].(map[string]interface{})["price"].(float64))).To(Equal(order1.Price))
	//			Expect(userOrderResponse[0].(map[string]interface{})["image"]).To(Equal(order1.Image))
	//			Expect(int(userOrderResponse[0].(map[string]interface{})["quantity"].(float64))).To(Equal(order1.Quantity))
	//			Expect(int(userOrderResponse[0].(map[string]interface{})["total_price_item"].(float64))).To(Equal(order1.TotalPriceItem))
	//
	//			Expect(userOrderResponse[1].(map[string]interface{})["code"]).To(Equal(order2.Code))
	//			Expect(userOrderResponse[1].(map[string]interface{})["name"]).To(Equal(order2.Name))
	//			Expect(int(userOrderResponse[1].(map[string]interface{})["price"].(float64))).To(Equal(order2.Price))
	//			Expect(userOrderResponse[1].(map[string]interface{})["image"]).To(Equal(order2.Image))
	//			Expect(int(userOrderResponse[1].(map[string]interface{})["quantity"].(float64))).To(Equal(order2.Quantity))
	//			Expect(int(userOrderResponse[1].(map[string]interface{})["total_price_item"].(float64))).To(Equal(order2.TotalPriceItem))
	//		})
	//	})
	//})

	Describe("Find All User Order By Order Id /user-order/:order_id", func() {
		When("the user order is not found", func() {
			It("should return error not found response", func() {
				// Find All User Order
				request := httptest.NewRequest(http.MethodGet, "/api/user-order/asasdw", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusNotFound))
				Expect(responseBody["status"]).To(Equal("data not found"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("the user order is exists", func() {
			It("should return successful find all user order by order id response", func() {
				// Find All User Order
				request := httptest.NewRequest(http.MethodGet, "/api/user-order/"+*paymentPayload.OrderId, nil)
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

				Expect(userOrderResponse[0].(map[string]interface{})["code"]).To(Equal("aXksCj2"))
				Expect(userOrderResponse[0].(map[string]interface{})["name"]).To(Equal("Bantal Biasa"))
				Expect(int64(userOrderResponse[0].(map[string]interface{})["price"].(float64))).To(Equal(int64(20000)))
				Expect(int(userOrderResponse[0].(map[string]interface{})["quantity"].(float64))).To(Equal(1))
				Expect(int64(userOrderResponse[0].(map[string]interface{})["total_price_item"].(float64))).To(Equal(int64(20000)))

				Expect(userOrderResponse[1].(map[string]interface{})["code"]).To(Equal("aXksCj2"))
				Expect(userOrderResponse[1].(map[string]interface{})["name"]).To(Equal("Guling"))
				Expect(int64(userOrderResponse[1].(map[string]interface{})["price"].(float64))).To(Equal(int64(10000)))
				Expect(int(userOrderResponse[1].(map[string]interface{})["quantity"].(float64))).To(Equal(2))
				Expect(int64(userOrderResponse[1].(map[string]interface{})["total_price_item"].(float64))).To(Equal(int64(20000)))
			})
		})
	})

	Describe("Find User Order By Id /user-order/single/:id", func() {
		When("the user order is not found", func() {
			It("should return error not found response", func() {
				// Find All User Order By Id
				request := httptest.NewRequest(http.MethodGet, "/api/user-order/single/12121", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusNotFound))
				Expect(responseBody["status"]).To(Equal("data not found"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("the user order is exists", func() {
			It("should return successful find user order by id response", func() {
				// Find All User Order By Id
				request := httptest.NewRequest(http.MethodGet, "/api/user-order/single/"+util.IntToStr(int(userOrder.IdOrder)), nil)
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

				Expect(responseBody["data"].(map[string]interface{})["code"]).To(Equal("aXksCj2"))
				Expect(responseBody["data"].(map[string]interface{})["name"]).To(Equal("Bantal Biasa"))
				Expect(int64(responseBody["data"].(map[string]interface{})["price"].(float64))).To(Equal(int64(20000)))
				Expect(int(responseBody["data"].(map[string]interface{})["quantity"].(float64))).To(Equal(1))
				Expect(int64(responseBody["data"].(map[string]interface{})["total_price_item"].(float64))).To(Equal(int64(20000)))
			})
		})
	})
})
