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
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.ORDER_INTERNAL_ERROR, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}
	var requirementOrder dtos.ToRequirement

	orderType := c.Query("order_type")
	switch orderType {
	case "WCB":
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
	// 获取门店城市与区域信息
	// var owner dtos.TuPethouse
	// db.DataBase.Model(&dtos.TuPethouse{}).Where("account_id = ?", accountID).First(&owner)
	// requirementOrder.City = owner.City
	// requirementOrder.Region = owner.Region
	requirementOrder.City = petHouseOrderReq.City
	requirementOrder.Region = petHouseOrderReq.Region

	// 转换请求数据，然后记录DB, 事务
	tx := db.DataBase.Begin()
	tx.Model(&dtos.ToRequirement{}).Create(&requirementOrder)
	tx.Commit()
	var orderResp dtos.PCOrderRespOnlyCreate
	err = orderResp.RespCreateOrder(requirementOrder)
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
	accountID := uint(tokenPayload["id"].(float64))
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
	tx.Model(&dtos.ToRequirement{}).Where("id = ? AND user_id = ?", uint(orderID), accountID).Count(&orderCount).First(&requirementOrder)
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
		requirementOrderStatus := dtos.CANCELORDER
		// RUNNING 取消
		// 十分钟校验
		var matchOrder dtos.ToMatch
		tx.Model(&dtos.ToMatch{}).Where("id = ?", requirementOrder.GroomerOrderID).First(&matchOrder)
		if (matchOrder.CreatedAt/1e3 + 600) < time.Now().UTC().Unix() {
			// 超出可取消时间
			// c.JSON(http.StatusBadRequest, gin.H{"code": dtos.ORDER_CANCEL_NOT_ALLOWED, "msg": "Sorry", "data": "", "detail": "被接单已经超过10分钟"})
			// tx.Commit()
			// return

			// 超时取消, 记录违约
			requirementOrderStatus = dtos.OUTTIME
		}

		// match 和 requirement 双表联动取消
		tx.Model(&dtos.ToMatch{}).Where("id = ?", requirementOrder.GroomerOrderID).UpdateColumns(dtos.ToMatch{
			UpdatedAt: time.Now().UTC().UnixNano() / 1e6,
			Status:    dtos.CANCELORDER})

		tx.Model(&dtos.ToRequirement{}).Where("id = ?", uint(orderID)).UpdateColumns(dtos.ToRequirement{
			UpdatedAt: time.Now().UTC().UnixNano() / 1e6,
			Status:    requirementOrderStatus,
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
	accountID := uint(tokenPayload["id"].(float64))
	if userType != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_EXPECTED_PETHOUSE_TOKEN, "msg": "Sorry", "data": "", "detail": "JWT_EXPECTED_PETHOUSE_TOKEN"})
		return
	}

	petHouseOrderID, err := strconv.ParseUint(c.Param("pethouseOrderID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.ORDER_BIZ_ID_WRONG, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}

	// 启动事务
	tx := db.DataBase.Begin()
	defer tx.Commit()
	count := 0
	var requirementOrder dtos.ToRequirement
	tx.Model(&dtos.ToRequirement{}).Where("id = ? AND user_id = ?", petHouseOrderID, accountID).Count(&count).First(&requirementOrder)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.ORDER_NOT_EXISTS, "msg": "Sorry", "data": "", "detail": "requirement中无该订单"})
		return
	}
	if requirementOrder.Status != dtos.RUNNING {
		// 不在可以deny的状态
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.ORDER_CANCEL_NOT_ALLOWED, "msg": "Sorry", "data": "", "detail": "订单不在RUNNING状态"})
		return
	}

	var matchOrder dtos.ToMatch
	tx.Model(&dtos.ToMatch{}).Where("id = ?", requirementOrder.GroomerOrderID).Count(&count).First(&matchOrder)
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
		"updated_at":       time.Now().UTC().UnixNano() / 1e6,
		"status":           dtos.NEW,
		"groomer_order_id": 0})

	tx.Model(&dtos.ToMatch{}).Where("id = ?", requirementOrder.GroomerOrderID).UpdateColumns(dtos.ToMatch{
		UpdatedAt: time.Now().UTC().UnixNano() / 1e6,
		Status:    dtos.CANCELGROOMER})
	c.JSON(http.StatusOK, gin.H{"code": dtos.OK, "msg": "OK", "data": "", "detail": "成功拒绝该美容师"})
}

