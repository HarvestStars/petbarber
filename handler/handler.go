package handler

import (
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/HarvestStars/petbarber/db"
	"github.com/HarvestStars/petbarber/setting"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

const uploadMaxBytes int64 = 1024 * 1024 // 1M
// ----------------------------------------------------- 普通用户 -----------------------------------------------------
// CreateOrUpdateAccount 更新账户信息:昵称，电话等文字信息
func CreateOrUpdateAccount(c *gin.Context) {
	var account db.AccountInfo
	err := c.Bind(&account)
	if err != nil {
		log.Print(err.Error())
		return
	}
	count := 0
	db.DataBase.Model(&db.AccountInfo{}).Where("account = ?", account.Account).Count(&count)
	if count != 0 {
		// exist
		db.DataBase.Model(&db.AccountInfo{}).Where("account = ?", account.Account).Update(&account)
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "update done."})
	} else {
		// create
		db.DataBase.Create(&account)
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "created."})
	}
}

// 美容师信息录入接口
// CreateOrUpdateGroomer 更新美容师信息:昵称，电话等文字信息
func CreateOrUpdateGroomer(c *gin.Context) {
	var groomer db.PetGroomer
	err := c.Bind(&groomer)
	if err != nil {
		log.Print(err.Error())
		return
	}
	count := 0
	db.DataBase.Model(&db.PetGroomer{}).Where("account_id = ?", groomer.AccountID).Count(&count)
	if count != 0 {
		// exist
		db.DataBase.Model(&db.PetGroomer{}).Where("account_id = ?", groomer.AccountID).Update(&groomer)
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "update done."})
	} else {
		// create
		db.DataBase.Create(&groomer)
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "created."})
	}
}

// UploadGroomerIDCard 上传美容师身份证正反面照片
func UploadGroomerIDCard(c *gin.Context) {
	AccountIDStr := c.Request.PostFormValue("account_id")
	AccountID, _ := strconv.ParseUint(AccountIDStr, 10, 32)
	var groomer db.PetGroomer
	db.DataBase.Where("account_id = ?", AccountID).First(&groomer)

	// FormFile方法会读取参数“upload”后面的文件名，返回值是一个File指针，和一个FileHeader指针，和一个err错误。
	fileFront, headerFront, err := c.Request.FormFile("id-front")
	fileBack, headerBack, err := c.Request.FormFile("id-back")
	IDCardNumber := c.Request.PostFormValue("idcardnumber")
	if err != nil {
		c.JSON(200, gin.H{"code": 0, "data": "错误请求", "error": err.Error()})
		return
	}

	fileNameFront, err := transferImage(fileFront, headerFront, setting.ImagePathSetting.GroomerIDCardPath)
	fileNameBack, err := transferImage(fileBack, headerBack, setting.ImagePathSetting.GroomerIDCardPath)
	if err != nil {
		c.JSON(200, gin.H{"code": 0, "data": "图片大小不能超过3M", "error": nil})
	}

	db.DataBase.Model(&groomer).Update(db.PetGroomer{
		IDCardNumber: IDCardNumber,
		IDCardFront:  setting.ImagePathSetting.GroomerIDCardPath + fileNameFront,
		IDCardBack:   setting.ImagePathSetting.GroomerIDCardPath + fileNameBack})
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "update done."})
}

// UploadGroomerAvatar 上传美容师头像
func UploadGroomerAvatar(c *gin.Context) {
	AccountIDStr := c.Request.PostFormValue("account_id")
	AccountID, _ := strconv.ParseUint(AccountIDStr, 10, 32)
	var groomer db.PetGroomer
	db.DataBase.Where("account_id = ?", AccountID).First(&groomer)

	// FormFile方法会读取参数“upload”后面的文件名，返回值是一个File指针，和一个FileHeader指针，和一个err错误。
	fileFront, headerFront, err := c.Request.FormFile("avatar")
	if err != nil {
		c.JSON(200, gin.H{"code": 0, "data": "错误请求", "error": err.Error()})
		return
	}

	fileNameFront, err := transferImage(fileFront, headerFront, setting.ImagePathSetting.AvatarPath)
	if err != nil {
		c.JSON(200, gin.H{"code": 0, "data": "图片大小不能超过3M", "error": nil})
	}

	db.DataBase.Model(&groomer).Update(db.PetGroomer{
		Avatar: setting.ImagePathSetting.AvatarPath + fileNameFront})
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "update done."})
}

// UploadGroomerCertificate 上传门美容师资格证
func UploadGroomerCertificate(c *gin.Context) {
	AccountIDStr := c.Request.PostFormValue("account_id")
	AccountID, _ := strconv.ParseUint(AccountIDStr, 10, 32)
	var groomer db.PetGroomer
	db.DataBase.Where("account_id = ?", AccountID).First(&groomer)

	// FormFile方法会读取参数“upload”后面的文件名，返回值是一个File指针，和一个FileHeader指针，和一个err错误。
	fileFront, headerFront, err := c.Request.FormFile("certifi-front")
	fileBack, headerBack, err := c.Request.FormFile("certifi-back")
	if err != nil {
		c.JSON(200, gin.H{"code": 0, "data": "错误请求", "error": err.Error()})
		return
	}

	fileNameFront, err := transferImage(fileFront, headerFront, setting.ImagePathSetting.GroomerCertificatePath)
	fileNameBack, err := transferImage(fileBack, headerBack, setting.ImagePathSetting.GroomerCertificatePath)
	if err != nil {
		c.JSON(200, gin.H{"code": 0, "data": "图片大小不能超过3M", "error": nil})
	}

	db.DataBase.Model(&groomer).Update(db.PetGroomer{
		CertificateFront: setting.ImagePathSetting.GroomerCertificatePath + fileNameFront,
		CertificateBack:  setting.ImagePathSetting.GroomerCertificatePath + fileNameBack})
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "update done."})
}

