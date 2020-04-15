package ctrl

import (
	"net/http"
	"reptile-go/args"
	"reptile-go/server"
	"reptile-go/util"
)

var messageService server.MessageService

// 获取消息记录
func ChatHistory(w http.ResponseWriter, r *http.Request) {
	var arg args.ContactArg
	util.Bind(r, &arg)
	if arg.Userid == 0 {
		util.RespOkList(w, nil, 0)
	}
	history := messageService.GetChatHistory(arg.Userid, arg.Dstid, arg.Cmd, arg.GetPageFrom(), arg.GetPageSize())
	util.RespOkList(w, history, len(history))
}
