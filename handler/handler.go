package handler

import (
	"errors"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/HarvestStars/petbarber/db"
	"github.com/HarvestStars/petbarber/setting"
	"github.com/gin-gonic/gin"
)

const uploadMaxBytes int64 = 1024 * 1024 // 1M

// ----------------------------------------------------- 普通用户 -----------------------------------------------------
// RegistUpdateAccount 更新账户信息:昵称，电话等文字信息
func RegistUpdateAccount(c *gin.Context) {
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
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "更新成功"})
	} else {
		// create
		db.DataBase.Create(&account)
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "创建成功"})
	}
}

// UploadGroomer 更新美容师非图片信息:昵称，电话等文字信息
func UploadGroomer(c *gin.Context) {
	auth := c.Request.Header.Get("authorization")
	tokenStr, err := extractTokenFromAuth(auth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": err.Error()})
		return
	}
	tokenPayload, err := ParseToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": err.Error()})
		return
	}
	accountIDStr := tokenPayload["account_id"].(string)
	accountID, _ := strconv.ParseUint(accountIDStr, 10, 32)

	var groomer db.PetGroomer
	err = c.Bind(&groomer)
	if err != nil {
		log.Print(err.Error())
		return
	}
	count := 0
	db.DataBase.Model(&db.PetGroomer{}).Where("account_id = ?", accountID).Count(&count)
	if count != 0 {
		// exist
		db.DataBase.Model(&db.PetGroomer{}).Where("account_id = ?", accountID).Update(&groomer)
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "更新成功"})
	} else {
		// create
		db.DataBase.Create(&groomer)
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "创建成功"})
	}
}

// UploadHouse 更新门店非图片类信息:昵称，电话等文字信息
func UploadHouse(c *gin.Context) {
	auth := c.Request.Header.Get("authorization")
	tokenStr, err := extractTokenFromAuth(auth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": err.Error()})
		return
	}
	tokenPayload, err := ParseToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": err.Error()})
		return
	}
	accountIDStr := tokenPayload["account_id"].(string)
	accountID, _ := strconv.ParseUint(accountIDStr, 10, 32)

	var house db.PetHouse
	err = c.Bind(&house)
	if err != nil {
		log.Print(err.Error())
		return
	}
	count := 0
	db.DataBase.Model(&db.PetHouse{}).Where("account_id = ?", accountID).Count(&count)
	if count != 0 {
		// exist
		db.DataBase.Model(&db.PetHouse{}).Where("account_id = ?", accountID).Update(&house)
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "更新成功"})
	} else {
		// create
		db.DataBase.Create(&house)
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "创建成功"})
	}
}

// UploadImage 上传图片功能
func UploadImage(c *gin.Context) {
	auth := c.Request.Header.Get("authorization")
	tokenStr, err := extractTokenFromAuth(auth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": err.Error()})
		return
	}
	tokenPayload, err := ParseToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": err.Error()})
		return
	}
	userType := tokenPayload["user_type"].(string)
	accountIDStr := tokenPayload["account_id"].(string)
	accountID, _ := strconv.ParseUint(accountIDStr, 10, 32)

	imageType := c.Query("image_type")
	switch imageType {
	case "avatar":
		fileFront, headerFront, err := c.Request.FormFile("avatar")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": err.Error()})
			return
		}
		err = UploadAvatar(accountID, fileFront, headerFront, userType)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 403, "msg": "Sorry", "data": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "更新成功"})
		return

	case "id_card":
		IDCardNumber := c.Request.PostFormValue("id_card_number")
		fileFront, headerFront, err := c.Request.FormFile("id-front")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": err.Error()})
			return
		}
		fileBack, headerBack, err := c.Request.FormFile("id-back")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": err.Error()})
			return
		}
		err = UploadIDCard(accountID, IDCardNumber, fileFront, headerFront, fileBack, headerBack, userType)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 403, "msg": "Sorry", "data": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "更新成功"})
		return

	case "certificate":
		var groomer db.PetGroomer
		groomerAccount := 0
		db.DataBase.Where("account_id = ?", accountID).First(&groomer).Count(&groomerAccount)
		if groomerAccount == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"code": 404, "msg": "Sorry", "data": "没有找到该美容师账户"})
			return
		}
		fileFront, headerFront, err := c.Request.FormFile("certifi-front")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": err.Error()})
			return
		}
		fileBack, headerBack, err := c.Request.FormFile("certifi-back")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": err.Error()})
			return
		}
		fileNameFront, err := transferImage(fileFront, headerFront, setting.ImagePathSetting.GroomerCertificatePath)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 403, "data": "图片大小不能超过5M", "error": nil})
			return
		}
		fileNameBack, err := transferImage(fileBack, headerBack, setting.ImagePathSetting.GroomerCertificatePath)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 403, "data": "图片大小不能超过5M", "error": nil})
			return
		}
		db.DataBase.Model(&groomer).Update(db.PetGroomer{
			CertificateFront: setting.ImagePathSetting.GroomerCertificatePath + fileNameFront,
			CertificateBack:  setting.ImagePathSetting.GroomerCertificatePath + fileNameBack})
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "更新成功"})
		return

	case "house_license":
		var house db.PetHouse
		houseAccount := 0
		db.DataBase.Where("account_id = ?", accountID).First(&house).Count(&houseAccount)
		if houseAccount == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"code": 404, "msg": "Sorry", "data": "没有找到该门店账户"})
			return
		}
		fileEnvFront, headerEnvFront, err := c.Request.FormFile("environment-front")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": err.Error()})
			return
		}
		fileEnvIn, headerEnvIn, err := c.Request.FormFile("environment-inside")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": err.Error()})
			return
		}
		fileFront, headerFront, err := c.Request.FormFile("license-front")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": err.Error()})
			return
		}
		fileNameEnvFront, err := transferImage(fileEnvFront, headerEnvFront, setting.ImagePathSetting.HouseEnvironmentPath)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 403, "data": "图片大小不能超过5M", "error": nil})
			return
		}
		fileNameEnvIn, err := transferImage(fileEnvIn, headerEnvIn, setting.ImagePathSetting.HouseEnvironmentPath)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 403, "data": "图片大小不能超过5M", "error": nil})
			return
		}
		fileNameFront, err := transferImage(fileFront, headerFront, setting.ImagePathSetting.HouseLicensePath)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 403, "data": "图片大小不能超过5M", "error": nil})
			return
		}
		db.DataBase.Model(&house).Update(db.PetHouse{
			EnvironmentFront:  setting.ImagePathSetting.HouseEnvironmentPath + fileNameEnvFront,
			EnvironmentInside: setting.ImagePathSetting.HouseEnvironmentPath + fileNameEnvIn,
			License:           setting.ImagePathSetting.HouseLicensePath + fileNameFront})
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "更新成功"})
		return

	default:
		c.JSON(http.StatusBadRequest, gin.H{"code": 404, "msg": "Sorry", "data": "未知请求"})
		return
	}
}

