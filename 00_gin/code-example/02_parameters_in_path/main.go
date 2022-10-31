package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	var r = gin.Default()

	// 请求示例 : /user/john
	r.GET("/user:name", func(c *gin.Context) {
		var name = c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)

	})

	r.GET("/user/:name/*action", func(c *gin.Context) {
		var name = c.Param("name")
		var action = c.Param("action")
		msg := name + "is" + action
		c.String(http.StatusOK, msg)
	})

	r.POST("/user/:name/*action", func(c *gin.Context) {
		b := c.FullPath() == "/user/:name/*action"
		c.String(http.StatusOK, "%t", b)
	})

	r.Run(":8080")
}
