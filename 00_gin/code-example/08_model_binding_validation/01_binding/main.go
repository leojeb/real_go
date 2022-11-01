package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func task1() {
	var r = gin.Default()
	r.GET("/log", func(c *gin.Context) {

	})
}

type login struct {
	User     string `form:"user" json:"user" xml:"user" binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

func main() {
	var r = gin.Default()
	r.POST("/loginJSON", func(c *gin.Context) {
		var json login
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if json.User != "manu" || json.Password != "123" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "未授权用户"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "logged in"})

	})

	r.POST("/loginXML", func(c *gin.Context) {
		var xml login
		if err := c.ShouldBindXML(&xml); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if xml.User != "manu" || xml.Password != "123" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "未授权用户"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "登录成功"})

	})

	r.POST("/loginFORM", func(c *gin.Context) {
		var form login
		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if form.User != "manu" || form.Password != "123" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "未授权用户"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "登录成功"})

	})

	r.Run(":8080")
}
