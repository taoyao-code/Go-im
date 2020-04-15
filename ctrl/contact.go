package ctrl

import (
	"net/http"
	"reptile-go/args"
	"reptile-go/model"
	"reptile-go/server"
	"reptile-go/util"
)

var contactService server.ContactService

// 添加好友
func Addfriend(w http.ResponseWriter, r *http.Request) {
	var arg args.ContactArg
	util.Bind(r, &arg)
	// 调用service
	err := contactService.AddFriend(arg.Userid, arg.Dstid)
	if err != nil {
		util.RespFail(w, err.Error())
	} else {
		util.RespOk(w, nil, "好友添加成功")
	}
}

// 加载好友列
func LoadFriend(w http.ResponseWriter, r *http.Request) {
	var arg args.ContactArg
	util.Bind(r, &arg)
	users := contactService.SearchFriend(arg.Userid)
	util.RespOkList(w, users, len(users))
}

// 创建群
func CreateCommunity(w http.ResponseWriter, r *http.Request) {
	var arg model.Community
	util.Bind(r, &arg)
	conn, err := contactService.CreateCommunity(arg)
	if err != nil {
		util.RespFail(w, err.Error())
	} else {
		util.RespOk(w, conn, "创建群成功")
	}
}

// 加入群
func JoinCommunity(w http.ResponseWriter, r *http.Request) {
	var arg args.ContactArg
	util.Bind(r, &arg)
	err := contactService.JoinCommunity(arg.Userid, arg.Dstid)
	//todo 刷新用户的群组信息
	AddGroupId(arg.Userid, arg.Dstid)
	if err != nil {
		util.RespFail(w, err.Error())
	} else {
		util.RespOk(w, nil, "success")
	}
}

// 获取群列表
func LoadCommunity(w http.ResponseWriter, r *http.Request) {
	var arg args.ContactArg
	util.Bind(r, &arg)
	comunitys := contactService.SearchComunity(arg.Userid)
	util.RespOkList(w, comunitys, len(comunitys))
}
