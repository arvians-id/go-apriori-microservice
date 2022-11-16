package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/arvians-id/go-apriori-microservice/adapter/config"
	"github.com/arvians-id/go-apriori-microservice/adapter/pb"
	"github.com/arvians-id/go-apriori-microservice/adapter/pkg/apriori"
	"github.com/arvians-id/go-apriori-microservice/adapter/tests/setup"
	"github.com/arvians-id/go-apriori-microservice/adapter/third-party/redis"
	. "github.com/onsi/ginkgo/v2"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/gomega"
)

var _ = Describe("Product API", func() {
	var server *gin.Engine
	var tokenJWT string
	var cookie *http.Cookie
	var aprioriService pb.AprioriServiceClient
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

		aprioriServiceClient := apriori.NewAprioriServiceClient(configuration)
		aprioriService = aprioriServiceClient
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

	Describe("Create Product /products", func() {
		When("the fields are correct", func() {
			When("the fields are filled", func() {
				It("should return successful create product response", func() {
					// Create Product
					requestBody := strings.NewReader(`{"code": "SK6","name": "Bantal Biasa","description": "Test","category": "Bantal, Kasur","mass": 1000,"price": 7000}`)
					request := httptest.NewRequest(http.MethodPost, "/api/products", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.XApiKey)
					request.AddCookie(cookie)
					request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					var responseBody map[string]interface{}
					_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
					Expect(responseBody["status"]).To(Equal("created"))
					Expect(responseBody["data"].(map[string]interface{})["code"]).To(Equal("SK6"))
					Expect(responseBody["data"].(map[string]interface{})["name"]).To(Equal("Bantal Biasa"))
					Expect(responseBody["data"].(map[string]interface{})["description"]).To(Equal("Test"))
					Expect(int(responseBody["data"].(map[string]interface{})["price"].(float64))).To(Equal(7000))
				})
			})
		})
	})

	Describe("Update Product /products/:code", func() {
		When("the product is not found", func() {
			It("should return error not found", func() {
				// Update Product
				requestBody := strings.NewReader(`{"code": "SK1","name": "Bantal Biasa","description": "Test","category": "Bantal, Kasur","mass":1000}`)
				request := httptest.NewRequest(http.MethodPatch, "/api/products/SK1", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusNotFound))
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("the fields are correct", func() {
			When("the fields are filled", func() {
				It("should return successful update product response", func() {
					// Create Product
					requestBody := strings.NewReader(`{"code": "SK6","name": "Guling","description": "Test","category": "Bantal, Kasur","mass": 1000,"price": 7000}`)
					request := httptest.NewRequest(http.MethodPost, "/api/products", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.XApiKey)
					request.AddCookie(cookie)
					request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					var responseBodyCreateProduct map[string]interface{}
					_ = json.NewDecoder(writer.Result().Body).Decode(&responseBodyCreateProduct)
					productCode := responseBodyCreateProduct["data"].(map[string]interface{})["code"].(string)

					// Update Product
					requestBody = strings.NewReader(`{"code": "SK1","name": "Guling Doti","description": "Test Bang","category": "Bantal, Kasur","mass":1000}`)
					request = httptest.NewRequest(http.MethodPatch, "/api/products/"+productCode, requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.XApiKey)
					request.AddCookie(cookie)
					request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

					writer = httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					var responseBody map[string]interface{}
					_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
					Expect(responseBody["status"]).To(Equal("updated"))
					Expect(responseBody["data"].(map[string]interface{})["name"]).ShouldNot(Equal("Guling"))
					Expect(responseBody["data"].(map[string]interface{})["description"]).ShouldNot(Equal("Test"))
				})
			})
		})
	})

	Describe("Delete Product /products/:code", func() {
		When("product is not found", func() {
			It("should return error not found", func() {
				// Delete Product
				request := httptest.NewRequest(http.MethodDelete, "/api/products/SK9", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusNotFound))
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("product is found", func() {
			It("should return a successful delete product response", func() {
				// Create Product
				requestBody := strings.NewReader(`{"code": "SK6","name": "Guling","description": "Test","category": "Bantal, Kasur","mass": 1000,"price": 7000}`)
				request := httptest.NewRequest(http.MethodPost, "/api/products", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBodyCreateProduct map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBodyCreateProduct)
				productCode := responseBodyCreateProduct["data"].(map[string]interface{})["code"].(string)

				// Delete Product
				request = httptest.NewRequest(http.MethodDelete, "/api/products/"+productCode, nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("deleted"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})
	})

	Describe("Find All Product /products", func() {
		When("the product is not present", func() {
			It("should return a successful but the data is null", func() {
				// Find All Product
				request := httptest.NewRequest(http.MethodGet, "/api/products", nil)
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

		When("the product is present", func() {
			It("should return a successful and show all products", func() {
				// Create Product One
				requestBody := strings.NewReader(`{"code": "SK6","name": "Guling","description": "Test","category": "Bantal, Kasur","mass": 1000,"price": 7000}`)
				request := httptest.NewRequest(http.MethodPost, "/api/products", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBodyCreateProductOne map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBodyCreateProductOne)
				productCodeOne := responseBodyCreateProductOne["data"].(map[string]interface{})["code"].(string)
				productNameOne := responseBodyCreateProductOne["data"].(map[string]interface{})["name"].(string)

				// Create Product Two
				requestBody = strings.NewReader(`{"code": "SK1","name": "Bantal","description": "Test Bang","category": "Bantal, Kasur","mass": 1000,"price": 7000}`)
				request = httptest.NewRequest(http.MethodPost, "/api/products", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBodyCreateProductTwo map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBodyCreateProductTwo)
				productCodeTwo := responseBodyCreateProductTwo["data"].(map[string]interface{})["code"].(string)
				productNameTwo := responseBodyCreateProductTwo["data"].(map[string]interface{})["name"].(string)

				// Find All Products
				request = httptest.NewRequest(http.MethodGet, "/api/products", nil)
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

				products := responseBody["data"].([]interface{})
				Expect(productCodeOne).To(Equal(products[1].(map[string]interface{})["code"]))
				Expect(productNameOne).To(Equal(products[1].(map[string]interface{})["name"]))

				Expect(productCodeTwo).To(Equal(products[0].(map[string]interface{})["code"]))
				Expect(productNameTwo).To(Equal(products[0].(map[string]interface{})["name"]))
			})
		})
	})

	Describe("Find All Product On Admin /products-admin", func() {
		When("the product is not present", func() {
			It("should return a successful but the data is null", func() {
				// Find All Product
				request := httptest.NewRequest(http.MethodGet, "/api/products-admin", nil)
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

		When("the product is present", func() {
			It("should return a successful and show all products with different status empty", func() {
				// Create Product One
				requestBody := strings.NewReader(`{"code": "SK6","name": "Guling","description": "Test","category": "Bantal, Kasur","mass": 1000,"price": 7000}`)
				request := httptest.NewRequest(http.MethodPost, "/api/products", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBodyCreateProductOne map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBodyCreateProductOne)
				productCodeOne := responseBodyCreateProductOne["data"].(map[string]interface{})["code"].(string)

				// Create Product Two
				requestBody = strings.NewReader(`{"code": "SK1","name": "Bantal","description": "Test Bang","category": "Bantal, Kasur","mass": 1000,"price": 7000}`)
				request = httptest.NewRequest(http.MethodPost, "/api/products", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBodyCreateProductTwo map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBodyCreateProductTwo)

				// Update Product
				requestBody = strings.NewReader(`{"name": "Guling Doti Bang","description": "Test Bang","category": "Bantal, Kasur","mass":1000}`)
				request = httptest.NewRequest(http.MethodPatch, "/api/products/"+productCodeOne, requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				// Find All Products
				request = httptest.NewRequest(http.MethodGet, "/api/products-admin", nil)
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

				products := responseBody["data"].([]interface{})
				Expect(products[1].(map[string]interface{})["code"]).To(Equal("SK6"))
				Expect(products[1].(map[string]interface{})["name"]).To(Equal("Guling Doti Bang"))

				Expect(products[0].(map[string]interface{})["code"]).To(Equal("SK1"))
				Expect(products[0].(map[string]interface{})["name"]).To(Equal("Bantal"))
			})
		})
	})

	Describe("Find All Product By Similar Category /products/:code/category", func() {
		When("the product similar is not present", func() {
			It("should return a successful but the data is null", func() {
				// Create Product One
				requestBody := strings.NewReader(`{"code": "SK6","name": "Guling","description": "Test","category": "Bantal, Kasur","mass": 1000,"price": 7000}`)
				request := httptest.NewRequest(http.MethodPost, "/api/products", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBodyCreateProductOne map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBodyCreateProductOne)
				productCodeOne := responseBodyCreateProductOne["data"].(map[string]interface{})["code"].(string)

				// Create Product Two
				requestBody = strings.NewReader(`{"code": "SK1","name": "Bantal","description": "Test Bang","category": "Elektronik, Guling","mass": 1000,"price": 7000}`)
				request = httptest.NewRequest(http.MethodPost, "/api/products", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBodyCreateProductTwo map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBodyCreateProductTwo)

				// Find All Products
				request = httptest.NewRequest(http.MethodGet, "/api/products/"+productCodeOne+"/category", nil)
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
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("the product is present", func() {
			It("should return a successful and show all products by similar category", func() {
				// Create Product One
				requestBody := strings.NewReader(`{"code": "SK6","name": "Guling","description": "Test","category": "Bantal, Kasur","mass": 1000,"price": 7000}`)
				request := httptest.NewRequest(http.MethodPost, "/api/products", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBodyCreateProductOne map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBodyCreateProductOne)
				productCodeOne := responseBodyCreateProductOne["data"].(map[string]interface{})["code"].(string)

				// Create Product Two
				requestBody = strings.NewReader(`{"code": "SK1","name": "Bantal","description": "Test Bang","category": "Bantal, Kasur","mass": 1000,"price": 7000}`)
				request = httptest.NewRequest(http.MethodPost, "/api/products", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBodyCreateProductTwo map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBodyCreateProductTwo)

				// Find All Products
				request = httptest.NewRequest(http.MethodGet, "/api/products/"+productCodeOne+"/category", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				products := responseBody["data"].([]interface{})
				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))
				Expect(products[0].(map[string]interface{})["code"]).To(Equal("SK1"))
				Expect(products[0].(map[string]interface{})["name"]).To(Equal("Bantal"))
				Expect(products[0].(map[string]interface{})["description"]).To(Equal("Test Bang"))
				Expect(products[0].(map[string]interface{})["category"]).To(Equal("Bantal, Kasur"))
			})
		})
	})

	Describe("Find By Code Product /products/:code", func() {
		When("product is not found", func() {
			It("should return error not found", func() {
				// Find By Code Product
				request := httptest.NewRequest(http.MethodGet, "/api/products/SK5", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusNotFound))
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("product is found", func() {
			It("should return a successful find product by code", func() {
				// Create Product
				requestBody := strings.NewReader(`{"code": "SK6","name": "Guling","description": "Test","category": "Bantal, Kasur","mass": 1000,"price": 7000}`)
				request := httptest.NewRequest(http.MethodPost, "/api/products", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBodyCreateProduct map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBodyCreateProduct)
				productCode := responseBodyCreateProduct["data"].(map[string]interface{})["code"].(string)

				// Find By Code Product
				request = httptest.NewRequest(http.MethodGet, "/api/products/"+productCode, nil)
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
				Expect(responseBody["data"].(map[string]interface{})["code"]).To(Equal("SK6"))
				Expect(responseBody["data"].(map[string]interface{})["name"]).To(Equal("Guling"))
				Expect(responseBody["data"].(map[string]interface{})["description"]).To(Equal("Test"))
			})
		})
	})

	Describe("Find All Recommendation Product By Code /products/:code/recommendation", func() {
		When("product recommendation is not found", func() {
			It("should return error not found", func() {
				// Create Product
				requestBody := strings.NewReader(`{"code": "SK6","name": "Guling Biasa","description": "Test","category": "Bantal, Kasur","mass": 1000,"price": 7000}`)
				request := httptest.NewRequest(http.MethodPost, "/api/products", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBodyCreateProduct map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBodyCreateProduct)
				productCode := responseBodyCreateProduct["data"].(map[string]interface{})["code"].(string)

				// Create Apriori
				var aprioriRequests []*pb.CreateAprioriRequest_Create
				aprioriRequests = append(aprioriRequests, &pb.CreateAprioriRequest_Create{
					Item:       "Guling Biasa",
					Discount:   25.00,
					Support:    50.00,
					Confidence: 71.43,
					RangeDate:  "2021-05-21 - 2022-05-21",
					IsActive:   true,
				})
				aprioriService.Create(context.Background(), &pb.CreateAprioriRequest{
					CreateAprioriRequest: aprioriRequests,
				})

				// Find All Recommendation
				request = httptest.NewRequest(http.MethodGet, "/api/products/"+productCode+"/recommendation", nil)
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
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("product recommendation is found", func() {
			It("should return a successful find recommendation product by code", func() {
				// Create Product
				requestBody := strings.NewReader(`{"code": "SK1","name": "Bantal Biasa","description": "Test Bang","category": "Bantal, Kasur","mass": 1000,"price": 7000}`)
				request := httptest.NewRequest(http.MethodPost, "/api/products", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBodyCreateProduct map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBodyCreateProduct)
				productCode := responseBodyCreateProduct["data"].(map[string]interface{})["code"].(string)

				// Create Apriori
				var aprioriRequests []*pb.CreateAprioriRequest_Create
				aprioriRequests = append(aprioriRequests, &pb.CreateAprioriRequest_Create{
					Item:       "bantal biasa",
					Discount:   25.00,
					Support:    50.00,
					Confidence: 71.43,
					RangeDate:  "2021-05-21 - 2022-05-21",
					IsActive:   true,
				})
				aprioriService.Create(context.Background(), &pb.CreateAprioriRequest{
					CreateAprioriRequest: aprioriRequests,
				})

				// Find All Recommendation
				request = httptest.NewRequest(http.MethodGet, "/api/products/"+productCode+"/recommendation", nil)
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

				products := responseBody["data"].([]interface{})
				Expect(products[0].(map[string]interface{})["apriori_item"]).To(Equal("bantal biasa"))
			})
		})
	})
})
