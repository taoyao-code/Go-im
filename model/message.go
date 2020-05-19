package model

type Message struct {
	Id       int64  `json:"id,omitempty" form:"id"`                //消息ID
	Userid   int64  `json:"userid,string,omitempty" form:"userid"` //谁发的
	Cmd      int    `json:"cmd,omitempty" form:"cmd"`              //群聊还是私聊
	Dstid    int64  `json:"dstid,omitempty" form:"dstid"`          //对端用户ID/群ID
	Media    int    `json:"media,omitempty" form:"media"`          //消息按照什么样式展示
	Content  string `json:"content,omitempty" form:"content"`      //消息的内容
	Pic      string `json:"pic,omitempty" form:"pic"`              //预览图片
	Url      string `json:"url,omitempty" form:"url"`              //服务的URL
	Memo     string `json:"memo,omitempty" form:"memo"`            //简单描述
	Amount   int    `json:"amount,omitempty" form:"amount"`        //其他和数字相关的/其他的附加数据，语音长度/红包金额
	Createat int64  `json:"createat,omitempty" form:"createat"`
	Type     int    `json:"type,omitempty" form:"type"`         // 消息类型/系统消息、用户消息
	Username string `json:"username,omitempty" form:"username"` // 用户昵称
	Face     string `json:"face,omitempty" form:"face"`         // 头像
}

const (
	CMD_NEW_FRIEND = 9  // 通知好友，有新朋友添加
	CMD_SINGLE_MSG = 10 // 点对点单聊,dstid是用户ID
	CMD_ROOM_MSG   = 11 // 群聊消息,dstid是群id
	CMD_HEART      = 0  // 心跳消息,不处理
	CMD_LOGIN      = 1  // 登录
	CMD_QUIT       = 2  // 退出
	CMD_FILTER     = 3  // 敏感信息
)

const (
	MSG_TYPE_SYSTEM = 1 // 系统消息
	MSG_TYPE_USER   = 2 // 用户消息
)

const (
	//文本样式
	MEDIA_TYPE_TEXT = 1
	//新闻样式,类比图文消息
	MEDIA_TYPE_News = 2
	//语音样式
	MEDIA_TYPE_VOICE = 3
	//图片样式
	MEDIA_TYPE_IMG = 4
	//红包样式
	MEDIA_TYPE_REDPACKAGR = 5
	//emoj表情样式
	MEDIA_TYPE_EMOJ = 6
	//超链接样式
	MEDIA_TYPE_LINK = 7
	//视频样式
	MEDIA_TYPE_VIDEO = 8
	//名片样式
	MEDIA_TYPE_CONCAT = 9
	//其他自己定义,前端做相应解析即可
	MEDIA_TYPE_UDEF = 100
)
