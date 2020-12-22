package handler

import (
	"net/http"

	"github.com/HarvestStars/petbarber/db"
	"github.com/HarvestStars/petbarber/dtos"
	"github.com/HarvestStars/petbarber/setting"
	"github.com/gin-gonic/gin"
)

// GetAvatarImage 获取头像图片
func GetAvatarImage(c *gin.Context) {
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

	path := c.Param("name")
	count := 0
	switch userType {
	case 1:
		// 门店头像
		var pethouse dtos.TuPethouse
		db.DataBase.Model(&dtos.TuPethouse{}).Where("account_id = ?", accountID).Count(&count).First(&pethouse)
		if count == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"code": dtos.IMAGE_TYPE_ERROR, "msg": "Sorry", "data": "", "detail": "该用户不存在"})
			return
		}
		avatarSRC := setting.ImagePathSetting.AvatarPath + path
		if pethouse.Avatar != avatarSRC {
			// 请求头像并不属于该用户
			c.JSON(http.StatusBadRequest, gin.H{"code": dtos.IMAGE_CANNOT_READ, "msg": "Sorry", "data": "", "detail": "请求头像并不属于该用户"})
			return
		}
		c.File(pethouse.Avatar)

	case 2:
		// 美容师头像
		var groomer dtos.TuGroomer
		db.DataBase.Model(&dtos.TuGroomer{}).Where("account_id = ?", accountID).Count(&count).First(&groomer)
		if count == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"code": dtos.IMAGE_TYPE_ERROR, "msg": "Sorry", "data": "", "detail": "该用户不存在"})
			return
		}
		avatarSRC := setting.ImagePathSetting.AvatarPath + path
		if groomer.Avatar != avatarSRC {
			// 请求头像并不属于该用户
			c.JSON(http.StatusBadRequest, gin.H{"code": dtos.IMAGE_CANNOT_READ, "msg": "Sorry", "data": "", "detail": "请求头像并不属于该用户"})
			return
		}
		c.File(groomer.Avatar)

	default:
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_VERIFY_RESULT_BAD_TOKEN, "msg": "Sorry", "data": "", "detail": "请求者类型不明"})
		return
	}
}

func GetIDCardImage(c *gin.Context) {
	// auth := c.Request.Header.Get("authorization")
	// tokenStr, err := extractTokenFromAuth(auth)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_TYPE_WRONG, "msg": "Sorry", "data": "", "detail": err.Error()})
	// 	return
	// }
	// tokenPayload, err := ParseToken(tokenStr)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_VERIFY_RESULT_BAD_TOKEN, "msg": "Sorry", "data": "", "detail": err.Error()})
	// 	return
	// }
	// userType := int(tokenPayload["utype"].(float64))
	// accountID := uint(tokenPayload["id"].(float64))

	path := c.Param("name")
	userType := c.Query("user_type")
	//count := 0
	switch userType {
	case "pethouse":
		// 门店身份证
		// var pethouse dtos.TuPethouse
		// db.DataBase.Model(&dtos.TuPethouse{}).Where("account_id = ?", accountID).Count(&count).First(&pethouse)
		// if count == 0 {
		// 	c.JSON(http.StatusBadRequest, gin.H{"code": dtos.IMAGE_TYPE_ERROR, "msg": "Sorry", "data": "", "detail": "该用户不存在"})
		// 	return
		// }
		idCardSRC := setting.ImagePathSetting.HouseIDCardPath + path
		// if pethouse.IDCardFront != idCardSRC && pethouse.IDCardBack != idCardSRC {
		// 	// 身份证正反面都不属于该用户
		// 	c.JSON(http.StatusBadRequest, gin.H{"code": dtos.IMAGE_CANNOT_READ, "msg": "Sorry", "data": "", "detail": "请求身份证并不属于该用户"})
		// 	return
		// }
		c.File(idCardSRC)

	case "groomer":
		// 美容师身份证
		// var groomer dtos.TuGroomer
		// db.DataBase.Model(&dtos.TuGroomer{}).Where("account_id = ?", accountID).Count(&count).First(&groomer)
		// if count == 0 {
		// 	c.JSON(http.StatusBadRequest, gin.H{"code": dtos.IMAGE_TYPE_ERROR, "msg": "Sorry", "data": "", "detail": "该用户不存在"})
		// 	return
		// }
		idCardSRC := setting.ImagePathSetting.GroomerIDCardPath + path
		// if groomer.IDCardFront != idCardSRC && groomer.IDCardBack != idCardSRC {
		// 	// 身份证正反面都不属于该用户
		// 	c.JSON(http.StatusBadRequest, gin.H{"code": dtos.IMAGE_CANNOT_READ, "msg": "Sorry", "data": "", "detail": "请求身份证并不属于该用户"})
		// 	return
		// }
		c.File(idCardSRC)

	default:
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_VERIFY_RESULT_BAD_TOKEN, "msg": "Sorry", "data": "", "detail": "请求者类型不明"})
		return
	}
}

