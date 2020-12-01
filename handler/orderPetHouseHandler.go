package handler

import (
	"net/http"
	"strconv"
	"time"

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
		requirementOrder.CreatedAt = time.Now().UTC().UnixNano() / 1e6
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

func PetHouseCancelOrder(c *gin.Context) {
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
	if userType != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_EXPECTED_PETHOUSE_TOKEN, "msg": "Sorry", "data": "", "detail": "JWT_EXPECTED_PETHOUSE_TOKEN"})
		return
	}
	orderIDStr := c.Param("orderID")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.ORDER_BIZ_ID_WRONG, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}
	orderCount := 0
	var requirementOrder dtos.ToRequirement
	tx := db.DataBase.Begin()
	tx.Model(&dtos.ToRequirement{}).Where("id = ?", uint(orderID)).Count(&orderCount).First(&requirementOrder)
	if orderCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.ORDER_NOT_EXISTS, "msg": "Sorry", "data": "", "detail": "目标requirement订单不存在"})
		return
	}

	switch requirementOrder.Status {
	case dtos.NEW:
		// NEW 取消
		// 直接取消requirement表
		tx.Model(&dtos.ToRequirement{}).Where("id = ?", uint(orderID)).UpdateColumns(dtos.ToRequirement{
			UpdatedAt: time.Now().UTC().UnixNano() / 1e6,
			Status:    dtos.CANCELORDER,
		})
		tx.Commit()
		c.JSON(http.StatusOK, gin.H{"code": dtos.OK, "msg": "OK", "data": "", "detail": "未被接单取消"})

	case dtos.RUNNING:
		// RUNNING 取消
		// 十分钟校验
		var matchOrder dtos.ToMatch
		tx.Model(&dtos.ToMatch{}).Where("id = ?", requirementOrder.MatchOrderID).First(&matchOrder)
		if (matchOrder.CreatedAt/1e3 + 600) < time.Now().UTC().Unix() {
			// 超出可取消时间
			c.JSON(http.StatusBadRequest, gin.H{"code": dtos.ORDER_CANCEL_NOT_ALLOWED, "msg": "Sorry", "data": "", "detail": "被接单已经超过10分钟"})
			tx.Commit()
			return
		}

		// match 和 requirement 双表联动取消
		tx.Model(&dtos.ToMatch{}).Where("id = ?", requirementOrder.MatchOrderID).UpdateColumns(dtos.ToMatch{
			UpdatedAt: time.Now().UTC().UnixNano() / 1e6,
			Status:    dtos.CANCELORDER})

		tx.Model(&dtos.ToRequirement{}).Where("id = ?", uint(orderID)).UpdateColumns(dtos.ToRequirement{
			UpdatedAt: time.Now().UTC().UnixNano() / 1e6,
			Status:    dtos.CANCELORDER,
		})
		tx.Commit()
		c.JSON(http.StatusOK, gin.H{"code": dtos.OK, "msg": "OK", "data": "", "detail": "10分钟内正常取消"})

	default:
		tx.Commit()
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.ORDER_CANCEL_NOT_ALLOWED, "msg": "Sorry", "data": "", "detail": "订单不为NEW或者RUNNING"})
	}
}

func PetHouseDenyUserOrder(c *gin.Context) {
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
	if userType != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_EXPECTED_PETHOUSE_TOKEN, "msg": "Sorry", "data": "", "detail": "JWT_EXPECTED_PETHOUSE_TOKEN"})
		return
	}

	petHouseOrderID, err := strconv.ParseUint(c.Param("pethouseOrderID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.ORDER_BIZ_ID_WRONG, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}
	// groomerID, err := strconv.ParseUint(c.Param("groomerUserID"), 10, 32)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"code": dtos.ORDER_GROOMER_ID_WRONG, "msg": "Sorry", "data": "", "detail": err.Error()})
	// 	return
	// }

	// 启动事务
	tx := db.DataBase.Begin()
	defer tx.Commit()
	count := 0
	var requirementOrder dtos.ToRequirement
	tx.Model(&dtos.ToRequirement{}).Where("id = ?", petHouseOrderID).Count(&count).First(&requirementOrder)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.ORDER_BIZ_ID_WRONG, "msg": "Sorry", "data": "", "detail": "requirement中无该订单"})
		return
	}
	if requirementOrder.Status != dtos.RUNNING {
		// 不在可以deny的状态
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.ORDER_CANCEL_NOT_ALLOWED, "msg": "Sorry", "data": "", "detail": "订单不在RUNNING状态"})
		return
	}

	var matchOrder dtos.ToMatch
	tx.Model(&dtos.ToMatch{}).Where("id = ?", requirementOrder.MatchOrderID).Count(&count).First(&matchOrder)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.ORDER_BIZ_ID_WRONG, "msg": "Sorry", "data": "", "detail": "match中无该订单"})
		return
	}
	if (matchOrder.CreatedAt/1e3 + 600) < time.Now().UTC().Unix() {
		// 已被接单超过10分钟
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.ORDER_CANCEL_NOT_ALLOWED, "msg": "Sorry", "data": "", "detail": "被该美容师接单已经超过10分钟"})
		return
	}
	tx.Model(&dtos.ToRequirement{}).Where("id = ?", petHouseOrderID).UpdateColumns(map[string]interface{}{
		"updated_at":     time.Now().UTC().UnixNano() / 1e6,
		"status":         dtos.NEW,
		"match_order_id": 0})

	tx.Model(&dtos.ToMatch{}).Where("id = ?", requirementOrder.MatchOrderID).UpdateColumns(dtos.ToMatch{
		UpdatedAt: time.Now().UTC().UnixNano() / 1e6,
		Status:    dtos.CANCELGROOMER})
	c.JSON(http.StatusOK, gin.H{"code": dtos.OK, "msg": "OK", "data": "", "detail": "成功拒绝该美容师"})
}

func PetHouseGetOrderList(c *gin.Context) {}

func PetHousGetOrder(c *gin.Context) {}

func PetHouseCloseOrder(c *gin.Context) {}
