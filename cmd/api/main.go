package main

import (
	"log"

	"github.com/savanrajyaguru/ecommerce-go-order-service/api"
	"github.com/savanrajyaguru/ecommerce-go-order-service/config"
	"github.com/savanrajyaguru/ecommerce-go-order-service/internal/cache"
	"github.com/savanrajyaguru/ecommerce-go-order-service/internal/database"
	"github.com/savanrajyaguru/ecommerce-go-order-service/internal/event"
	"github.com/savanrajyaguru/ecommerce-go-order-service/migrations"
	"github.com/savanrajyaguru/ecommerce-go-order-service/pkg/logger"
	"github.com/savanrajyaguru/ecommerce-go-order-service/repository"
	"github.com/savanrajyaguru/ecommerce-go-order-service/services"
)

func main() {
	// 1. Initialize Logger
	logger.InitLogger()
	defer logger.Log.Sync()

	// 2. Load Configuration
	config.LoadConfig()

	// 3. Initialize Database
	database.ConnectDB()

	// 4. Run Migrations
	migrations.RunMigrations()

	// 5. Initialize Redis
	cache.ConnectRedis()

	// 6. Initialize Kafka Producer
	event.InitProducer()
	defer event.CloseProducer()

	// 7. Initialize Services
	repo := repository.NewOrderRepository(database.DB)
	orderService, err := services.NewOrderService(repo)
	if err != nil {
		log.Fatalf("Failed to initialize order service: %v", err)
	}

	// 8. Setup Router
	r := api.SetupRouter(orderService)

	// 9. Start Server
	port := config.AppConfig.AppPort
	log.Printf("Starting Order Service on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
