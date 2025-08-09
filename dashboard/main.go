package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
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

	http.ListenAndServe(os.Getenv("DASHBOARD_PORT"), nil)
}
