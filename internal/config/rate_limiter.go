package config

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/vgmrs/goexpert-rate-limiter/internal/limiter"
	"github.com/vgmrs/goexpert-rate-limiter/internal/redis"
)

type RateLimiterConfig struct {
	RedisAddress  string
	RedisPassword string
	RedisDB       int

	IPMaxRequests      int
	IPBlockDuration    time.Duration
	TokenMaxRequests   int
	TokenBlockDuration time.Duration
}

func loadRateLimiterConfig() (*RateLimiterConfig, error) {
	redisAddr := os.Getenv("REDIS_ADDRESS")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")

	redisDB := 0
	if db := os.Getenv("REDIS_DB"); db != "" {
		if dbVal, err := strconv.Atoi(db); err == nil {
			redisDB = dbVal
		}
	}

	ipMaxRequests := 10
	if v := os.Getenv("RATE_LIMITER_IP_MAX_REQUESTS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			ipMaxRequests = n
		}
	}

	ipBlockDuration := 5 * time.Minute
	if v := os.Getenv("RATE_LIMITER_IP_BLOCK_DURATION"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			ipBlockDuration = d
		}
	}

	tokenMaxRequests := 100
	if v := os.Getenv("RATE_LIMITER_TOKEN_MAX_REQUESTS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			tokenMaxRequests = n
		}
	}

	tokenBlockDuration := time.Hour
	if v := os.Getenv("RATE_LIMITER_TOKEN_BLOCK_DURATION"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			tokenBlockDuration = d
		}
	}

	return &RateLimiterConfig{
		RedisAddress:       redisAddr,
		RedisPassword:      redisPassword,
		RedisDB:            redisDB,
		IPMaxRequests:      ipMaxRequests,
		IPBlockDuration:    ipBlockDuration,
		TokenMaxRequests:   tokenMaxRequests,
		TokenBlockDuration: tokenBlockDuration,
	}, nil
}

func SetupRateLimiter() (func(http.Handler) http.Handler, error) {
	cfg, err := loadRateLimiterConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading rate limiter config: %w", err)
	}

	redisStorage, err := redis.NewRedisStorage(cfg.RedisAddress, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		return nil, err
	}

	limiterInstance := limiter.NewLimiter(redisStorage)

	ipConfig := limiter.LimitConfig{
		MaxRequests:   cfg.IPMaxRequests,
		BlockDuration: cfg.IPBlockDuration,
	}

	tokenConfigs := make(map[string]limiter.LimitConfig)
	tokenConfigs["*"] = limiter.LimitConfig{
		MaxRequests:   cfg.TokenMaxRequests,
		BlockDuration: cfg.TokenBlockDuration,
	}

	middleware := limiter.NewRateLimiterMiddleware(limiterInstance, ipConfig, tokenConfigs)
	return middleware.Handler, nil
}
