package main

import (
	"database/sql"
	"github.com/arvians-id/go-apriori-microservice/adapter/middleware"
	"github.com/arvians-id/go-apriori-microservice/adapter/pkg/apriori"
	"github.com/arvians-id/go-apriori-microservice/config"
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

	// Services
	apriori.RegisterRoutes(router, configuration)

	err = router.Run(configuration.Port)
	if err != nil {
		log.Fatalln("Failed at running", err)
	}
}
