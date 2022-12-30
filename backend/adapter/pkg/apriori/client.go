package apriori

import (
	"github.com/arvians-id/go-apriori-microservice/adapter/config"
	"github.com/arvians-id/go-apriori-microservice/adapter/middleware"
	"github.com/arvians-id/go-apriori-microservice/adapter/model"
	"github.com/arvians-id/go-apriori-microservice/adapter/pb"
	"github.com/arvians-id/go-apriori-microservice/adapter/response"
	"github.com/arvians-id/go-apriori-microservice/adapter/third-party/aws"
	"github.com/arvians-id/go-apriori-microservice/adapter/util"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"strconv"
	"strings"
)

type ServiceClient struct {
	AprioriService pb.AprioriServiceClient
	StorageS3      *aws.StorageS3
}

func NewAprioriServiceClient(configuration *config.Config) pb.AprioriServiceClient {
	connection, err := grpc.Dial(configuration.AprioriSvcUrl, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}

	return pb.NewAprioriServiceClient(connection)
}

func RegisterRoutes(router *gin.Engine, configuration *config.Config, storageS3 *aws.StorageS3) *ServiceClient {
	serviceClient := &ServiceClient{
		AprioriService: NewAprioriServiceClient(configuration),
		StorageS3:      storageS3,
	}

	authorized := router.Group("/api", middleware.AuthJwtMiddleware(configuration))
	{
		authorized.PATCH("/apriori/:code", serviceClient.UpdateStatus)
		authorized.POST("/apriori", serviceClient.Create)
		authorized.DELETE("/apriori/:code", serviceClient.Delete)
		authorized.PATCH("/apriori/:code/update/:id", serviceClient.Update)
		authorized.POST("/apriori/generate", serviceClient.Generate)
	}

	unauthorized := router.Group("/api")
	{
		unauthorized.GET("/apriori", serviceClient.FindAll)
		unauthorized.GET("/apriori/:code", serviceClient.FindAllByCode)
		unauthorized.GET("/apriori/:code/detail/:id", serviceClient.FindByCodeAndId)
		unauthorized.GET("/apriori/actives", serviceClient.FindAllByActive)
	}

	return serviceClient
}

func (client *ServiceClient) FindAll(c *gin.Context) {
	apriories, err := client.AprioriService.FindAll(c.Request.Context(), new(emptypb.Empty))
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", apriories.GetApriori())
}

func (client *ServiceClient) FindAllByActive(c *gin.Context) {
	apriories, err := client.AprioriService.FindAllByActive(c.Request.Context(), new(emptypb.Empty))
	if err != nil {
		if err.Error() == util.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", apriories.GetApriori())
}

func (client *ServiceClient) FindAllByCode(c *gin.Context) {
	codeParam := c.Param("code")
	apriories, err := client.AprioriService.FindAllByCode(c.Request.Context(), &pb.GetAprioriByCodeRequest{
		Code: codeParam,
	})
	if err != nil {
		if err.Error() == util.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", apriories.GetApriori())
}

func (client *ServiceClient) FindByCodeAndId(c *gin.Context) {
	codeParam := c.Param("code")
	idParam, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	apriori, err := client.AprioriService.FindByCodeAndId(c.Request.Context(), &pb.GetAprioriByCodeAndIdRequest{
		Code: codeParam,
		Id:   idParam,
	})
	if err != nil {
		if err.Error() == util.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", apriori.GetProductRecommendation())
}

func (client *ServiceClient) Update(c *gin.Context) {
	codeParam := c.Param("code")
	idParam, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	description := c.PostForm("description")
	file, header, err := c.Request.FormFile("image")
	var filePath string
	if err == nil {
		path, fileName := client.StorageS3.GenerateNewFile(header.Filename)
		go func() {
			err = client.StorageS3.UploadToAWS(file, fileName, header.Header.Get("Content-Type"))
			if err != nil {
				log.Println("[Apriori][Update][UploadToAWS] error upload file to S3, err: ", err.Error())
			}
		}()
		filePath = path
	}

	apriories, err := client.AprioriService.Update(c.Request.Context(), &pb.UpdateAprioriRequest{
		IdApriori:   idParam,
		Code:        codeParam,
		Description: description,
		Image:       filePath,
	})
	if err != nil {
		if err.Error() == util.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", apriories.GetApriori())
}

func (client *ServiceClient) UpdateStatus(c *gin.Context) {
	codeParam := c.Param("code")
	_, err := client.AprioriService.UpdateStatus(c.Request.Context(), &pb.GetAprioriByCodeRequest{
		Code: codeParam,
	})
	if err != nil {
		if err.Error() == util.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", nil)
}

func (client *ServiceClient) Create(c *gin.Context) {
	var generateRequests []*model.GenerateApriori
	err := c.ShouldBindJSON(&generateRequests)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	var aprioriRequests []*pb.CreateAprioriRequest_Create
	for _, generateRequest := range generateRequests {
		ItemSet := strings.Join(generateRequest.ItemSet, ", ")
		aprioriRequests = append(aprioriRequests, &pb.CreateAprioriRequest_Create{
			Item:       ItemSet,
			Discount:   generateRequest.Discount,
			Support:    generateRequest.Support,
			Confidence: generateRequest.Confidence,
			RangeDate:  generateRequest.RangeDate,
		})
	}

	_, err = client.AprioriService.Create(c.Request.Context(), &pb.CreateAprioriRequest{
		CreateAprioriRequest: aprioriRequests,
	})
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", nil)
}
func (client *ServiceClient) Delete(c *gin.Context) {
	codeParam := c.Param("code")
	_, err := client.AprioriService.Delete(c.Request.Context(), &pb.GetAprioriByCodeRequest{
		Code: codeParam,
	})
	if err != nil {
		if err.Error() == util.ErrorNotFound {
			response.ReturnErrorNotFound(c, err, nil)
			return
		}

		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", nil)

}
func (client *ServiceClient) Generate(c *gin.Context) {
	var requestGenerate GenerateAprioriRequest
	err := c.ShouldBindJSON(&requestGenerate)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	apriori, err := client.AprioriService.Generate(c.Request.Context(), &pb.GenerateAprioriRequest{
		MinimumSupport:    requestGenerate.MinimumSupport,
		MinimumConfidence: requestGenerate.MinimumConfidence,
		MinimumDiscount:   requestGenerate.MinimumDiscount,
		MaximumDiscount:   requestGenerate.MaximumDiscount,
		StartDate:         requestGenerate.StartDate,
		EndDate:           requestGenerate.EndDate,
	})
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", apriori.GetGenerateApriori())
}
