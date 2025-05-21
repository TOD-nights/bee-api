package service

import (
	"gitee.com/stuinfer/bee-api/db"
	"gitee.com/stuinfer/bee-api/model"
	"gitee.com/stuinfer/bee-api/proto"
	"time"
)

type memberCardService struct {
}

var MemberCardService = &memberCardService{}

func (s *memberCardService) ListAll() ([]proto.MemberCardDto, error) {
	var list []proto.MemberCardDto
	var err = db.GetDB().Model(&model.BeeMemberCard{}).Joins("left join sys_dictionaries b on b.type='member_card_type' left join sys_dictionary_details a on( a.status = 1 and a.value = bee_member_card.valid_month and a.sys_dictionary_id = b.id) ").Select("bee_member_card.*,a.label as dictLabel").Where("bee_member_card.delete_flag = ?", 0).Order("bee_member_card.sort_num desc").Find(&list).Error
	return list, err
}

// 领取
func (s *memberCardService) SaveUserMemberCardLog(log model.BeeUserMemberCardUseLog) error {
	var userMemberCard model.BeeUserMemberCard
	if err := db.GetDB().Model(&model.BeeUserMemberCard{}).First(&userMemberCard, log.UserCardID).Error; err != nil {
		return err
	} else if err := db.GetDB().Model(&model.BeeUserMemberCard{ID: log.UserCardID}).Update("left_count", userMemberCard.LeftCount-1).Error; err != nil {
		return err
	} else {
		log.UseTime = time.Now()
		return db.GetDB().Save(&log).Error
	}
}
