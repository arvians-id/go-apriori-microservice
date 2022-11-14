package product

import (
	"github.com/arvians-id/go-apriori-microservice/adapter/middleware"
	"github.com/arvians-id/go-apriori-microservice/adapter/pb"
	"github.com/arvians-id/go-apriori-microservice/adapter/response"
	"github.com/arvians-id/go-apriori-microservice/config"
	"github.com/arvians-id/go-apriori-microservice/third-party/aws"
	"github.com/arvians-id/go-apriori-microservice/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"log"
	"strings"
)

type ServiceClient struct {
	ProductService pb.ProductServiceClient
	StorageS3      *aws.StorageS3
}

func NewCommentServiceClient(configuration *config.Config) pb.ProductServiceClient {
	connection, err := grpc.Dial(configuration.ProductSvcUrl, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}

	return pb.NewProductServiceClient(connection)
}

func RegisterRoutes(router *gin.Engine, configuration *config.Config, storageS3 *aws.StorageS3) *ServiceClient {
	serviceClient := &ServiceClient{
		ProductService: NewCommentServiceClient(configuration),
		StorageS3:      storageS3,
	}

	authorized := router.Group("/api", middleware.AuthJwtMiddleware(configuration))
	{
		authorized.GET("/products-admin", serviceClient.FindAllByAdmin)
		authorized.POST("/products", serviceClient.Create)
		authorized.PATCH("/products/:code", serviceClient.Update)
		authorized.DELETE("/products/:code", serviceClient.Delete)
	}

	unauthorized := router.Group("/api")
	{
		unauthorized.GET("/products", serviceClient.FindAllByUser)
		unauthorized.GET("/products/:code/category", serviceClient.FindAllSimilarCategory)
		unauthorized.GET("/products/:code/recommendation", serviceClient.FindAllRecommendation)
		unauthorized.GET("/products/:code", serviceClient.FindByCode)
	}

	return serviceClient
}

func (client *ServiceClient) FindAllByAdmin(c *gin.Context) {
	products, err := client.ProductService.FindAllByAdmin(c.Request.Context(), new(empty.Empty))
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", products.GetProduct())
}

func (client *ServiceClient) FindAllSimilarCategory(c *gin.Context) {
	codeParam := c.Param("code")
	products, err := client.ProductService.FindAllBySimilarCategory(c.Request.Context(), &pb.GetProductByProductCodeRequest{
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

	response.ReturnSuccessOK(c, "OK", products.GetProduct())
}

func (client *ServiceClient) FindAllByUser(c *gin.Context) {
	searchQuery := strings.ToLower(c.Query("search"))
	categoryQuery := strings.ToLower(c.Query("category"))
	products, err := client.ProductService.FindAll(c.Request.Context(), &pb.GetProductByFiltersRequest{
		Search:   searchQuery,
		Category: categoryQuery,
	})
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "OK", products.GetProduct())
}

func (client *ServiceClient) FindAllRecommendation(c *gin.Context) {
	codeParam := c.Param("code")
	products, err := client.ProductService.FindAllRecommendation(c.Request.Context(), &pb.GetProductByProductCodeRequest{
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

	response.ReturnSuccessOK(c, "OK", products.GetProductRecommendation())
}

func (client *ServiceClient) FindByCode(c *gin.Context) {
	codeParam := c.Param("code")
	product, err := client.ProductService.FindByCode(c.Request.Context(), &pb.GetProductByProductCodeRequest{
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

	response.ReturnSuccessOK(c, "OK", product.GetProduct())
}

func (client *ServiceClient) Create(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	filePath := client.StorageS3.DefaultPath()
	if err == nil {
		path, fileName := client.StorageS3.GenerateNewFile(header.Filename)
		go func() {
			err = client.StorageS3.UploadToAWS(file, fileName, header.Header.Get("Content-Type"))
			if err != nil {
				log.Println("[Product][Create][UploadFileS3Test] error upload file S3, err: ", err.Error())
			}
		}()
		filePath = path
	}

	var requestCreate CreateProductRequest
	err = c.ShouldBind(&requestCreate)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	product, err := client.ProductService.Create(c.Request.Context(), &pb.CreateProductRequest{
		Code:        requestCreate.Code,
		Name:        requestCreate.Name,
		Description: requestCreate.Description,
		Price:       requestCreate.Price,
		Category:    requestCreate.Category,
		Mass:        requestCreate.Mass,
		Image:       filePath,
	})
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	response.ReturnSuccessOK(c, "created", product.GetProduct())
}

func (client *ServiceClient) Update(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	var filePath string
	if err == nil {
		path, fileName := client.StorageS3.GenerateNewFile(header.Filename)
		go func() {
			err = client.StorageS3.UploadToAWS(file, fileName, header.Header.Get("Content-Type"))
			if err != nil {
				log.Println("[Product][Create][UploadFileS3Test] error upload file S3, err: ", err.Error())
			}
		}()
		filePath = path
	}

	var requestUpdate UpdateProductRequest
	err = c.ShouldBind(&requestUpdate)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	product, err := client.ProductService.Update(c.Request.Context(), &pb.UpdateProductRequest{
		Code:        c.Param("code"),
		Name:        requestUpdate.Name,
		Description: requestUpdate.Description,
		Price:       requestUpdate.Price,
		Category:    requestUpdate.Category,
		IsEmpty:     requestUpdate.IsEmpty,
		Mass:        requestUpdate.Mass,
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

	response.ReturnSuccessOK(c, "updated", product.GetProduct())
}

func (client *ServiceClient) Delete(c *gin.Context) {
	codeParam := c.Param("code")
	_, err := client.ProductService.Delete(c.Request.Context(), &pb.GetProductByProductCodeRequest{
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

	response.ReturnSuccessOK(c, "deleted", nil)
}
