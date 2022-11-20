package setup

import (
	"database/sql"
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
)

func ModuleSetup(configuration *config.Config) (*gin.Engine, *sql.DB) {
	db, _ := config.NewPostgresSQL(configuration)
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
		NsqdAddress: "nsqd:4150",
	})

	// Services
	user.RegisterRoutes(router, configuration)
	apriori.RegisterRoutes(router, configuration, storageS3)
	auth.RegisterRoutes(router, configuration, jwtAuth, messagingProducer)
	category.RegisterRoutes(router, configuration)
	comment.RegisterRoutes(router, configuration)
	notification.RegisterRoutes(router, configuration)
	payment.RegisterRoutes(router, configuration, messagingProducer)
	product.RegisterRoutes(router, configuration, storageS3)
	raja_ongkir.RegisterRoutes(router, configuration)
	transaction.RegisterRoutes(router, configuration, storageS3)
	user_order.RegisterRoutes(router, configuration)

	return router, db
}
