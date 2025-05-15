package main

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go-basics/config"
	"go-basics/internal/handler"
	"go-basics/internal/middleware"
	"go-basics/internal/middleware/auth"
	httpresponseencoder "go-basics/pkg/http"
	"net/http"
	"time"
)

func main() {
	// Load configuration
	cfg, err := config.LoadDefaultConfig()
	if err != nil {
		log.Fatal().Msgf("failed to load config %v", err)
	}

	// Configure logging
	configureLogging(cfg.Server.Debug)

	// create dependencies
	responseEncoder := httpresponseencoder.NewJSONEncoder()
	authenticator := auth.NewTokenAuthenticator(cfg.Auth.Token)

	// Create HTTP handlers
	apiHandler := handler.NewAPIHandler(responseEncoder)
	helloHandler := handler.NewHelloHandler(responseEncoder)

	// Set up router
	mux := http.NewServeMux()
	mux.Handle("/api", apiHandler)
	mux.Handle("/api/hello", helloHandler)

	// Apply middleware
	var muxHandler http.Handler = mux
	muxHandler = middleware.LoggingMiddleware(muxHandler)
	muxHandler = middleware.AuthMiddleware(authenticator)(muxHandler)

	// Apply rate limiting middleware if enabled
	if cfg.RateLimit.Enabled {
		windowSize, err := time.ParseDuration(cfg.RateLimit.WindowSize)
		if err != nil {
			log.Fatal().Msgf("failed to parse rate limit interval: %v", err)
		}
		limiter := middleware.NewRateLimit(cfg.RateLimit.MaxRequests, windowSize)
		muxHandler = middleware.RateLimiterMiddleware(limiter)(muxHandler)
	}

	// Start server
	log.Info().Msgf("Starting server on port: %d", cfg.Server.Port)
	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	if err := http.ListenAndServe(serverAddr, muxHandler); err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}

func configureLogging(debug bool) {
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
