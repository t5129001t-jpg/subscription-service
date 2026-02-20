package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/t5129001t-jpg/subscription-service/internal/model"
	"github.com/t5129001t-jpg/subscription-service/internal/service"
)

// Мок для репозитория
type MockSubscriptionRepository struct {
	mock.Mock
}

func (m *MockSubscriptionRepository) Create(sub *model.Subscription) error {
	args := m.Called(sub)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) GetByID(id string) (*model.Subscription, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) Update(id string, updates map[string]interface{}) error {
	args := m.Called(id, updates)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) List(filter model.SubscriptionFilter) ([]model.Subscription, int, error) {
	args := m.Called(filter)
	return args.Get(0).([]model.Subscription), args.Int(1), args.Error(2)
}

func (m *MockSubscriptionRepository) GetTotalPrice(userID, serviceName, startMonth, endMonth string) (int, error) {
	args := m.Called(userID, serviceName, startMonth, endMonth)
	return args.Int(0), args.Error(1)
}

// Тесты для сервиса
func TestCreateSubscription_ValidData(t *testing.T) {
	mockRepo := new(MockSubscriptionRepository)
	svc := service.NewSubscriptionService(mockRepo)

	req := &model.CreateSubscriptionRequest{
		ServiceName: "Netflix",
		Price:       1000,
		UserID:      "123e4567-e89b-12d3-a456-426614174000",
		StartDate:   "01-2024",
	}

	expectedSub := &model.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   req.StartDate,
	}

	mockRepo.On("Create", mock.AnythingOfType("*model.Subscription")).Return(nil)

	sub, err := svc.Create(req)

	assert.NoError(t, err)
	assert.Equal(t, expectedSub.ServiceName, sub.ServiceName)
	assert.Equal(t, expectedSub.Price, sub.Price)
	mockRepo.AssertExpectations(t)
}

func TestCreateSubscription_InvalidDate(t *testing.T) {
	mockRepo := new(MockSubscriptionRepository)
	svc := service.NewSubscriptionService(mockRepo)

	req := &model.CreateSubscriptionRequest{
		ServiceName: "Netflix",
		Price:       1000,
		UserID:      "123e4567-e89b-12d3-a456-426614174000",
		StartDate:   "13-2024", // Неверный месяц
	}

	sub, err := svc.Create(req)

	assert.Error(t, err)
	assert.Nil(t, sub)
	assert.Contains(t, err.Error(), "invalid start_date format")
}
