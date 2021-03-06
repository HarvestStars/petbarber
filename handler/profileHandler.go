package handler

import (
	"errors"
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
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_TYPE_WRONG, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}
	tokenPayload, err := ParseToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_VERIFY_RESULT_BAD_TOKEN, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}
	accountID := uint(tokenPayload["id"].(float64))

	var groomer dtos.TuGroomer
	err = c.Bind(&groomer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.PROFILE_UPLOAD_ERROR, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}
	count := 0
	db.DataBase.Model(&dtos.TuGroomer{}).Where("account_id = ?", accountID).Count(&count)
	if count != 0 {
		// exist
		groomer.UpdatedAt = time.Now().UTC().UnixNano() / 1e6
		db.DataBase.Model(&dtos.TuGroomer{}).Where("account_id = ?", accountID).UpdateColumns(groomer)
		c.JSON(http.StatusOK, gin.H{"code": dtos.OK, "msg": "OK", "data": "", "detail": "更新成功"})
	} else {
		// create
		groomer.AccountID = accountID
		groomer.CreatedAt = time.Now().UTC().UnixNano() / 1e6
		groomer.UpdatedAt = time.Now().UTC().UnixNano() / 1e6
		db.DataBase.Create(&groomer)
		c.JSON(http.StatusOK, gin.H{"code": dtos.OK, "msg": "OK", "data": "", "detail": "创建成功"})
	}
}

// UploadHouse 更新门店非图片类信息:昵称，电话等文字信息
func UploadHouse(c *gin.Context) {
	auth := c.Request.Header.Get("authorization")
	tokenStr, err := extractTokenFromAuth(auth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_TYPE_WRONG, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}
	tokenPayload, err := ParseToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_VERIFY_RESULT_BAD_TOKEN, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}
	accountID := uint(tokenPayload["id"].(float64))

	var house dtos.TuPethouse
	err = c.Bind(&house)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.PROFILE_UPLOAD_ERROR, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}
	count := 0
	db.DataBase.Model(&dtos.TuPethouse{}).Where("account_id = ?", accountID).Count(&count)
	if count != 0 {
		// exist
		house.UpdatedAt = time.Now().UTC().UnixNano() / 1e6
		db.DataBase.Model(&dtos.TuPethouse{}).Where("account_id = ?", accountID).UpdateColumns(house)
		c.JSON(http.StatusOK, gin.H{"code": dtos.OK, "msg": "OK", "data": "", "detail": "更新成功"})
	} else {
		// create
		house.AccountID = accountID
		house.CreatedAt = time.Now().UTC().UnixNano() / 1e6
		house.UpdatedAt = time.Now().UTC().UnixNano() / 1e6
		db.DataBase.Create(&house)
		c.JSON(http.StatusOK, gin.H{"code": dtos.OK, "msg": "OK", "data": "", "detail": "创建成功"})
	}
}

type IDCardResp struct {
	Front string `json:"front"`
	Back  string `json:"back"`
}

