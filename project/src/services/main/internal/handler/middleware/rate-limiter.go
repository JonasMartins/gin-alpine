package middleware

import (
	"log"
	"net"
	"net/http"
	"sync"

	"gin-alpine/src/services/main/internal/bootstrap"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var ipLimiter sync.Map

//	RateLimitMiddleware limit: max req per second
//
// burst: max number of events happening at once
func RateLimitMiddleware(b *bootstrap.Bootstrap) gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter := handleGetLimiter(c.Request, b)

		if limiter != nil && !limiter.Allow() {
			b.Logger.Error("TOO_MANY_REQUESTS")

			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
			})
			return
		}

		c.Next()
	}
}

func handleGetLimiter(r *http.Request, b *bootstrap.Bootstrap) *rate.Limiter {
	ip := getIP(r)
	if ip == "" {
		return nil
	}
	limiterAny, _ := ipLimiter.LoadOrStore(ip, b.Config.NewLimiter())
	return limiterAny.(*rate.Limiter)
}

func getIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Printf("error parsing ip %v", err)
		return ""
	}
	return host
}
