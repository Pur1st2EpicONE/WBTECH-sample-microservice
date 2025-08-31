// Package handler contains HTTP handlers for working with orders.
// It handles routing, template rendering, and JSON responses.
// Swagger annotations are used for automatic API documentation generation.
package handler

import (
	"fmt"
	"net/http"
	"strings"

	_ "github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/models"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/service"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/logger"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Handler wraps the service layer and logger for HTTP handling.
//
// Provides routes for serving API requests, static files, and HTML templates.
type Handler struct {
	service      service.ServiceProvider // service layer interface
	logger       logger.Logger           // structured logger
	TemplatePath string                  // path pattern to HTML templates
}

// NewHandler constructs a new Handler with the given service and logger.
func NewHandler(service service.ServiceProvider, logger logger.Logger) *Handler {
	return &Handler{
		service:      service,
		logger:       logger,
		TemplatePath: "web/templates/*", // default template path
	}
}

// InitRoutes initializes all HTTP routes for the service.
//
// Includes:
// - Swagger documentation at /swagger/*any
// - Static files under /static
// - API endpoints under /api/v1
// - HTML pages at root and /orders/:orderId
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Static("/static", "./web/static")
	if h.TemplatePath != "" {
		router.LoadHTMLGlob(h.TemplatePath)
	}

	api := router.Group("/api/v1")
	{
		api.GET("/orders/:orderId", h.getOrder)
	}

	basePath := router.Group("/")
	basePath.GET("/", h.showHomePage)
	basePath.GET("/orders/:orderId", h.showOrderPage)

	return router
}

// ErrorResponse defines the standard structure for JSON error messages.
//
// Used primarily for Swagger documentation to display errors in a structured and readable way.
// Without this struct, error responses may not be properly documented in Swagger.
type ErrorResponse struct {
	Error string `json:"error"`
}

// getOrder handles GET /api/v1/orders/:orderId.
//
// Returns order data in JSON with cache status indicated in the X-Cache header.
// - HIT: order retrieved from cache
// - MISS: order retrieved from database
//
// Responds with:
// - 200 OK + order JSON
// - 404 Not Found if order does not exist
// - 500 Internal Server Error on unexpected failures
//
// @Summary Get order by UID with cache status indication
// @Description Returns order details in JSON format.<br>Check <strong>X-Cache</strong> header for cache status: <strong>HIT</strong> (from cache) or <strong>MISS</strong> (from database)
// @Tags Orders
// @Produce json
// @Param orderId path string true "Order ID (UUID)"
// @Success 200 {object} models.Order "Order data"
// @Failure 404 {object} ErrorResponse "Order not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Header 200 {string} X-Cache "Cache status: HIT or MISS"
// @Router /api/v1/orders/{orderId} [get]
func (h *Handler) getOrder(c *gin.Context) {
	orderID := c.Param("orderId")
	order, fromCache, err := h.service.GetOrder(orderID, h.logger)
	if err != nil {
		h.logger.Debug("handler — failed to get order", "orderUID", orderID, "layer", "handler")
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("%s — order not found", orderID)})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "something broke on our end, sorry :("})
		return
	}
	if fromCache {
		c.Header("X-Cache", "HIT")
	} else {
		c.Header("X-Cache", "MISS") // I guess they never miss, huh? 💀
	}
	c.JSON(http.StatusOK, order)
}
