package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func login(c *gin.Context) {
	//b := c.FullPath() == "/user/:name/*action"
	var msg string
	if strings.Contains(c.FullPath(), "/v1") {
		msg = "v1"
	} else if strings.Contains(c.FullPath(), "/v2") {
		msg = "v2"
	}
	c.String(http.StatusOK, "登录成功"+msg)
}

func submit(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{
		"a": "b",
	})
}

func read(c *gin.Context) {
	c.String(http.StatusOK, "已读不回")
}

func main() {
	var r = gin.Default()
	var v1 = r.Group("/v1")

	{
		v1.GET("/login", login)
		v1.POST("/submit", submit)
		v1.POST("/read", read)
	}

	var v2 = r.Group("/v2")
	{
		v2.GET("/login", login)
		v2.POST("/submit", submit)
		v2.POST("/read", read)
	}

	r.Run(":8080")

}