func PetHouseGetOrderList(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_EXPECTED_PETHOUSE_TOKEN, "msg": "Sorry", "data": "", "detail": "JWT_EXPECTED_PETHOUSE_TOKEN"})
		return
	}

	pageSize, err := strconv.Atoi(c.Query("page_size"))
	pageIndex, err := strconv.Atoi(c.Query("page_index"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.URL_ERROR, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}
	orderStatusList := c.QueryArray("order_status")

	var requirementOrders []dtos.ToRequirement
	count := 0
	db.DataBase.Model(&dtos.ToRequirement{}).Where("status in (?) AND user_id = ?", orderStatusList, accountID).Count(&count).Limit(pageSize).Offset((pageIndex - 1) * pageSize).Find(&requirementOrders)

	var listResp []dtos.PCOrderResp
	for _, order := range requirementOrders {
		var matchOrder dtos.ToMatch
		db.DataBase.Model(&dtos.ToMatch{}).Where("id = ?", order.GroomerOrderID).First(&matchOrder)
		var groomer dtos.TuGroomer
		db.DataBase.Model(&dtos.TuGroomer{}).Where("account_id = ?", matchOrder.UserID).First(&groomer)
		var orderResp dtos.PCOrderResp
		err = orderResp.RespTransfer(order, matchOrder, groomer)
		if err == nil {
			listResp = append(listResp, orderResp)
		}
	}

	var orderListResp dtos.PCOrderListResp
	orderListResp.List = listResp
	orderListResp.PageInfo = dtos.PageInfo{
		TotalItems: count,
		TotalPages: count/pageSize + 1,
		PageSize:   pageSize,
		PageIndex:  pageIndex,
	}
	c.JSON(http.StatusOK, gin.H{"code": dtos.OK, "msg": "OK", "data": orderListResp, "detail": ""})
}

func PetHousGetOrder(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_EXPECTED_PETHOUSE_TOKEN, "msg": "Sorry", "data": "", "detail": "JWT_EXPECTED_PETHOUSE_TOKEN"})
		return
	}
	orderID, err := strconv.ParseUint(c.Param("orderID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.ORDER_BIZ_ID_WRONG, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}

	var requirementOrder dtos.ToRequirement
	count := 0
	db.DataBase.Model(&dtos.ToRequirement{}).Where("id = ? AND user_id = ?", orderID, accountID).Count(&count).First(&requirementOrder)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.ORDER_NOT_EXISTS, "msg": "Sorry", "data": "", "detail": "requirement中无该订单"})
		return
	}
	var matchOrder dtos.ToMatch
	db.DataBase.Model(&dtos.ToMatch{}).Where("id = ?", requirementOrder.GroomerOrderID).First(&matchOrder)
	var groomer dtos.TuGroomer
	db.DataBase.Model(&dtos.TuGroomer{}).Where("account_id = ?", matchOrder.UserID).First(&groomer)
	var orderResp dtos.PCOrderResp
	err = orderResp.RespTransfer(requirementOrder, matchOrder, groomer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.ORDER_PAYMENT_DATA_MISSION, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": dtos.OK, "msg": "OK", "data": orderResp, "detail": ""})
}

func PetHouseCloseOrder(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_EXPECTED_PETHOUSE_TOKEN, "msg": "Sorry", "data": "", "detail": "JWT_EXPECTED_PETHOUSE_TOKEN"})
		return
	}
	orderID, err := strconv.ParseUint(c.Param("orderID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.ORDER_BIZ_ID_WRONG, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}
	totalPayment, err := strconv.ParseFloat(c.Query("total_payment"), 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.ORDER_NOT_FINISHED, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}
	var requirementOrder dtos.ToRequirement
	count := 0
	db.DataBase.Model(&dtos.ToRequirement{}).Where("id = ? AND user_id = ?", orderID, accountID).Count(&count).First(&requirementOrder)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.ORDER_NOT_EXISTS, "msg": "Sorry", "data": "", "detail": "requirement中无该订单"})
		return
	}
	if requirementOrder.Status != dtos.RUNNING {
		// 未在可完成状态
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.ORDER_NOT_FINISHED, "msg": "Sorry", "data": "", "detail": "订单不在RUNNING状态"})
		return
	}

	tx := db.DataBase.Begin()
	defer tx.Commit()
	tx.Model(&dtos.ToRequirement{}).Where("id = ?", orderID).UpdateColumns(dtos.ToRequirement{
		UpdatedAt:    time.Now().UTC().UnixNano() / 1e6,
		TotalPayment: float32(totalPayment),
		Status:       dtos.FINISHED,
	})
	tx.Model(&dtos.ToMatch{}).Where("id = ?", requirementOrder.GroomerOrderID).UpdateColumns(dtos.ToMatch{
		UpdatedAt: time.Now().UTC().UnixNano() / 1e6,
		Status:    dtos.FINISHED,
	})
	c.JSON(http.StatusOK, gin.H{"code": dtos.OK, "msg": "OK", "data": "", "detail": "订单完成"})
}
