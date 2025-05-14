package middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// Rate Limiter Middleware
type RateLimiter struct {
	mu              sync.Mutex             // mutex for thread safety
	requests        map[string][]time.Time // map of IP addresses to slice of request times
	maxRequests     int                    // maximum requests per interval
	windowSize      time.Duration          // interval for rate limit (e.g., 1 minute)
	cleanupInterval time.Duration          // interval for cleanup goroutine
}

func NewRateLimit(maxRequests int, windowSize time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests:        make(map[string][]time.Time),
		maxRequests:     maxRequests,
		windowSize:      windowSize,
		cleanupInterval: windowSize * 2, // Cleanup every 2x the window size
	}

	// start cleanup goroutine
	go rl.cleanup()

	return rl
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.windowSize)
	defer ticker.Stop()

	for range ticker.C {
		log.Info().Msg("Cleaning up rate limiter")
		rl.mu.Lock()

		now := time.Now()

		for ip, timestamps := range rl.requests {
			var valid []time.Time
			for _, ts := range timestamps {
				if now.Sub(ts) <= rl.windowSize {
					valid = append(valid, ts)
				}
			}

			if len(valid) == 0 {
				delete(rl.requests, ip)
			} else {
				rl.requests[ip] = valid
			}
		}

		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) isAllowed(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	timestamps, exists := rl.requests[ip]

	if !exists {
		rl.requests[ip] = []time.Time{now}
		return true
	}

	// Filter out timestamps outside the window
	var validTimestamps []time.Time
	for _, ts := range timestamps {
		if now.Sub(ts) <= rl.windowSize {
			validTimestamps = append(validTimestamps, ts)
		}
	}

	if len(validTimestamps) >= rl.maxRequests {
		rl.requests[ip] = validTimestamps // update with filtered timestamps
		return false
	}

	validTimestamps = append(validTimestamps, now)
	rl.requests[ip] = validTimestamps

	return true
}

func RateLimiterMiddleware(rl *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getClientIP(r)
			log.Info().Msg(ip)
			if !rl.isAllowed(ip) {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-Ip header
	if xri := r.Header.Get("X-Real-Ip"); xri != "" {
		return xri
	}

	// If neither header is set, fallback to the remote address
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
