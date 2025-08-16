package handler

import (
	"net/http"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
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
	order, fromCache, err := h.service.GetOrder(orderID)
	if err != nil {
		logger.LogError("error getting the order", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request data"})
		return
	}
	if fromCache {
		c.Header("X-Cache", "HIT")
	} else {
		c.Header("X-Cache", "MISS") // I guess they never miss, huh? 💀
	}
	c.JSON(http.StatusOK, order)
}
