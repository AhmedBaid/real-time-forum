package middleware

import (
	"net/http"
	"real_time/backend/config"
	"sync"
	"time"
)

var (
	requests = make(map[string]int)
	lastSeen = make(map[string]time.Time)
	mu       sync.Mutex
)

func RateLimit(HandlerFunc http.HandlerFunc) http.HandlerFunc  {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		mu.Lock()
		defer mu.Unlock()

		if time.Since(lastSeen[ip]) > time.Minute {
			requests[ip] = 0
		}

		requests[ip]++
		lastSeen[ip] = time.Now()

		if requests[ip] > 30 {
			config.ResponseJSON(w, http.StatusTooManyRequests, map[string]any{
				"message": "too many request , try again after 1  min ",
				"status":  http.StatusTooManyRequests,
			})
			return
		}
		HandlerFunc.ServeHTTP(w, r)
	}
}
