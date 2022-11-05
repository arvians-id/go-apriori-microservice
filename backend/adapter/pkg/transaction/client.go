package transaction

import (
	"github.com/arvians-id/go-apriori-microservice/adapter/middleware"
	"github.com/arvians-id/go-apriori-microservice/adapter/pkg/transaction/pb"
	"github.com/arvians-id/go-apriori-microservice/adapter/response"
	"github.com/arvians-id/go-apriori-microservice/config"
	"github.com/arvians-id/go-apriori-microservice/third-party/aws"
	"github.com/arvians-id/go-apriori-microservice/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"log"
)

type ServiceClient struct {
	TransactionService pb.TransactionServiceClient
	StorageS3          aws.StorageS3
}

func NewTransactionServiceClient(configuration *config.Config) pb.TransactionServiceClient {
	connection, err := grpc.Dial(configuration.TransactionSvcUrl, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}

	return pb.NewTransactionServiceClient(connection)
}

func RegisterRoutes(router *gin.Engine, configuration *config.Config, storageS3 aws.StorageS3) *ServiceClient {
	serviceClient := &ServiceClient{
		TransactionService: NewTransactionServiceClient(configuration),
		StorageS3:          storageS3,
	}

	authorized := router.Group("/api", middleware.AuthJwtMiddleware())
	{
		authorized.GET("/transactions", serviceClient.FindAll)
		authorized.GET("/transactions/:number_transaction", serviceClient.FindByNoTransaction)
		authorized.POST("/transactions", serviceClient.Create)
		authorized.POST("/transactions/csv", serviceClient.CreateByCSV)
		authorized.PATCH("/transactions/:number_transaction", serviceClient.Update)
		authorized.DELETE("/transactions/:number_transaction", serviceClient.Delete)
		authorized.DELETE("/transactions/truncate", serviceClient.Truncate)
	}

	return serviceClient
}

func (client *ServiceClient) FindAll(c *gin.Context) {
	transaction, err := client.TransactionService.FindAll(c.Request.Context(), new(empty.Empty))
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", transaction)
}

func (client *ServiceClient) FindByNoTransaction(c *gin.Context) {
	noTransactionParam := c.Param("number_transaction")
	transactions, err := client.TransactionService.FindByNoTransaction(c.Request.Context(), &pb.GetTransactionByNoTransactionRequest{
		NoTransaction: noTransactionParam,
	})
	if err != nil {
		if err.Error() == util.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", transactions)
}

func (client *ServiceClient) Create(c *gin.Context) {
	var requestCreate CreateTransactionRequest
	err := c.ShouldBindJSON(&requestCreate)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	transaction, err := client.TransactionService.Create(c.Request.Context(), &pb.CreateTransactionRequest{
		ProductName:  requestCreate.ProductName,
		CustomerName: requestCreate.CustomerName,
	})
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "created", transaction)
}

func (client *ServiceClient) CreateByCSV(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	filePath, fileName := client.StorageS3.GenerateNewFile(header.Filename)
	err = client.StorageS3.UploadToAWS(file, fileName, header.Header.Get("Content-Type"))
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	data, err := util.OpenCsvFile(filePath)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	_, err = client.TransactionService.CreateByCSV(c.Request.Context(), &pb.CreateTransactionByCSVRequest{
		Request: data,
	})
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	defer func(StorageS3 *aws.StorageS3, filePath string) {
		_ = StorageS3.DeleteFromAWS(filePath)
	}(&client.StorageS3, filePath)

	response.ReturnSuccessOK(c, "created", nil)
}

func (client *ServiceClient) Update(c *gin.Context) {
	var requestUpdate UpdateTransactionRequest
	err := c.ShouldBindJSON(&requestUpdate)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	noTransactionParam := c.Param("number_transaction")
	requestUpdate.NoTransaction = noTransactionParam
	transaction, err := client.TransactionService.Update(c.Request.Context(), &pb.UpdateTransactionRequest{
		ProductName:   requestUpdate.ProductName,
		CustomerName:  requestUpdate.CustomerName,
		NoTransaction: requestUpdate.NoTransaction,
	})
	if err != nil {
		if err.Error() == util.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "updated", transaction)
}

func (client *ServiceClient) Delete(c *gin.Context) {
	noTransactionParam := c.Param("number_transaction")
	_, err := client.TransactionService.Delete(c.Request.Context(), &pb.GetTransactionByNoTransactionRequest{
		NoTransaction: noTransactionParam,
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

func (client *ServiceClient) Truncate(c *gin.Context) {
	_, err := client.TransactionService.Truncate(c.Request.Context(), new(empty.Empty))
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "deleted", nil)
}
