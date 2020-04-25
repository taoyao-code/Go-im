package ctrl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reptile-go/model"
	"reptile-go/server"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
)

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

var log = logrus.New()

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
		log.WithFields(logrus.Fields{
			"animal": "walrus",
			"size":   10,
		}).Warn(err.Error())
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
				log.WithFields(logrus.Fields{
					"animal": "walrus",
					"size":   10,
				}).Warn(err.Error())
				return
			}
			//default:
			//	continue
		}
	}
}

// 接收协程(读)
func recvproc(node *Node) {
	defer node.Conn.Close()
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			log.WithFields(logrus.Fields{
				"animal": "walrus",
				"size":   10,
			}).Warn(err.Error())
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
		log.WithFields(logrus.Fields{
			"animal": "walrus",
			"size":   10,
		}).Warn(err.Error())
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
	case model.CMD_NEW_FRIEND:
		// 通知新朋友添加
		sendMsg(msg.Dstid, data)
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
		log.WithFields(logrus.Fields{
			"animal": "walrus",
			"size":   10,
		}).Warn(err.Error())
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
		log.WithFields(logrus.Fields{
			"animal": "walrus",
			"size":   10,
		}).Warn("定时打开，请关闭这个")
		log.Printf("定时打开，请关闭这个")
		<-timeTickerChan
	}
}

//func init() {
//	fmt.Println("redis timing LPop")
//	//go TickGetRedisLPop()
//}
