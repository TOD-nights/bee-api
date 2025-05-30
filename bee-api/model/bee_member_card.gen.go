package model

import (
	"time"
)

const TableNameBeeMemberCard = "bee_member_card"

// BeeMemberCard 会员卡
type BeeMemberCard struct {
	ID         int32     `gorm:"column:id;primaryKey" json:"id"`
	Name       string    `gorm:"column:name;comment:会员卡名称" json:"name"`                            // 会员卡名称
	SortNum    int32     `gorm:"column:sort_num;comment:排序" json:"sortNum"`                        // 排序
	Amount     float64   `gorm:"column:amount;comment:金额" json:"amount"`                           // 金额
	CreateTime time.Time `gorm:"column:create_time;comment:创建时间" json:"create_time"`               // 创建时间
	UpdateTime time.Time `gorm:"column:update_time;comment:更新时间" json:"update_time"`               // 更新时间
	DeleteFlag bool      `gorm:"column:delete_flag;comment:删除标识   0未删除  1已删除" json:"deleteFlag"`   // 删除标识   0未删除  1已删除
	DeleteTime time.Time `gorm:"column:delete_time;comment:删除时间" json:"delete_time"`               // 删除时间
	ValidMonth int32     `gorm:"column:valid_month;comment:有效期月数 用于区分月卡,年卡,半年卡" json:"validMonth"` // 有效期月数 用于区分月卡,年卡,半年卡
}

// TableName BeeMemberCard's table name
func (*BeeMemberCard) TableName() string {
	return TableNameBeeMemberCard
}
