package handler

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"

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

func parseJWTPayload(tokenStr string, tokenPayload *map[string]interface{}) error {
	parts := strings.Split(tokenStr, ".")
	payloadByte, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return err
	}
	err = json.Unmarshal(payloadByte, tokenPayload)
	if err != nil {
		return err
	}
	return nil
}
