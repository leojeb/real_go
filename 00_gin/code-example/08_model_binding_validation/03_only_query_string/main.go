package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type person struct {
	Name    string `form:"name"`
	Address string `form:"address""`
}

func main() {
	var r = gin.Default()
	r.Any("/startPage", startPage)
	r.Run(":8080")
}

func startPage(c *gin.Context) {
	var p1 person
	if err := c.ShouldBindQuery(&p1); err != nil {
		c.String(http.StatusBadRequest, err.Error()) // 不会走这里, 因为没有任何字段要求
	} else {
		log.Println("---only bind query string---")
		log.Println(p1.Name)
		log.Println(p1.Address)
	}
	c.String(http.StatusOK, "验证成功")

}
