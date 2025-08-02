package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	pb "gRPC-rating/gen/github.com/ArtemKVD/gRPC-rating/gen"

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
	client := pb.NewProductViewServiceClient(conn)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("client/web/templates/index.html"))
		tmpl.Execute(w, nil)
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
	http.ListenAndServe(":8081", nil)
}
