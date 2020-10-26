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

	// for general users
	// 提交注册，修改资料
	r.POST("/api/infoload/v2/createorupdateaccount", handler.CreateOrUpdateAccount) // 上层表测试接口 上线时会移除本接口

	r.POST("/api/infoload/v2/createorupdategroomer", handler.CreateOrUpdateGroomer)
	r.POST("/api/infoload/v2/uploadgroomeridcard", handler.UploadGroomerIDCard)
	r.POST("/api/infoload/v2/uploadgroomeravatar", handler.UploadGroomerAvatar)
	r.POST("/api/infoload/v2/uploadgroomercertificate", handler.UploadGroomerCertificate)

	r.POST("/api/infoload/v2/createorupdatehouse", handler.CreateOrUpdateHouse)
	r.POST("/api/infoload/v2/uploadhouseidcard", handler.UploadHouseIDCard)
	r.POST("/api/infoload/v2/uploadhouseavatar", handler.UploadHouseAvatar)
	r.POST("/api/infoload/v2/uploadhouselicense", handler.UploadHouseLicense)

	// for super users
	// 审核, 封禁, 查阅，删除

	r.Run(setting.ServerSetting.Host)
}
