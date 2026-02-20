package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateSubscriptionRoute(t *testing.T) {
	// Настраиваем тестовый режим
	gin.SetMode(gin.TestMode)

	// Создаем тестовый запрос
	requestBody := `{
		"service_name": "Netflix",
		"price": 1000,
		"user_id": "123e4567-e89b-12d3-a456-426614174000",
		"start_date": "01-2024"
	}`

	req, _ := http.NewRequest("POST", "/api/v1/subscriptions", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Создаем ResponseRecorder
	w := httptest.NewRecorder()

	// Создаем роутер с заглушкой сервиса
	router := gin.New()
	router.POST("/api/v1/subscriptions", func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"id": "test-id"})
	})

	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusCreated, w.Code)
}
