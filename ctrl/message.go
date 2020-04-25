package ctrl

import (
	"encoding/json"
	"net/http"
	"reptile-go/args"
	"reptile-go/model"
	"reptile-go/server"
	"reptile-go/util"
	"time"
)

var messageService server.MessageService

// 获取消息记录
/**
@api {post} /message/chathistory 获取消息记录
@apiName 消息记录
@apiGroup 聊天消息

@apiParam {Number} userid 用户ID
@apiParam {Number} dstid 好友ID/群ID
@apiParam {Number} cmd 单/群聊

@apiSuccessExample Success-Response:
HTTP/1.1 200 OK
{
	"code": 0,
	"data": "",
	"msg": "xxx"
}

@apiErrorExample Error-Response:
HTTP/1.1 404 Not Found
{
	"code": -1,
	"msg": "xxx"
}
@apiUse CommonError
*/
func ChatHistory(w http.ResponseWriter, r *http.Request) {
	var arg args.ContactArg
	util.Bind(r, &arg)
	if arg.Userid == 0 || arg.Dstid == 0 || arg.Cmd == 0 {
		util.RespFail(w, "参数错误")
		return
	}
	var chat string
	if arg.Cmd == model.CMD_ROOM_MSG {
		chat = "chat_11"
	} else {
		chat = "chat_10"
	}
	lRange, b2 := server.LRange(chat, int64(arg.GetPageFrom()), int64(arg.GetPageSize()))
	if b2 {
		var msg model.Message
		msgList := make([]model.Message, 0)
		for _, data := range lRange {
			json.Unmarshal([]byte(data), &msg)
			if arg.Cmd == model.CMD_ROOM_MSG {
				if msg.Userid != 0 && arg.Dstid == msg.Dstid {
					msg.Createat = time.Now().Unix()
					//msgList = append(msgList, msg)
					msgList = append([]model.Message{msg}, msgList...)
				}
			} else {
				//(userid = ? and dstid = ?) or (dstid = ? and userid = ?)
				if msg.Userid != 0 && (msg.Dstid == arg.Userid && msg.Userid == arg.Dstid) || (arg.Userid == msg.Userid && arg.Dstid == msg.Dstid) {
					msg.Createat = time.Now().Unix()
					//msgList = append(msgList, msg)
					msgList = append([]model.Message{msg}, msgList...)
				}
			}
		}
		util.RespOkList(w, msgList, len(msgList))
		return
	} else {
		history := messageService.GetChatHistory(arg.Userid, arg.Dstid, arg.Cmd, arg.GetPageFrom(), arg.GetPageSize())
		util.RespOkList(w, history, len(history))
	}
}
