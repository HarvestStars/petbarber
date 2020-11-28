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
	r.GET("/api/v1/account/smscode/signin_or_signup", handler.SigninOrSignup)
	// 订单业务

	// 上传资料
	r.POST("/api/v1/account/uploadgroomer", handler.UploadGroomer) // 美容师信息页 非图片类信息上传
	r.POST("/api/v1/account/uploadhouse", handler.UploadHouse)     // 门店信息页 非图片类信息上传
	r.POST("/api/v1/account/uploadimage", handler.UploadImage)     // 门店与美容师 图片类信息上传
	r.Run(setting.ServerSetting.Host)
}
