package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RefreshAccessToken(c *gin.Context) {
	auth := c.Request.Header.Get("authorization")
	tokenStr, err := extractTokenFromAuth(auth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}
	tokenPayload, err := ParseToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}
	JwtToken, err := RefreshJwtToken(tokenPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": JwtToken, "detail": ""})
}
