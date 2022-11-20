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

var _ = Describe("Comment API", func() {
	var server *gin.Engine
	var tokenJWT string
	var cookie *http.Cookie
	var order *pb.UserOrder
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
		name := "Bantal"
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

		order = userOrderOneResponse.UserOrder
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

	Describe("Create Comment /comments", func() {
		When("the fields are filled", func() {
			It("should return successful create comment response", func() {
				// Create Comment
				stringBody := fmt.Sprintf(`{"user_order_id": %v,"product_code": "Lfanp","description": "mantap bang","tag": "keren, mantap","rating": 4}`, order.IdOrder)
				requestBody := strings.NewReader(stringBody)
				request := httptest.NewRequest(http.MethodPost, "/api/comments", requestBody)
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

				Expect(responseBody["data"].(map[string]interface{})["product_code"]).To(Equal("Lfanp"))
				Expect(responseBody["data"].(map[string]interface{})["description"]).To(Equal("mantap bang"))
				Expect(responseBody["data"].(map[string]interface{})["tag"]).To(Equal("keren, mantap"))
				Expect(int(responseBody["data"].(map[string]interface{})["rating"].(float64))).To(Equal(4))
			})
		})
	})

	Describe("Find By Id Comment /comments/:id", func() {
		When("comment is not found", func() {
			It("should return error not found", func() {
				// Find By Id Comment
				request := httptest.NewRequest(http.MethodGet, "/api/comments/1", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusNotFound))
				Expect(responseBody["status"]).To(Equal("data not found"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("comment is found", func() {
			It("should return a successful find comment by id", func() {
				// Create Comment
				stringBody := fmt.Sprintf(`{"user_order_id": %v,"product_code": "Lfanp","description": "mantap bang","tag": "keren, mantap","rating": 4}`, order.IdOrder)
				requestBody := strings.NewReader(stringBody)
				request := httptest.NewRequest(http.MethodPost, "/api/comments", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBodyCreateComment map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBodyCreateComment)
				log.Println(responseBodyCreateComment)
				idComment := int(responseBodyCreateComment["data"].(map[string]interface{})["id_comment"].(float64))

				// Find By Id Comment
				request = httptest.NewRequest(http.MethodGet, "/api/comments/"+util.IntToStr(idComment), nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))
				Expect(responseBody["data"].(map[string]interface{})["product_code"]).To(Equal("Lfanp"))
				Expect(responseBody["data"].(map[string]interface{})["description"]).To(Equal("mantap bang"))
				Expect(responseBody["data"].(map[string]interface{})["tag"]).To(Equal("keren, mantap"))
				Expect(int(responseBody["data"].(map[string]interface{})["rating"].(float64))).To(Equal(4))
			})
		})
	})

	Describe("Find Comment By User Order Id /comments/user-order/:user_order_id", func() {
		When("comment is not found", func() {
			It("should return error not found", func() {
				// Find Comment By User Order Id
				request := httptest.NewRequest(http.MethodGet, "/api/comments/user-order/1", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusNotFound))
				Expect(responseBody["status"]).To(Equal("data not found"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("comment is found", func() {
			It("should return a successful find comment by user order id", func() {
				// Create Comment
				stringBody := fmt.Sprintf(`{"user_order_id": %v,"product_code": "Lfanp","description": "mantap bang","tag": "keren, mantap","rating": 4}`, order.IdOrder)
				requestBody := strings.NewReader(stringBody)
				request := httptest.NewRequest(http.MethodPost, "/api/comments", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBodyCreateComment map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBodyCreateComment)
				userOrderId := int(responseBodyCreateComment["data"].(map[string]interface{})["user_order_id"].(float64))

				// Find Comment By User Order Id
				request = httptest.NewRequest(http.MethodGet, "/api/comments/user-order/"+util.IntToStr(userOrderId), nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))
				Expect(responseBody["data"].(map[string]interface{})["product_code"]).To(Equal("Lfanp"))
				Expect(responseBody["data"].(map[string]interface{})["description"]).To(Equal("mantap bang"))
				Expect(responseBody["data"].(map[string]interface{})["tag"]).To(Equal("keren, mantap"))
				Expect(int(responseBody["data"].(map[string]interface{})["rating"].(float64))).To(Equal(4))
			})
		})
	})

	Describe("Find All Rating By Product Code /comments/rating/:product_code", func() {
		When("rating's comment by product code is not found", func() {
			It("should return error not found", func() {
				// Find All Rating By Product Code
				request := httptest.NewRequest(http.MethodGet, "/api/comments/rating/XX1", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusNotFound))
				Expect(responseBody["status"]).To(Equal("data not found"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("rating's comment is exists", func() {
			It("should return a successful find comment by user order id", func() {
				// Create Comment
				stringBody := fmt.Sprintf(`{"user_order_id": %v,"product_code": "Lfanp","description": "mantap bang","tag": "keren, mantap","rating": 4}`, order.IdOrder)
				requestBody := strings.NewReader(stringBody)
				request := httptest.NewRequest(http.MethodPost, "/api/comments", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBodyCreateCommentOne map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBodyCreateCommentOne)
				productCode := responseBodyCreateCommentOne["data"].(map[string]interface{})["product_code"].(string)

				// Create Comment
				stringBody = fmt.Sprintf(`{"user_order_id": %v,"product_code": "Lfanp","tag": "keren, mantap","rating": 3}`, order.IdOrder)
				requestBody = strings.NewReader(stringBody)
				request = httptest.NewRequest(http.MethodPost, "/api/comments", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				// Create Comment
				stringBody = fmt.Sprintf(`{"user_order_id": %v,"product_code": "Lfanp","tag": "keren, mantap","rating": 4}`, order.IdOrder)
				requestBody = strings.NewReader(stringBody)
				request = httptest.NewRequest(http.MethodPost, "/api/comments", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				// Find All Rating By Product Code
				request = httptest.NewRequest(http.MethodGet, "/api/comments/rating/"+productCode, nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))

				commentResponse := responseBody["data"].([]interface{})
				Expect(int(commentResponse[0].(map[string]interface{})["rating"].(float64))).To(Equal(4))
				Expect(int(commentResponse[0].(map[string]interface{})["result_comment"].(float64))).To(Equal(1))
				Expect(int(commentResponse[0].(map[string]interface{})["result_rating"].(float64))).To(Equal(8))

				Expect(int(commentResponse[1].(map[string]interface{})["rating"].(float64))).To(Equal(3))
				Expect(int(commentResponse[1].(map[string]interface{})["result_comment"].(float64))).To(Equal(0))
				Expect(int(commentResponse[1].(map[string]interface{})["result_rating"].(float64))).To(Equal(3))
			})
		})
	})

	Describe("Find All Comment By Product Code /comments/product/:product_code", func() {
		When("comment by product code is not found", func() {
			It("should return error not found", func() {
				// Find All Comment By Product Code
				request := httptest.NewRequest(http.MethodGet, "/api/comments/product/XX1", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusNotFound))
				Expect(responseBody["status"]).To(Equal("data not found"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("comment is exists", func() {
			It("should return a successful find all comment by product code", func() {
				// Create Comment
				stringBody := fmt.Sprintf(`{"user_order_id": %v,"product_code": "Lfanp","description": "mantap bang","tag": "keren, mantap","rating": 4}`, order.IdOrder)
				requestBody := strings.NewReader(stringBody)
				request := httptest.NewRequest(http.MethodPost, "/api/comments", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBodyCreateCommentOne map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBodyCreateCommentOne)
				productCode := responseBodyCreateCommentOne["data"].(map[string]interface{})["product_code"].(string)

				// Create Comment
				stringBody = fmt.Sprintf(`{"user_order_id": %v,"product_code": "Lfanp","tag": "jelek, tidak memuaskan","rating": 2}`, order.IdOrder)
				requestBody = strings.NewReader(stringBody)
				request = httptest.NewRequest(http.MethodPost, "/api/comments", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				// Find All Comment By Product Code
				request = httptest.NewRequest(http.MethodGet, "/api/comments/product/"+productCode, nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))

				commentResponse := responseBody["data"].([]interface{})
				Expect(commentResponse[1].(map[string]interface{})["product_code"]).To(Equal("Lfanp"))
				Expect(commentResponse[1].(map[string]interface{})["description"]).To(Equal("mantap bang"))
				Expect(commentResponse[1].(map[string]interface{})["tag"]).To(Equal("keren, mantap"))
				Expect(int(commentResponse[1].(map[string]interface{})["rating"].(float64))).To(Equal(4))

				Expect(commentResponse[0].(map[string]interface{})["product_code"]).To(Equal("Lfanp"))
				Expect(commentResponse[0].(map[string]interface{})["description"]).To(Equal(""))
				Expect(commentResponse[0].(map[string]interface{})["tag"]).To(Equal("jelek, tidak memuaskan"))
				Expect(int(commentResponse[0].(map[string]interface{})["rating"].(float64))).To(Equal(2))
			})
		})
	})
})
