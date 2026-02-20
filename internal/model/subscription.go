package model

import (
	"time"
)

type Subscription struct {
	ID          string     `json:"id" db:"id"`
	ServiceName string     `json:"service_name" db:"service_name" binding:"required"`
	Price       int        `json:"price" db:"price" binding:"required,min=0"`
	UserID      string     `json:"user_id" db:"user_id" binding:"required,uuid"`
	StartDate   string     `json:"start_date" db:"start_date" binding:"required"`
	EndDate     *string    `json:"end_date,omitempty" db:"end_date"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"-" db:"deleted_at"`
}

type CreateSubscriptionRequest struct {
	ServiceName string  `json:"service_name" binding:"required"`
	Price       int     `json:"price" binding:"required,min=0"`
	UserID      string  `json:"user_id" binding:"required,uuid"`
	StartDate   string  `json:"start_date" binding:"required"`
	EndDate     *string `json:"end_date,omitempty"`
}

type UpdateSubscriptionRequest struct {
	ServiceName *string `json:"service_name,omitempty"`
	Price       *int    `json:"price,omitempty" binding:"omitempty,min=0"`
	UserID      *string `json:"user_id,omitempty" binding:"omitempty,uuid"`
	StartDate   *string `json:"start_date,omitempty"`
	EndDate     *string `json:"end_date,omitempty"`
}

type SubscriptionFilter struct {
	UserID      string `form:"user_id"`
	ServiceName string `form:"service_name"`
	Month       string `form:"month"`
	StartMonth  string `form:"start_month"`
	EndMonth    string `form:"end_month"`
	Limit       int    `form:"limit,default=10"`
	Offset      int    `form:"offset,default=0"`
}
