package cache

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/savanrajyaguru/ecommerce-go-order-service/config"
)

var Rdb *redis.Client
var Ctx = context.Background()

func ConnectRedis() {
	cfg := config.AppConfig.Redis
	if !cfg.Enabled {
		log.Println("Redis is disabled in configuration")
		return
	}

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	Rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       0, // Use default DB
	})

	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		log.Printf("Failed to connect to Redis at %s: %v", addr, err)
		// We might not want to panic if Redis is optional, but for order service idempotency it's critical.
		// For now, log fatal as per reliability requirements.
		log.Fatal("Redis connection failed")
	}

	log.Println("Connected to Redis successfully")
}
