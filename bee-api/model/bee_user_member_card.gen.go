package model

import (
	"time"
)

const TableNameBeeUserMemberCard = "bee_user_member_card"

// BeeUserMemberCard 前端用户-会员卡关联表
type BeeUserMemberCard struct {
	ID         int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UserID     int64     `gorm:"column:user_id;comment:前端用户id" json:"user_id"`       // 前端用户id
	CardID     int32     `gorm:"column:card_id;comment:会员卡id" json:"card_id"`        // 会员卡id
	CreateTime time.Time `gorm:"column:create_time;comment:添加时间" json:"create_time"` // 添加时间
	Amount     float64   `gorm:"column:amount;comment:支付金额" json:"amount"`           // 支付金额
	Name       string    `gorm:"column:name;comment:卡片名称" json:"name"`               // 卡片名称
	ValidMonth int32     `gorm:"column:valid_month;comment:有效月数" json:"valid_month"` // 有效月数
	TotalCount int32     `gorm:"column:total_count;comment:总消费次数" json:"totalCount"` // 总消费次数
	LeftCount  int32     `gorm:"column:left_count;comment:剩余消费次数" json:"leftCount"`  // 剩余消费次数
	OutTradeNo string    `gorm:"column:out_trade_no;comment:订单号" json:"outTradeNo"`  // 订单号
}

// TableName BeeUserMemberCard's table name
func (*BeeUserMemberCard) TableName() string {
	return TableNameBeeUserMemberCard
}
