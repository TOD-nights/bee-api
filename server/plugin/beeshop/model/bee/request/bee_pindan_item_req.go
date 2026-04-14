package request

import (
	"gitee.com/stuinfer/bee-api/enum"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/beeshop/dto"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/beeshop/model/bee"
)

type BeePindanItemReq struct {
	Id       int64            `json:"id" form:"id"`             // 订单ID
	Status   enum.OrderStatus `json:"status" form:"status"`     // 订单状态
	PageSize int              `json:"pageSize" form:"pageSize"` // 每页条数
	Page     int              `json:"page" form:"page"`         // 页码
}

type BeePindanItemResp struct {
	bee.BeePindanOrderItem
	ShopInfo  dto.PindanShopInfo  `json:"shopInfo" gorm:"embedded;embeddedPrefix:shopInfo_;"`
	GoodsInfo dto.PindanItemGoods `json:"goodsInfo" gorm:"embedded;embeddedPrefix:goodsInfo_;"`
}
