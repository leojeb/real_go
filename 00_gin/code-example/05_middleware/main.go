package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 按组分配middleware

func main() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	//r.GET("/benchmark", MyBenchLogger(), benchEndpoint)

	r.GET("/login", LoginMiddleWare, Login)

	authorized := r.Group("/")
	{
		authorized.POST("/login")
	}

	r.Run(":8080")
}

func Login(c *gin.Context) {
	//b := c.FullPath() == "/user/:name/*action"
	var msg string
	println(c.Keys)
	fmt.Println(c.Keys)
	c.String(http.StatusOK, "登录成功"+msg)
}

func LoginMiddleWare(c *gin.Context) {
	println(c.Keys)
	fmt.Println(c.Keys)
	c.Keys = map[string]any{
		"a": "c",
	}
	fmt.Println(c.Keys)
	c.Keys["a"] = "b"
	println(c.Keys)
	fmt.Println(c.Keys)
	c.String(http.StatusOK, "中间件执行")
}
