package main

import (
	"github.com/arvians-id/go-apriori-microservice/adapter/config"
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
	"github.com/arvians-id/go-apriori-microservice/adapter/third-party/aws"
	"github.com/arvians-id/go-apriori-microservice/adapter/third-party/jwt"
	messaging "github.com/arvians-id/go-apriori-microservice/adapter/third-party/message-queue"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	jwtAuth := jwt.NewJsonWebToken(configuration)
	storageS3 := aws.NewStorageS3(configuration)
	messagingProducer := messaging.NewProducer(messaging.ProducerConfig{
		NsqdAddress: "nsq-release-nsqd:4150",
	})

	// Other routes
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Apriori Algorithm API. Created By https://github.com/arvians-id",
		})
	})

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Services
	payment.RegisterRoutes(router, configuration, messagingProducer)
	router.Use(middleware.SetupXApiKeyMiddleware(configuration))
	user.RegisterRoutes(router, configuration)
	apriori.RegisterRoutes(router, configuration, storageS3)
	auth.RegisterRoutes(router, configuration, jwtAuth, messagingProducer)
	category.RegisterRoutes(router, configuration)
	comment.RegisterRoutes(router, configuration)
	notification.RegisterRoutes(router, configuration)
	product.RegisterRoutes(router, configuration, storageS3)
	raja_ongkir.RegisterRoutes(router, configuration)
	transaction.RegisterRoutes(router, configuration, storageS3)
	user_order.RegisterRoutes(router, configuration)

	err = router.Run(configuration.Port)
	if err != nil {
		log.Fatalln("Failed at running", err)
	}
}
