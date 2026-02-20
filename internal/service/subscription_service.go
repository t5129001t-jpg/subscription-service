package service

import (
	"errors"
	"regexp"
	"time"

	"github.com/t5129001t-jpg/subscription-service/internal/model"
	"github.com/t5129001t-jpg/subscription-service/internal/repository"
)

type SubscriptionService interface {
	Create(req *model.CreateSubscriptionRequest) (*model.Subscription, error)
	GetByID(id string) (*model.Subscription, error)
	Update(id string, req *model.UpdateSubscriptionRequest) error
	Delete(id string) error
	List(filter model.SubscriptionFilter) ([]model.Subscription, int, error)
	GetTotalPrice(filter model.SubscriptionFilter) (int, error)
}

type subscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{repo: repo}
}

func (s *subscriptionService) Create(req *model.CreateSubscriptionRequest) (*model.Subscription, error) {
	if !isValidDateFormat(req.StartDate) {
		return nil, errors.New("invalid start_date format, expected MM-YYYY")
	}

	if req.EndDate != nil && !isValidDateFormat(*req.EndDate) {
		return nil, errors.New("invalid end_date format, expected MM-YYYY")
	}

	if req.EndDate != nil {
		if !isEndDateAfterStartDate(req.StartDate, *req.EndDate) {
			return nil, errors.New("end_date must be after or equal to start_date")
		}
	}

	sub := &model.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	err := s.repo.Create(sub)
	return sub, err
}

func (s *subscriptionService) GetByID(id string) (*model.Subscription, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	return s.repo.GetByID(id)
}

func (s *subscriptionService) Update(id string, req *model.UpdateSubscriptionRequest) error {
	if id == "" {
		return errors.New("id is required")
	}

	updates := make(map[string]interface{})

	if req.ServiceName != nil {
		updates["service_name"] = *req.ServiceName
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.UserID != nil {
		updates["user_id"] = *req.UserID
	}
	if req.StartDate != nil {
		if !isValidDateFormat(*req.StartDate) {
			return errors.New("invalid start_date format")
		}
		updates["start_date"] = *req.StartDate
	}
	if req.EndDate != nil {
		if !isValidDateFormat(*req.EndDate) {
			return errors.New("invalid end_date format")
		}
		updates["end_date"] = *req.EndDate
	}

	return s.repo.Update(id, updates)
}

func (s *subscriptionService) Delete(id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	return s.repo.Delete(id)
}

func (s *subscriptionService) List(filter model.SubscriptionFilter) ([]model.Subscription, int, error) {
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	return s.repo.List(filter)
}

func (s *subscriptionService) GetTotalPrice(filter model.SubscriptionFilter) (int, error) {
	if filter.Month != "" {
		if !isValidDateFormat(filter.Month) {
			return 0, errors.New("invalid month format")
		}
		return s.repo.GetTotalPrice(filter.UserID, filter.ServiceName, filter.Month, filter.Month)
	}

	if filter.StartMonth == "" || filter.EndMonth == "" {
		return 0, errors.New("both start_month and end_month must be provided")
	}

	if !isValidDateFormat(filter.StartMonth) {
		return 0, errors.New("invalid start_month format")
	}
	if !isValidDateFormat(filter.EndMonth) {
		return 0, errors.New("invalid end_month format")
	}

	return s.repo.GetTotalPrice(filter.UserID, filter.ServiceName, filter.StartMonth, filter.EndMonth)
}

func isValidDateFormat(date string) bool {
	match, _ := regexp.MatchString(`^(0[1-9]|1[0-2])-[0-9]{4}$`, date)
	return match
}

func isEndDateAfterStartDate(start, end string) bool {
	startTime, _ := time.Parse("01-2006", start)
	endTime, _ := time.Parse("01-2006", end)
	return !endTime.Before(startTime)
}
