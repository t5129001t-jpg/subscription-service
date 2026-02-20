package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/t5129001t-jpg/subscription-service/internal/model"
)

type SubscriptionRepository interface {
	Create(sub *model.Subscription) error
	GetByID(id string) (*model.Subscription, error)
	Update(id string, updates map[string]interface{}) error
	Delete(id string) error
	List(filter model.SubscriptionFilter) ([]model.Subscription, int, error)
	GetTotalPrice(userID, serviceName, startMonth, endMonth string) (int, error)
}

type subscriptionRepository struct {
	db *sqlx.DB
}

func NewSubscriptionRepository(db *sqlx.DB) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) Create(sub *model.Subscription) error {
	query := `
		INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
	).Scan(&sub.ID, &sub.CreatedAt, &sub.UpdatedAt)

	return err
}

func (r *subscriptionRepository) GetByID(id string) (*model.Subscription, error) {
	var sub model.Subscription
	query := `SELECT * FROM subscriptions WHERE id = $1 AND deleted_at IS NULL`
	err := r.db.Get(&sub, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &sub, err
}

func (r *subscriptionRepository) Update(id string, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	setClauses := make([]string, 0, len(updates))
	args := make([]interface{}, 0, len(updates)+1)
	i := 1

	for field, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", field, i))
		args = append(args, value)
		i++
	}

	args = append(args, id)
	query := fmt.Sprintf(`
		UPDATE subscriptions 
		SET %s, updated_at = NOW()
		WHERE id = $%d AND deleted_at IS NULL
	`, strings.Join(setClauses, ", "), i)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *subscriptionRepository) Delete(id string) error {
	query := `UPDATE subscriptions SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *subscriptionRepository) List(filter model.SubscriptionFilter) ([]model.Subscription, int, error) {
	var subscriptions []model.Subscription
	var total int

	baseQuery := `FROM subscriptions WHERE deleted_at IS NULL`
	args := make([]interface{}, 0)
	conditions := make([]string, 0)
	argCount := 1

	if filter.UserID != "" {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argCount))
		args = append(args, filter.UserID)
		argCount++
	}

	if filter.ServiceName != "" {
		conditions = append(conditions, fmt.Sprintf("service_name = $%d", argCount))
		args = append(args, filter.ServiceName)
		argCount++
	}

	if filter.Month != "" {
		conditions = append(conditions, fmt.Sprintf("(start_date <= $%d AND (end_date IS NULL OR end_date >= $%d))", argCount, argCount))
		args = append(args, filter.Month)
		argCount++
	}

	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(*) " + baseQuery
	err := r.db.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	dataQuery := "SELECT * " + baseQuery + " ORDER BY start_date DESC"
	
	if filter.Limit > 0 {
		dataQuery += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filter.Limit)
		argCount++
	}
	
	if filter.Offset > 0 {
		dataQuery += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, filter.Offset)
		argCount++
	}

	err = r.db.Select(&subscriptions, dataQuery, args...)
	return subscriptions, total, err
}

func (r *subscriptionRepository) GetTotalPrice(userID, serviceName, startMonth, endMonth string) (int, error) {
	var total int

	query := `
		SELECT COALESCE(SUM(price), 0)
		FROM subscriptions
		WHERE deleted_at IS NULL
		AND ($1 = '' OR $1 IS NULL OR user_id = $1)
		AND ($2 = '' OR $2 IS NULL OR service_name = $2)
		AND start_date >= $3
		AND (end_date IS NULL OR end_date <= $4)
	`

	err := r.db.Get(&total, query, userID, serviceName, startMonth, endMonth)
	return total, err
}
