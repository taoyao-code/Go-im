package server

import (
	"reptile-go/model"
	"time"
)

type MessageService struct {
}

// 添加消息
func (service MessageService) AddMessage(
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
