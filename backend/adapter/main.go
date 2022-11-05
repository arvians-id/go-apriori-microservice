package main

import (
	"database/sql"
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
	"github.com/arvians-id/go-apriori-microservice/third-party/jwt"
	"github.com/gin-gonic/gin"
	"log"
)

func NewInitializedDatabase(configuration *config.Config) (*sql.DB, error) {
	db, err := config.NewPostgresSQL(configuration)
	if err != nil {
		return nil, err
	}

	return db, nil
}

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

	// Services
	userService := user.RegisterRoutes(router, configuration)
	apriori.RegisterRoutes(router, configuration)
	auth.RegisterRoutes(router, configuration, userService, jwtAuth)
	category.RegisterRoutes(router, configuration)
	comment.RegisterRoutes(router, configuration)
	notification.RegisterRoutes(router, configuration)
	payment.RegisterRoutes(router, configuration)
	product.RegisterRoutes(router, configuration)
	raja_ongkir.RegisterRoutes(router)
	transaction.RegisterRoutes(router, configuration)
	user_order.RegisterRoutes(router, configuration)

	err = router.Run(configuration.Port)
	if err != nil {
		log.Fatalln("Failed at running", err)
	}
}
