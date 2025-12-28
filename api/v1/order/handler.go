package order

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/savanrajyaguru/ecommerce-go-order-service/internal/auth"
	"github.com/savanrajyaguru/ecommerce-go-order-service/pkg/logger"
	"github.com/savanrajyaguru/ecommerce-go-order-service/pkg/utils"
	"github.com/savanrajyaguru/ecommerce-go-order-service/services"
	"go.uber.org/zap"
)

type OrderHandler struct {
	service *services.OrderService
}

func NewOrderHandler(service *services.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req services.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	// Get UserID from context (set by Auth Middleware)
	userID, exists := c.Get("userID")
	if !exists {
		utils.JSONError(c, http.StatusUnauthorized, "Unixuthenticated", "")
		return
	}
	req.UserID = userID.(string)

	// Get Token from Header for forwarding
	req.Token = c.GetHeader("Authorization")[7:] // Remove "Bearer "

	order, err := h.service.CreateOrder(c.Request.Context(), req)
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "Failed to create order", err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusCreated, true, "Order created successfully", order)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.JSONError(c, http.StatusBadRequest, "Order ID required", "")
		return
	}

	order, err := h.service.GetOrder(c.Request.Context(), id)
	if err != nil {
		utils.JSONError(c, http.StatusNotFound, "Order not found", err.Error())
		return
	}

	// Auth Check: User can only see own orders, Admin sees all
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	if role.(string) != auth.RoleAdmin && order.UserID != userID.(string) {
		utils.JSONError(c, http.StatusForbidden, "Unauthorized access to order", "")
		return
	}

	utils.JSONResponse(c, http.StatusOK, true, "Order retrieved", order)
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	// If Admin -> All Orders? Or specific filtering?
	// Prompt: "Get Order history (USER)" and "Get Order by ID (USER, ADMIN)".
	// Prompt implies User gets own history. Admin validation implies they might want all.
	// But let's stick to: If User -> GetUserOrders. If Admin -> GetAllOrders?

	role, _ := c.Get("role")
	userID, _ := c.Get("userID")

	if role.(string) == auth.RoleAdmin {
		// Admin fetch all?
		orders, err := h.service.GetAllOrders(c.Request.Context())
		if err != nil {
			utils.JSONError(c, http.StatusInternalServerError, "Failed to fetch orders", err.Error())
			return
		}
		utils.JSONResponse(c, http.StatusOK, true, "All orders retrieved", orders)
		return
	}

	// Regular User
	orders, err := h.service.GetUserOrders(c.Request.Context(), userID.(string))
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "Failed to fetch order history", err.Error())
		return
	}
	utils.JSONResponse(c, http.StatusOK, true, "Order history retrieved", orders)
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")

	if err := h.service.CancelOrder(c.Request.Context(), id, userID.(string)); err != nil {
		// Could distinguish 400 vs 500 based on error
		logger.Log.Error("Cancel order failed", zap.Error(err))
		utils.JSONError(c, http.StatusBadRequest, "Failed to cancel order", err.Error())
		return
	}

	utils.JSONResponse(c, http.StatusOK, true, "Order cancelled successfully", nil)
}
