package raja_ongkir

import (
	"encoding/json"
	"fmt"
	"github.com/arvians-id/go-apriori-microservice/adapter/response"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strings"
)

type ServiceClient struct {
}

func RegisterRoutes(router *gin.Engine) *ServiceClient {
	serviceClient := &ServiceClient{}

	router.GET("/api/raja-ongkir/:place", serviceClient.FindAll)
	router.POST("/api/raja-ongkir/cost", serviceClient.GetCost)

	return serviceClient
}

func (client *ServiceClient) FindAll(c *gin.Context) {
	placeParam := c.Param("place")
	if placeParam == "province" {
		placeParam = "province"
	} else if placeParam == "city" {
		placeParam = "city?province=" + c.Query("province")
	}

	url := "https://api.rajaongkir.com/starter/" + placeParam
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("[RajaOngkirController][FindAll] problem on request to raja ongkir, err: ", err.Error())
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	req.Header.Add("key", os.Getenv("RAJA_ONGKIR_SECRET_KEY"))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("[RajaOngkirController][FindAll] problem on send http request to raja ongkir, err: ", err.Error())
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	defer res.Body.Close()

	var rajaOngkirModel interface{}
	err = json.NewDecoder(res.Body).Decode(&rajaOngkirModel)
	if err != nil {
		log.Println("[RajaOngkirController][FindAll] unable to decode data err: ", err.Error())
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": rajaOngkirModel,
	})
}

func (client *ServiceClient) GetCost(c *gin.Context) {
	var requestDelivery GetDeliveryRequest
	err := c.ShouldBind(&requestDelivery)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	payload := fmt.Sprintf(
		"origin=%v&destination=%v&weight=%v&courier=%v",
		requestDelivery.Origin,
		requestDelivery.Destination,
		requestDelivery.Weight,
		requestDelivery.Courier,
	)
	data := strings.NewReader(payload)
	req, err := http.NewRequest("POST", "https://api.rajaongkir.com/starter/cost", data)
	if err != nil {
		log.Println("[RajaOngkirController][GetCost] problem on request to raja ongkir, err: ", err.Error())
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	req.Header.Add("key", os.Getenv("RAJA_ONGKIR_SECRET_KEY"))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("[RajaOngkirController][GetCost] problem on send http request to raja ongkir, err: ", err.Error())
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	defer res.Body.Close()

	var rajaOngkirModel interface{}
	err = json.NewDecoder(res.Body).Decode(&rajaOngkirModel)
	if err != nil {
		log.Println("[RajaOngkirController][GetCost] unable to decode data err: ", err.Error())
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": rajaOngkirModel,
	})
}
