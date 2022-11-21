package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/arvians-id/go-apriori-microservice/adapter/config"
	"github.com/arvians-id/go-apriori-microservice/adapter/pb"
	"github.com/arvians-id/go-apriori-microservice/adapter/pkg/payment"
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

var _ = Describe("Payment API", func() {
	var server *gin.Engine
	var tokenJWT string
	var cookie *http.Cookie
	var userId int64
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

		userId = idUser
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

	Describe("Find All Payment /payments", func() {
		When("payment is exists", func() {
			It("should return a successful find all payment", func() {
				// Create payment
				paymentService := payment.NewCommentServiceClient(configuration)
				_, _ = paymentService.OnlyCreate(context.Background(), &pb.OnlyCreatePaymentRequest{
					UserId:         userId,
					Address:        "Jl Merdeka",
					Courier:        "Widdy",
					CourierService: "TokopediaCourier",
				})

				// Find All Payment
				request := httptest.NewRequest(http.MethodGet, "/api/payments", nil)
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

				paymentResponse := responseBody["data"].([]interface{})
				Expect(paymentResponse[0].(map[string]interface{})["address"]).To(Equal("Jl Merdeka"))
				Expect(paymentResponse[0].(map[string]interface{})["courier"]).To(Equal("Widdy"))
				Expect(paymentResponse[0].(map[string]interface{})["courier_service"]).To(Equal("TokopediaCourier"))
			})
		})
	})

	Describe("Find Payment By Order Id /payments/:order_id", func() {
		When("payment is not exists", func() {
			It("should return error not found", func() {
				// Find Payment By Order Id
				request := httptest.NewRequest(http.MethodGet, "/api/payments/asdasdsa", nil)
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

		When("payment is exists", func() {
			It("should return a successful find all payment", func() {
				// Create payment
				paymentService := payment.NewCommentServiceClient(configuration)
				paymentResponse, _ := paymentService.OnlyCreate(context.Background(), &pb.OnlyCreatePaymentRequest{
					UserId:         userId,
					Address:        "Jl Merdeka",
					Courier:        "Widdy",
					CourierService: "TokopediaCourier",
				})

				// Find Payment By Order Id
				request := httptest.NewRequest(http.MethodGet, "/api/payments/"+*paymentResponse.Payment.OrderId, nil)
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

				Expect(responseBody["data"].(map[string]interface{})["address"]).To(Equal("Jl Merdeka"))
				Expect(responseBody["data"].(map[string]interface{})["courier"]).To(Equal("Widdy"))
				Expect(responseBody["data"].(map[string]interface{})["courier_service"]).To(Equal("TokopediaCourier"))
			})
		})
	})

	Describe("Update Receipt Number Payment By Order Id /payments/:order_id", func() {
		When("payment is not exists", func() {
			It("should return error not found", func() {
				// Update Receipt Number Payment By Order Id
				stringBody := fmt.Sprintf(`{"receipt_number": "asfxxd"}`)
				requestBody := strings.NewReader(stringBody)
				request := httptest.NewRequest(http.MethodPatch, "/api/payments/asdasdsa", requestBody)
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

		When("payment is exists", func() {
			It("should return a successful update receipt number payment", func() {
				// Create payment
				paymentService := payment.NewCommentServiceClient(configuration)
				paymentResponse, _ := paymentService.OnlyCreate(context.Background(), &pb.OnlyCreatePaymentRequest{
					UserId:         userId,
					Address:        "Jl Merdeka",
					Courier:        "Widdy",
					CourierService: "TokopediaCourier",
				})

				// Update Receipt Number Payment By Order Id
				stringBody := fmt.Sprintf(`{"receipt_number": "asfxxd"}`)
				requestBody := strings.NewReader(stringBody)
				request := httptest.NewRequest(http.MethodPatch, "/api/payments/"+*paymentResponse.Payment.OrderId, requestBody)
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
				Expect(responseBody["data"]).To(BeNil())
			})
		})
	})

	Describe("Send Notification After Paid /payments/notification", func() {
		When("payment is exists", func() {
			It("should return a successful update user's payment", func() {
				// Create payment
				paymentService := payment.NewCommentServiceClient(configuration)
				paymentResponse, _ := paymentService.OnlyCreate(context.Background(), &pb.OnlyCreatePaymentRequest{
					UserId:            userId,
					TransactionStatus: "pending",
					Address:           "Jl Merdeka",
					Courier:           "Widdy",
					CourierService:    "TokopediaCourier",
				})

				// Update Receipt Number Payment By Order Id
				stringBody := fmt.Sprintf(`
{
    "va_numbers": [
        {
            "va_number": "62888785966",
            "bank": "bca"
        }
    ],
    "transaction_time": "2022-07-11 20:11:46",
    "transaction_status": "settlement",
    "transaction_id": "becb3ef8-0a3b-4eda-b97a-a147214a8b32",
    "status_message": "midtrans payment notification",
    "status_code": "200",
    "signature_key": "e079aa332e704bd5cc7a9466144fec6fd0a0c788fc2ea1b3c0666aba4ba2a87aec7b352bf23a55cbfae7f7a6bcc0c80f678e1417028b39ff9923ac913e9e91cd",
    "settlement_time": "2022-07-11 20:13:13",
    "payment_type": "bank_transfer",
    "payment_amounts": [],
    "order_id": "%s",
    "merchant_id": "G271262888",
    "gross_amount": "59996.00",
    "fraud_status": "accept",
    "custom_field1": "%s",
    "custom_field2": "26",
    "custom_field3": "WIddy",
    "currency": "IDR"
}
`, *paymentResponse.Payment.OrderId, util.Int64ToStr(userId))
				requestBody := strings.NewReader(stringBody)
				request := httptest.NewRequest(http.MethodPost, "/api/payments/notification", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)
				log.Println(responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))
				Expect(responseBody["data"]).To(BeNil())

				// Find Payment By Order Id
				request = httptest.NewRequest(http.MethodGet, "/api/payments/"+*paymentResponse.Payment.OrderId, nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBodyFindPayment map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBodyFindPayment)

				Expect(int(responseBodyFindPayment["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBodyFindPayment["status"]).To(Equal("OK"))

				Expect(responseBodyFindPayment["data"].(map[string]interface{})["transaction_status"]).To(Equal("settlement"))
				Expect(responseBodyFindPayment["data"].(map[string]interface{})["address"]).To(Equal("Jl Merdeka"))
				Expect(responseBodyFindPayment["data"].(map[string]interface{})["courier"]).To(Equal("Widdy"))
				Expect(responseBodyFindPayment["data"].(map[string]interface{})["courier_service"]).To(Equal("TokopediaCourier"))
			})
		})
	})

	Describe("Delete Payment By Order Id /payments/:order_id", func() {
		When("payment is not exists", func() {
			It("should return error not found", func() {
				// Delete Payment By Order Id
				request := httptest.NewRequest(http.MethodDelete, "/api/payments/asdasdsa", nil)
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

		When("payment is exists", func() {
			It("should return a successful find all payment", func() {
				// Create payment
				paymentService := payment.NewCommentServiceClient(configuration)
				paymentResponse, _ := paymentService.OnlyCreate(context.Background(), &pb.OnlyCreatePaymentRequest{
					UserId:         userId,
					Address:        "Jl Merdeka",
					Courier:        "Widdy",
					CourierService: "TokopediaCourier",
				})

				// Delete Payment By Order Id
				request := httptest.NewRequest(http.MethodDelete, "/api/payments/"+*paymentResponse.Payment.OrderId, nil)
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
				Expect(responseBody["data"]).To(BeNil())
			})
		})
	})
})