// UploadImage 上传图片功能
func UploadImage(c *gin.Context) {
	auth := c.Request.Header.Get("authorization")
	tokenStr, err := extractTokenFromAuth(auth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_TYPE_WRONG, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}
	tokenPayload, err := ParseToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_VERIFY_RESULT_BAD_TOKEN, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}
	userType := int(tokenPayload["utype"].(float64))
	accountID := uint(tokenPayload["id"].(float64))

	imageType := c.Query("image_type")
	switch imageType {
	case "avatar":
		fileFront, headerFront, err := c.Request.FormFile("avatar")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": dtos.IMAGE_FETCH_ERROR, "msg": "Sorry", "data": "", "detail": err.Error()})
			return
		}
		err = UploadAvatar(accountID, fileFront, headerFront, userType)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": dtos.IMAGE_UPLOAD_ERROR, "msg": "Sorry", "data": "", "detail": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": dtos.OK, "msg": "OK", "data": "", "detail": "更新成功"})
		return

	case "id_card":
		IDCardUrl := IDCardResp{}
		cardFlag := 3 // 01=front, 10=back, 11=both
		IDCardNumber := c.Request.PostFormValue("id_card_number")
		name := c.Request.PostFormValue("name")

		fileFront, headerFront, err := c.Request.FormFile("id_front")
		if err != nil {
			cardFlag = cardFlag & 2
		}
		fileBack, headerBack, err := c.Request.FormFile("id_back")
		if err != nil {
			cardFlag = cardFlag & 1
		}
		err = UploadIDCard(accountID, name, IDCardNumber, fileFront, headerFront, fileBack, headerBack, userType, cardFlag, &IDCardUrl)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": dtos.IMAGE_UPLOAD_ERROR, "msg": "Sorry", "data": "", "detail": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": IDCardUrl, "detail": "更新成功"})
		return

	case "certificate":
		type certifiUrl struct {
			Front string `json:"front"`
		}
		certificate := certifiUrl{}
		if userType != 2 {
			c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_EXPECTED_PETGROOMER_TOKEN, "msg": "Sorry", "data": "", "detail": "jwt usertype error"})
			return
		}
		var groomer dtos.TuGroomer
		groomerAccount := 0
		db.DataBase.Model(&dtos.TuGroomer{}).Where("account_id = ?", accountID).First(&groomer).Count(&groomerAccount)
		if groomerAccount == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_EXPECTED_PETGROOMER_TOKEN, "msg": "Sorry", "data": "", "detail": "没有找到该美容师账户"})
			return
		}
		fileFront, headerFront, err := c.Request.FormFile("certifi_front")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": dtos.IMAGE_FETCH_ERROR, "msg": "Sorry", "data": "", "detail": err.Error()})
			return
		}
		// fileBack, headerBack, err := c.Request.FormFile("certifi_back")
		// if err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"code": dtos.IMAGE_FETCH_ERROR, "msg": "Sorry", "data": "", "detail": err.Error()})
		// 	return
		// }
		fileNameFront, err := transferImage(fileFront, headerFront, setting.ImagePathSetting.GroomerCertificatePath)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": dtos.IMAGE_UPLOAD_ERROR, "msg": "Sorry", "data": "", "detail": "图片大小不能超过5M"})
			return
		}
		// fileNameBack, err := transferImage(fileBack, headerBack, setting.ImagePathSetting.GroomerCertificatePath)
		// if err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"code": dtos.IMAGE_UPLOAD_ERROR, "msg": "Sorry", "data": "", "detail": "图片大小不能超过5M"})
		// 	return
		// }
		db.DataBase.Model(&groomer).UpdateColumns(dtos.TuGroomer{
			UpdatedAt:        time.Now().UTC().UnixNano() / 1e6,
			CertificateFront: setting.ImagePathSetting.GroomerCertificatePath + fileNameFront,
		})
		certificate.Front = "/api/v1/images/certifi/" + fileNameFront
		c.JSON(http.StatusOK, gin.H{"code": dtos.OK, "msg": "OK", "data": certificate, "detail": "更新成功"})
		return

	case "house_license":
		type houseImageResp struct {
			Front   string `json:"front"`
			Inside  string `json:"inside"`
			License string `json:"license"`
		}
		houseImageUrl := houseImageResp{}

		licenseFlag := 7 // 111 = 门店正面，门店环境，营业执照
		if userType != 1 {
			c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_EXPECTED_PETHOUSE_TOKEN, "msg": "Sorry", "data": "", "detail": "jwt usertype error"})
			return
		}
		var house dtos.TuPethouse
		houseAccount := 0
		db.DataBase.Model(&dtos.TuPethouse{}).Where("account_id = ?", accountID).First(&house).Count(&houseAccount)
		if houseAccount == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_EXPECTED_PETHOUSE_TOKEN, "msg": "Sorry", "data": "", "detail": "没有找到该门店账户"})
			return
		}
		fileEnvFront, headerEnvFront, err := c.Request.FormFile("environment_front")
		if err != nil {
			licenseFlag = licenseFlag & 3 // 011
		}
		fileEnvIn, headerEnvIn, err := c.Request.FormFile("environment_inside")
		if err != nil {
			licenseFlag = licenseFlag & 5 // 101
		}
		fileFront, headerFront, err := c.Request.FormFile("license_front")
		if err != nil {
			licenseFlag = licenseFlag & 6 // 110
		}

		// 正面 100
		var fileNameEnvFront string
		if licenseFlag&4 == 4 {
			fileNameEnvFront, err = transferImage(fileEnvFront, headerEnvFront, setting.ImagePathSetting.HouseEnvironmentPath)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"code": dtos.IMAGE_UPLOAD_ERROR, "msg": "Sorry", "data": "", "detail": "图片大小不能超过5M"})
				return
			}
			db.DataBase.Model(&house).UpdateColumns(dtos.TuPethouse{
				UpdatedAt:        time.Now().UTC().UnixNano() / 1e6,
				EnvironmentFront: setting.ImagePathSetting.HouseEnvironmentPath + fileNameEnvFront,
			})
			houseImageUrl.Front = "/api/v1/images/envir/" + fileNameEnvFront
		}

		// 环境 010
		var fileNameEnvIn string
		if licenseFlag&2 == 2 {
			fileNameEnvIn, err = transferImage(fileEnvIn, headerEnvIn, setting.ImagePathSetting.HouseEnvironmentPath)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"code": dtos.IMAGE_UPLOAD_ERROR, "msg": "Sorry", "data": "", "detail": "图片大小不能超过5M"})
				return
			}
			db.DataBase.Model(&house).UpdateColumns(dtos.TuPethouse{
				UpdatedAt:         time.Now().UTC().UnixNano() / 1e6,
				EnvironmentInside: setting.ImagePathSetting.HouseEnvironmentPath + fileNameEnvIn,
			})
			houseImageUrl.Inside = "/api/v1/images/envir/" + fileNameEnvIn
		}

		// 营业执照 001
		var fileNameLicense string
		if licenseFlag&1 == 1 {
			fileNameLicense, err = transferImage(fileFront, headerFront, setting.ImagePathSetting.HouseLicensePath)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"code": dtos.IMAGE_UPLOAD_ERROR, "msg": "Sorry", "data": "", "detail": "图片大小不能超过5M"})
				return
			}
			db.DataBase.Model(&house).UpdateColumns(dtos.TuPethouse{
				UpdatedAt: time.Now().UTC().UnixNano() / 1e6,
				License:   setting.ImagePathSetting.HouseLicensePath + fileNameLicense,
			})
			houseImageUrl.License = "/api/v1/images/license/" + fileNameLicense
		}
		c.JSON(http.StatusOK, gin.H{"code": dtos.OK, "msg": "OK", "data": houseImageUrl, "detail": "更新成功"})
		return

	default:
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.UNKNOW_REQUEST, "msg": "Sorry", "data": "", "detail": "未知请求"})
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
		db.DataBase.Model(&dtos.TuPethouse{}).Where("account_id = ?", accountID).Count(&houseAccount).First(&house)
		if houseAccount == 0 {
			return errors.New("没有找到该门店账户")
		}
		fileNameFront, err := transferImage(fileFront, headerFront, setting.ImagePathSetting.AvatarPath)
		if err != nil {
			return errors.New("头像大小不能超过3M")
		}
		db.DataBase.Model(&house).UpdateColumns(dtos.TuPethouse{
			UpdatedAt: time.Now().UTC().UnixNano() / 1e6,
			Avatar:    setting.ImagePathSetting.AvatarPath + fileNameFront,
		})
		return nil
	case 2:
		var groomer dtos.TuGroomer
		groomerAccount := 0
		db.DataBase.Model(&dtos.TuGroomer{}).Where("account_id = ?", accountID).First(&groomer).Count(&groomerAccount)
		if groomerAccount == 0 {
			return errors.New("没有找到该美容师账户")
		}
		fileNameFront, err := transferImage(fileFront, headerFront, setting.ImagePathSetting.AvatarPath)
		if err != nil {
			return errors.New("头像大小不能超过3M")
		}
		db.DataBase.Model(&groomer).UpdateColumns(map[string]interface{}{
			"updated_at": time.Now().UTC().UnixNano() / 1e6,
			"avatar":     setting.ImagePathSetting.AvatarPath + fileNameFront,
		})
		return nil
	default:
		return errors.New("头像类型错误")
	}
}

