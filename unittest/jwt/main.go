package main

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	SECRETKEY = "SECRETKEY" //私钥
)

//自定义Claims
type CustomClaims struct {
	UserId int64
	jwt.StandardClaims
}

func main() {
	//生成token
	maxAge := 60 * 60 * 24
	// Create the Claims
	//claims := &jwt.StandardClaims{
	//    //    ExpiresAt: time.Now().Add(time.Duration(maxAge)*time.Second).Unix(), // 过期时间，必须设置,
	//    //    Issuer:    "jerry",// 非必须，也可以填充用户名，
	//    //}

	//或者用下面自定义claim
	claims := jwt.MapClaims{
		"key":   "testing",
		"id":    1,
		"phone": 13988888888,
		"utype": 1,
		"exp":   time.Now().Add(time.Duration(maxAge) * time.Second).Unix(), // 过期时间，必须设置
		"iat":   time.Now().Unix(),                                          // 当前时间
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(SECRETKEY))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("token: %v\n", tokenString)

	//解析token
	ret, err := ParseToken(tokenString)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("userinfo: %v\n", ret)
}

//解析token
func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(SECRETKEY), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
