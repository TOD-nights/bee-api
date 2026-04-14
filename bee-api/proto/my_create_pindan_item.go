package proto

import (
	"gitee.com/stuinfer/bee-api/common"
	"gitee.com/stuinfer/bee-api/enum"
	"github.com/shopspring/decimal"
)

type MyCreatePindanItem struct {
	Id              int64           `json:"id" mapstructure:"id"`
	UserId          int64           `json:"userId" mapstructure:"userId"` //主账号uid
	DateAdd         common.JsonTime `json:"dateAdd" mapstructure:"dateAdd"`
	DateUpdate      common.JsonTime `json:"dateUpdate" mapstructure:"dateUpdate"`
	Amount          decimal.Decimal `json:"amount" mapstructure:"amount"` // 总价格
	AmountCoupons   decimal.Decimal `json:"amountCoupons" mapstructure:"amountCoupons"`
	AmountLogistics decimal.Decimal `json:"amountLogistics" mapstructure:"amountLogistics"`
	AmountBalance   decimal.Decimal `json:"amountBalance" mapstructure:"amountBalance"`
	AmountReal      decimal.Decimal `json:"amountReal" mapstructure:"amountReal"` // 除去优惠、余额抵扣，实际支付金额

	DatePay     common.JsonTime  `json:"datePay" mapstructure:"datePay"`         // 支付时间
	GoodsNumber int64            `json:"goodsNumber" mapstructure:"goodsNumber"` //商品数量
	Ip          string           `json:"ip" mapstructure:"ip"`
	IsPay       bool             `json:"isPay" mapstructure:"isPay"`
	OrderNumber string           `json:"orderNumber" mapstructure:"orderNumber"`
	ShopId      int64            `json:"shopId" mapstructure:"shopId"`
	ShopIdZt    int64            `json:"shopIdZt" mapstructure:"shopIdZt"`
	ShopNameZt  string           `json:"shopNameZt" mapstructure:"shopNameZt"`
	Status      enum.OrderStatus `json:"status" mapstructure:"status"` //0未支付  1待核销  2已完成
	ExtJsonStr  string           `json:"extJsonStr" mapstructure:"extJsonStr"`
	PeisongType int8             `json:"peisongType" mapstructure:"peisongType"` //1.自提 2.配送
	IsVip       bool             `json:"isVip" mapstructure:"isVip"`

	UserInfo PinDanUserInfo `json:"userInfo" mapstructure:"userInfo"` // 拼单发起人信息
	ShopInfo PindanShopInfo `json:"shopInfo" mapstructure:"shopInfo"` // 拼单门店信息
}
