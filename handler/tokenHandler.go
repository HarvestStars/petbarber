package handler

import (
	"net/http"

	"github.com/HarvestStars/petbarber/dtos"
	"github.com/gin-gonic/gin"
)

func RefreshAccessToken(c *gin.Context) {
	auth := c.Request.Header.Get("authorization")
	tokenStr, err := extractTokenFromAuth(auth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_TYPE_WRONG, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}
	tokenPayload, err := ParseToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_VERIFY_RESULT_BAD_TOKEN, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}
	JwtToken, err := RefreshJwtToken(tokenPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_CREATE_WRONG, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": dtos.OK, "msg": "OK", "data": JwtToken, "detail": ""})
}
