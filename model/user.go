package model

import "time"

const (
	SEX_WOMEN  = "W" // 女
	SEX_MAN    = "M" // 男
	SEX_UNKNOW = "U" // 未知
)

type User struct {
	Id       int64     `xorm:"pk autoincr bigint(20)" form:"id" json:"id"`  // 用户ID
	Mobile   string    `xorm:"varchar(20)" form:"mobile" json:"mobile"`     // 手机号
	Passwd   string    `xorm:"varchar(40)" form:"passwd" json:"-"`          // 用户密码 = f(plainpwd + salt),MD5,
	Avatar   string    `xorm:"varchar(150)" form:"avatar" json:"avatar"`    // 头像
	Sex      string    `xorm:"varchar(2)" form:"sex" json:"sex"`            // 性别
	Nickname string    `xorm:"varchar(20)" form:"nickname" json:"nickname"` // 昵称
	Salt     string    `xorm:"varchar(10)" form:"salt" json:"-"`            // 加盐随机字符串6
	Online   int       `xorm:"int(10)" form:"online" json:"online"`         // 是否在线
	Token    string    `xorm:"varchar(40)" form:"token" json:"token"`       // Token
	Memo     string    `xorm:"varchar(140)" form:"memo" json:"memo"`        // 备注
	Createat time.Time `xorm:"datetime" form:"createat" json:"createat"`    // 创建时间
}
