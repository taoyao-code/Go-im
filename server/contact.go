package server

import (
	"errors"
	"reptile-go/model"
	"time"
)

type ContactService struct {
}

// 添加好友
func (service *ContactService) AddFriend(userid, dstid int64) error {
	// 如果添加自己为好友
	if userid == dstid {
		return errors.New("无法添加自己为好友")
	}
	// 判断是否已是好友
	tmp := model.Contact{}
	tmpUser := model.User{}
	// 判断将要添加的好友是否存在
	DbEngin.Where("id = ?", dstid).Get(&tmpUser)
	// 用户不存在
	if tmpUser.Id == 0 {
		return errors.New("用户不存在!")
	}
	// 1.查询是否已经是好友了
	// 这里是条件链式操作
	// 获取1条数据
	DbEngin.Where("ownerid = ?", userid).
		And("dstobj = ?", dstid).
		And("cate = ?", model.CONCAT_CATE_USER).
		Get(&tmp)
	if tmp.Id > 0 {
		return errors.New("请勿重复添加好友!")
	}
	// 开启事务
	session := DbEngin.NewSession()
	session.Begin()
	// 插入自己的好友数据
	_, e1 := session.InsertOne(model.Contact{
		Ownerid:  userid,
		Dstobj:   dstid,
		Cate:     model.CONCAT_CATE_USER,
		Createat: time.Now(),
	})
	// 插入对方的数据
	_, e2 := session.InsertOne(model.Contact{
		Ownerid:  dstid,
		Dstobj:   userid,
		Cate:     model.CONCAT_CATE_USER,
		Createat: time.Now(),
	})
	// 如果没有错误
	if e1 == nil && e2 == nil {
		// 提交事务
		session.Commit()
		return nil
	} else {
		// 回滚事务
		session.Rollback()
		if e1 != nil {
			return e1
		} else {
			return e2
		}
	}
}

//查找好友列表
func (service *ContactService) SearchFriend(userId int64) []model.User {
	conconts := make([]model.Contact, 0)
	objIds := make([]int64, 0)
	// 查询好友列表
	DbEngin.Where("ownerid = ? and cate = ?", userId, model.CONCAT_CATE_USER).Find(&conconts)
	for _, v := range conconts {
		objIds = append(objIds, v.Dstobj)
	}
	coms := make([]model.User, 0)
	// 未查询到好友
	if len(objIds) == 0 {
		return coms
	}
	// 根据ID,查询符合条件的好友
	DbEngin.In("id", objIds).Find(&coms)
	return coms
}

// 创建群
func (service *ContactService) CreateCommunity(comm model.Community) (ret model.Community, err error) {
	if len(comm.Name) == 0 {
		return ret, errors.New("请输入群名称")
	}
	if comm.Ownerid == 0 {
		return ret, errors.New("请先登录")
	}
	com := model.Community{
		Ownerid: comm.Ownerid,
	}
	// 判断建群数量
	num, err := DbEngin.Count(&com)
	if num > 5 {
		return ret, errors.New("一个用户最多创建5个群")
	} else {
		comm.Createat = time.Now()
		// 开启事务
		session := DbEngin.NewSession()
		session.Begin()
		_, err = session.InsertOne(&comm)
		if err != nil {
			// 回滚事务
			session.Rollback()
			return com, err
		}
		_, err := session.InsertOne(
			model.Contact{
				Ownerid:  comm.Ownerid,
				Dstobj:   comm.Id,
				Cate:     model.CONCAT_CATE_COMUNITY,
				Createat: time.Now(),
			})
		if err != nil {
			session.Rollback()
		} else {
			session.Commit()
		}
		return com, err
	}
}

//获取群列表
func (service *ContactService) SearchComunity(userId int64) []model.Community {
	conconts := make([]model.Contact, 0)
	comIds := make([]int64, 0)
	DbEngin.Where("ownerid = ? and cate = ?", userId, model.CONCAT_CATE_COMUNITY).Find(&conconts)
	for _, v := range conconts {
		comIds = append(comIds, v.Dstobj)
	}
	coms := make([]model.Community, 0)
	if len(comIds) == 0 {
		return coms
	}
	DbEngin.In("id", comIds).Find(&coms)
	return coms
}

// 加入群
func (service *ContactService) JoinCommunity(userId, comId int64) error {
	_, err := service.ShowCommunityID(comId)
	if err != nil {
		return err
	}
	cot := model.Contact{
		Ownerid: userId,
		Dstobj:  comId,
		Cate:    model.CONCAT_CATE_COMUNITY,
	}
	DbEngin.Get(&cot)
	if cot.Id == 0 {
		cot.Createat = time.Now()
		_, err := DbEngin.InsertOne(cot)
		return err
	} else if cot.Id > 0 {
		return errors.New("已在该群")
	} else {
		return nil
	}
}

// 获取群信息
func (service *ContactService) SearchComunityIds(userId int64) (comIds []int64) {
	//	TODO 获取用户全部群ID
	conconts := make([]model.Contact, 0)
	comIds = make([]int64, 0)
	DbEngin.Where("ownerid = ? and cate = ?", userId, model.CONCAT_CATE_COMUNITY).Find(&conconts)
	for _, v := range conconts {
		comIds = append(comIds, v.Dstobj)
	}
	return comIds
}

// 获取单个群信息
func (server *ContactService) ShowCommunityID(dstId int64) (ret model.Community, err error) {
	com := model.Community{
		Id: dstId,
	}
	b, _ := DbEngin.Get(&com)
	if b == false {
		return ret, errors.New("该群不存在")
	} else {
		return com, nil
	}
}
