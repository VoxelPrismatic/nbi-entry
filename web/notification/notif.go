package notif

import (
	"nbientry/web"
	"nbientry/web/common"
)

var _ = web.Migrate(Notification{})

type Notification struct {
	BusinessUnitId int
	StageId        int
	UserId         int
}

func (n Notification) Exists(u common.User) bool {
	return web.GetFirst(Notification{UserId: u.Id, BusinessUnitId: n.BusinessUnitId, StageId: n.StageId}).UserId != 0
}

func (n Notification) Users() []common.User {
	user_ids := []int{}
	for _, u := range web.GetSorted(Notification{BusinessUnitId: n.BusinessUnitId, StageId: n.StageId}, "user_id ASC") {
		user_ids = append(user_ids, u.UserId)
	}

	ret := []common.User{}
	web.Db().Model(&common.User{}).Where("id in (?)", user_ids).Find(&ret)
	return ret
}

func (n Notification) NegateUsers() []common.User {
	user_ids := []int{}
	for _, u := range web.GetSorted(Notification{BusinessUnitId: n.BusinessUnitId, StageId: n.StageId}, "user_id ASC") {
		user_ids = append(user_ids, u.UserId)
	}

	ret := []common.User{}
	web.Db().Model(&common.User{}).Where("id not in (?)", user_ids).Find(&ret)
	return ret
}

func (n Notification) GetBusinessUnit() BusinessUnit {
	return web.GetFirst(BusinessUnit{Id: n.BusinessUnitId})
}

func (n Notification) GetStage() Stage {
	return web.GetFirst(Stage{Id: n.StageId})
}

func (n Notification) GetUser() common.User {
	return web.GetFirst(common.User{Id: n.UserId})
}

func GetNotifications(s Stage, bu BusinessUnit) []common.User {
	return web.GetSorted(common.User{}, "name ASC")
}
