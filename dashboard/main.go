package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	http.HandleFunc("/ratings", func(w http.ResponseWriter, r *http.Request) {
		ratings, err := rdb.ZRevRangeWithScores(context.Background(), "product_ratings", 0, 10).Result()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, item := range ratings {
			fmt.Fprintf(w, "%s: %v \n", item.Member, item.Score)
		}
	})

	http.ListenAndServe(":8080", nil)
}
