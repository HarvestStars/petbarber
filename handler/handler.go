package handler

import (
	"errors"
	"io"
	"log"
	"mime/multipart"
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

// CreateOrUpdateAccount 更新账户信息:昵称，电话等文字信息
func CreateOrUpdateAccount(c *gin.Context) {
	var account db.AccountInfo
	err := c.Bind(&account)
	if err != nil {
		log.Print(err.Error())
		return
	}
	count := 0
	db.DataBase.Model(&db.AccountInfo{}).Where("tel = ?", account.Tel).Count(&count)
	if count != 0 {
		// exist
		db.DataBase.Model(&db.AccountInfo{}).Where("tel = ?", account.Tel).Update(&account)
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "update done."})
	} else {
		// create
		db.DataBase.Create(&account)
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "created."})
	}
}

// UploadIDCard 上传身份证正反面照片
func UploadIDCard(c *gin.Context) {
	tel := c.Request.PostFormValue("tel")
	var account db.AccountInfo
	db.DataBase.Where("tel = ?", tel).First(&account)

	// FormFile方法会读取参数“upload”后面的文件名，返回值是一个File指针，和一个FileHeader指针，和一个err错误。
	fileFront, headerFront, err := c.Request.FormFile("front")
	fileBack, headerBack, err := c.Request.FormFile("back")
	IDCardNumber := c.Request.PostFormValue("idcardnumber")
	if err != nil {
		c.JSON(200, gin.H{"code": 0, "data": "错误请求", "error": err.Error()})
		return
	}

	fileNameFront, err := transferImage(fileFront, headerFront)
	fileNameBack, err := transferImage(fileBack, headerBack)
	if err != nil {
		c.JSON(200, gin.H{"code": 0, "data": "图片大小不能超过3M", "error": nil})
	}

	db.DataBase.Model(&account).Update(db.AccountInfo{IDCardNumber: IDCardNumber, IDCardFront: setting.ImagePathSetting.IDCardPath + fileNameFront, IDCardBack: setting.ImagePathSetting.IDCardPath + fileNameBack})
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "update done."})
}

func transferImage(file multipart.File, header *multipart.FileHeader) (string, error) {
	// header调用Filename方法，就可以得到文件名
	fileName := header.Filename
	filesuffix := path.Ext(fileName)
	u1, _ := uuid.NewV4()
	fileName = u1.String()
	fileName += filesuffix
	if header.Size > 3*uploadMaxBytes {
		return "", errors.New("over size")
	}

	// 创建一个文件，文件名为filename，这里的返回值out也是一个File指针
	_, err := os.Stat(setting.ImagePathSetting.IDCardPath)
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

	out, err := os.Create(setting.ImagePathSetting.IDCardPath + fileName)
	if err != nil {
		log.Fatal(err)
	}

	defer out.Close()

	// 将file的内容拷贝到out
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}
	return fileName, nil
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
