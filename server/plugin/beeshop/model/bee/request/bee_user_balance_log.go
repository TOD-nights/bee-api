package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type BeeUserBalanceLogSearch struct {
	StartDateAdd *time.Time `json:"startDateAdd" form:"startDateAdd"`
	EndDateAdd   *time.Time `json:"endDateAdd" form:"endDateAdd"`
	Uid          *int       `json:"uid" form:"uid" `
	Type         string     `json:"type" form:"type"` // 新增字段：payment 或 recharge

	request.PageInfo
	Sort   string `json:"sort" form:"sort"`
	Order  string `json:"order" form:"order"`
	ShopId int    `json:"shopId" form:"shopId"`
}
