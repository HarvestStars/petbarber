package handler

import "github.com/gin-gonic/gin"

// ----------------------------------------------------- 普通用户 -----------------------------------------------------
// 综合信息接口
// GetAccount 获取账户信息:昵称，头像和身份证图片路径等
func GetAccount(c *gin.Context) {}

// UpdateAccount 更新账户信息:昵称，电话，身份证图片正反面等信息
func UpdateAccount(c *gin.Context) {}

// 美容师专业信息接口
// GetGroomer 获取宠物美容师专业信息，等级，星级，证书照片正反面等
func GetGroomer(c *gin.Context) {}

// UpdateGroomer 更新宠物美容师专业信息
func UpdateGroomer(c *gin.Context) {}

// 门店信息接口
// GetHouse 获取门店营业执照信息
func GetHouse(c *gin.Context) {}

// UpdateHouse 更新门店营业执照信息
func UpdateHouse(c *gin.Context) {}

// ----------------------------------------------------- 超级管理员 -----------------------------------------------------
