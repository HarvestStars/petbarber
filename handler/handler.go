package handler

import "github.com/gin-gonic/gin"

func Register(c *gin.Context) {
	baseInfo := &BaseInfo{}
	c.Bind(baseInfo)
	if baseInfo.Account == "123" {
		c.JSON(200, gin.H{"code": 200, "msg": "OK", "data": "test"})
	}
}
