package integration

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	model2 "github.com/arvians-id/go-apriori-microservice/adapter/model"
	"github.com/arvians-id/go-apriori-microservice/config"
	commentRepository "github.com/arvians-id/go-apriori-microservice/services/comment-service/repository"
	paymentRepository "github.com/arvians-id/go-apriori-microservice/services/payment-service/repository"
	productRepository "github.com/arvians-id/go-apriori-microservice/services/product-service/repository"
	userOrderRepository "github.com/arvians-id/go-apriori-microservice/services/user-order-service/repository"
	userRepository "github.com/arvians-id/go-apriori-microservice/services/user-service/repository"
	"github.com/arvians-id/go-apriori-microservice/tests/setup"
	"github.com/arvians-id/go-apriori-microservice/third-party/redis"
	"github.com/arvians-id/go-apriori-microservice/util"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

var _ = Describe("Comment API", func() {
	var server *gin.Engine
	var database *sql.DB
	var tokenJWT string
	var cookie *http.Cookie
	var order *model2.UserOrder
	configuration, err := config.LoadConfig("../../config/envs")
	if err != nil {
		log.Fatal(err)
	}

	BeforeEach(func() {
		// Setup Configuration
		router, db := setup.ModuleSetup(configuration)

		database = db
		server = router

		// User Authentication
		// Create user
		tx, _ := database.Begin()
		userRepository := userRepository.NewUserRepository()
		password, _ := bcrypt.GenerateFromPassword([]byte("Rahasia123"), bcrypt.DefaultCost)
		user, _ := userRepository.Create(context.Background(), tx, &model2.User{
			Role:      1,
			Name:      "Widdy",
			Email:     "widdy@gmail.com",
			Password:  string(password),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		_ = tx.Commit()

		// Login
		requestBody := strings.NewReader(`{"email": "widdy@gmail.com","password":"Rahasia123"}`)
		request := httptest.NewRequest(http.MethodPost, "/api/auth/login", requestBody)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("X-API-KEY", configuration.XApiKey)

		writer := httptest.NewRecorder()
		server.ServeHTTP(writer, request)

		var responseBody map[string]interface{}
		_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

		tokenJWT = responseBody["data"].(map[string]interface{})["access_token"].(string)
		for _, c := range writer.Result().Cookies() {
			if c.Name == "token" {
				cookie = c
			}
		}

		// Create product
		tx, _ = database.Begin()
		productRepository := productRepository.NewProductRepository()
		description := "Test Bang"
		_, _ = productRepository.Create(context.Background(), tx, &model2.Product{
			Code:        "Lfanp",
			Name:        "Bantal Biasa",
			Description: &description,
			Category:    "Bantal, Kasur",
			Mass:        1000,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})

		// Create payload
		payloadRepository := paymentRepository.NewPaymentRepository()
		payload, _ := payloadRepository.Create(context.Background(), tx, &model2.Payment{
			UserId: user.IdUser,
		})

		// Create User Order
		userOrderRepository := userOrderRepository.NewUserOrderRepository()
		code := "aXksCj2"
		name := "Bantal"
		price := int64(20000)
		image := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/assets/%s", configuration.AwsBucket, configuration.AwsRegion, "no-image.png")
		var quantity int32 = 1
		totalPriceItem := int64(20000)
		userOrder, _ := userOrderRepository.Create(context.Background(), tx, &model2.UserOrder{
			PayloadId:      payload.IdPayload,
			Code:           &code,
			Name:           &name,
			Price:          &price,
			Image:          &image,
			Quantity:       &quantity,
			TotalPriceItem: &totalPriceItem,
		})
		_ = tx.Commit()

		order = userOrder
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
				tx, _ := database.Begin()
				commentRepository := commentRepository.NewCommentRepository()
				description := "mantap bang"
				tag := "keren, mantap"
				comment, _ := commentRepository.Create(context.Background(), tx, &model2.Comment{
					UserOrderId: order.IdOrder,
					ProductCode: "Lfanp",
					Description: &description,
					Tag:         &tag,
					Rating:      4,
				})
				_ = tx.Commit()

				// Find By Id Comment
				request := httptest.NewRequest(http.MethodGet, "/api/comments/"+util.IntToStr(int(comment.IdComment)), nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)

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
				tx, _ := database.Begin()
				commentRepository := commentRepository.NewCommentRepository()
				description := "mantap bang"
				tag := "keren, mantap"
				comment, _ := commentRepository.Create(context.Background(), tx, &model2.Comment{
					UserOrderId: order.IdOrder,
					ProductCode: "Lfanp",
					Description: &description,
					Tag:         &tag,
					Rating:      4,
				})
				_ = tx.Commit()

				// Find Comment By User Order Id
				request := httptest.NewRequest(http.MethodGet, "/api/comments/user-order/"+util.IntToStr(int(comment.UserOrderId)), nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)

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
				tx, _ := database.Begin()
				commentRepository := commentRepository.NewCommentRepository()
				description := "mantap bang"
				tag := "keren, mantap"
				comment1, _ := commentRepository.Create(context.Background(), tx, &model2.Comment{
					UserOrderId: order.IdOrder,
					ProductCode: "Lfanp",
					Description: &description,
					Tag:         &tag,
					Rating:      4,
				})

				_, _ = commentRepository.Create(context.Background(), tx, &model2.Comment{
					UserOrderId: order.IdOrder,
					ProductCode: "Lfanp",
					Rating:      3,
				})

				_, _ = commentRepository.Create(context.Background(), tx, &model2.Comment{
					UserOrderId: order.IdOrder,
					ProductCode: "Lfanp",
					Rating:      4,
				})
				_ = tx.Commit()

				// Find All Rating By Product Code
				request := httptest.NewRequest(http.MethodGet, "/api/comments/rating/"+comment1.ProductCode, nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)

				writer := httptest.NewRecorder()
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
				tx, _ := database.Begin()
				commentRepository := commentRepository.NewCommentRepository()
				description := "mantap bang"
				tag := "keren, mantap"
				comment1, _ := commentRepository.Create(context.Background(), tx, &model2.Comment{
					UserOrderId: order.IdOrder,
					ProductCode: "Lfanp",
					Description: &description,
					Tag:         &tag,
					Rating:      4,
				})

				tag = "jelek, tidak memuaskan"
				_, _ = commentRepository.Create(context.Background(), tx, &model2.Comment{
					UserOrderId: order.IdOrder,
					ProductCode: "Lfanp",
					Tag:         &tag,
					Rating:      2,
				})
				tx.Commit()

				// Find All Comment By Product Code
				request := httptest.NewRequest(http.MethodGet, "/api/comments/product/"+comment1.ProductCode, nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)

				writer := httptest.NewRecorder()
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
				Expect(commentResponse[0].(map[string]interface{})["description"]).To(BeNil())
				Expect(commentResponse[0].(map[string]interface{})["tag"]).To(Equal("jelek, tidak memuaskan"))
				Expect(int(commentResponse[0].(map[string]interface{})["rating"].(float64))).To(Equal(2))
			})
		})
	})
})
