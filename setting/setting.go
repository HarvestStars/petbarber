package setting

import (
	"log"

	"github.com/go-ini/ini"
)

type MySQLConf struct {
	Host     string
	User     string
	PassWord string
	DataBase string
}

var MySQLSetting = &MySQLConf{}

// Setup 启动配置
func Setup() {
	cfg, err := ini.Load("./conf/my.ini")
	if err != nil {
		log.Fatalf("Fail to parse '../my.ini': %v", err)
	}

	mapTo(cfg, "mysql", MySQLSetting)
}

func mapTo(cfg *ini.File, section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo RedisSetting err: %v", err)
	}
}
