package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	pb "gRPC-rating/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(
		os.Getenv("GRPC_SERVER_ADDRESS"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewProductViewServiceClient(conn)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "templates/index.html")
	})

	http.HandleFunc("/view", func(w http.ResponseWriter, r *http.Request) {

		productID := strings.TrimSpace(r.FormValue("product_id"))

		_, err := client.SendProductView(r.Context(), &pb.ProductViewRequest{
			ProductId: productID,
		})

		if err != nil {
			log.Printf("error send production view %v", err)
		}
	})
	http.ListenAndServe(os.Getenv("WEB_PORT"), nil)
}
