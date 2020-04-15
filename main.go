package main

import (
	"html/template"
	"log"
	"net/http"
	"reptile-go/ctrl"
)

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
	http.Handle("/asset/", http.FileServer(http.Dir(".")))
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
	http.HandleFunc("/attach/upload", cors(ctrl.Upload))           //上传文件
	http.HandleFunc("/user/updateUser", cors(ctrl.UpdateUserInfo)) // 更新用户数据

	// 记录
	http.HandleFunc("/message/chathistory", cors(ctrl.ChatHistory)) // 获取聊天记录

	RegisterView()

	//	https://www.hi-linux.com/posts/42176.html 配置反向代理
}

//func handleFunc() {
//	// 1. 提供静态资源目录支持
//	http.Handle("/asset/", http.FileServer(http.Dir(".")))
//	http.Handle("/mnt/", http.FileServer(http.Dir(".")))
//	// 绑定请求的处理函数
//	http.HandleFunc("/user/register", ctrl.UserRegister)          // 注册
//	http.HandleFunc("/user/login", ctrl.UserLogin)                // 登录
//	http.HandleFunc("/contact/addfriend", ctrl.Addfriend)         // 添加好友
//	http.HandleFunc("/contact/loadfriend", ctrl.LoadFriend) // 加载好友列表
//
//	http.HandleFunc("/contact/createcommunity", ctrl.CreateCommunity) // 创建群
//	http.HandleFunc("/contact/joincommunity", ctrl.JoinCommunity)     // 添加群
//	http.HandleFunc("/contact/loadcommunity", ctrl.LoadCommunity)     // 获取群列表
//
//	http.HandleFunc("/chat", ctrl.Chat)      // ws
//	http.HandleFunc("/attach/upload", ctrl.Upload) //上传文件
//	RegisterView()
//}

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

func main() {
	handleFunc()
	http.ListenAndServe(":8081", nil)
}
