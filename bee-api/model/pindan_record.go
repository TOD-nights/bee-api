package model

import (
	"gitee.com/stuinfer/bee-api/common"
	"gitee.com/stuinfer/bee-api/enum"
	"github.com/shopspring/decimal"
)

// 拼单记录表
type PinDanRecord struct {
	common.BaseModel
	Amount          decimal.Decimal `gorm:"column:amount;type:decimal(10,2);comment:商品金额" json:"amount"` // 总价格
	AmountCoupons   decimal.Decimal `gorm:"column:amount_coupons;type:decimal(10,2);comment:优惠券抵扣" json:"amountCoupons"`
	AmountLogistics decimal.Decimal `gorm:"column:amount_logistics;type:decimal(10,2);comment:运费" json:"amountLogistics"`
	AmountBalance   decimal.Decimal `gorm:"column:amount_balance;type:decimal(10,2);comment:余额抵扣" json:"amountBalance"`
	AmountReal      decimal.Decimal `gorm:"column:amount_real;type:decimal(10,2);comment:付款金额" json:"amountReal"` // 除去优惠、余额抵扣，实际支付金额

	DateClose   common.JsonTime  `gorm:"column:date_close;type:datetime(3);comment:关闭订单时间" json:"dateClose"`
	DatePay     common.JsonTime  `gorm:"column:date_pay;type:datetime(3);comment:支付订单时间" json:"datePay"`
	GoodsNumber int64            `gorm:"column:goods_number;type:bigint(11);comment:商品总数量" json:"goodsNumber"`
	Ip          string           `gorm:"column:ip;type:varchar(100);comment:下单ip" json:"ip"`
	IsDelUser   bool             `gorm:"column:is_del_user;type:tinyint(1);comment:用户删除" json:"isDelUser"`
	IsEnd       bool             `gorm:"column:is_end;type:tinyint(1);comment:订单已经结束" json:"isEnd"`
	IsPay       bool             `gorm:"column:is_pay;type:tinyint(1);comment:是否已经支付" json:"isPay"`
	OrderNumber string           `gorm:"uniqueIndex;column:order_number;type:varchar(100);comment:订单号" json:"orderNumber"`
	OrderType   enum.OrderType   `gorm:"column:order_type;type:bigint(11);comment:订单类型" json:"orderType"`
	Remark      string           `gorm:"column:remark;type:varchar(100);comment:备注" json:"remark"`
	ShopId      int64            `gorm:"column:shop_id;type:bigint(11);comment:商店id" json:"shopId"`
	ShopIdZt    int64            `gorm:"column:shop_id_zt;type:bigint(11);comment:自提商店id" json:"shopIdZt"`
	ShopNameZt  string           `gorm:"column:shop_name_zt;type:varchar(100);comment:自提商店名称" json:"shopNameZt"`
	Status      enum.OrderStatus `gorm:"column:status;type:bigint(11);comment:订单状态，0未发货" json:"status"`
	ExtJsonStr  string           `gorm:"column:ext_json_str;type:varchar(1000);comment:扩展信息" json:"extJsonStr"`
	PeisongType int8             `gorm:"column:peisong_type;type:int(2);comment:配送类型, 1.自提 2.配送" json:"peisongType"`
	IsVip       bool             `gorm:"column:is_vip;type:tinyint(1);comment:是否是vip订单" json:"isVip"`
}

func (b *PinDanRecord) TableName() string {
	return "bee_pindan_record"
}
