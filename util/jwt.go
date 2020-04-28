package util

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

const TokenExpireDuration = time.Hour * 2 // 2小时

var MySecret = []byte("夏天夏天悄悄过去")

// GenToken 生成JWT
func GenToken(username string) (string, error) {
	// 创建一个我们自己的声明
	c := MyClaims{
		username, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			Issuer:    "my-project",                               // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	if len(MySecret) == 0 {
		return "", errors.New("token_key为空")
	}
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(MySecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
		//fmt.Println(claims)
		return claims, nil
	}
	return nil, errors.New("令牌无效")
}

// AuthHandler: 获取Token
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	r.ParseForm()
	// 检查提供的凭据-如果将这些凭据存储在数据库中，则查询将在此处进行检查。
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	if username != "myusername" || password != "mypassword" {
		io.WriteString(w, `{"error":"账号或密码错误"}`)
		return
	}
	tokenString, _ := GenToken(username)
	io.WriteString(w, `{"token":"`+tokenString+`"}`)
	return
}

// JWTAuthMiddleware : 验证Token
func JWTAuthMiddleware(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	r.ParseForm()
	//authHeader := r.Header.Get("Authorization")	// 头信息中的
	authHeader := r.Form.Get("Authorization") // 路由中的
	if authHeader == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"code":-2,"msg":"token不存在"}`)
		return
	}
	mc, err := ParseToken(authHeader)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, `{"code":-3,"msg":"无效的Token"}`)
		return
	}
	fmt.Println(mc.Username)
}
