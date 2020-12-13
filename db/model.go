package db

import (
	"fmt"
	"time"

	"github.com/HarvestStars/petbarber/dtos"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// DataBase 是orm实例
var DataBase *gorm.DB

// Setup 启动mysql配置
func Setup(user string, pwd string, host string, db string) {
	url := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=10s", user, pwd, host, db)
	var err error
	DataBase, err = gorm.Open("mysql", url)
	if err != nil {
		panic("failed to connect database")
	}
	DataBase.DB().SetConnMaxLifetime(2 * 3600 * time.Second) // 2小时空闲链接超时
	DataBase.SingularTable(true)                             // 可以取消表名的复数形式，使得表名和结构体名称一致
	DataBase.AutoMigrate(&dtos.TComment{}, &dtos.ToMatch{}, &dtos.ToRequirement{}, &dtos.TuAccount{}, &dtos.TuGroomer{}, &dtos.TuPethouse{}, &dtos.CLaw{})
}
