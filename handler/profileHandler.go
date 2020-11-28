package handler

import (
	"errors"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/HarvestStars/petbarber/db"
	"github.com/HarvestStars/petbarber/dtos"
	"github.com/HarvestStars/petbarber/setting"
	"github.com/gin-gonic/gin"
)

const uploadMaxBytes int64 = 1024 * 1024 // 1M

// ----------------------------------------------------- 普通用户 -----------------------------------------------------
// UploadGroomer 更新美容师非图片信息:昵称，电话等文字信息
func UploadGroomer(c *gin.Context) {
	auth := c.Request.Header.Get("authorization")
	tokenStr, err := extractTokenFromAuth(auth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}
	tokenPayload, err := ParseToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}
	accountID := uint(tokenPayload["id"].(float64))

	var groomer dtos.TuGroomer
	err = c.Bind(&groomer)
	if err != nil {
		log.Print(err.Error())
		return
	}
	count := 0
	db.DataBase.Model(&dtos.TuGroomer{}).Where("account_id = ?", accountID).Count(&count)
	if count != 0 {
		// exist
		groomer.UpdatedAt = time.Now().UTC().UnixNano() / 1e6
		db.DataBase.Model(&dtos.TuGroomer{}).Where("account_id = ?", accountID).Update(&groomer)
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "", "detail": "更新成功"})
	} else {
		// create
		groomer.AccountID = accountID
		groomer.CreatedAt = time.Now().UTC().Unix() / 1e6
		groomer.UpdatedAt = time.Now().UTC().Unix() / 1e6
		db.DataBase.Create(&groomer)
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "", "detail": "创建成功"})
	}
}

// UploadHouse 更新门店非图片类信息:昵称，电话等文字信息
func UploadHouse(c *gin.Context) {
	auth := c.Request.Header.Get("authorization")
	tokenStr, err := extractTokenFromAuth(auth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}
	tokenPayload, err := ParseToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}
	accountID := uint(tokenPayload["id"].(float64))

	var house dtos.TuPethouse
	err = c.Bind(&house)
	if err != nil {
		log.Print(err.Error())
		return
	}
	count := 0
	db.DataBase.Model(&dtos.TuPethouse{}).Where("account_id = ?", accountID).Count(&count)
	if count != 0 {
		// exist
		house.UpdatedAt = time.Now().UTC().Unix() / 1e6
		db.DataBase.Model(&dtos.TuPethouse{}).Where("account_id = ?", accountID).Update(&house)
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "", "detail": "更新成功"})
	} else {
		// create
		house.AccountID = accountID
		house.CreatedAt = time.Now().UTC().Unix() / 1e6
		house.UpdatedAt = time.Now().UTC().Unix() / 1e6
		db.DataBase.Create(&house)
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "", "detail": "创建成功"})
	}
}

// UploadImage 上传图片功能
func UploadImage(c *gin.Context) {
	auth := c.Request.Header.Get("authorization")
	tokenStr, err := extractTokenFromAuth(auth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}
	tokenPayload, err := ParseToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}
	userType := int(tokenPayload["utype"].(float64))
	accountID := uint(tokenPayload["id"].(float64))

	imageType := c.Query("image_type")
	switch imageType {
	case "avatar":
		fileFront, headerFront, err := c.Request.FormFile("avatar")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": "", "detail": err.Error()})
			return
		}
		err = UploadAvatar(accountID, fileFront, headerFront, userType)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 403, "msg": "Sorry", "data": "", "detail": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "", "detail": "更新成功"})
		return

	case "id_card":
		IDCardNumber := c.Request.PostFormValue("id_card_number")
		fileFront, headerFront, err := c.Request.FormFile("id-front")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": "", "detail": err.Error()})
			return
		}
		fileBack, headerBack, err := c.Request.FormFile("id-back")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": "", "detail": err.Error()})
			return
		}
		err = UploadIDCard(accountID, IDCardNumber, fileFront, headerFront, fileBack, headerBack, userType)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 403, "msg": "Sorry", "data": "", "detail": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "更新成功"})
		return

	case "certificate":
		if userType != 2 {
			c.JSON(http.StatusBadRequest, gin.H{"code": 404, "msg": "Sorry", "data": "", "detail": "jwt usertype error"})
			return
		}
		var groomer dtos.TuGroomer
		groomerAccount := 0
		db.DataBase.Where("account_id = ?", accountID).First(&groomer).Count(&groomerAccount)
		if groomerAccount == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"code": 404, "msg": "Sorry", "data": "", "detail": "没有找到该美容师账户"})
			return
		}
		fileFront, headerFront, err := c.Request.FormFile("certifi-front")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": "", "detail": err.Error()})
			return
		}
		fileBack, headerBack, err := c.Request.FormFile("certifi-back")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": "", "detail": err.Error()})
			return
		}
		fileNameFront, err := transferImage(fileFront, headerFront, setting.ImagePathSetting.GroomerCertificatePath)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 403, "msg": "Sorry", "data": "", "detail": "图片大小不能超过5M"})
			return
		}
		fileNameBack, err := transferImage(fileBack, headerBack, setting.ImagePathSetting.GroomerCertificatePath)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 403, "msg": "Sorry", "data": "", "detail": "图片大小不能超过5M"})
			return
		}
		db.DataBase.Model(&groomer).Update(dtos.TuGroomer{
			UpdatedAt:        time.Now().UTC().Unix() / 1e6,
			CertificateFront: setting.ImagePathSetting.GroomerCertificatePath + fileNameFront,
			CertificateBack:  setting.ImagePathSetting.GroomerCertificatePath + fileNameBack})
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "", "detail": "更新成功"})
		return

	case "house_license":
		if userType != 1 {
			c.JSON(http.StatusBadRequest, gin.H{"code": 404, "msg": "Sorry", "data": "", "detail": "jwt usertype error"})
			return
		}
		var house dtos.TuPethouse
		houseAccount := 0
		db.DataBase.Where("account_id = ?", accountID).First(&house).Count(&houseAccount)
		if houseAccount == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"code": 404, "msg": "Sorry", "data": "", "detail": "没有找到该门店账户"})
			return
		}
		fileEnvFront, headerEnvFront, err := c.Request.FormFile("environment-front")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": "", "detail": err.Error()})
			return
		}
		fileEnvIn, headerEnvIn, err := c.Request.FormFile("environment-inside")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": "", "detail": err.Error()})
			return
		}
		fileFront, headerFront, err := c.Request.FormFile("license-front")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": "", "detail": err.Error()})
			return
		}
		fileNameEnvFront, err := transferImage(fileEnvFront, headerEnvFront, setting.ImagePathSetting.HouseEnvironmentPath)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 403, "msg": "Sorry", "data": "", "detail": "图片大小不能超过5M"})
			return
		}
		fileNameEnvIn, err := transferImage(fileEnvIn, headerEnvIn, setting.ImagePathSetting.HouseEnvironmentPath)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 403, "msg": "Sorry", "data": "", "detail": "图片大小不能超过5M"})
			return
		}
		fileNameFront, err := transferImage(fileFront, headerFront, setting.ImagePathSetting.HouseLicensePath)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 403, "msg": "Sorry", "data": "", "detail": "图片大小不能超过5M"})
			return
		}
		db.DataBase.Model(&house).Update(dtos.TuPethouse{
			UpdatedAt:         time.Now().UTC().Unix() / 1e6,
			EnvironmentFront:  setting.ImagePathSetting.HouseEnvironmentPath + fileNameEnvFront,
			EnvironmentInside: setting.ImagePathSetting.HouseEnvironmentPath + fileNameEnvIn,
			License:           setting.ImagePathSetting.HouseLicensePath + fileNameFront})
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "", "detail": "更新成功"})
		return

	default:
		c.JSON(http.StatusBadRequest, gin.H{"code": 404, "msg": "Sorry", "data": "", "detail": "未知请求"})
		return
	}
}

