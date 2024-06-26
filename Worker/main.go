// main.go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nikhilsiwach28/MyCode.git/consumer"
	"github.com/nikhilsiwach28/MyCode.git/producer"
	"github.com/nikhilsiwach28/MyCode.git/redis"
)

func main() {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	redisClient := redis.NewRedisService("localhost:6379", "", 0)

	// Initialize Kafka consumer
	consumer, err := consumer.NewConsumerFromEnv()
	if err != nil {
		fmt.Printf("Error initializing Kafka consumer: %v\n", err)
		os.Exit(1)
	}
	defer consumer.Close()

	// Initialize Kafka producer
	producer, err := producer.NewProducerFromEnv()
	if err != nil {
		fmt.Printf("Error initializing Kafka producer: %v\n", err)
		os.Exit(1)
	}
	defer producer.Close()

	// Start consuming messages by passing Producer instance
	fmt.Println("Consuming Messages")
	go consumer.ConsumeMessages(sigterm, producer, redisClient)

	// Wait for shutdown signal
	<-sigterm
	fmt.Println("Shutting down worker service")
}