// UploadAvatar 上传头像
func UploadAvatar(accountID uint64, fileFront multipart.File, headerFront *multipart.FileHeader, userType string) error {
	switch userType {
	case "groomer":
		var groomer db.PetGroomer
		groomerAccount := 0
		db.DataBase.Where("account_id = ?", accountID).First(&groomer).Count(&groomerAccount)
		if groomerAccount == 0 {
			return errors.New("没有找到该美容师账户")
		}
		fileNameFront, err := transferImage(fileFront, headerFront, setting.ImagePathSetting.AvatarPath)
		if err != nil {
			return errors.New("头像大小不能超过3M")
		}
		db.DataBase.Model(&groomer).Update(db.PetGroomer{
			Avatar: setting.ImagePathSetting.AvatarPath + fileNameFront})
		return nil
	case "house":
		var house db.PetHouse
		houseAccount := 0
		db.DataBase.Where("account_id = ?", accountID).First(&house).Count(&houseAccount)
		if houseAccount == 0 {
			return errors.New("没有找到该门店账户")
		}
		fileNameFront, err := transferImage(fileFront, headerFront, setting.ImagePathSetting.AvatarPath)
		if err != nil {
			return errors.New("头像大小不能超过3M")
		}
		db.DataBase.Model(&house).Update(db.PetHouse{
			Avatar: setting.ImagePathSetting.AvatarPath + fileNameFront})
		return nil
	default:
		return errors.New("头像类型错误")
	}
}

// UploadIDCard 上传身份证正反面照片
func UploadIDCard(accountID uint64, IDCardNumber string, fileFront multipart.File, headerFront *multipart.FileHeader, fileBack multipart.File, headerBack *multipart.FileHeader, userType string) error {
	switch userType {
	case "groomer":
		var groomer db.PetGroomer
		groomerAccount := 0
		db.DataBase.Where("account_id = ?", accountID).First(&groomer).Count(&groomerAccount)
		if groomerAccount == 0 {
			return errors.New("没有找到该美容师账户")
		}
		fileNameFront, err := transferImage(fileFront, headerFront, setting.ImagePathSetting.GroomerIDCardPath)
		if err != nil {
			return errors.New("图片大小不能超过5M")
		}
		fileNameBack, err := transferImage(fileBack, headerBack, setting.ImagePathSetting.GroomerIDCardPath)
		if err != nil {
			return errors.New("图片大小不能超过5M")
		}
		db.DataBase.Model(&groomer).Update(db.PetGroomer{
			IDCardNumber: IDCardNumber,
			IDCardFront:  setting.ImagePathSetting.GroomerIDCardPath + fileNameFront,
			IDCardBack:   setting.ImagePathSetting.GroomerIDCardPath + fileNameBack})
		return nil
	case "house":
		var house db.PetHouse
		houseAccount := 0
		db.DataBase.Where("account_id = ?", accountID).First(&house).Count(&houseAccount)
		if houseAccount == 0 {
			return errors.New("没有找到该门店账户")
		}
		fileNameFront, err := transferImage(fileFront, headerFront, setting.ImagePathSetting.HouseIDCardPath)
		if err != nil {
			return errors.New("图片大小不能超过5M")
		}
		fileNameBack, err := transferImage(fileBack, headerBack, setting.ImagePathSetting.HouseIDCardPath)
		if err != nil {
			return errors.New("图片大小不能超过5M")
		}
		db.DataBase.Model(&house).Update(db.PetHouse{
			IDCardNumber: IDCardNumber,
			IDCardFront:  setting.ImagePathSetting.HouseIDCardPath + fileNameFront,
			IDCardBack:   setting.ImagePathSetting.HouseIDCardPath + fileNameBack})
		return nil
	default:
		return errors.New("身份证类型错误")
	}
}
