package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/HarvestStars/petbarber/db"
	"github.com/HarvestStars/petbarber/dtos"
	"github.com/gin-gonic/gin"
)

func GroomerCreateOrder(c *gin.Context) {
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
	userType := int(tokenPayload["utype"].(float64))
	accountID := uint(tokenPayload["id"].(float64))
	if userType != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_EXPECTED_PETHOUSE_TOKEN, "msg": "Sorry", "data": "", "detail": "jwt类型错误"})
		return
	}

	requirementOrderID, err := strconv.Atoi(c.Param("bizOrderID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.ORDER_BIZ_ID_WRONG, "msg": "Sorry", "data": "", "detail": "ORDER_BIZ_ID_WRONG" + err.Error()})
		return
	}
	var requirementOrder dtos.ToRequirement
	requireCount := 0
	db.DataBase.Model(&dtos.ToRequirement{}).Where("id = ?", requirementOrderID).Count(&requireCount).First(&requirementOrder)
	if requireCount == 0 {
		// 没有该需求订单
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.ORDER_BIZ_ID_WRONG, "msg": "Sorry", "data": "", "detail": "DB中无该需求订单"})
		return
	}
	if requirementOrder.Status != dtos.NEW {
		// 订单已被接单
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.ORDER_NOT_ACTIVE, "msg": "Sorry", "data": "", "detail": "该需求订单不为等待接单状态"})
		return
	}

	var matchOrder dtos.ToMatch
	matchOrder.CreatedAt = time.Now().UTC().UnixNano() / 1e6
	matchOrder.Status = dtos.RUNNING
	matchOrder.PethouseOrderID = requirementOrder.ID
	matchOrder.UserID = accountID

	// 双表事务
	tx := db.DataBase.Begin()
	tx.Model(&dtos.ToMatch{}).Create(&matchOrder)
	updatedTime := time.Now().UTC().UnixNano() / 1e6
	tx.Model(&dtos.ToRequirement{}).Where("id = ?", requirementOrder.ID).UpdateColumns(dtos.ToRequirement{
		UpdatedAt:    updatedTime,
		Status:       dtos.RUNNING,
		MatchOrderID: matchOrder.ID,
	}).First(&requirementOrder)
	tx.Commit()
	var matchResp dtos.PCMatchResp
	err = matchResp.RespTransfer(matchOrder, requirementOrder)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.ORDER_PAYMENT_DATA_MISSION, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": dtos.OK, "msg": "OK", "data": matchResp, "detail": ""})
}

func GroomerCancelOrder(c *gin.Context) {}

func GroomerGetOrderList(c *gin.Context) {}

func GroomerGetActivePethouseOrder(c *gin.Context) {}

func GroomerGetOrder(c *gin.Context) {}
