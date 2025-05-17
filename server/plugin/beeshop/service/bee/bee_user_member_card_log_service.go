package bee

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/beeshop/dto"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/beeshop/model/bee"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/beeshop/model/bee/request"
)

type userMemberCardLogService struct {
}

var UserMemberCardLogService = &userMemberCardLogService{}

func (s *userMemberCardLogService) Page(info request.MemberCardSearchInfo) ([]dto.BeeUserMemberCardUseLogDto, int64, error) {

	tx := global.GVA_DB.Model(&bee.BeeUserMemberCardUseLog{}).Joins("left join bee_user_member_card  a on a.id  = bee_user_member_card_use_log.user_card_id ")
	if info.Name != "" {
		tx = tx.Where("a.name like ?", "%"+info.Name+"%")
	}
	if info.ValidMonth > 0 {
		tx = tx.Where("a.valid_month = ?", info.ValidMonth)
	}
	var count int64 = 0
	var list []dto.BeeUserMemberCardUseLogDto
	tx.Count(&count)

	if info.Page > 0 {
		tx = tx.Offset((info.Page - 1) * info.PageSize).Limit(info.PageSize).Select("bee_user_member_card_use_log.*,a.name,a.valid_month ").Find(&list)
	}

	return list, count, tx.Error
}
