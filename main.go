package main

import (
	"html/template"
	"log"
	"net/http"
	"reptile-go/ctrl"
	"reptile-go/router"
	"time"

	"github.com/gorilla/mux"
)

/**
*	@apiDefine CommonError
*
*   @apiError (客户端错误) 400-BadRequest 请求信息有误，服务器不能或不会处理该请求
*   @apiError (服务端错误) 500-ServerError 服务器遇到了一个未曾预料的状况，导致了它无法完成对请求的处理。
*   @apiErrorExample {json} BadRequest
*	HTTP/1.1 401 BadRequest
*	{
*		"msg": "请求信息有误",
*		"code": -1,
*	}
*   @apiErrorExample {json} ServerError
*	HTTP/1.1 500 Internal Server Error
*	{
*		"message": "系统错误，请稍后再试",
*		"code": -1,
*		"data":[]
*	}
 */

func RegisterView() {
	//一次解析出全部模板
	tpl, err := template.ParseGlob("view/**/*")
	if err != nil {
		// 打印并直接退出
		log.Fatal(err.Error())
	}
	//通过for循环做好映射
	for _, v := range tpl.Templates() {
		tplname := v.Name()
		http.HandleFunc(tplname, func(w http.ResponseWriter, r *http.Request) {
			tpl.ExecuteTemplate(w, tplname, nil)
		})
	}
}

func handleFunc() {
	// 1. 提供静态资源目录支持
	//http.Handle("/asset/", http.FileServer(http.Dir(".")))
	http.Handle("/apidoc/", http.FileServer(http.Dir(".")))

	http.Handle("/mnt/", http.FileServer(http.Dir(".")))
	// 绑定请求的处理函数
	http.HandleFunc("/getCaptcha", cors(ctrl.GetCaptcha))         // 获取验证码
	http.HandleFunc("/user/register", cors(ctrl.UserRegister))    // 注册
	http.HandleFunc("/user/login", cors(ctrl.UserLogin))          // 登录
	http.HandleFunc("/contact/addfriend", cors(ctrl.Addfriend))   // 添加好友
	http.HandleFunc("/contact/loadfriend", cors(ctrl.LoadFriend)) // 加载好友列表

	http.HandleFunc("/contact/createcommunity", cors(ctrl.CreateCommunity)) // 创建群
	http.HandleFunc("/contact/joincommunity", cors(ctrl.JoinCommunity))     // 添加群
	http.HandleFunc("/contact/loadcommunity", cors(ctrl.LoadCommunity))     // 获取群列表

	http.HandleFunc("/chat", cors(ctrl.Chat))                      // ws
	http.HandleFunc("/attach/upload", cors(ctrl.UploadLocal))      //上传文件
	http.HandleFunc("/user/updateUser", cors(ctrl.UpdateUserInfo)) // 更新用户数据
	// 记录
	http.HandleFunc("/message/chathistory", cors(ctrl.ChatHistory)) // 获取聊天记录

	//RegisterView()
	//http.HandleFunc("/auth", util.AuthHandler) // 获取token
	//http.HandleFunc("/", util.JWTAuthMiddleware) // 验证tokne

}
func cors(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")                                                            // 允许访问所有域，可以换成具体url，注意仅具体url才能带cookie信息
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token") //header的类型
		w.Header().Add("Access-Control-Allow-Credentials", "true")                                                    //设置为true，允许ajax异步请求带cookie信息
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")                             //允许请求方法
		w.Header().Set("content-type", "application/json;charset=UTF-8")                                              //返回数据格式是json
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		f(w, r)
	}
}

type Middleware func(http.HandlerFunc) http.HandlerFunc

// 记录每个URL请求的执行时长
func Logging() mux.MiddlewareFunc {
	//	创建中间件
	return func(f http.Handler) http.Handler {
		//	创建一个新的handler包装http.HandlerFunc
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//	中间件的处理逻辑
			start := time.Now()
			defer func() {
				log.Println(r.URL.Path, time.Since(start))
			}()
			// 调用下一个中间件或者最终的handler处理程序
			f.ServeHTTP(w, r)
		})
	}
}

// 验证请求用的是否是指定的HTTP Method，不是则返回 400 Bad Request
func Method(m string) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != m {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			f(w, r)
		}
	}
}

func RegisterRoutes(r *mux.Router) {
	// apply Logging middleware
	r.Use(Logging(), router.AccessLogging)
}

func main() {
	handleFunc()
	// 将logrus的Logger转换为io.Writer
	//errorWriter := vlog.ErrorLog.Writer()
	// 关闭io.Writer
	//defer errorWriter.Close()
	//muxRouter := mux.NewRouter()
	//RegisterRoutes(muxRouter)

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		panic(err)
	}
}
