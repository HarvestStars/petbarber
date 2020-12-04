package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/HarvestStars/petbarber/db"
	"github.com/HarvestStars/petbarber/dtos"
	"github.com/gin-gonic/gin"
)

func CreateOrderComment(c *gin.Context) {
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
	var commentReq dtos.CommentReq
	err = c.Bind(&commentReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.COMMENT_ERROR_TYPE, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}

	var comment dtos.TComment
	commentType := c.Query("comment_type")
	switch commentType {
	case "ToGroomerOrder":
		// 对groomer评论
		if userType != 1 {
			c.JSON(http.StatusBadRequest, gin.H{"code": dtos.COMMENT_CANT_CREATE_COMMENT, "msg": "Sorry", "data": "", "detail": "jwt不是门店用户, 无权评论美容师"})
			return
		}
		comment.CreatedAt = time.Now().UTC().UnixNano() / 1e6
		comment.Status = 1
		comment.FromUserID = accountID
		var requirementOrder dtos.ToRequirement
		count := 0
		db.DataBase.Model(&dtos.ToRequirement{}).Where("id = ? AND user_id = ?", commentReq.OrderID, accountID).Count(&count).First(&requirementOrder)
		if count == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"code": dtos.ORDER_NOT_EXISTS, "msg": "Sorry", "data": "", "detail": "目标requirement订单不存在"})
			return
		}

		if requirementOrder.Status != dtos.FINISHED {
			c.JSON(http.StatusBadRequest, gin.H{"code": dtos.COMMENT_CANT_CREATE_COMMENT, "msg": "Sorry", "data": "", "detail": "订单未完成, 无法评论"})
			return
		}
		var matchOrder dtos.ToMatch
		db.DataBase.Model(&dtos.ToMatch{}).Where("id = ?", requirementOrder.GroomerOrderID).First(&matchOrder)
		var groomer dtos.TuGroomer
		db.DataBase.Model(&dtos.TuGroomer{}).Where("account_id = ?", matchOrder.UserID).First(&groomer)
		comment.ToUserID = groomer.AccountID
		comment.CommentType = 1
		comment.Favor = commentReq.Favor
		comment.Content = commentReq.Content
		comment.PethouseOrderID = commentReq.OrderID
		comment.GroomerOrderID = matchOrder.ID
		db.DataBase.Model(&dtos.TComment{}).Where("pethouse_order_id = ? AND groomer_order_id = ?", commentReq.OrderID, matchOrder.ID).Count(&count)
		if count == 0 {
			// 未评论记录
			db.DataBase.Create(&comment)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"code": dtos.COMMENT_CANT_CREATE_COMMENT, "msg": "Sorry", "data": "", "detail": "不可重复评论"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": dtos.OK, "msg": "OK", "data": "", "detail": ""})

	case "ToPetHouseOrder":
		// 对pethouse评论
		if userType != 2 {
			c.JSON(http.StatusBadRequest, gin.H{"code": dtos.COMMENT_CANT_CREATE_COMMENT, "msg": "Sorry", "data": "", "detail": "jwt不是美容师用户, 无权评论门店"})
			return
		}
		comment.CreatedAt = time.Now().UTC().UnixNano() / 1e6
		comment.Status = 1
		comment.FromUserID = accountID
		var matchOrder dtos.ToMatch
		count := 0
		db.DataBase.Model(&dtos.ToMatch{}).Where("id = ? AND user_id = ?", commentReq.OrderID, accountID).Count(&count).First(&matchOrder)
		if count == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"code": dtos.ORDER_NOT_EXISTS, "msg": "Sorry", "data": "", "detail": "目标match订单不存在"})
			return
		}

		if matchOrder.Status != dtos.FINISHED {
			c.JSON(http.StatusBadRequest, gin.H{"code": dtos.COMMENT_CANT_CREATE_COMMENT, "msg": "Sorry", "data": "", "detail": "订单未完成, 无法评论"})
			return
		}
		var requirementOrder dtos.ToRequirement
		db.DataBase.Model(&dtos.ToRequirement{}).Where("id = ?", matchOrder.PethouseOrderID).First(&requirementOrder)
		var petHouse dtos.TuPethouse
		db.DataBase.Model(&dtos.TuPethouse{}).Where("account_id = ?", requirementOrder.UserID).First(&petHouse)
		comment.ToUserID = petHouse.AccountID
		comment.CommentType = 2
		comment.Favor = commentReq.Favor
		comment.Content = commentReq.Content
		comment.PethouseOrderID = requirementOrder.ID
		comment.GroomerOrderID = commentReq.OrderID
		db.DataBase.Model(&dtos.TComment{}).Where("pethouse_order_id = ? AND groomer_order_id = ?", requirementOrder.ID, commentReq.OrderID).Count(&count)
		if count == 0 {
			// 未评论记录
			db.DataBase.Create(&comment)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"code": dtos.COMMENT_CANT_CREATE_COMMENT, "msg": "Sorry", "data": "", "detail": "不可重复评论"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": dtos.OK, "msg": "OK", "data": "", "detail": ""})

	default:
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.COMMENT_ERROR_TYPE, "msg": "Sorry", "data": "", "detail": "评论者身份不明"})
	}
}

// GetComment 根据目标对象id和两份订单的订单号,返回对应的评论
func GetComment(c *gin.Context) {
	// 获取其他人对to_user_id用户的评论信息
	accountID, err := strconv.ParseUint(c.Query("to_user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.URL_ERROR, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}
	requirementOrderID, err := strconv.ParseUint(c.Query("pethouse_order_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.URL_ERROR, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}
	matchOrderID, err := strconv.ParseUint(c.Query("groomer_order_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.URL_ERROR, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}

	var comment dtos.TComment
	count := 0
	db.DataBase.Model(&dtos.TComment{}).Where("to_user_id = ? AND pethouse_order_id = ? AND groomer_order_id = ?", accountID, requirementOrderID, matchOrderID).
		Count(&count).First(&comment)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.COMMENT_CANT_READ, "msg": "Sorry", "data": "", "detail": "没有这条评论"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": dtos.OK, "msg": "OK", "data": comment, "detail": ""})
}

// GetCommentList 目标对象id所有的被评论信息
func GetCommentList(c *gin.Context) {
	// 其他人对to_user_id用户的所有评论信息
	accountID, err := strconv.ParseUint(c.Query("to_user_id"), 10, 32)
	pageSize, err := strconv.Atoi(c.Query("page_size"))
	pageIndex, err := strconv.Atoi(c.Query("page_index"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.URL_ERROR, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}

	var comments []dtos.TComment
	count := 0
	db.DataBase.Model(&dtos.TComment{}).Where("to_user_id = ?", accountID).
		Count(&count).Limit(pageSize).Offset((pageIndex - 1) * pageSize).Find(&comments)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.COMMENT_CANT_READ, "msg": "Sorry", "data": "", "detail": "该用户没有评论"})
		return
	}

	var commentList dtos.CommentResp
	commentList.List = comments
	commentList.PageInfo = dtos.PageInfo{
		TotalItems: count,
		TotalPages: count/pageSize + 1,
		PageSize:   pageSize,
		PageIndex:  pageIndex,
	}
	c.JSON(http.StatusOK, gin.H{"code": dtos.OK, "msg": "OK", "data": commentList, "detail": ""})
}
