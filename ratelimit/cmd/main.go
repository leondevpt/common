package main

import (
	"github.com/leondevpt/common/ratelimit"
	"net/http"
)
import "github.com/gin-gonic/gin"

func BusinessLogic(c *gin.Context)  {
	  c.JSON(http.StatusOK, gin.H{
		"code": 0 ,
		"message": "ok",
	})
}

func main()  {
	rateLimit := ratelimit.NewTokenBuket(10,1)

	rateLimit = rateLimit
	router := gin.New()
	router.GET("/", ratelimit.RateLimiter(10,1)(ratelimit.ParamLimiter(5, 1, "name")(BusinessLogic)))
	router.Run(":8080")
}