// 门店信息录入接口
// CreateOrUpdateHouse 更新美容师信息:昵称，电话等文字信息
func CreateOrUpdateHouse(c *gin.Context) {
	var house db.PetHouse
	err := c.Bind(&house)
	if err != nil {
		log.Print(err.Error())
		return
	}
	count := 0
	db.DataBase.Model(&db.PetHouse{}).Where("account_id = ?", house.AccountID).Count(&count)
	if count != 0 {
		// exist
		db.DataBase.Model(&db.PetHouse{}).Where("account_id = ?", house.AccountID).Update(&house)
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "update done."})
	} else {
		// create
		db.DataBase.Create(&house)
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "created."})
	}
}

// UploadHouseIDCard 上传门店主身份证正反面照片
func UploadHouseIDCard(c *gin.Context) {
	AccountIDStr := c.Request.PostFormValue("account_id")
	AccountID, _ := strconv.ParseUint(AccountIDStr, 10, 32)
	var house db.PetHouse
	db.DataBase.Where("account_id = ?", AccountID).First(&house)

	// FormFile方法会读取参数“upload”后面的文件名，返回值是一个File指针，和一个FileHeader指针，和一个err错误。
	fileFront, headerFront, err := c.Request.FormFile("id-front")
	fileBack, headerBack, err := c.Request.FormFile("id-back")
	IDCardNumber := c.Request.PostFormValue("idcardnumber")
	if err != nil {
		c.JSON(200, gin.H{"code": 0, "data": "错误请求", "error": err.Error()})
		return
	}

	fileNameFront, err := transferImage(fileFront, headerFront, setting.ImagePathSetting.HouseIDCardPath)
	fileNameBack, err := transferImage(fileBack, headerBack, setting.ImagePathSetting.HouseIDCardPath)
	if err != nil {
		c.JSON(200, gin.H{"code": 0, "data": "图片大小不能超过3M", "error": nil})
	}

	db.DataBase.Model(&house).Update(db.PetHouse{
		IDCardNumber: IDCardNumber,
		IDCardFront:  setting.ImagePathSetting.HouseIDCardPath + fileNameFront,
		IDCardBack:   setting.ImagePathSetting.HouseIDCardPath + fileNameBack})
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "update done."})
}

// UploadHouseAvatar 上传门店主头像
func UploadHouseAvatar(c *gin.Context) {
	AccountIDStr := c.Request.PostFormValue("account_id")
	AccountID, _ := strconv.ParseUint(AccountIDStr, 10, 32)
	var house db.PetHouse
	db.DataBase.Where("account_id = ?", AccountID).First(&house)

	// FormFile方法会读取参数“upload”后面的文件名，返回值是一个File指针，和一个FileHeader指针，和一个err错误。
	fileFront, headerFront, err := c.Request.FormFile("avatar")
	if err != nil {
		c.JSON(200, gin.H{"code": 0, "data": "错误请求", "error": err.Error()})
		return
	}

	fileNameFront, err := transferImage(fileFront, headerFront, setting.ImagePathSetting.AvatarPath)
	if err != nil {
		c.JSON(200, gin.H{"code": 0, "data": "图片大小不能超过3M", "error": nil})
	}

	db.DataBase.Model(&house).Update(db.PetHouse{
		Avatar: setting.ImagePathSetting.AvatarPath + fileNameFront})
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "update done."})
}

// UploadHouseLicense 上传门店执照, 环境
func UploadHouseLicense(c *gin.Context) {
	AccountIDStr := c.Request.PostFormValue("account_id")
	AccountID, _ := strconv.ParseUint(AccountIDStr, 10, 32)
	var house db.PetHouse
	db.DataBase.Where("account_id = ?", AccountID).First(&house)

	fileEnvFront, headerEnvFront, err := c.Request.FormFile("environment-front")
	fileEnvIn, headerEnvIn, err := c.Request.FormFile("environment-inside")
	fileFront, headerFront, err := c.Request.FormFile("license-front")
	if err != nil {
		c.JSON(200, gin.H{"code": 0, "data": "错误请求", "error": err.Error()})
		return
	}

	fileNameEnvFront, err := transferImage(fileEnvFront, headerEnvFront, setting.ImagePathSetting.HouseEnvironmentPath)
	fileNameEnvIn, err := transferImage(fileEnvIn, headerEnvIn, setting.ImagePathSetting.HouseEnvironmentPath)
	fileNameFront, err := transferImage(fileFront, headerFront, setting.ImagePathSetting.HouseLicensePath)
	if err != nil {
		c.JSON(200, gin.H{"code": 0, "data": "图片大小不能超过3M", "error": nil})
	}

	db.DataBase.Model(&house).Update(db.PetHouse{
		EnvironmentFront:  setting.ImagePathSetting.HouseEnvironmentPath + fileNameEnvFront,
		EnvironmentInside: setting.ImagePathSetting.HouseEnvironmentPath + fileNameEnvIn,
		License:           setting.ImagePathSetting.HouseLicensePath + fileNameFront})
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "update done."})
}

func transferImage(file multipart.File, header *multipart.FileHeader, rootPath string) (string, error) {
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
	_, err := os.Stat(rootPath)
	if err != nil {
		if os.IsExist(err) {
			// 文件夹存在
		} else {
			err = os.Mkdir(rootPath, os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	out, err := os.Create(rootPath + fileName)
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

// ----------------------------------------------------- 超级管理员 -----------------------------------------------------
// 人工审核接口
