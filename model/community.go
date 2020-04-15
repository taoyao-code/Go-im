package model

import "time"

const (
	COMMUNITY_CATE_COM = 0x01
)

// 群组
type Community struct {
	Id       int64     `xorm:"pk autoincr bigint(20)" form:"id" json:"id"`
	Name     string    `xorm:"varchar(100)" form:"name",json:"name"`     // 名称
	Ownerid  int64     `xorm:"bigint(20)" form:"ownerid" json:"ownerid"` // 群主ID
	Icon     string    `xorm:"varchar(120)" form:"icon" json:"icon"`     // 群logo
	Cate     int       `xorm:"int(11)" form:"cate" json:"cate"`          // como、 群类型
	Memo     string    `xorm:"varchar(120)" form:"memo" json:"memo"`     // 描述
	Createat time.Time `xorm:"datetime" form:"createat" json:"createat"`
}
