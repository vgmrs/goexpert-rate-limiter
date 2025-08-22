package limiter

import (
	"net/http"
)

type RateLimiterMiddleware struct {
	limiter      *Limiter
	ipConfig     LimitConfig
	tokenConfigs map[string]LimitConfig
}

func NewRateLimiterMiddleware(limiter *Limiter, ipConfig LimitConfig, tokenConfigs map[string]LimitConfig) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		limiter:      limiter,
		ipConfig:     ipConfig,
		tokenConfigs: tokenConfigs,
	}
}

func (m *RateLimiterMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := GetIPFromRequest(r)
		token := GetTokenFromRequest(r)

		config := m.ipConfig

		if token != "" {
			if tokenConfig, exists := m.tokenConfigs[token]; exists {
				config = tokenConfig
			} else if wildcardConfig, exists := m.tokenConfigs["*"]; exists {
				config = wildcardConfig
			}
		}

		key := GetKey(ip, token)

		allowed, err := m.limiter.CheckLimit(r.Context(), key, config)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		if !allowed {
			http.Error(
				w,
				"you have reached the maximum number of requests allowed",
				http.StatusTooManyRequests,
			)
			return
		}

		next.ServeHTTP(w, r)
	})
}
