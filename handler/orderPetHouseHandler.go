package handler

import (
	"net/http"

	"github.com/HarvestStars/petbarber/db"
	"github.com/HarvestStars/petbarber/dtos"
	"github.com/gin-gonic/gin"
)

func PetHouseCreateOrder(c *gin.Context) {
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
	if userType != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_EXPECTED_PETHOUSE_TOKEN, "msg": "Sorry", "data": "", "detail": "jwt类型错误"})
		return
	}
	var petHouseOrderReq dtos.CreatePetHousePCOrderReq
	err = c.Bind(&petHouseOrderReq)
	var requirementOrder dtos.ToRequirement

	orderType := c.Query("order_type")
	switch orderType {
	case "WCB":
		requirementOrder.UserID = accountID
		requirementOrder.CreatedAt = petHouseOrderReq.RequestedAt
		requirementOrder.StartedAt = petHouseOrderReq.StartedAt
		requirementOrder.FinishedAt = petHouseOrderReq.FinishedAt
		requirementOrder.ServiceBits = dtos.ToServiceBits(petHouseOrderReq.ServiceItems)
		requirementOrder.ServiceItemsDesc = dtos.ToServiceDesc(petHouseOrderReq.ServiceItems)

		requirementOrder.Basic = petHouseOrderReq.Basic
		requirementOrder.Commission = petHouseOrderReq.Commission
		payModeInt, err := dtos.ToPayMode(petHouseOrderReq.Basic, petHouseOrderReq.Commission)
		payModeDesc, err := dtos.ToPayModeDesc(petHouseOrderReq.Basic, petHouseOrderReq.Commission)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": dtos.ORDER_PAYMENT_DATA_MISSION, "msg": "Sorry", "data": "", "detail": err.Error()})
			return
		}
		requirementOrder.PayMode = payModeInt
		if payModeInt == 1 {
			requirementOrder.TotalPayment = petHouseOrderReq.Basic
		}
		requirementOrder.PayModeDesc = payModeDesc
		requirementOrder.OrderType = dtos.WCB
		requirementOrder.Status = dtos.NEW
		requirementOrder.UserID = accountID

	case "WalkTheDog":
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.ORDER_UNKNOWN_ORDER_TYPE, "msg": "Sorry", "data": "", "detail": "该业务尚未开放"})
		return
	case "PickUp":
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.ORDER_UNKNOWN_ORDER_TYPE, "msg": "Sorry", "data": "", "detail": "该业务尚未开放"})
		return
	default:
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.ORDER_UNKNOWN_ORDER_TYPE, "msg": "Sorry", "data": "", "detail": "无业务计划"})
		return
	}
	// 转换请求数据，然后记录DB, 事务
	var matchOrder dtos.ToMatch
	var groomer dtos.TuGroomer
	tx := db.DataBase.Begin()
	tx.Model(&dtos.ToRequirement{}).Create(&requirementOrder)
	//tx.Model(&dtos.ToMatch{}).Where("pethouse_order_id = ?", requirementOrder.ID).First(&matchOrder)
	//tx.Model(&dtos.TuGroomer{}).Where("pethouse_order_id = ?", requirementOrder.ID).First(&matchOrder)
	tx.Commit()
	var orderResp dtos.PCOrderResp

	err = orderResp.RespTransfer(requirementOrder, matchOrder, groomer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.ORDER_PAYMENT_DATA_MISSION, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": dtos.OK, "msg": "OK", "data": orderResp, "detail": ""})
}

func PetHouseCancelOrder(c *gin.Context) {}

func PetHouseDenyUserOrder(c *gin.Context) {}

func PetHouseGetOrderList(c *gin.Context) {}

func PetHousGetOrder(c *gin.Context) {}

func PetHouseCloseOrder(c *gin.Context) {}
