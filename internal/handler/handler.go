package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/service"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.ServiceProvider
	logger  logger.Logger
}

func NewHandler(service service.ServiceProvider, logger logger.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	basePath := router.Group("/")
	{
		orders := basePath.Group("/orders")
		{
			orders.GET("/:orderId", h.getOrder)
		}
	}
	return router
}

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
