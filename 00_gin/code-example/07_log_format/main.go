package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"time"
)

func task1() {
	var r = gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.Run(":8080")
}
func task2() {
	// 自定义日志格式

	r1 := gin.New()
	r1.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 自定义格式
		return fmt.Sprintf(
			"%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	r1.GET("/log", func(c *gin.Context) {
		c.String(http.StatusOK, "logging:")
	})
	r1.Run(":9090")
}

func main() {
	//gin.DisableConsoleColor() // 日志总是无颜色
	gin.ForceConsoleColor() // 日志总是有颜色
	// 日志写到指定文件
	f, _ := os.Create("gin.log")
	//gin.DefaultWriter = io.MultiWriter(f)
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout) //同时输出到文件和控制台

	task1()
	//task2()

}