// UploadIDCard 上传身份证正反面照片
func UploadIDCard(accountID uint, name string, IDCardNumber string, fileFront multipart.File, headerFront *multipart.FileHeader,
	fileBack multipart.File, headerBack *multipart.FileHeader, userType int, cardFlag int, idUrl *IDCardResp) error {
	var fileNameFront = ""
	var fileNameBack = ""
	err := errors.New("")

	switch userType {
	case 0:
		return errors.New("请先确定职业身份")

	case 1:
		var house dtos.TuPethouse
		houseAccount := 0
		db.DataBase.Model(&dtos.TuPethouse{}).Where("account_id = ?", accountID).First(&house).Count(&houseAccount)
		if houseAccount == 0 {
			return errors.New("没有找到该门店账户")
		}

		if cardFlag == 1 || cardFlag == 3 {
			fileNameFront, err = transferImage(fileFront, headerFront, setting.ImagePathSetting.HouseIDCardPath)
			if err != nil {
				return errors.New("图片大小不能超过5M")
			}
		}
		if cardFlag == 2 || cardFlag == 3 {
			fileNameBack, err = transferImage(fileBack, headerBack, setting.ImagePathSetting.HouseIDCardPath)
			if err != nil {
				return errors.New("图片大小不能超过5M")
			}
		}
		if cardFlag == 1 {
			db.DataBase.Model(&house).UpdateColumns(map[string]interface{}{
				"updated_at":     time.Now().UTC().UnixNano() / 1e6,
				"name":           name,
				"id_card_number": IDCardNumber,
				"id_card_front":  setting.ImagePathSetting.HouseIDCardPath + fileNameFront})
		}
		if cardFlag == 2 {
			db.DataBase.Model(&house).UpdateColumns(map[string]interface{}{
				"updated_at":     time.Now().UTC().UnixNano() / 1e6,
				"name":           name,
				"id_card_number": IDCardNumber,
				"id_card_back":   setting.ImagePathSetting.HouseIDCardPath + fileNameBack})
		}
		if cardFlag == 3 {
			db.DataBase.Model(&house).UpdateColumns(map[string]interface{}{
				"updated_at":     time.Now().UTC().UnixNano() / 1e6,
				"name":           name,
				"id_card_number": IDCardNumber,
				"id_card_front":  setting.ImagePathSetting.HouseIDCardPath + fileNameFront,
				"id_card_back":   setting.ImagePathSetting.HouseIDCardPath + fileNameBack})
		}
		idUrl.Front = "/api/v1/images/idcard/" + fileNameFront
		idUrl.Back = "/api/v1/images/idcard/" + fileNameBack
		return nil

	case 2:
		var groomer dtos.TuGroomer
		groomerAccount := 0
		db.DataBase.Model(&dtos.TuGroomer{}).Where("account_id = ?", accountID).First(&groomer).Count(&groomerAccount)
		if groomerAccount == 0 {
			return errors.New("没有找到该美容师账户")
		}
		if cardFlag == 1 || cardFlag == 3 {
			fileNameFront, err = transferImage(fileFront, headerFront, setting.ImagePathSetting.GroomerIDCardPath)
			if err != nil {
				return errors.New("图片大小不能超过5M")
			}
		}
		if cardFlag == 2 || cardFlag == 3 {
			fileNameBack, err = transferImage(fileBack, headerBack, setting.ImagePathSetting.GroomerIDCardPath)
			if err != nil {
				return errors.New("图片大小不能超过5M")
			}
		}
		if cardFlag == 1 {
			db.DataBase.Model(&groomer).UpdateColumns(map[string]interface{}{
				"updated_at":     time.Now().UTC().UnixNano() / 1e6,
				"id_card_number": IDCardNumber,
				"id_card_front":  setting.ImagePathSetting.GroomerIDCardPath + fileNameFront})
		}
		if cardFlag == 2 {
			db.DataBase.Model(&groomer).UpdateColumns(map[string]interface{}{
				"updated_at":     time.Now().UTC().UnixNano() / 1e6,
				"id_card_number": IDCardNumber,
				"id_card_back":   setting.ImagePathSetting.GroomerIDCardPath + fileNameBack})
		}
		if cardFlag == 3 {
			db.DataBase.Model(&groomer).UpdateColumns(map[string]interface{}{
				"updated_at":     time.Now().UTC().UnixNano() / 1e6,
				"id_card_number": IDCardNumber,
				"id_card_front":  setting.ImagePathSetting.GroomerIDCardPath + fileNameFront,
				"id_card_back":   setting.ImagePathSetting.GroomerIDCardPath + fileNameBack})
		}
		idUrl.Front = "/api/v1/images/idcard/" + fileNameFront
		idUrl.Back = "/api/v1/images/idcard/" + fileNameBack
		return nil

	default:
		return errors.New("身份证类型错误")
	}
}
