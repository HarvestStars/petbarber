package handler

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/HarvestStars/petbarber/db"
	"github.com/HarvestStars/petbarber/setting"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

const uploadMaxBytes int64 = 1024 * 1024 // 1M
// ----------------------------------------------------- 普通用户 -----------------------------------------------------
// 综合信息接口
// GetAccount 获取账户信息:昵称，头像和身份证图片路径等
func GetAccount(c *gin.Context) {
	tel := c.Query("tel")
	var account db.AccountInfo
	db.DataBase.Where("tel = ?", tel).First(&account)
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": account})
}

// UpdateAccount 更新账户信息:昵称，电话等文字信息
func UpdateAccount(c *gin.Context) {
	var account db.AccountInfo
	err := c.Bind(&account)
	if err != nil {
		log.Print(err.Error())
		return
	}
	db.DataBase.Model(&db.AccountInfo{}).Where("tel = ?", account.Tel).Update("usertype", 11)
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "update done."})
}

// UploadIDCard 上传身份证正反面照片
func UploadIDCard(c *gin.Context) {
	tel := c.Request.PostFormValue("tel")
	var account db.AccountInfo
	db.DataBase.Where("tel = ?", tel).First(&account)

	// FormFile方法会读取参数“upload”后面的文件名，返回值是一个File指针，和一个FileHeader指针，和一个err错误。
	file, header, err := c.Request.FormFile("front")
	if err != nil {
		c.JSON(200, gin.H{"code": 0, "data": "错误请求", "error": err.Error()})
		return
	}

	// header调用Filename方法，就可以得到文件名
	filename := header.Filename
	filesuffix := path.Ext(filename)
	u1, _ := uuid.NewV4()
	filename = u1.String()
	filename += filesuffix
	if header.Size > 3*uploadMaxBytes {
		c.JSON(200, gin.H{"code": 0, "data": "图片大小不能超过3M", "error": nil})
		return
	}

	// 创建一个文件，文件名为filename，这里的返回值out也是一个File指针
	_, err = os.Stat(setting.ImagePathSetting.IDCardPath)
	if err != nil {
		if os.IsExist(err) {
			// 文件夹存在
		} else {
			err = os.Mkdir(setting.ImagePathSetting.IDCardPath, os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	out, err := os.Create(setting.ImagePathSetting.IDCardPath + filename)
	if err != nil {
		log.Fatal(err)
	}

	defer out.Close()

	// 将file的内容拷贝到out
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}

	account.IDCardFront = setting.ImagePathSetting.IDCardPath + filename
	db.DataBase.Model(&account).Update("id_card_front", setting.ImagePathSetting.IDCardPath+filename)
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "update done."})
}

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
// 人工审核接口
