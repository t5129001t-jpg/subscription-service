package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"github.com/t5129001t-jpg/subscription-service/internal/config"
	"github.com/t5129001t-jpg/subscription-service/internal/handler"
	"github.com/t5129001t-jpg/subscription-service/internal/repository"
	"github.com/t5129001t-jpg/subscription-service/internal/service"
)

func main() {
	cfg := config.LoadConfig()

	db, err := setupDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := runMigrations(cfg); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	subscriptionRepo := repository.NewSubscriptionRepository(db)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService)

	router := setupRouter(subscriptionHandler)

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

func setupDatabase(cfg *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", cfg.Database.GetDBConnString())
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Database connected successfully")
	return db, nil
}

func runMigrations(cfg *config.Config) error {
	db, err := sqlx.Connect("postgres", cfg.Database.GetDBConnString())
	if err != nil {
		return err
	}
	defer db.Close()

	if err := goose.Up(db.DB, "migrations"); err != nil {
		return err
	}

	log.Println("Migrations applied successfully")
	return nil
}

func setupRouter(subscriptionHandler *handler.SubscriptionHandler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	
	r := gin.Default()

	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	api := r.Group("/api/v1")
	{
		subscriptions := api.Group("/subscriptions")
		{
			subscriptions.POST("/", subscriptionHandler.CreateSubscription)
			subscriptions.GET("/", subscriptionHandler.ListSubscriptions)
			subscriptions.GET("/total", subscriptionHandler.GetTotalPrice)
			subscriptions.GET("/:id", subscriptionHandler.GetSubscription)
			subscriptions.PUT("/:id", subscriptionHandler.UpdateSubscription)
			subscriptions.DELETE("/:id", subscriptionHandler.DeleteSubscription)
		}
	}

	return r
}
