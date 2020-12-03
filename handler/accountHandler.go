package handler

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"crypto/hmac"
	"crypto/sha512"

	"github.com/HarvestStars/petbarber/db"
	"github.com/HarvestStars/petbarber/dtos"
	"github.com/HarvestStars/petbarber/setting"
	"github.com/gin-gonic/gin"
)

func SendSmsCode(c *gin.Context) {
	phone := c.Param("phone")
	smid, expireAt := createSmsCode(phone)
	res := dtos.SmsToken{Smsid: smid, ExpireAt: expireAt}
	c.JSON(http.StatusOK, gin.H{"code": dtos.OK, "msg": "OK", "data": res, "detail": ""})
}

func createSmsCode(phone string) (string, int64) {
	code := 1234 // random.randrange(1000,9999)
	codeStr := strconv.Itoa(code)
	expireAt := time.Now().Add(time.Duration(setting.JwtSetting.SmsExpireTimeSec) * time.Second).UTC().Unix()
	expireAtStr := strconv.FormatInt(expireAt, 10)
	expireAtHexStr := strconv.FormatInt(expireAt, 16)
	msg := fmt.Sprintf("%s.%s.%s", phone, expireAtStr, codeStr)
	hmac := hmac.New(sha512.New384, []byte(setting.JwtSetting.SmsKey))
	hmac.Write([]byte(msg))
	hmacByte := hmac.Sum(nil)
	sign := hex.EncodeToString(hmacByte)[:50]
	smid := sign + expireAtHexStr
	return smid, expireAt
}

func SigninOrSignup(c *gin.Context) {
	var signInReq dtos.UserSigninReq
	c.Bind(&signInReq)
	if len(signInReq.Smsid) != 58 {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.LOGIN_SMS_CODE_INVALID, "msg": "Sorry", "data": "", "detail": "smsid length wrong"})
		return
	}
	// get expire time
	expireAtHexStr := signInReq.Smsid[50:]
	expireAt, err := strconv.ParseInt(expireAtHexStr, 16, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.LOGIN_SMS_CODE_MISSMATCH, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}
	expireAtStr := strconv.FormatInt(expireAt, 10)
	if time.Now().UTC().Unix() > expireAt {
		// smsid超时
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.LOGIN_SMS_CODE_EXPIRED, "msg": "Sorry", "data": "", "detail": "smsid time out"})
		return
	}

	// check smsid
	msg := fmt.Sprintf("%s.%s.%s", signInReq.Phone, expireAtStr, signInReq.Code)
	hmac := hmac.New(sha512.New384, []byte(setting.JwtSetting.SmsKey))
	hmac.Write([]byte(msg))
	hmacByte := hmac.Sum(nil)
	sign := hex.EncodeToString(hmacByte)[:50]
	if sign != signInReq.Smsid[:50] {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.LOGIN_SMS_CODE_MISSMATCH, "msg": "Sorry", "data": "", "detail": "smsid check error"})
		return
	}

	// orm action
	// 需要添加角色切换功能，目前角色选定后不可更改
	var account dtos.TuAccount
	count := 0
	db.DataBase.Model(&dtos.TuAccount{}).Where("account = ?", signInReq.Phone).Count(&count).First(&account)
	if count == 0 {
		// create
		account.Account = signInReq.Phone
		account.IsActive = true
		account.IsSuperuser = false
		account.UserType = 0
		account.CreatedAt = time.Now().UTC().UnixNano() / 1e6
		db.DataBase.Create(&account)
	}

	// create jwt token
	var signinRes dtos.UserSigninRep
	signinRes.User = dtos.User{UserID: account.ID, Phone: account.Account, UserType: account.UserType}
	JwtToken, err := CreateJwtToken(signinRes.User)
	if err != nil {

	}
	signinRes.Token = JwtToken
	c.JSON(http.StatusOK, gin.H{"code": dtos.OK, "msg": "OK", "data": signinRes, "detail": ""})
}

func GetUserProfile(c *gin.Context) {
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
	jwtUserType := int(tokenPayload["utype"].(float64))
	phone := tokenPayload["phone"].(string)

	reqUserType := c.Query("utype")
	count := 0
	switch reqUserType {
	case "PetHouse":
		var house dtos.TuPethouse
		db.DataBase.Model(&dtos.TuPethouse{}).Where("account_id = ?", accountID).Count(&count).First(&house)
		if count == 0 {
			// create
			house.AccountID = accountID
			house.CreatedAt = time.Now().UTC().UnixNano() / 1e6
			db.DataBase.Create(&house)
		} else {
			// 转换图片URL
			house.Avatar = GenImageURL("/api/v1/images/avatar/", house.Avatar)
			house.IDCardFront = GenImageURL("/api/v1/images/idcard/", house.IDCardFront)
			house.IDCardBack = GenImageURL("/api/v1/images/idcard/", house.IDCardBack)
			house.EnvironmentFront = GenImageURL("/api/v1/images/envir/", house.EnvironmentFront)
			house.EnvironmentInside = GenImageURL("/api/v1/images/envir/", house.EnvironmentInside)
			house.License = GenImageURL("/api/v1/images/license/", house.License)
		}
		if jwtUserType != 1 {
			// 该身份的第一次登陆, 同时更新account表
			db.DataBase.Model(&dtos.TuAccount{}).Where("id = ?", accountID).UpdateColumns(dtos.TuAccount{UserType: 1, UpdatedAt: time.Now().UTC().UnixNano() / 1e6})
		}

		token, err := CreateJwtToken(dtos.User{UserID: accountID, Phone: phone, UserType: 1})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_CREATE_WRONG, "msg": "Sorry", "data": "", "detail": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": dtos.OK, "msg": "OK", "data": dtos.UserPetHouseProfileRep{User: house, Token: token}, "detail": ""})

	case "Groomer":
		var groomer dtos.TuGroomer
		db.DataBase.Model(&dtos.TuGroomer{}).Where("account_id = ?", accountID).Count(&count).First(&groomer)
		if count == 0 {
			// create
			groomer.AccountID = accountID
			groomer.CreatedAt = time.Now().UTC().UnixNano() / 1e6
			db.DataBase.Create(&groomer)
		} else {
			// 转换图片URL
			groomer.Avatar = GenImageURL("/api/v1/images/avatar/", groomer.Avatar)
			groomer.IDCardFront = GenImageURL("/api/v1/images/idcard/", groomer.IDCardFront)
			groomer.IDCardBack = GenImageURL("/api/v1/images/idcard/", groomer.IDCardBack)
			groomer.CertificateFront = GenImageURL("/api/v1/images/certifi/", groomer.CertificateFront)
			groomer.CertificateBack = GenImageURL("/api/v1/images/certifi/", groomer.CertificateBack)
		}
		if jwtUserType != 2 {
			// 该身份的第一次登陆, 同时更新account表
			db.DataBase.Model(&dtos.TuAccount{}).Where("id = ?", accountID).UpdateColumns(dtos.TuAccount{UserType: 2, UpdatedAt: time.Now().UTC().UnixNano() / 1e6})
		}

		token, err := CreateJwtToken(dtos.User{UserID: accountID, Phone: phone, UserType: 2})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_CREATE_WRONG, "msg": "Sorry", "data": "", "detail": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": dtos.OK, "msg": "OK", "data": dtos.UserGroomerProfileRep{User: groomer, Token: token}, "detail": ""})

	default:
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.JWT_TYPE_WRONG, "msg": "Sorry", "data": "", "detail": err.Error()})
	}
}
