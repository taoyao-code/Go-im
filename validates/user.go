package validates

import (
	"reptile-go/model"

	"github.com/gookit/validate"
)

type UserValidate struct {
}

func (validatec *ContactValidate) UserValidates(userid, dstid int64) (string, error) {
	contact := &model.User{
		Mobile: "",
		Passwd: "",
		Avatar: "",
		Token:  "",
	}
	// 创建 Validation 实例
	v := validate.Struct(contact)
	if v.Validate() { // 验证成功
		return "", nil
	} else {
		//fmt.Println(v.Errors)                  // 所有的错误消息
		//fmt.Println(v.Errors.One()) 			 // 返回随机一条错误消息
		//fmt.Println(v.Errors.Field("Ownerid")) // 返回该字段的错误消息
		return v.Errors.One(), v.Errors
	}
}
