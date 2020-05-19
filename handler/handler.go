package handler

import (
	"log"
	"net/http"
	"reptile-go/ctrl"
	"reptile-go/middleware"
	"reptile-go/util"
	"text/template"

	"github.com/gorilla/mux"
)

// RegisterRoutes 路由
func RegisterRoutes(r *mux.Router) {
	// 1. 提供静态资源目录支持
	r.Handle("/apidoc/", http.FileServer(http.Dir(".")))
	r.Handle("/mnt/", http.FileServer(http.Dir(".")))

	// 不需要验证
	indexRouter := r.PathPrefix("/index").Subrouter()
	// 绑定请求的处理函数
	indexRouter.HandleFunc("/getCaptcha", ctrl.GetCaptcha).Methods(http.MethodPost)      // 获取验证码
	indexRouter.HandleFunc("/user/register", ctrl.UserRegister).Methods(http.MethodPost) // 注册
	indexRouter.HandleFunc("/user/login", ctrl.UserLogin).Methods(http.MethodPost)       // 登录
	indexRouter.HandleFunc("/auth", util.AuthHandler).Methods(http.MethodPost)           // 获取token
	// 文件
	fileRouter := r.PathPrefix("/attach").Subrouter()
	fileRouter.HandleFunc("/upload", ctrl.UploadLocal).Methods(http.MethodPost, http.MethodOptions) //上传文件

	indexRouter.HandleFunc("/chat", ctrl.Chat) // ws

	// 需要验证Token
	authRouter := r.PathPrefix("/").Subrouter()
	authRouter.Use(middleware.JWTAuthMiddleware, middleware.AccessLogging)

	authRouter.HandleFunc("/contact/addfriend", ctrl.Addfriend).Methods(http.MethodPost)   // 添加好友
	authRouter.HandleFunc("/contact/loadfriend", ctrl.LoadFriend).Methods(http.MethodPost) // 加载好友列表

	authRouter.HandleFunc("/contact/createcommunity", ctrl.CreateCommunity).Methods(http.MethodPost) // 创建群
	authRouter.HandleFunc("/contact/joincommunity", ctrl.JoinCommunity).Methods(http.MethodPost)     // 添加群
	authRouter.HandleFunc("/contact/loadcommunity", ctrl.LoadCommunity).Methods(http.MethodPost)     // 获取群列表

	authRouter.HandleFunc("/user/updateUser", ctrl.UpdateUserInfo).Methods(http.MethodPost) // 更新用户数据
	// 记录
	authRouter.HandleFunc("/message/chathistory", ctrl.ChatHistory).Methods(http.MethodPost) // 获取聊天记录

	//RegisterView()
	authRouter.HandleFunc("/", util.JWTAuthMiddleware).Methods("POST") // 验证tokne
}

/**
*	@apiDefine CommonError
*
*   @apiError (客户端错误) 400-BadRequest 请求信息有误，服务器不能或不会处理该请求
*   @apiError (服务端错误) 500-ServerError 服务器遇到了一个未曾预料的状况，导致了它无法完成对请求的处理。
*   @apiErrorExample {json} BadRequest
*	HTTP/1.1 400 BadRequest
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
