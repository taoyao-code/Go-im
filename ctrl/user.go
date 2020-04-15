package ctrl

import (
	"fmt"
	"math/rand"
	"net/http"
	"reptile-go/model"
	"reptile-go/server"
	"reptile-go/util"
	"strconv"
	"time"
)

var userService server.UserService

// 登录
func UserLogin(w http.ResponseWriter, r *http.Request) {
	// 解析参数
	if r.Method == http.MethodPost {
		r.ParseForm()
		//1.获取前端传递过来的参数
		mobile := r.PostForm.Get("mobile")
		plainpwd := r.PostForm.Get("passwd")
		user, err := userService.Login(mobile, plainpwd)
		if err != nil {
			util.RespFail(w, err.Error())
		} else {

			util.RespOk(w, user, "")
		}
	}
}

// 注册
func UserRegister(w http.ResponseWriter, r *http.Request) {
	//1.获取前端传递过来的参数
	// 解析参数
	if r.Method == http.MethodPost {
		r.ParseForm()
		mobile := r.PostForm.Get("mobile")
		plainpwd := r.PostForm.Get("passwd")

		uuid := r.PostForm.Get("uuid")
		code := r.PostForm.Get("code")
		// 检验验证码
		err := util.CaptchaVerifyHandle(uuid, code)
		if err != nil {
			util.RespFail(w, err.Error())
			return
		}
		rand.Seed(time.Now().UnixNano()) // 设置种子数为当前时间
		nickname := fmt.Sprintf("user%06d", rand.Int31())
		avatar := ""
		sex := model.SEX_UNKNOW
		user, err := userService.Register(mobile, plainpwd, nickname, avatar, sex)
		if err != nil {
			util.RespFail(w, err.Error())
		} else {
			util.RespOk(w, user, "")
		}
	}
}

//修改用户数据
func UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userid := r.PostForm.Get("userid")
	avatar := r.PostForm.Get("avatar")
	id, _ := strconv.Atoi(userid)
	userService.UserInfo(int64(id), avatar)
	util.RespOk(w, nil, "")
}

// 获取验证码
func GetCaptcha(w http.ResponseWriter, r *http.Request) {
	util.GenerateCaptchaHandler(w, r)
}
