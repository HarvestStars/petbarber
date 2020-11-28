package handler

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"crypto/hmac"
	"crypto/sha512"

	"github.com/HarvestStars/petbarber/setting"
	"github.com/gin-gonic/gin"
)

func SendSmsCode(c *gin.Context) {
	phone := c.Param("phone")
	smid, expireAt, err := createSmsCode(phone)
	if err != nil {

	}
	c.JSON(http.StatusOK, gin.H{"smid": smid, "expireAt": expireAt})
}

func createSmsCode(phone string) (string, int64, error) {
	code := 1234 // random.randrange(1000,9999)
	codeStr := strconv.Itoa(code)
	expireAt := time.Now().Add(300).UTC().Unix()
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

}

func GetUserProfile(c *gin.Context) {

}
