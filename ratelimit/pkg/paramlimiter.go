package pkg

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 根据参数key 进行限流
func ParamLimiter(cap int64, rate int64, key string)func(handler gin.HandlerFunc) gin.HandlerFunc  {
	limiter := NewTokenBuket(cap, rate)
	return func(handler gin.HandlerFunc) gin.HandlerFunc {
		return func(context *gin.Context) {
			if context.Query(key) != "" {
				if limiter.IsAccept() {
					handler(context)
				} else {
					context.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
						"code": http.StatusTooManyRequests,
						"message": "too many requests",
					})
				}
			} else {  // 不需要进行限流
				handler(context)
			}
		}
	}
}
