package ctrl

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reptile-go/model"
	"reptile-go/server"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
)

//type Message struct {
//	Id int64 `json:"id,omitempty" form:"id"` //消息ID
//	//谁发的
//	Userid int64 `json:"userid,string,omitempty" form:"userid"` //谁发的
//	//什么业务
//	Cmd int `json:"cmd,omitempty" form:"cmd"` //群聊还是私聊
//	//发给谁
//	Dstid int64 `json:"dstid,omitempty" form:"dstid"` //对端用户ID/群ID
//	//怎么展示
//	Media int `json:"media,omitempty" form:"media"` //消息按照什么样式展示
//	//内容是什么
//	Content string `json:"content,omitempty" form:"content"` //消息的内容
//	//图片是什么
//	Pic string `json:"pic,omitempty" form:"pic"` //预览图片
//	//连接是什么
//	Url string `json:"url,omitempty" form:"url"` //服务的URL
//	//简单描述
//	Memo string `json:"memo,omitempty" form:"memo"` //简单描述
//	//其他的附加数据，语音长度/红包金额
//	Amount int `json:"amount,omitempty" form:"amount"` //其他和数字相关的
//}

/**
消息发送结构体
1、MEDIA_TYPE_TEXT
{id:1,userid:2,dstid:3,cmd:10,media:1,content:"hello"}
2、MEDIA_TYPE_News
{id:1,userid:2,dstid:3,cmd:10,media:2,content:"标题",pic:"http://www.baidu.com/a/log,jpg",url:"http://www.a,com/dsturl","memo":"这是描述"}
3、MEDIA_TYPE_VOICE，amount单位秒
{id:1,userid:2,dstid:3,cmd:10,media:3,url:"http://www.a,com/dsturl.mp3",anount:40}
4、MEDIA_TYPE_IMG
{id:1,userid:2,dstid:3,cmd:10,media:4,url:"http://www.baidu.com/a/log,jpg"}
5、MEDIA_TYPE_REDPACKAGR //红包amount 单位分
{id:1,userid:2,dstid:3,cmd:10,media:5,url:"http://www.baidu.com/a/b/c/redpackageaddress?id=100000","amount":300,"memo":"恭喜发财"}
6、MEDIA_TYPE_EMOJ 6
{id:1,userid:2,dstid:3,cmd:10,media:6,"content":"cry"}
7、MEDIA_TYPE_Link 6
{id:1,userid:2,dstid:3,cmd:10,media:7,"url":"http://www.a,com/dsturl.html"}

7、MEDIA_TYPE_Link 6
{id:1,userid:2,dstid:3,cmd:10,media:7,"url":"http://www.a,com/dsturl.html"}

8、MEDIA_TYPE_VIDEO 8
{id:1,userid:2,dstid:3,cmd:10,media:8,pic:"http://www.baidu.com/a/log,jpg",url:"http://www.a,com/a.mp4"}

9、MEDIA_TYPE_CONTACT 9
{id:1,userid:2,dstid:3,cmd:10,media:9,"content":"10086","pic":"http://www.baidu.com/a/avatar,jpg","memo":"胡大力"}

*/

// 本核心在与形成userId 和 Node的映射关系
type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte //并行转串行
	GroupSets set.Interface
}

// 映射关系表
var clientMap = make(map[int64]*Node, 0)

// 读写锁
var rwlocker sync.RWMutex

// ws://ip/chat?id=1&token=xxxx
func Chat(w http.ResponseWriter, r *http.Request) {
	//TODO 校验Token合法性
	// checkToken
	// 获取路由参数
	query := r.URL.Query()
	id := query.Get("id")
	token := query.Get("token")
	userId, _ := strconv.ParseInt(id, 10, 64) // 将字符串转换为int64类型
	isvalida := checkToken(userId, token)
	//如果isvalida=true
	//isvalida=false
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(w, r, nil)
	if err != nil {
		log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
		log.Printf(err.Error())
		return
	}
	//	TODO 获得 conn
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}

	// TODO 获取用户全部群id
	comIds := contactService.SearchComunityIds(userId)
	for _, v := range comIds {
		node.GroupSets.Add(v)
	}
	// todo userid 和 node 形成绑定关系
	rwlocker.Lock() // 锁
	clientMap[userId] = node
	rwlocker.Unlock() // 释放锁
	// todo 完成发送逻辑，conn
	go sendproc(node)
	// todo 完成接收逻辑
	go recvproc(node)
}

// 发送协程(写)
func sendproc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
				log.Printf(err.Error())
				return
			}
		default:
		}
	}
}

// 接收协程(读)
func recvproc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
			log.Printf(err.Error())
			return
		}
		// todo 对data进一步处理
		dispatch(data)
		fmt.Printf("recv <= %s\n", data)
	}
}

