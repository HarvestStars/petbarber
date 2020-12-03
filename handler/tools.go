package handler

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"

	"github.com/HarvestStars/petbarber/dtos"
	"github.com/HarvestStars/petbarber/setting"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
)

type BaseInfo struct {
	Account     string `json:"account"`
	IsActive    bool   `json:"isactive"`
	IsSuperUser bool   `json:"issuperuser"`
	UserType    int    `json:"usertype"`
}

// -----------------------------------------------------------------------------profile-----------------------------------------------------------------------------
// image
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

// GenImageURL 生成图片访问url
func GenImageURL(baseURL string, imagePath string) string {
	fileName := getFileNameWithSuffix(imagePath)
	URL := baseURL + fileName
	return URL
}

func getFileNameWithSuffix(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			return path[i+1:]
		}
	}
	return ""
}

// -----------------------------------------------------------------------------jwt-----------------------------------------------------------------------------
// 生成Jwt
func CreateJwtToken(user dtos.User) (dtos.Token, error) {
	claims := jwt.MapClaims{
		"key":   "testing",
		"id":    user.UserID,
		"phone": user.Phone,
		"utype": user.UserType,
		"exp":   time.Now().Add(time.Duration(setting.JwtSetting.JwtExpireTimeSec) * time.Second).Unix(), // 过期时间
		"iat":   time.Now().Unix(),                                                                       // 当前时间
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(setting.JwtSetting.JwtKey))
	if err != nil {
		return dtos.Token{}, err
	}
	JwtToken := dtos.Token{AccessToken: tokenString, TokenType: "Bearer", ExpireAt: claims["exp"].(int64)}
	return JwtToken, nil
}

// 更新Jwt
func RefreshJwtToken(claims jwt.MapClaims) (dtos.Token, error) {
	claims["exp"] = time.Now().Add(time.Duration(setting.JwtSetting.JwtExpireTimeSec) * time.Second).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(setting.JwtSetting.JwtKey))
	if err != nil {
		return dtos.Token{}, err
	}
	JwtToken := dtos.Token{AccessToken: tokenString, TokenType: "Bearer", ExpireAt: claims["exp"].(int64)}
	return JwtToken, nil
}

// 解析Jwt
func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(setting.JwtSetting.JwtKey), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func extractTokenFromAuth(auth string) (string, error) {
	if strings.HasPrefix(auth, "Bearer ") {
		token := auth[len("Bearer "):]
		return token, nil
	}
	return auth, errors.New("token 格式不合法")
}
