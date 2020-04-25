package model

import (
	"time"

	"github.com/gookit/validate"
)

const (
	CONCAT_CATE_USER     = 0x01 //用户
	CONCAT_CATE_COMUNITY = 0x02 //群组
)

// 好友和群都存在这个表中
// 可根据具体业务做拆分
type Contact struct {
	Id       int64     `xorm:"pk autoincr bigint(20)" form:"id" json:"id"`
	Ownerid  int64     `xorm:"bigint(20)" form:"ownerid" json:"ownerid" validate:"required|int"` // 谁的10000
	Dstobj   int64     `xorm:"bigint(20)" form:"dstobj" json:"dstobj" validate:"required|int"`   // 对端，10001
	Cate     int       `xorm:"int(11)" form:"cate" json:"cate"`                                  // 什么角色： 用户/群组
	Memo     string    `xorm:"varchar(120)" form:"memo" json:"memo"`                             // 介绍
	Createat time.Time `xorm:"datetime" form:"createat" json:"createat"`                         // 时间
}

// Messages 您可以自定义验证器错误消息
func (f Contact) Messages() map[string]string {
	return validate.MS{
		"required": "{field} 不能为空",
		"int":      "{field} 必须是数字",
	}
}

// Translates 你可以自定义字段翻译
func (f Contact) Translates() map[string]string {
	return validate.MS{
		"Ownerid": "userid",
		"Dstobj":  "dstid",
	}
}
