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
	r.POST("/api/stoserver/v2/admin/backend/register", handler.UpdateAccount)

	// for super users
	// 审核, 封禁, 查阅，删除

	r.Run(setting.ServerSetting.Host)
}
