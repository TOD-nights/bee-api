package dto

import "github.com/flipped-aurora/gin-vue-admin/server/model/bee"

type BeePayLogInfoDto struct {
	bee.BeePayLog
	ShopName    string `json:"shopName" gorm:"column:shopName;"`
	GoodsName   string `json:"goodsName" gorm:"column:goodsName;"`
	GoodsNumber string `json:"goods_number" gorm:"column:goods_number;"`
}
