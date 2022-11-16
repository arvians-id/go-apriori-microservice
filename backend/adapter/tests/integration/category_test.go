package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/arvians-id/go-apriori-microservice/adapter/config"
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

var _ = Describe("Category API", func() {
	var server *gin.Engine
	var tokenJWT string
	var cookie *http.Cookie
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

	Describe("Create Category /categories", func() {
		When("the fields are filled", func() {
			It("should return successful create category response", func() {
				// Create Category
				requestBody := strings.NewReader(`{"name": "Produk Kasur"}`)
				request := httptest.NewRequest(http.MethodPost, "/api/categories", requestBody)
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

				Expect(responseBody["data"].(map[string]interface{})["name"]).To(Equal("Produk Kasur"))
			})
		})
	})

	Describe("Find All Category /categories", func() {
		When("the categories is not found", func() {
			It("should return successful find all category response with empty data", func() {
				// Find All Category
				request := httptest.NewRequest(http.MethodGet, "/api/categories", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("the categories is exist", func() {
			It("should return successful find all category response", func() {
				// Create Category
				requestBody := strings.NewReader(`{"name": "Produk Kasur"}`)
				request := httptest.NewRequest(http.MethodPost, "/api/categories", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				// Find All Category
				request = httptest.NewRequest(http.MethodGet, "/api/categories", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))

				categoryResponse := responseBody["data"].([]interface{})
				Expect(categoryResponse[0].(map[string]interface{})["name"]).To(Equal("Produk Kasur"))
			})
		})
	})

	Describe("Find Category By Id /categories/:id", func() {
		When("the categories is not found", func() {
			It("should return error not found", func() {
				// Find Category By Id
				request := httptest.NewRequest(http.MethodGet, "/api/categories/1", nil)
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

		When("the categories is exist", func() {
			It("should return successful find category by id response", func() {
				// Create Category
				requestBody := strings.NewReader(`{"name": "Produk Kasur"}`)
				request := httptest.NewRequest(http.MethodPost, "/api/categories", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBodyCreate map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBodyCreate)
				idCategory := int64(responseBodyCreate["data"].(map[string]interface{})["id_category"].(float64))

				// Find Category By Id
				request = httptest.NewRequest(http.MethodGet, "/api/categories/"+util.IntToStr(int(idCategory)), nil)
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
				Expect(responseBody["data"].(map[string]interface{})["name"]).To(Equal("Produk Kasur"))
			})
		})
	})

	Describe("Update Category By Id /categories/:id", func() {
		When("the categories is not found", func() {
			It("should return error not found", func() {
				// Update Category By Id
				requestBody := strings.NewReader(`{"name": "Produk Bantal"}`)
				request := httptest.NewRequest(http.MethodPatch, "/api/categories/1", requestBody)
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

		When("the categories is exist", func() {
			It("should return successful update category response", func() {
				// Create Category
				requestBody := strings.NewReader(`{"name": "Produk Kasur"}`)
				request := httptest.NewRequest(http.MethodPost, "/api/categories", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBodyCreate map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBodyCreate)
				idCategory := int64(responseBodyCreate["data"].(map[string]interface{})["id_category"].(float64))

				// Update Category By Id
				requestBody = strings.NewReader(`{"name": "Produk Bantal"}`)
				request = httptest.NewRequest(http.MethodPatch, "/api/categories/"+util.IntToStr(int(idCategory)), requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				// Find Category By Id
				request = httptest.NewRequest(http.MethodGet, "/api/categories/"+util.IntToStr(int(idCategory)), nil)
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
				Expect(responseBody["data"].(map[string]interface{})["name"]).To(Equal("Produk Bantal"))
			})
		})
	})

	Describe("Delete Category By Id /categories/:id", func() {
		When("the categories is not found", func() {
			It("should return error not found", func() {
				// Delete Category By Id
				request := httptest.NewRequest(http.MethodDelete, "/api/categories/1", nil)
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

		When("the categories is exist", func() {
			It("should return successful delete category by id and set not found after deleted", func() {
				// Create Category
				requestBody := strings.NewReader(`{"name": "Produk Kasur"}`)
				request := httptest.NewRequest(http.MethodPost, "/api/categories", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBodyCreate map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBodyCreate)
				idCategory := int64(responseBodyCreate["data"].(map[string]interface{})["id_category"].(float64))

				// Delete Category By Id
				request = httptest.NewRequest(http.MethodDelete, "/api/categories/"+util.IntToStr(int(idCategory)), nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				// Find Category By Id
				request = httptest.NewRequest(http.MethodGet, "/api/categories/"+util.IntToStr(int(idCategory)), nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusNotFound))
				Expect(responseBody["status"]).To(Equal("data not found"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})
	})
})
