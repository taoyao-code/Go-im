package ctrl

import (
	"net/http"
	"reptile-go/args"
	"reptile-go/model"
	"reptile-go/server"
	"reptile-go/util"
	"reptile-go/validates"
)

var contactService server.ContactService
var contactValidate validates.ContactValidate

/**
@api {post} /contact/addfriend 添加好友
@apiName 添加好友
@apiGroup 好友
@apiParam {Number} userid Users unique ID.
@apiParam {Number} dstid 好友ID.
@apiSuccessExample Success-Response:
HTTP/1.1 200 OK
{
	"code": 0,
	"data": "",
	"msg": "xxx"
}
@apiError UserNotFound The id of the User was not found.
@apiErrorExample Error-Response:
HTTP/1.1 404 Not Found
{
	"code": -1,
	"msg": "xxx"
}
*/
func Addfriend(w http.ResponseWriter, r *http.Request) {
	var arg args.ContactArg
	util.Bind(r, &arg)
	contactValidates, err := contactValidate.ContactValidates(arg.Userid, arg.Dstid)
	if err != nil {
		util.RespFail(w, contactValidates)
		return
	}
	//调用service
	err = contactService.AddFriend(arg.Userid, arg.Dstid)
	if err != nil {
		util.RespFail(w, err.Error())
	} else {
		util.RespOk(w, nil, "好友添加成功")
	}
}

/**
@api {post} /contact/loadfriend 加载好友列表
@apiName 加载好友列表
@apiGroup 好友
@apiParam {Number} userid Users unique ID.
@apiSuccessExample Success-Response:
HTTP/1.1 200 OK
{
	"code": 0,
	"msg": "xxx"
	"rows": [
		{
			avatar: "头像"
			createat: "2020-04-02T17:32:57+08:00"
			id: 2
			memo: ""
			mobile: "账号"
			nickname: "昵称"
			online: 0
			sex: "U"
			token: "C6628AEC84FB7713FA9BE3A28A25BA50"
		},
	],
	"total": total,
}
@apiError UserNotFound The id of the User was not found.
@apiErrorExample Error-Response:
HTTP/1.1 404 Not Found
{
	"code": -1,
	"msg": "xxx"
}
*/
func LoadFriend(w http.ResponseWriter, r *http.Request) {
	var arg args.ContactArg
	util.Bind(r, &arg)
	if arg.Userid == 0 {
		util.RespFail(w, "参数错误")
		return
	}
	users := contactService.SearchFriend(arg.Userid)
	util.RespOkList(w, users, len(users))
}

/**
@api {post} /contact/createcommunity 创建群
@apiName 创建群
@apiGroup 群
@apiParam {String} name 群昵称
@apiParam {String} ownerid  用户ID
@apiParam {String} icon  群logo
@apiParam {String{可选}} icon  群logo

@apiSuccessExample Success-Response:
HTTP/1.1 200 OK
{
	"code": 0,
	"data": "",
	"msg": "xxx"
}
@apiError UserNotFound The id of the User was not found.

@apiErrorExample Error-Response:
HTTP/1.1 404 Not Found
{
	"code": -1,
	"msg": "xxx"
}
*/
func CreateCommunity(w http.ResponseWriter, r *http.Request) {
	var arg model.Community
	util.Bind(r, &arg)
	if arg.Ownerid == 0 || len(arg.Name) == 0 || len(arg.Icon) == 0 {
		util.RespFail(w, "参数错误")
		return
	}
	conn, err := contactService.CreateCommunity(arg)
	if err != nil {
		util.RespFail(w, err.Error())
	} else {
		util.RespOk(w, conn, "创建群成功")
	}
}

// 加入群
/**
@api {post} /contact/joincommunity 加入群
@apiName 加入群
@apiGroup 群

@apiParam {Number} userid 用户ID
@apiParam {Number} dstid 群ID

@apiSuccessExample Success-Response:
HTTP/1.1 200 OK
{
	"code": 0,
	"data": "",
	"msg": "xxx"
}
@apiError UserNotFound The id of the User was not found.

@apiErrorExample Error-Response:
HTTP/1.1 404 Not Found
{
	"code": -1,
	"msg": "xxx"
}
*/
func JoinCommunity(w http.ResponseWriter, r *http.Request) {
	var arg args.ContactArg
	util.Bind(r, &arg)
	if arg.Userid == 0 || arg.Dstid == 0 {
		util.RespFail(w, "参数错误")
		return
	}
	err := contactService.JoinCommunity(arg.Userid, arg.Dstid)
	//todo 刷新用户的群组信息
	AddGroupId(arg.Userid, arg.Dstid)
	if err != nil {
		util.RespFail(w, err.Error())
	} else {
		util.RespOk(w, nil, "success")
	}
}

/**
@api {post} /contact/loadcommunity 获取群列表
@apiName 获取群列表
@apiGroup 群

@apiParam {Number} userid 用户ID
@apiSuccessExample Success-Response:
HTTP/1.1 200 OK
{
"code": 0,
"msg": "xxx"
"rows": [
	{
		"id": 1,
		"Name": "1群",
		"ownerid": 1,
		"icon": "/mnt/1584786931964894457.jpg",
		"cate": 1,
		"memo": "1231",
		"createat": "2020-03-21T18:35:39+08:00"
	},
],
"total": total,
}
@apiErrorExample Error-Response:
HTTP/1.1 404 Not Found
{
	"code": -1,
	"msg": "xxx"
}
@apiUse CommonError
*/
func LoadCommunity(w http.ResponseWriter, r *http.Request) {
	var arg args.ContactArg
	util.Bind(r, &arg)
	if arg.Userid == 0 {
		util.RespFail(w, "参数错误")
		return
	}
	comunitys := contactService.SearchComunity(arg.Userid)
	util.RespOkList(w, comunitys, len(comunitys))
}
