package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

func RateLimiter(redisClient *redis.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get client IP address
			ip := r.RemoteAddr
			key := fmt.Sprintf("ratelimit:%s", ip)

			// Get current count from Redis
			ctx := context.Background()
			count, err := redisClient.Get(ctx, key).Int()
			if err != nil && err != redis.Nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// If count is 0 or key doesn't exist, set it to 1 with 1 minute expiry
			if err == redis.Nil {
				err = redisClient.Set(ctx, key, 1, time.Minute).Err()
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			} else if count >= 5 {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				fmt.Fprintf(w, `{"status":"error","message":"Rate limit exceeded. Please try again in 1 minute."}`)
				return
			} else {
				// Increment the counter
				err = redisClient.Incr(ctx, key).Err()
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
