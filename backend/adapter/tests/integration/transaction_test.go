package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/arvians-id/go-apriori-microservice/adapter/config"
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

/*
	Error :
		- Create Transactions By CSV File /transactions/csv
*/
var _ = Describe("Transaction API", func() {
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

	Describe("Create Transaction /transactions", func() {
		When("the fields are incorrect", func() {
			When("the product name field is incorrect", func() {
				It("should return error required", func() {
					// Create Transaction
					requestBody := strings.NewReader(`{"customer_name": "Wids"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/transactions", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.XApiKey)
					request.AddCookie(cookie)
					request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					var responseBody map[string]interface{}
					_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'CreateTransactionRequest.ProductName' Error:Field validation for 'ProductName' failed on the 'required' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})

			When("the customer name field is incorrect", func() {
				It("should return error required", func() {
					// Create Transaction
					requestBody := strings.NewReader(`{"product_name": "Kasur cinta, Bantal memori"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/transactions", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.XApiKey)
					request.AddCookie(cookie)
					request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					var responseBody map[string]interface{}
					_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'CreateTransactionRequest.CustomerName' Error:Field validation for 'CustomerName' failed on the 'required' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})
		})

		When("the fields are correct", func() {
			It("should return successful create transaction response", func() {
				// Create Transaction
				requestBody := strings.NewReader(`{"product_name": "Kasur cinta, Bantal memori","customer_name": "Wids"}`)
				request := httptest.NewRequest(http.MethodPost, "/api/transactions", requestBody)
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
				Expect(responseBody["data"].(map[string]interface{})["product_name"]).To(Equal("kasur cinta, bantal memori"))
				Expect(responseBody["data"].(map[string]interface{})["customer_name"]).To(Equal("Wids"))
			})
		})
	})

	//Describe("Create Transactions By CSV File /transactions/csv", func() {
	//	When("file exist", func() {
	//		It("should return error no such file", func() {
	//			path := "./assets/example1.csv"
	//			body := new(bytes.Buffer)
	//			writer := multipart.NewWriter(body)
	//			part, _ := writer.CreateFormFile("file", path)
	//			sample, _ := os.Open(path)
	//
	//			_, _ = io.Copy(part, sample)
	//			writer.Close()
	//
	//			// Create Transaction
	//			request := httptest.NewRequest(http.MethodPost, "/api/transactions/csv", body)
	//			request.Header.Add("Content-Type", writer.FormDataContentType())
	//			request.AddCookie(cookie)
	//			request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))
	//
	//			rec := httptest.NewRecorder()
	//			server.ServeHTTP(rec, request)
	//
	//			var responseBody map[string]interface{}
	//			_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)
	//
	//			Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
	//			Expect(responseBody["status"]).To(Equal("created"))
	//			Expect(responseBody["data"]).To(BeNil())
	//		})
	//	})
	//})

	Describe("Update Transaction /transactions/:number_transaction", func() {
		When("the fields are incorrect", func() {
			When("the product name field is incorrect", func() {
				It("should return error required", func() {
					// Create Transaction
					requestBody := strings.NewReader(`{"product_name": "kasur cinta, bantal memori", "customer_name": "Wids", "no_transaction": "202320"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/transactions", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.XApiKey)
					request.AddCookie(cookie)
					request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					var responseBodyCreateTransaction map[string]interface{}
					_ = json.NewDecoder(writer.Result().Body).Decode(&responseBodyCreateTransaction)
					noTransaction := responseBodyCreateTransaction["data"].(map[string]interface{})["no_transaction"].(string)

					// Update Transaction
					requestBody = strings.NewReader(`{"customer_name": "Wids"}`)
					request = httptest.NewRequest(http.MethodPatch, "/api/transactions/"+noTransaction, requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.XApiKey)
					request.AddCookie(cookie)
					request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

					writer = httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					var responseBody map[string]interface{}
					_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'UpdateTransactionRequest.ProductName' Error:Field validation for 'ProductName' failed on the 'required' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})

			When("the customer name field is incorrect", func() {
				It("should return error required", func() {
					// Create Transaction
					requestBody := strings.NewReader(`{"product_name": "kasur cinta, bantal memori", "customer_name": "Wids", "no_transaction": "202320"}`)
					request := httptest.NewRequest(http.MethodPost, "/api/transactions", requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.XApiKey)
					request.AddCookie(cookie)
					request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

					writer := httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					var responseBodyCreateTransaction map[string]interface{}
					_ = json.NewDecoder(writer.Result().Body).Decode(&responseBodyCreateTransaction)
					noTransaction := responseBodyCreateTransaction["data"].(map[string]interface{})["no_transaction"].(string)

					// Update Transaction
					requestBody = strings.NewReader(`{"product_name": "Kasur cinta, Bantal memori"}`)
					request = httptest.NewRequest(http.MethodPatch, "/api/transactions/"+noTransaction, requestBody)
					request.Header.Add("Content-Type", "application/json")
					request.Header.Add("X-API-KEY", configuration.XApiKey)
					request.AddCookie(cookie)
					request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

					writer = httptest.NewRecorder()
					server.ServeHTTP(writer, request)

					var responseBody map[string]interface{}
					_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("Key: 'UpdateTransactionRequest.CustomerName' Error:Field validation for 'CustomerName' failed on the 'required' tag"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})
		})

		When("the fields are correct", func() {
			It("should return successful update transaction response", func() {
				// Create Transaction
				requestBody := strings.NewReader(`{"product_name": "kasur cinta, bantal memori", "customer_name": "Wids", "no_transaction": "202320"}`)
				request := httptest.NewRequest(http.MethodPost, "/api/transactions", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBodyCreateTransaction map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBodyCreateTransaction)
				noTransaction := responseBodyCreateTransaction["data"].(map[string]interface{})["no_transaction"].(string)

				// Update Transaction
				requestBody = strings.NewReader(`{"product_name": "Guling cinta, Guling memori","customer_name": "Goengs"}`)
				request = httptest.NewRequest(http.MethodPatch, "/api/transactions/"+noTransaction, requestBody)
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
				Expect(responseBody["data"].(map[string]interface{})["product_name"]).To(Equal("guling cinta, guling memori"))
				Expect(responseBody["data"].(map[string]interface{})["product_name"]).ShouldNot(Equal("kasur cinta, bantal memori"))
				Expect(responseBody["data"].(map[string]interface{})["customer_name"]).ShouldNot(Equal("Wids"))
			})
		})
	})

	Describe("Delete Transaction /transactions/:number_transaction", func() {
		When("transaction is not found", func() {
			It("should return error not found", func() {
				// Delete Transaction
				request := httptest.NewRequest(http.MethodDelete, "/api/transactions/32412", nil)
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

		When("transaction is found", func() {
			It("should return a successful delete transaction response", func() {
				// Create Transaction
				requestBody := strings.NewReader(`{"product_name": "kasur cinta, bantal memori", "customer_name": "Wids", "no_transaction": "202320"}`)
				request := httptest.NewRequest(http.MethodPost, "/api/transactions", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBodyCreateTransaction map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBodyCreateTransaction)
				noTransaction := responseBodyCreateTransaction["data"].(map[string]interface{})["no_transaction"].(string)

				// Delete Transaction
				request = httptest.NewRequest(http.MethodDelete, "/api/transactions/"+noTransaction, nil)
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

	Describe("Find All Transaction /transactions", func() {
		When("the transaction is not present", func() {
			It("should return a successful but the data is null", func() {
				// Find All Transaction
				request := httptest.NewRequest(http.MethodGet, "/api/transactions", nil)
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

		When("the transactions is present", func() {
			It("should return a successful and show all transactions", func() {
				// Create Transaction
				requestBody := strings.NewReader(`{"product_name": "kasur cinta, bantal memori", "customer_name": "Wids", "no_transaction": "202320"}`)
				request := httptest.NewRequest(http.MethodPost, "/api/transactions", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				// Create Transaction
				requestBody = strings.NewReader(`{"product_name": "guling cinta, guling memori", "customer_name": "Goengs", "no_transaction": "202320"}`)
				request = httptest.NewRequest(http.MethodPost, "/api/transactions", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBodyCreateTransaction map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBodyCreateTransaction)

				// Find All Transaction
				request = httptest.NewRequest(http.MethodGet, "/api/transactions", nil)
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

				transactions := responseBody["data"].([]interface{})
				Expect(transactions[1].(map[string]interface{})["product_name"]).To(Equal("kasur cinta, bantal memori"))
				Expect(transactions[1].(map[string]interface{})["customer_name"]).To(Equal("Wids"))

				Expect(transactions[0].(map[string]interface{})["product_name"]).To(Equal("guling cinta, guling memori"))
				Expect(transactions[0].(map[string]interface{})["customer_name"]).To(Equal("Goengs"))
			})
		})
	})

	Describe("Find By No Transaction /transactions/:number_transaction", func() {
		When("transaction is not found", func() {
			It("should return error not found", func() {
				// Find By No Transaction
				request := httptest.NewRequest(http.MethodGet, "/api/transactions/52324", nil)
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

		When("transaction is found", func() {
			It("should return a successful find transaction by no transaction", func() {
				// Create Transaction
				requestBody := strings.NewReader(`{"product_name": "kasur cinta, bantal memori", "customer_name": "Wids", "no_transaction": "202320"}`)
				request := httptest.NewRequest(http.MethodPost, "/api/transactions", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBodyCreateTransaction map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBodyCreateTransaction)
				noTransaction := responseBodyCreateTransaction["data"].(map[string]interface{})["no_transaction"].(string)

				// Find By No Transaction
				request = httptest.NewRequest(http.MethodGet, "/api/transactions/"+noTransaction, nil)
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
				Expect(responseBody["data"].(map[string]interface{})["product_name"]).To(Equal("kasur cinta, bantal memori"))
				Expect(responseBody["data"].(map[string]interface{})["customer_name"]).To(Equal("Wids"))
			})
		})
	})

	Describe("Truncate Transaction /transactions/truncate", func() {
		When("transaction is found", func() {
			It("should return successful delete all transactions", func() {
				// Create Transaction
				requestBody := strings.NewReader(`{"product_name": "kasur cinta, bantal memori", "customer_name": "Wids", "no_transaction": "202320"}`)
				request := httptest.NewRequest(http.MethodPost, "/api/transactions", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBodyCreateTransaction map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBodyCreateTransaction)

				// Delete Transaction
				request = httptest.NewRequest(http.MethodDelete, "/api/transactions/truncate", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Add("X-API-KEY", configuration.XApiKey)
				request.AddCookie(cookie)
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenJWT))

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				var responseBody map[string]interface{}
				_ = json.NewDecoder(writer.Result().Body).Decode(&responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["data"]).To(BeNil())
			})
		})
	})

	Describe("Access Transaction Endpoint", func() {
		When("the user is not logged in", func() {
			It("should return error unauthorized response", func() {
				request := httptest.NewRequest(http.MethodGet, "/api/transactions", nil)
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
	})
})
