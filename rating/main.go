package main

import (
	"context"
	"errors"
	"log"
	"time"

	pb "gRPC-rating/gen/github.com/ArtemKVD/gRPC-rating/gen"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RatingService struct {
	redisClient     *redis.Client
	dashboardClient pb.DashboardServiceClient
	grpcConn        *grpc.ClientConn
}

func main() {

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	conn, err := grpc.NewClient(
		"localhost:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	service := &RatingService{
		redisClient:     rdb,
		dashboardClient: pb.NewDashboardServiceClient(conn),
		grpcConn:        conn,
	}

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "rating-group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	err = consumer.SubscribeTopics([]string{"product-views"}, nil)
	if err != nil {
		log.Fatal(err)
	}

	for {
		msg, err := consumer.ReadMessage(10 * time.Second)
		if msg != nil {
			log.Printf("Delivered message: %v", msg)
		}
		if err != nil {
			log.Printf("Consumer error: %v", err)
			continue
		}

		err = service.RatingUpdate(context.Background(), msg.Value)
		if err != nil {
			log.Printf("Processing failed: %v", err)
		}
	}
}

func (s *RatingService) RatingUpdate(ctx context.Context, msgValue []byte) error {
	productID := string(msgValue)
	if productID == "" {
		return errors.New("invalid product id")
	}

	_, err := s.redisClient.ZIncrBy(ctx, "product_ratings", 1, productID).Result()
	return err
}