func GetEnvironmentImage(c *gin.Context) {
	// auth := c.Request.Header.Get("authorization")
	// tokenStr, err := extractTokenFromAuth(auth)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_TYPE_WRONG, "msg": "Sorry", "data": "", "detail": err.Error()})
	// 	return
	// }
	// tokenPayload, err := ParseToken(tokenStr)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_VERIFY_RESULT_BAD_TOKEN, "msg": "Sorry", "data": "", "detail": err.Error()})
	// 	return
	// }
	// userType := int(tokenPayload["utype"].(float64))
	// accountID := uint(tokenPayload["id"].(float64))

	path := c.Param("name")
	//count := 0
	//switch userType {
	//case 1:

	// 门店环境
	// var pethouse dtos.TuPethouse
	// db.DataBase.Model(&dtos.TuPethouse{}).Where("account_id = ?", accountID).Count(&count).First(&pethouse)
	// if count == 0 {
	// 	c.JSON(http.StatusBadRequest, gin.H{"code": dtos.IMAGE_TYPE_ERROR, "msg": "Sorry", "data": "", "detail": "该用户不存在"})
	// 	return
	// }
	envirSRC := setting.ImagePathSetting.HouseEnvironmentPath + path
	// if pethouse.EnvironmentFront != envirSRC && pethouse.EnvironmentInside != envirSRC {
	// 	// 环境图都不属于该用户
	// 	c.JSON(http.StatusBadRequest, gin.H{"code": dtos.IMAGE_CANNOT_READ, "msg": "Sorry", "data": "", "detail": "jwt不是门店用户"})
	// 	return
	// }
	c.File(envirSRC)

	//default:
	//	c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_VERIFY_RESULT_BAD_TOKEN, "msg": "Sorry", "data": "", "detail": "请求者类型不明"})
	//	return
	//}
}

func GetLicenseImage(c *gin.Context) {
	// auth := c.Request.Header.Get("authorization")
	// tokenStr, err := extractTokenFromAuth(auth)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_TYPE_WRONG, "msg": "Sorry", "data": "", "detail": err.Error()})
	// 	return
	// }
	// tokenPayload, err := ParseToken(tokenStr)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_VERIFY_RESULT_BAD_TOKEN, "msg": "Sorry", "data": "", "detail": err.Error()})
	// 	return
	// }
	// userType := int(tokenPayload["utype"].(float64))
	// accountID := uint(tokenPayload["id"].(float64))

	path := c.Param("name")
	//count := 0
	// switch userType {
	// case 1:
	// 门店执照
	// var pethouse dtos.TuPethouse
	// db.DataBase.Model(&dtos.TuPethouse{}).Where("account_id = ?", accountID).Count(&count).First(&pethouse)
	// if count == 0 {
	// 	c.JSON(http.StatusBadRequest, gin.H{"code": dtos.IMAGE_TYPE_ERROR, "msg": "Sorry", "data": "", "detail": "该用户不存在"})
	// 	return
	// }
	licenseSRC := setting.ImagePathSetting.HouseLicensePath + path
	// if pethouse.License != licenseSRC {
	// 	// 营业执照不属于该用户
	// 	c.JSON(http.StatusBadRequest, gin.H{"code": dtos.IMAGE_CANNOT_READ, "msg": "Sorry", "data": "", "detail": "jwt不是门店用户"})
	// 	return
	// }
	c.File(licenseSRC)

	// default:
	// 	c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_VERIFY_RESULT_BAD_TOKEN, "msg": "Sorry", "data": "", "detail": "请求者类型不明"})
	// 	return
	// }
}

func GetCertificateImage(c *gin.Context) {
	// auth := c.Request.Header.Get("authorization")
	// tokenStr, err := extractTokenFromAuth(auth)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_TYPE_WRONG, "msg": "Sorry", "data": "", "detail": err.Error()})
	// 	return
	// }
	// tokenPayload, err := ParseToken(tokenStr)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_VERIFY_RESULT_BAD_TOKEN, "msg": "Sorry", "data": "", "detail": err.Error()})
	// 	return
	// }
	// userType := int(tokenPayload["utype"].(float64))
	// accountID := uint(tokenPayload["id"].(float64))

	path := c.Param("name")
	//count := 0
	// switch userType {
	// case 2:
	// 美容师证书
	// var groomer dtos.TuGroomer
	// db.DataBase.Model(&dtos.TuGroomer{}).Where("account_id = ?", accountID).Count(&count).First(&groomer)
	// if count == 0 {
	// 	c.JSON(http.StatusBadRequest, gin.H{"code": dtos.IMAGE_TYPE_ERROR, "msg": "Sorry", "data": "", "detail": "该用户不存在"})
	// 	return
	// }
	certifiSRC := setting.ImagePathSetting.GroomerCertificatePath + path
	// if groomer.CertificateFront != certifiSRC && groomer.CertificateBack != certifiSRC {
	// 	// 身份证正反面都不属于该用户
	// 	c.JSON(http.StatusBadRequest, gin.H{"code": dtos.IMAGE_CANNOT_READ, "msg": "Sorry", "data": "", "detail": "jwt不是门店用户"})
	// 	return
	// }
	c.File(certifiSRC)

	// default:
	// 	c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_VERIFY_RESULT_BAD_TOKEN, "msg": "Sorry", "data": "", "detail": "请求者类型不明"})
	// 	return
	// }
}
