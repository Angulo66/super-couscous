package main

import (
	"flag"
	"fmt"
	"go-basics/config"
	"go-basics/internal/handler"
	"go-basics/internal/middleware"
	http_response_encoder "go-basics/pkg/http"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Use the default config loader
	cfg, err := config.LoadDefaultConfig()
	if err != nil {
		log.Fatal().Msgf("failed to load config %v", err)
	}

	// Or use a custom loader
	// loader := config.NewViperLoader("custom-config", "json", []string{"./config"}, true)
	// cfg, err := config.LoadConfig(loader)
	// if err != nil {
	// 	log.Fatal().Msgf("failed to load config %v", err)
	// }

	log.Info().Msgf("Starting server on port: %d", cfg.Server.Port)

	//zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	debug := flag.Bool("debug", false, "enable debug logging")
	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	mux := http.NewServeMux()

	var muxHandler http.Handler = mux

	muxHandler = middleware.LoggingMiddleware(muxHandler)
	muxHandler = middleware.AuthMiddleware(cfg, muxHandler)

	if cfg.RateLimit.Enabled {
		windowSize, err := time.ParseDuration(cfg.RateLimit.WindowSize)
		if err != nil {
			log.Fatal().Msgf("failed to parse rate limit interval: %v", err)
		}
		limiter := middleware.NewRateLimit(cfg.RateLimit.MaxRequests, windowSize)
		muxHandler = middleware.RateLimiterMiddleware(limiter)(muxHandler)
	}

	// Initialize the response encoder
	responseEncoder := http_response_encoder.NewJSONEncoder()
	// Initialize handlers with dependency injection
	apiHandler := handler.NewAPIHandler(responseEncoder)
	helloHandler := handler.NewHelloHandler(responseEncoder)

	mux.Handle("/api", apiHandler)
	mux.Handle("/api/hello", helloHandler)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), muxHandler); err != nil {
		log.Fatal().Err(err).Msg("server failed to start")
	}

}
