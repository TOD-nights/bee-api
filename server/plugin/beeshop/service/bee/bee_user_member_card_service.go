package bee

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/beeshop/model/bee"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/beeshop/model/bee/request"
)

type userMemberCardService struct {
}

var UserMemberCardService = &userMemberCardService{}

func (s *userMemberCardService) Page(info request.MemberCardSearchInfo) ([]bee.BeeUserMemberCard, int64, error) {

	tx := global.GVA_DB.Model(&bee.BeeUserMemberCard{})
	if info.Name != "" {
		tx = tx.Where("name like ?", "%"+info.Name+"%")
	}
	if info.ValidMonth > 0 {
		tx = tx.Where("valid_month = ?", info.ValidMonth)
	}

	var count int64 = 0
	var list []bee.BeeUserMemberCard
	tx.Count(&count)

	if info.Page > 0 {
		tx = tx.Offset((info.Page - 1) * info.PageSize).Limit(info.PageSize).Find(&list)
	}

	return list, count, tx.Error
}
