package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Set("example", "123567")
		// before request
		c.Next()
		// after request, 相当于响应前的中间件
		latency := time.Since(t)
		log.Println(latency)

		// 获取将要发送的status
		status := c.Writer.Status()
		log.Println(status)
	}
}
func main() {
	var r = gin.New()
	r.Use(Logger())

	r.GET("/test", func(c *gin.Context) {
		var example = c.MustGet("example")
		log.Println(example)
		//c.String(http.StatusOK, "执行完成")
	})

	r.Run(":8080")
}
