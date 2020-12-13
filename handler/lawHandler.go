package handler

import (
	"net/http"

	"github.com/HarvestStars/petbarber/db"
	"github.com/HarvestStars/petbarber/dtos"
	"github.com/gin-gonic/gin"
)

func GetAgreement(c *gin.Context) {
	var law dtos.CLaw
	db.DataBase.Where("id = ?", 1).First(&law)
	c.JSON(http.StatusOK, gin.H{"code": dtos.OK, "msg": "Sorry", "data": law, "detail": ""})
	return
}
