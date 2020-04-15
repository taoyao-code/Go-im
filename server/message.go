package server

import (
	"reptile-go/model"
	"time"
)

type MessageService struct {
}

// 添加消息
func (service *MessageService) AddMessage(
	userId int64, cmd int, dstId int64, media int,
	content, pic, url, memo string, amount int, types int,
	username, face string) error {
	_, err := DbEngin.InsertOne(model.Message{
		Userid:   userId,
		Cmd:      cmd,
		Dstid:    dstId,
		Media:    media,
		Content:  content,
		Pic:      pic,
		Url:      url,
		Memo:     memo,
		Amount:   amount,
		Createat: time.Now().Unix(),
		Type:     types,
		Username: username,
		Face:     face,
	})
	if err != nil {
		return err
	}
	return nil
}

//获取聊天记录
func (service *MessageService) GetChatHistory(userId, dstId int64, cmd, pageForm, pageSize int) []model.Message {
	message := make([]model.Message, 0)
	if cmd == model.CMD_ROOM_MSG {
		DbEngin.Where("dstid = ? and cmd = ?", dstId, cmd).Desc("id").Limit(pageSize, pageForm).Find(&message)
		return message
	}
	DbEngin.Where("(userid = ? and dstid = ?) or (dstid = ? and userid = ?) and cmd = ?", userId, dstId, userId, dstId, cmd).
		Desc("id").Limit(pageSize, pageForm).Find(&message)
	return message
}
