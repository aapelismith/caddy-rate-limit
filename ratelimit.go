package ratelimit

import (
	"net/http"

	"github.com/mholt/caddy/caddyhttp/httpserver"
)

type RateLimit struct {
	Next  httpserver.Handler
	Rules []Rule
}

type Rule struct {
	Rate      float64
	Burst     int
	Resources []string
}

var (
	customLimiter *CustomLimiter
)

func init() {

	customLimiter = NewCustomLimiter()
}

func (rl RateLimit) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {

	for _, rule := range rl.Rules {
		for _, res := range rule.Resources {
			if !httpserver.Path(r.URL.Path).Matches(res) {
				continue
			}

			sliceKeys := buildKeys(res, r)
			for _, keys := range sliceKeys {
				ret := customLimiter.Allow(keys, rule)
				if !ret {
					return http.StatusTooManyRequests, nil
				}
			}
		}
	}

	return rl.Next.ServeHTTP(w, r)
}