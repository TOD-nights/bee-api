package model

import (
	"time"
)

const TableNameBeeUserMemberCardUseLog = "bee_user_member_card_use_log"

// BeeUserMemberCardUseLog 用户-会员卡使用记录
type BeeUserMemberCardUseLog struct {
	ID         int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UserCardID int32     `gorm:"column:user_card_id;comment:用户-会员卡 表id" json:"user_card_id" form:"user_card_id"` // 用户-会员卡 表id
	UseTime    time.Time `gorm:"column:use_time;comment:使用时间" json:"use_time" form:"use_time"`                   // 使用时间
	ShopId     int32     `gorm:"column:shop_id;comment:门店id" json:"shopId" form:"shopId"`                        // 门店id
	GoodsId    int32     `gorm:"column:goods_id;comment:商品id" json:"goodsId" form:"goodsId"`                     // 商品id
	GoodsName  string    `gorm:"column:goods_name;comment:商品名" json:"goodsName" form:"goodsName"`                // 商品名
	GoodsPrice float64   `gorm:"column:goods_price;comment:商品价格" json:"goodsPrice" form:"goodsPrice"`            // 商品价格
	ShopName   string    `gorm:"column:shop_name;comment:门店名" json:"shopName" form:"shopName"`                   // 门店名
}

// TableName BeeUserMemberCardUseLog's table name
func (*BeeUserMemberCardUseLog) TableName() string {
	return TableNameBeeUserMemberCardUseLog
}
