package ratelimit

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RateLimiter(cap int64, rate int64)  func(handlerFunc gin.HandlerFunc)  gin.HandlerFunc{
	rateLimiter := NewTokenBuket(cap, rate)
	return func(handlerFunc gin.HandlerFunc) gin.HandlerFunc {
		return func(context *gin.Context) {
			if rateLimiter.IsAccept(){
				handlerFunc(context)
			} else {
				context.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
					"code":http.StatusTooManyRequests,
					"message": "too many requests",
				})
			}
		}
	}
}
