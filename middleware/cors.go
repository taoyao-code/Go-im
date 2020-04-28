package middleware

import (
	"io"
	"log"
	"net/http"
	"reptile-go/util"
)

// 跨域
func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//	中间件的处理逻辑
		w.Header().Set("Access-Control-Allow-Origin", "*")                                                            // 允许访问所有域，可以换成具体url，注意仅具体url才能带cookie信息
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token") //header的类型
		w.Header().Add("Access-Control-Allow-Credentials", "true")                                                    //设置为true，允许ajax异步请求带cookie信息
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")                             //允许请求方法
		w.Header().Set("content-type", "application/json;charset=UTF-8")                                              //返回数据格式是json
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		// 调用下一个中间件或者最终的handler处理程序
		next.ServeHTTP(w, r)
	})
}

// JWTAuthMiddleware : 验证Token
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		r.ParseForm()
		//authHeader := r.Header.Get("Authorization") // 头信息中的
		authHeader := r.Form.Get("Authorization") // 路由中的
		if authHeader == "" {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, `{"code":-2,"msg":"token不存在"}`)
			return
		}
		_, err := util.ParseToken(authHeader)
		if err != nil {
			//w.WriteHeader(http.StatusUnauthorized)
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, `{"code":-3,"msg":`+err.Error()+`}`)
			return
		}
		// 调用下一个中间件或者最终的handler处理程序
		next.ServeHTTP(w, r)
	})
}

// 保存Token
type authenticationMiddleware struct {
	tokenUsers map[string]string
}

// Initialize it somewhere
func (amw *authenticationMiddleware) Populate() {
	amw.tokenUsers["00000000"] = "user0"
	amw.tokenUsers["aaaaaaaa"] = "userA"
	amw.tokenUsers["05f717e5"] = "randomUser"
	amw.tokenUsers["deadbeef"] = "user0"
}

// Middleware function, which will be called for each request
func (amw *authenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Session-Token")

		if user, found := amw.tokenUsers[token]; found {
			// We found the token in our map
			log.Printf("Authenticated user %s\n", user)
			// Pass down the request to the next middleware (or final handler)
			next.ServeHTTP(w, r)
		} else {
			// Write an error and stop the handler chain
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}
