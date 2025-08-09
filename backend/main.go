package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "gRPC-rating/gen"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type server struct {
	pb.UnimplementedProductViewServiceServer
	producer *kafka.Producer
}

func (s *server) SendProductView(ctx context.Context, req *pb.ProductViewRequest) (*pb.ProductViewResponse, error) {
	topic := "product-views"
	event := &pb.KafkaProductEvent{
		ProductId: req.ProductId,
	}

	protoData, err := proto.Marshal(event)
	if err != nil {
		log.Printf("kafka produce error")
		return nil, err
	}

	err = s.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: protoData,
	}, nil)

	if err != nil {
		log.Printf("kafka produce error")
		return nil, err
	}

	log.Printf("message delivered: %v", req.ProductId)
	return &pb.ProductViewResponse{Success: true}, nil
}

func main() {
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	config := &kafka.ConfigMap{
		"bootstrap.servers": kafkaBrokers,
		"client.id":         "gRPC-backend",
		"acks":              "1",
	}

	producer, err := kafka.NewProducer(config)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	lis, err := net.Listen("tcp", os.Getenv("GRPC_PORT"))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterProductViewServiceServer(s, &server{producer: producer})

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := s.Serve(lis)
		if err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	<-stop
	s.GracefulStop()
	log.Println("Server stopped")
}
