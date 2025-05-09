package model

import (
	"gitee.com/stuinfer/bee-api/common"
	"gitee.com/stuinfer/bee-api/enum"
	"github.com/shopspring/decimal"
)

type BeePayLog struct {
	common.BaseModel
	Money        decimal.Decimal   `gorm:"column:money;type:decimal(10,2);comment:总金额" json:"money"`
	NextAction   string            `gorm:"column:next_action;type:varchar(100);comment:next_action" json:"nextAction"`
	OrderNo      string            `gorm:"column:order_no;type:varchar(100);comment:订单号" json:"orderNo"`
	OrderNumber  string            `gorm:"column:order_number;type:varchar(255);comment:order 表的orderNumber订单号" json:"orderNumber"`
	ThirdOrderNo string            `gorm:"column:third_order_no;type:varchar(200);comment:第三方订单号" json:"thirdOrderNo"`
	PayGate      enum.PayGate      `gorm:"column:pay_gate;type:varchar(100);comment:pay_gate" json:"payGate"`
	PayGateStr   string            `gorm:"-" json:"payGateStr"`
	Remark       string            `gorm:"column:remark;type:varchar(100);comment:备注" json:"remark"`
	Status       enum.PayLogStatus `gorm:"column:status;type:int(11);comment:状态" json:"status"`
	StatusStr    string            `gorm:"-" json:"statusStr"`
	Uid          int64             `gorm:"column:uid;type:bigint(11);comment:用户id" json:"uid"`
	ShopId       int64             `gorm:"column:shop_id;type:bigint(11);comment:所属门店" json:"shopId"`
	OrderType    int32             `gorm:"column:order_type;type:int(11);comment:订单类型,默认值0   0消费订单   1会员卡订单" json:"orderType"`
	MemberCardId int64             `gorm:"column:member_card_id;type:bigint(11);comment:会员卡id,默认值0 " json:"memberCardId"`
}

func (b *BeePayLog) TableName() string {
	return "bee_pay_log"
}

func (b *BeePayLog) FillData() {
	b.StatusStr = enum.PayLogStatusMap[b.Status]
	b.PayGateStr = enum.PayGateMap[b.PayGate]
}
