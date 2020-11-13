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

func extractTokenFromAuth(auth string) (string, error) {
	if strings.HasPrefix(auth, "Bearer ") {
		token := auth[len("Bearer "):]
		return token, nil
	}
	return auth, errors.New("token 格式不合法")
}

//解析token
func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(setting.JwtSetting.SecretKey), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