// 调度逻辑处理
func dispatch(data []byte) {
	// TODO 解析data为message
	msg := model.Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
		log.Printf(err.Error())
		return
	}
	// TODO 根据cmd对逻辑进行处理
	switch msg.Cmd {
	case model.CMD_HEART:
		// 心跳 TODO 一般都不做
	case model.CMD_SINGLE_MSG:
		// 单聊
		sendMsg(msg.Dstid, data)
		//TODO 添加聊天记录
		// 例如key: chat_10_uid_1_tid_2
		//server.Rpush("chat_10_uid_"+strconv.FormatInt(msg.Userid, 10)+"_tid_"+strconv.FormatInt(msg.Dstid, 10), data)
		server.Rpush("chat_10", data)
		//go AddMessagesChat(msg)
	case model.CMD_ROOM_MSG:
		// 群聊
		// TODO 群聊转发逻辑
		for _, v := range clientMap {
			if v.GroupSets.Has(msg.Dstid) {
				v.DataQueue <- data
			}
		}
		//TODO 添加聊天记录
		//server.Rpush("chat_11_dstid_"+strconv.FormatInt(msg.Dstid, 10), data)
		server.Rpush("chat_11", data)
		//go AddMessagesChat(msg)
	case model.CMD_QUIT:
		//	退出
		DelClientMapID(msg.Userid)

	}
}

//todo 发送消息
func sendMsg(userId int64, msg []byte) {
	rwlocker.RLock() // 锁
	node, ok := clientMap[userId]
	rwlocker.RUnlock() //释放锁
	if ok {
		node.DataQueue <- msg
	}
}

// 检查 Token 是否有效
func checkToken(userId int64, token string) bool {
	// 从数据库中查询 并 比对
	user := userService.Find(userId)
	return user.Token == token
}

//todo 添加新的群ID到用户的groupset中
func AddGroupId(userId, gid int64) {
	//取得node
	rwlocker.Lock()
	node, ok := clientMap[userId]
	if ok {
		node.GroupSets.Add(gid)
	}
	//clientMap[userId] = node
	rwlocker.Unlock()
	//添加gid到set
}

// todo 用户退出删除连接
func DelClientMapID(userId int64) {
	rwlocker.Lock()
	_, ok := clientMap[userId]
	if ok {
		delete(clientMap, userId) //将userId:值,从map中删除
	}
	rwlocker.Unlock()
}

// 添加记录
func AddMessagesChat(msg model.Message) {
	var messageService server.MessageService
	if msg.Userid == 0 {
		return
	}
	err := messageService.AddMessage(msg.Userid, msg.Cmd, msg.Dstid, msg.Media, msg.Content, msg.Pic, msg.Url, msg.Memo, msg.Amount, msg.Type, msg.Username, msg.Face)
	if err != nil {
		log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
		log.Printf(err.Error())
	}
}

// 从redis队列中取出数据
func GetLRangeMessage(key string, start, stop int64) {
	lRange, b2 := server.LRange(key, start, stop)
	if b2 {
		fmt.Println(lRange)
	} else {
		fmt.Println("lRange: is not")
	}
}

// 定时器
// 定时从redis中取出数据
func TickGetRedisLPop() {
	// 创建一个计时器
	var msg model.Message
	var message server.MessageService
	timeTickerChan := time.Tick(time.Second * 600)
	for {
		lRange, b2 := server.LRange("chat_10", 0, -1)
		if b2 {
			msgList := make([]model.Message, 0)
			for _, data := range lRange {
				json.Unmarshal([]byte(data), &msg)
				if msg.Userid != 0 {
					server.Lpop("chat_10")
					msg.Createat = time.Now().Unix()
					msgList = append(msgList, msg)
				}
			}
			if len(msgList) > 0 {
				message.AddMessageList(msgList)
			}
			//AddMessagesChat(msg)
		} else {
			fmt.Println("lRange: is not")
		}
		// 群
		lRange, b2 = server.LRange("chat_11", 0, -1)
		if b2 {
			msgList := make([]model.Message, 0)
			for _, data := range lRange {
				json.Unmarshal([]byte(data), &msg)
				if msg.Userid != 0 {
					server.Lpop("chat_11")
					msg.Createat = time.Now().Unix()
					msgList = append(msgList, msg)
				}
			}
			if len(msgList) > 0 {
				message.AddMessageList(msgList)
			}
			//AddMessagesChat(msg)
		} else {
			fmt.Println("lRange: is not")
		}
		<-timeTickerChan
	}
}
func init() {
	fmt.Println("redis timing LPop")
	//go TickGetRedisLPop()
}
