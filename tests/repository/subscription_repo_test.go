package repository_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/t5129001t-jpg/subscription-service/internal/model"
	"github.com/t5129001t-jpg/subscription-service/internal/repository"
)

func TestCreateSubscription(t *testing.T) {
	// Создаем мок БД
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewSubscriptionRepository(sqlxDB)

	// Подготовка данных
	sub := &model.Subscription{
		ServiceName: "Netflix",
		Price:       1000,
		UserID:      "123e4567-e89b-12d3-a456-426614174000",
		StartDate:   "01-2024",
	}

	// Ожидаем запрос к БД
	mock.ExpectQuery(`INSERT INTO subscriptions`).
		WithArgs(sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow("123e4567-e89b-12d3-a456-426614174001", time.Now(), time.Now()))

	// Выполняем тестируемую функцию
	err = repo.Create(sub)

	// Проверяем результаты
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
