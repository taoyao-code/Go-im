package ctrl

import (
	"net/http"
	"reptile-go/args"
	"reptile-go/server"
	"reptile-go/util"
)

var messageService server.MessageService

// 添加消息记录
func AddMessageLog(w http.ResponseWriter, r *http.Request) {
	var arg args.ContactArg
	util.Bind(r, &arg)
}