// UploadAvatar 上传头像
func UploadAvatar(accountID uint, fileFront multipart.File, headerFront *multipart.FileHeader, userType int) error {
	switch userType {
	case 0:
		return errors.New("请先确定职业身份")
	case 1:
		var house dtos.TuPethouse
		houseAccount := 0
		db.DataBase.Where("account_id = ?", accountID).First(&house).Count(&houseAccount)
		if houseAccount == 0 {
			return errors.New("没有找到该门店账户")
		}
		fileNameFront, err := transferImage(fileFront, headerFront, setting.ImagePathSetting.AvatarPath)
		if err != nil {
			return errors.New("头像大小不能超过3M")
		}
		db.DataBase.Model(&house).Update(dtos.TuPethouse{
			UpdatedAt: time.Now().UTC().Unix() / 1e6,
			Avatar:    setting.ImagePathSetting.AvatarPath + fileNameFront})
		return nil
	case 2:
		var groomer dtos.TuGroomer
		groomerAccount := 0
		db.DataBase.Where("account_id = ?", accountID).First(&groomer).Count(&groomerAccount)
		if groomerAccount == 0 {
			return errors.New("没有找到该美容师账户")
		}
		fileNameFront, err := transferImage(fileFront, headerFront, setting.ImagePathSetting.AvatarPath)
		if err != nil {
			return errors.New("头像大小不能超过3M")
		}
		db.DataBase.Model(&groomer).Update(dtos.TuGroomer{
			UpdatedAt: time.Now().UTC().Unix() / 1e6,
			Avatar:    setting.ImagePathSetting.AvatarPath + fileNameFront})
		return nil
	default:
		return errors.New("头像类型错误")
	}
}

// UploadIDCard 上传身份证正反面照片
func UploadIDCard(accountID uint, IDCardNumber string, fileFront multipart.File, headerFront *multipart.FileHeader, fileBack multipart.File, headerBack *multipart.FileHeader, userType int) error {
	switch userType {
	case 0:
		return errors.New("请先确定职业身份")
	case 1:
		var house dtos.TuPethouse
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
		db.DataBase.Model(&house).Update(dtos.TuPethouse{
			UpdatedAt:    time.Now().UTC().Unix() / 1e6,
			IDCardNumber: IDCardNumber,
			IDCardFront:  setting.ImagePathSetting.HouseIDCardPath + fileNameFront,
			IDCardBack:   setting.ImagePathSetting.HouseIDCardPath + fileNameBack})
		return nil
	case 2:
		var groomer dtos.TuGroomer
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
		db.DataBase.Model(&groomer).Update(dtos.TuGroomer{
			UpdatedAt:    time.Now().UTC().Unix() / 1e6,
			IDCardNumber: IDCardNumber,
			IDCardFront:  setting.ImagePathSetting.GroomerIDCardPath + fileNameFront,
			IDCardBack:   setting.ImagePathSetting.GroomerIDCardPath + fileNameBack})
		return nil
	default:
		return errors.New("身份证类型错误")
	}
}
