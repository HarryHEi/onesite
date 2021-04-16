package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"

	"onesite/common/rest"
)

type RateLimiter struct {
	Bucket *ratelimit.Bucket
}

func InitRateLimiter(qps int64) *RateLimiter {
	return &RateLimiter{
		ratelimit.NewBucketWithRate(float64(qps), qps),
	}
}

func (r *RateLimiter) Take() bool {
	return r.Bucket.TakeAvailable(1) > 0
}

func (r *RateLimiter) Middleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		if !r.Take() {
			rest.BadRequest(c, errors.New("rate limited"))
			return
		}
	}
}
