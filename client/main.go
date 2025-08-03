package main

import (
	"context"
	"log"
	"time"

	pb "gRPC-rating/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client := pb.NewProductViewServiceClient(conn)

	for {
		_, err = client.SendProductView(ctx, &pb.ProductViewRequest{
			ProductId: "test",
		})
		if err != nil {
			log.Printf("Failed to send view %v", err)
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}
}
