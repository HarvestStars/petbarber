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
	smid, expireAt, err := createSmsCode(phone)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": "", "detail": err.Error()})
	}
	res := dtos.SmsToken{Smsid: smid, ExpireAt: expireAt}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "OK", "data": res, "detail": ""})
}

func createSmsCode(phone string) (string, int64, error) {
	code := 1234 // random.randrange(1000,9999)
	codeStr := strconv.Itoa(code)
	expireAt := time.Now().Add(time.Duration(300) * time.Second).UTC().Unix()
	expireAtStr := strconv.FormatInt(expireAt, 10)
	expireAtHexStr := strconv.FormatInt(expireAt, 16)
	msg := fmt.Sprintf("%s.%s.%s", phone, expireAtStr, codeStr)
	hmac := hmac.New(sha512.New384, []byte(setting.JwtSetting.SmsKey))
	hmac.Write([]byte(msg))
	hmacByte := hmac.Sum(nil)
	sign := hex.EncodeToString(hmacByte)[:50]
	smid := sign + expireAtHexStr
	return smid, expireAt, nil
}

func SigninOrSignup(c *gin.Context) {
	var signInReq dtos.UserSigninReq
	c.Bind(&signInReq)
	if len(signInReq.Smsid) != 58 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": "", "detail": "smsid length wrong"})
		return
	}
	// get expire time
	expireAtHexStr := signInReq.Smsid[50:]
	expireAt, err := strconv.ParseInt(expireAtHexStr, 16, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": "", "detail": err.Error()})
		return
	}
	expireAtStr := strconv.FormatInt(expireAt, 10)
	if time.Now().UTC().Unix() > expireAt {
		// smsid超时
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": "", "detail": "smsid time out"})
		return
	}

	// check smsid
	msg := fmt.Sprintf("%s.%s.%s", signInReq.Phone, expireAtStr, signInReq.Code)
	hmac := hmac.New(sha512.New384, []byte(setting.JwtSetting.SmsKey))
	hmac.Write([]byte(msg))
	hmacByte := hmac.Sum(nil)
	sign := hex.EncodeToString(hmacByte)[:50]
	if sign != signInReq.Smsid[:50] {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "msg": "Sorry", "data": "", "detail": "smsid check error"})
		return
	}

	// orm action
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
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "OK", "data": signinRes, "detail": ""})
}

func GetUserProfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "OK", "data": "校验通过", "detail": ""})
}
