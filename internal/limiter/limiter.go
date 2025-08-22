package limiter

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type LimitConfig struct {
	MaxRequests   int
	BlockDuration time.Duration
}

type Limiter struct {
	storage Storage
}

func NewLimiter(storage Storage) *Limiter {
	return &Limiter{
		storage: storage,
	}
}

func (l *Limiter) CheckLimit(ctx context.Context, key string, config LimitConfig) (bool, error) {
	blockKey := fmt.Sprintf("block:%s", key)
	countKey := fmt.Sprintf("count:%s", key)

	exists, err := l.storage.Exists(ctx, blockKey)
	if err != nil {
		return false, err
	}
	if exists {
		return false, nil
	}

	count, err := l.storage.Increment(ctx, countKey, int64(config.BlockDuration.Seconds()))
	if err != nil {
		return false, err
	}

	if count > int64(config.MaxRequests) {
		err = l.storage.Set(ctx, blockKey, "1", int64(config.BlockDuration.Seconds()))
		if err != nil {
			return false, err
		}

		_ = l.storage.Delete(ctx, countKey)
		return false, nil
	}

	return true, nil
}

func GetKey(ip, token string) string {
	if token != "" {
		return fmt.Sprintf("token:%s", token)
	}
	return fmt.Sprintf("ip:%s", ip)
}

func GetIPFromRequest(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")

	if ip != "" {
		ips := strings.Split(ip, ",")
		ip = strings.TrimSpace(ips[0])
	} else {
		ip = r.RemoteAddr

		if strings.Contains(ip, ":") {
			ip = strings.Split(ip, ":")[0]
		}
	}

	return ip
}

func GetTokenFromRequest(r *http.Request) string {
	token := r.Header.Get("API_KEY")
	if token == "" {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			return ""
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			token = parts[1]
		}
	}

	return token
}
