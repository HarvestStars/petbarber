package handler

import (
	"net/http"
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
	case "CommentToPetGroomerOrder":
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
		var matchOrder dtos.ToMatch
		db.DataBase.Model(&dtos.ToMatch{}).Where("id = ?", requirementOrder.MatchOrderID).First(&matchOrder)
		var groomer dtos.TuGroomer
		db.DataBase.Model(&dtos.TuGroomer{}).Where("account_id = ?", matchOrder.UserID).First(&groomer)
		comment.ToUserID = groomer.AccountID
		comment.CommentType = 2
		comment.Favor = commentReq.Favor
		comment.Content = commentReq.Content
		db.DataBase.Create(&comment)
		c.JSON(http.StatusOK, gin.H{"code": dtos.OK, "msg": "OK", "data": "", "detail": ""})

	case "CommentToPetHouseOrder":
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
		var requirementOrder dtos.ToRequirement
		db.DataBase.Model(&dtos.ToRequirement{}).Where("id = ?", matchOrder.PethouseOrderID).First(&requirementOrder)
		var petHouse dtos.TuPethouse
		db.DataBase.Model(&dtos.TuPethouse{}).Where("account_id = ?", requirementOrder.UserID).First(&petHouse)
		comment.ToUserID = petHouse.AccountID
		comment.CommentType = 2
		comment.Favor = commentReq.Favor
		comment.Content = commentReq.Content
		db.DataBase.Create(&comment)
		c.JSON(http.StatusOK, gin.H{"code": dtos.OK, "msg": "OK", "data": "", "detail": ""})

	default:
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.COMMENT_ERROR_TYPE, "msg": "Sorry", "data": "", "detail": "评论者身份不明"})
	}
}
