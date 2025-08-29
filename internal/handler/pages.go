package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) showHomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "layout.html", nil)
}

func (h *Handler) showOrderPage(c *gin.Context) {
	c.HTML(200, "order.html", gin.H{
		"orderId": c.Param("orderId"),
	})
}
