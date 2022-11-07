package main

import (
	"github.com/arvians-id/go-apriori-microservice/adapter/middleware"
	"github.com/arvians-id/go-apriori-microservice/adapter/pkg/apriori"
	"github.com/arvians-id/go-apriori-microservice/adapter/pkg/auth"
	"github.com/arvians-id/go-apriori-microservice/adapter/pkg/category"
	"github.com/arvians-id/go-apriori-microservice/adapter/pkg/comment"
	"github.com/arvians-id/go-apriori-microservice/adapter/pkg/notification"
	"github.com/arvians-id/go-apriori-microservice/adapter/pkg/payment"
	"github.com/arvians-id/go-apriori-microservice/adapter/pkg/product"
	raja_ongkir "github.com/arvians-id/go-apriori-microservice/adapter/pkg/raja-ongkir"
	"github.com/arvians-id/go-apriori-microservice/adapter/pkg/transaction"
	"github.com/arvians-id/go-apriori-microservice/adapter/pkg/user"
	user_order "github.com/arvians-id/go-apriori-microservice/adapter/pkg/user-order"
	"github.com/arvians-id/go-apriori-microservice/config"
	"github.com/arvians-id/go-apriori-microservice/third-party/aws"
	"github.com/arvians-id/go-apriori-microservice/third-party/jwt"
	messaging "github.com/arvians-id/go-apriori-microservice/third-party/message-queue"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	configuration, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	router := gin.Default()

	// Main Middlewares
	middleware.RegisterPrometheusMetrics()
	router.Use(middleware.SetupCorsMiddleware())
	router.Use(middleware.GinContextToContextMiddleware())
	router.Use(middleware.PrometheusMetricsMiddleware())

	// Third Party
	jwtAuth := jwt.NewJsonWebToken()
	storageS3 := aws.NewStorageS3(configuration)
	messagingProducer := messaging.NewProducer(messaging.ProducerConfig{
		NsqdAddress: "nsqd:4150",
	})

	// Services
	user.RegisterRoutes(router, configuration)
	apriori.RegisterRoutes(router, configuration)
	auth.RegisterRoutes(router, configuration, jwtAuth, messagingProducer)
	category.RegisterRoutes(router, configuration)
	comment.RegisterRoutes(router, configuration)
	notification.RegisterRoutes(router, configuration)
	payment.RegisterRoutes(router, configuration, messagingProducer)
	product.RegisterRoutes(router, configuration, storageS3)
	raja_ongkir.RegisterRoutes(router)
	transaction.RegisterRoutes(router, configuration, storageS3)
	user_order.RegisterRoutes(router, configuration)

	err = router.Run(configuration.Port)
	if err != nil {
		log.Fatalln("Failed at running", err)
	}
}
