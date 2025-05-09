package bee

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/beeshop/model/bee"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/beeshop/model/bee/request"
	"time"
)

type memberCardService struct {
}

var MemberCardService = &memberCardService{}

func (s *memberCardService) Page(info request.MemberCardSearchInfo) ([]bee.BeeMemberCard, int64, error) {

	tx := global.GVA_DB.Model(&bee.BeeMemberCard{})
	if info.Name != "" {
		tx = tx.Where("name like ?", "%"+info.Name+"%")
	}
	if info.ValidMonth > 0 {
		tx = tx.Where("valid_month = ?", info.ValidMonth)
	}

	if info.DeleteFlag >= 0 {
		tx = tx.Where("delete_flag = ?", info.DeleteFlag)
	}

	var count int64 = 0
	var list []bee.BeeMemberCard
	tx.Count(&count)

	if info.Page > 0 {
		tx = tx.Offset((info.Page - 1) * info.PageSize).Limit(info.PageSize).Find(&list)
	}

	return list, count, tx.Error
}

func (s *memberCardService) SaveOrUpdate(memberCard bee.BeeMemberCard) error {

	if memberCard.ID == 0 {
		memberCard.DeleteFlag = false
		memberCard.CreateTime = time.Now()
		memberCard.UpdateTime = time.Now()
		return global.GVA_DB.Model(&bee.BeeMemberCard{}).Omit("delete_time").Create(&memberCard).Error
	}
	return global.GVA_DB.Debug().Model(&bee.BeeMemberCard{ID: memberCard.ID}).Omit("delete_time").Save(&memberCard).Error
}

func (s *memberCardService) DeleteOneById(id int32) {
	global.GVA_DB.Model(&bee.BeeMemberCard{ID: id}).UpdateColumns(map[string]interface{}{"delete_time": time.Now(), "delete_flag": true})
}

func (s *memberCardService) Info(id int) (bee.BeeMemberCard, error) {
	var memberCard bee.BeeMemberCard
	err := global.GVA_DB.Where("id = ?", id).First(&memberCard).Error
	return memberCard, err
}

func (s *memberCardService) RecoverOneById(id int32) {
	global.GVA_DB.Model(&bee.BeeMemberCard{ID: id}).UpdateColumns(map[string]interface{}{"delete_time": nil, "delete_flag": false})
}
