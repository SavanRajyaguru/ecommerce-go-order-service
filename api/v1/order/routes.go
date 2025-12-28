package order

import (
	"github.com/gin-gonic/gin"
	"github.com/savanrajyaguru/ecommerce-go-order-service/internal/auth"
	"github.com/savanrajyaguru/ecommerce-go-order-service/services"
)

func RegisterRoutes(router *gin.RouterGroup, service *services.OrderService) {
	handler := NewOrderHandler(service)

	orderGroup := router.Group("/orders")
	orderGroup.Use(auth.AuthMiddleware()) // Protect all routes
	{
		orderGroup.POST("", handler.CreateOrder)
		orderGroup.GET("", handler.GetOrders)
		orderGroup.GET("/:id", handler.GetOrder)
		orderGroup.POST("/:id/cancel", handler.CancelOrder)
	}
}
