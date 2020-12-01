package main

import (
	"github.com/HarvestStars/petbarber/db"
	"github.com/HarvestStars/petbarber/handler"
	"github.com/HarvestStars/petbarber/setting"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 解析配置文件
	setting.Setup()

	// 注册数据库
	db.Setup(setting.MySQLSetting.User, setting.MySQLSetting.PassWord, setting.MySQLSetting.Host, setting.MySQLSetting.DataBase)

	// 开启服务
	r := gin.Default()

	// 注册与登录
	r.GET("/api/v1/account/smscode/:phone", handler.SendSmsCode)
	r.POST("/api/v1/account/smscode/signin_or_signup", handler.SigninOrSignup)
	r.GET("/api/v1/account/profile", handler.GetUserProfile)

	// 更新jwt
	r.GET("/api/v1/jwt/refresh", handler.RefreshAccessToken)

	// 上传资料
	r.POST("/api/v1/account/uploadgroomer", handler.UploadGroomer) // 美容师信息页 非图片类信息上传
	r.POST("/api/v1/account/uploadhouse", handler.UploadHouse)     // 门店信息页 非图片类信息上传
	r.POST("/api/v1/account/uploadimage", handler.UploadImage)     // 门店与美容师 图片类信息上传

	// 订单业务
	// pethouse
	r.POST("/api/v1/order/pethouse/create/", handler.PetHouseCreateOrder)
	r.DELETE("/api/v1/order/pethouse/cancel/:orderID", handler.PetHouseCancelOrder)
	r.DELETE("/api/v1/order/pethouse/deny/:pethouseOrderID/:groomerUserID", handler.PetHouseDenyUserOrder)
	r.GET("/api/v1/order/pethouse/list", handler.PetHouseGetOrderList)
	r.GET("/api/v1/order/pethouse/getorder/:orderID", handler.PetHousGetOrder)
	r.GET("/api/v1/order/pethouse/close/:orderID", handler.PetHouseCloseOrder)
	// groomer
	r.POST("/api/v1/order/groomer/create/:bizOrderID", handler.GroomerCreateOrder)
	r.DELETE("/api/v1/order/groomer/cancel/:orderID", handler.GroomerCancelOrder)
	r.GET("/api/v1/order/groomer/list", handler.GroomerGetOrderList)
	r.GET("/api/v1/order/groomer/active", handler.GroomerGetActivePethouseOrder)
	r.GET("/api/v1/order/groomer/getorder/:orderID", handler.GroomerGetOrder)

	// 评论
	r.POST("/api/v1/comment/", handler.CreateOrderComment)

	// 支付

	r.Run(setting.ServerSetting.Host)
}
