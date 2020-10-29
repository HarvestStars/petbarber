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
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "update done."})
	} else {
		// create
		db.DataBase.Create(&account)
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "created."})
	}
}

// RegistUpdateGroomer 更新美容师信息:昵称，电话等文字信息
func RegistUpdateGroomer(c *gin.Context) {
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

// RegistUpdateHouse 更新美容师信息:昵称，电话等文字信息
func RegistUpdateHouse(c *gin.Context) {
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

// UploadImage 上传图片的总路由
func UploadImage(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	accountID, _ := strconv.ParseUint(accountIDStr, 10, 32)
	userType := c.Query("user_type")
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
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "update done."})
		return

	case "id_card":
		IDCardNumber := c.Query("id_card_number")
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
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "update done."})
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
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "update done."})
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
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "OK", "data": "update done."})
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

func transferImage(file multipart.File, header *multipart.FileHeader, rootPath string) (string, error) {
	// header调用Filename方法，就可以得到文件名
	fileName := header.Filename
	filesuffix := path.Ext(fileName)
	u1, _ := uuid.NewV4()
	fileName = u1.String()
	fileName += filesuffix
	if header.Size > 5*uploadMaxBytes {
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
