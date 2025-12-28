package migrations

import (
	"log"

	"github.com/savanrajyaguru/ecommerce-go-order-service/internal/database"
	"github.com/savanrajyaguru/ecommerce-go-order-service/models"
)

func RunMigrations() {
	log.Println("Running database migrations...")
	database.Migrate(&models.Order{}, &models.OrderItem{})
	log.Println("Migrations completed successfully.")
}
