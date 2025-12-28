package api

import (
	"github.com/gin-gonic/gin"
	middleware "github.com/savanrajyaguru/ecommerce-go-order-service/api/v1/middleware"
	"github.com/savanrajyaguru/ecommerce-go-order-service/api/v1/order"
	"github.com/savanrajyaguru/ecommerce-go-order-service/config"
	"github.com/savanrajyaguru/ecommerce-go-order-service/services"
)

func SetupRouter(orderService *services.OrderService) *gin.Engine {
	if config.AppConfig.AppPort == "" && config.AppConfig.GrpcPort == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(middleware.Cors())

	v1 := r.Group("/v1")
	{
		order.RegisterRoutes(v1, orderService)
	}

	return r
}
