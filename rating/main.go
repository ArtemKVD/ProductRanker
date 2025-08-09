package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "gRPC-rating/gen"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

type RatingService struct {
	redisClient     *redis.Client
	dashboardClient pb.DashboardServiceClient
	grpcConn        *grpc.ClientConn
}

func main() {

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	conn, err := grpc.NewClient(
		os.Getenv("DASHBOARD_ADDRESS"),
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
		"bootstrap.servers": "kafka:9092",
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
	var event pb.KafkaProductEvent
	err := proto.Unmarshal(msgValue, &event)
	if err != nil {
		return err
	}

	_, err = s.redisClient.ZIncrBy(context.Background(), "product_ratings", 1, event.ProductId).Result()

	return err
}
