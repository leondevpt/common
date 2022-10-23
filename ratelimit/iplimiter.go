package ratelimit


import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sync"
	"time"
)

type LimiterCache struct {
	data sync.Map // key ==ip+端口  value==>bucket
}
var IpCache *LimiterCache
var IpCache2 *Cache
func init() {
	IpCache=&LimiterCache{}
	IpCache2=NewCache(WithMaxSize(10000))

}
func IPLimiter(cap int64,rate int64) func(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(handler gin.HandlerFunc) gin.HandlerFunc {
		return func(c *gin.Context) {
			ip:=c.Request.RemoteAddr
			var limiter *TokenBuket

			//if v,ok:=IpCache.data.Load(ip);ok{
			//	limiter=v.(*Bucket)
			//}else{
			//	limiter=NewBucket(cap,rate )
			//	IpCache.data.Store(ip,limiter)
			//}
			if v:=IpCache2.Get(ip);v!=nil{
				limiter=v.(*TokenBuket)
			}else{
				limiter=NewTokenBuket(cap,rate )
				log.Print("from cache")
				IpCache2.Set(ip,limiter,time.Second*5)
			}

			if limiter.IsAccept(){
				handler(c)
			}else{
				c.AbortWithStatusJSON(http.StatusTooManyRequests,gin.H{"message":"too many requests"})
			}
		}
	}
}
