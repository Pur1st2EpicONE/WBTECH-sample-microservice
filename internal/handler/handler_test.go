package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/models"
	mock_service "github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/service/mocks"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/logger"
	mock_logger "github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/logger/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestInitRoutes_OrderRoute(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := mock_service.NewMockServiceProvider(controller)
	mockLogger, _ := logger.NewLogger(configs.Logger{LogDir: "/tmp", Debug: false})

	h := NewHandler(mockService, mockLogger)
	h.TemplatePath = ""

	gin.SetMode(gin.ReleaseMode)
	router := h.InitRoutes()
	order := &models.Order{OrderUID: "orderAbobaId"}
	mockService.EXPECT().GetOrder("orderAbobaId", gomock.Any()).Return(order, false, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/orders/orderAbobaId", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.NotEqual(t, 404, w.Code)
}

func setupHandlerWithMock(t *testing.T) (*Handler, *mock_service.MockServiceProvider, *gin.Engine) {
	controller := gomock.NewController(t)

	mockService := mock_service.NewMockServiceProvider(controller)
	mockLogger, _ := logger.NewLogger(configs.Logger{LogDir: "/tmp", Debug: false})

	gin.SetMode(gin.ReleaseMode)
	h := NewHandler(mockService, mockLogger)
	router := gin.New()
	router.GET("/orders/:orderId", h.getOrder)

	return h, mockService, router
}

func TestGetOrder_Miss(t *testing.T) {
	_, mockService, router := setupHandlerWithMock(t)

	order := &models.Order{OrderUID: "test_aboba"}
	mockService.EXPECT().GetOrder("test_aboba", gomock.Any()).Return(order, false, nil)

	req := httptest.NewRequest(http.MethodGet, "/orders/test_aboba", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "MISS", w.Header().Get("X-Cache"))
}

func TestGetOrder_Hit(t *testing.T) {
	_, mockService, router := setupHandlerWithMock(t)

	order := &models.Order{OrderUID: "squid_aboba456"}
	mockService.EXPECT().GetOrder("squid_aboba456", gomock.Any()).Return(order, true, nil)

	req := httptest.NewRequest(http.MethodGet, "/orders/squid_aboba456", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "HIT", w.Header().Get("X-Cache"))
}

func TestGetOrder_Error(t *testing.T) {
	_, mockService, router := setupHandlerWithMock(t)

	mockService.EXPECT().GetOrder("game_of_abobas", gomock.Any()).Return(nil, false, errors.New("not found"))

	req := httptest.NewRequest(http.MethodGet, "/orders/game_of_abobas", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetOrder_NotFound(t *testing.T) {
	h, mockService, router := setupHandlerWithMock(t)
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLogger := mock_logger.NewMockLogger(controller)
	h.logger = mockLogger

	orderID := "a1b2o3b4a5"
	mockService.EXPECT().GetOrder(orderID, gomock.Any()).Return(nil, false, errors.New("sql: no rows in result set"))
	mockLogger.EXPECT().Debug("handler — failed to get order", "orderUID", orderID, "layer", "handler")

	req := httptest.NewRequest(http.MethodGet, "/orders/"+orderID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "order not found")
}
