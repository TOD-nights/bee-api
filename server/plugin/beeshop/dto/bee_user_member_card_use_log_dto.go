package dto

import "github.com/flipped-aurora/gin-vue-admin/server/plugin/beeshop/model/bee"

type BeeUserMemberCardUseLogDto struct {
	bee.BeeUserMemberCardUseLog
	Name       string `gorm:"column:name" json:"name"`
	ValidMonth int32  `gorm:"column:valid_month" json:"valid_month"`
}
