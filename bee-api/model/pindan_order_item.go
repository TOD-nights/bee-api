package model

import (
	"encoding/json"
	"fmt"

	"gitee.com/stuinfer/bee-api/common"
	"gitee.com/stuinfer/bee-api/enum"
	"github.com/shopspring/decimal"
)

// 拼单订单项
type BeePindanOrderItem struct {
	common.BaseModel
	Amount             decimal.Decimal  `gorm:"column:amount;type:decimal(10,2);comment:商品金额" json:"amount"`
	AmountVip          decimal.Decimal  `gorm:"column:amount_vip;type:decimal(10,2);comment:vip价格" json:"amountVip"`
	DatePay            common.JsonTime  `gorm:"column:date_pay;type:datetime(3);default:null;<-:create" json:"datePay"`
	GoodsNumber        int64            `gorm:"column:goods_number;type:bigint(11);comment:商品总数量" json:"goodsNumber"`
	HxNumber           string           `gorm:"uniqueIndex;column:hx_number;type:varchar(100);comment:核销码" json:"hxNumber"`
	Ip                 string           `gorm:"column:ip;type:varchar(100);comment:下单ip" json:"ip"`
	IsCanHx            bool             `gorm:"column:is_can_hx;type:tinyint(1);comment:是否可以核销" json:"isCanHx"`
	IsPay              bool             `gorm:"column:is_pay;type:tinyint(1);comment:是否已经支付" json:"isPay"`
	OrderNumber        string           `gorm:"uniqueIndex;column:order_number;type:varchar(100);comment:订单号" json:"orderNumber"`
	OrderType          enum.OrderType   `gorm:"column:order_type;type:bigint(11);comment:订单类型" json:"orderType"`
	Qudanhao           string           `gorm:"column:qudanhao;type:varchar(100);comment:取单号" json:"qudanhao"`
	Remark             string           `gorm:"column:remark;type:varchar(100);comment:备注" json:"remark"`
	ShopId             int64            `gorm:"column:shop_id;type:bigint(11);comment:商店id" json:"shopId"`
	ShopIdZt           int64            `gorm:"column:shop_id_zt;type:bigint(11);comment:自提商店id" json:"shopIdZt"`
	ShopNameZt         string           `gorm:"column:shop_name_zt;type:varchar(100);comment:自提商店名称" json:"shopNameZt"`
	Status             enum.OrderStatus `gorm:"column:status;type:bigint(11);comment:订单状态，1未发货" json:"status"`
	ExtJsonStr         string           `gorm:"column:ext_json_str;type:varchar(1000);comment:扩展信息" json:"extJsonStr"`
	PeisongType        int8             `gorm:"column:peisong_type;type:int(2);comment:配送类型  1.自提 2.配送" json:"peisongType"`
	PindanId           int64            `gorm:"column:pindan_id;type:bigint(20);comment:拼单id" json:"pindanId"`
	GoodsPropertyIds   string           `gorm:"column:goods_property_ids;type:varchar(255);comment:商品属性ids组合" json:"goodsPropertyIds"`
	GoodsPropertyNames string           `gorm:"column:goods_property_names;type:varchar(255);comment:商品属性名称" json:"goodsPropertyNames"`
	GoodsId            int64            `gorm:"column:goods_id;type:bigint(20);comment:商品id" json:"goodsId"`
}

func (m *BeePindanOrderItem) TableName() string {
	return "bee_pindan_order_item"
}

func (m *BeePindanOrderItem) ExtJsonToString() string {
	if m.ExtJsonStr == "" {
		return ""
	}
	var ma = make(map[string]interface{})
	_ = json.Unmarshal([]byte(m.ExtJsonStr), &ma)
	str := ""
	for k, v := range ma {
		str = str + k + ":" + fmt.Sprintf("%v", v) + ","
	}
	return str
}
