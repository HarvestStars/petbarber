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

type Server struct {
	Host string
}

var ServerSetting = &Server{}

type ImagePathConf struct {
	AvatarPath             string
	GroomerIDCardPath      string
	GroomerCertificatePath string
	HouseIDCardPath        string
	HouseEnvironmentPath   string
	HouseLicensePath       string
}

var ImagePathSetting = &ImagePathConf{}

type JwtConf struct {
	JwtKey           string
	SmsKey           string
	JwtExpireTimeSec int64
	SmsExpireTimeSec int
}

var JwtSetting = &JwtConf{}

// Setup 启动配置
func Setup() {
	cfg, err := ini.Load("./conf/my.ini")
	if err != nil {
		log.Fatalf("Fail to parse '../my.ini': %v", err)
	}

	mapTo(cfg, "mysql", MySQLSetting)
	mapTo(cfg, "server", ServerSetting)
	mapTo(cfg, "image", ImagePathSetting)
	mapTo(cfg, "keys", JwtSetting)
}

func mapTo(cfg *ini.File, section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo RedisSetting err: %v", err)
	}
}
