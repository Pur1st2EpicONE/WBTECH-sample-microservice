package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// showHomePage renders the main page of the application.
func (h *Handler) showHomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "layout.html", nil)
}

// showOrderPage renders the page for a specific order.
// The order ID is passed to the template context.
func (h *Handler) showOrderPage(c *gin.Context) {
	c.HTML(http.StatusOK, "order.html", gin.H{
		"orderId": c.Param("orderId"),
	})
}
