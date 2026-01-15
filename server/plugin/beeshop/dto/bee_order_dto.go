package dto

import (
	"fmt"

	"gitee.com/stuinfer/bee-api/common"
	"gitee.com/stuinfer/bee-api/enum"
	"gitee.com/stuinfer/bee-api/model"
	"github.com/shopspring/decimal"
)

type BeeOrderDto struct {
	model.BeeOrder
	PageNum  int `json:"pageNum" form:"pageNum"`
	PageSize int `json:"pageSize" form:"pageSize"`
}

// / 发起拼单请求
type CreatePindanReq struct {
	PindanId    int64                `json:"pindanId" form:"pindanId"`       // pindanId
	GoodsId     int64                `json:"goodsId" form:"goodsId"`         // 商品ID
	Number      int64                `json:"number" form:"number"`           // 商品数量
	Addition    string               `json:"addition" form:"addition"`       // 其他附加信息
	Sku         []CreatePindanSkuReq `json:"sku" form:"sku"`                 // 商品规格
	ShopId      int64                `json:"shopId" form:"shopId"`           // 门店ID
	PeisongType int8                 `json:"peisongType" form:"peisongType"` // 配送类型，1.自提 2.配送
}

// / 发起拼单请求 sku
type CreatePindanSkuReq struct {
	OptionId      int64 `json:"optionId" form:"optionId"`           // 规格ID
	OptionValueId int64 `json:"optionValueId" form:"optionValueId"` // 规格值Id
}

// / 将sku 结构体参数转化为字符串描述
func (m CreatePindanSkuReq) Desc() string {
	return fmt.Sprintf("%d:%d", m.OptionId, m.OptionValueId)
}

// / 查询拼单详情请求
type GetPindanInfoReq struct {
	ShopID int64 `json:"shopId" form:"shopId"` //门店ID
	UserID int64 `json:"userId" form:"userId"` //用户ID
	ID     int64 `json:"id" form:"id"`         //拼单ID
}

// / 拼单详细
type PindanResp struct {
	// 拼单信息
	Amount          decimal.Decimal `json:"amount"` // 总价格
	AmountCoupons   decimal.Decimal `json:"amountCoupons"`
	AmountLogistics decimal.Decimal `json:"amountLogistics"`
	AmountBalance   decimal.Decimal `json:"amountBalance"`
	AmountReal      decimal.Decimal `json:"amountReal"` // 除去优惠、余额抵扣，实际支付金额

	DatePay        common.JsonTime  `json:"datePay"`
	GoodsNumber    int64            `json:"goodsNumber"`
	DateClose      common.JsonTime  `json:"dateClose"`
	Ip             string           `json:"ip"`
	IsDelUser      bool             `json:"isDelUser"`
	IsEnd          bool             `json:"isEnd"`
	IsPay          bool             `json:"isPay"`
	OrderNumber    string           `json:"orderNumber"`
	OrderType      enum.OrderType   `json:"orderType"`
	Qudanhao       string           `json:"qudanhao"`
	Remark         string           `json:"remark"`
	ShopId         int64            `json:"shopId"`
	ShopIdZt       int64            `json:"shopIdZt"`
	ShopNameZt     string           `json:"shopNameZt"`
	Status         enum.OrderStatus `json:"status"`
	Uid            int64            `json:"uid"`
	ExtJsonStr     string           `json:"extJsonStr"`
	PeisongType    int8             `json:"peisongType"`
	IsVip          bool             `json:"isVip"`
	Id             int64            `json:"id"`
	UserId         int64            `json:"userId"` //主账号uid
	DateAdd        common.JsonTime  `son:"dateAdd"`
	DateUpdate     common.JsonTime  `json:"dateUpdate"`
	PindanShopInfo PindanShopInfo   `json:"shopInfo"` // 拼单店铺信息
	// 拼单项信息
	Items []*PindanItemResp `json:"items"`
}

type PindanItemResp struct {
	Id                 int64            `json:"id"`
	Amount             decimal.Decimal  `json:"amount"`
	AmountVip          decimal.Decimal  `json:"amountVip"`
	DatePay            common.JsonTime  `json:"datePay"`
	GoodsNumber        int64            `json:"goodsNumber"`
	HxNumber           string           `json:"hxNumber"`
	Ip                 string           `json:"ip"`
	IsCanHx            bool             `json:"isCanHx"`
	IsPay              bool             `json:"isPay"`
	OrderNumber        string           `json:"orderNumber"`
	OrderType          enum.OrderType   `json:"orderType"`
	Qudanhao           string           `json:"qudanhao"`
	Remark             string           `json:"remark"`
	ShopId             int64            `json:"shopId"`
	ShopIdZt           int64            `json:"shopIdZt"`
	ShopNameZt         string           `json:"shopNameZt"`
	Status             enum.OrderStatus `json:"status"`
	ExtJsonStr         string           `json:"extJsonStr"`
	PeisongType        int8             `json:"peisongType"`
	PindanId           int64            `json:"pindanId"`
	GoodsPropertyIds   string           `json:"goodsPropertyIds"`
	GoodsPropertyNames string           `json:"goodsPropertyNames"`
	GoodsId            int64            `json:"goodsId"`
	UserId             int64            `json:"userId"`
	GoodsInfo          PindanItemGoods  `json:"goodsInfo"`
	UserInfo           PinDanUserInfo   `json:"userInfo"`
}

type PindanItemGoods struct {
	Id   int64  `json:"id" mapstructure:"id"`
	Name string `json:"name" mapstructure:"name"`
	Pic  string `json:"pic" mapstructure:"pic"`
}

type PinDanUserInfo struct {
	Id        int64  `gorm:"column:id;" json:"id" mapstructure:"id"`                       // 用户ID
	Nick      string `gorm:"column:nick;" json:"nick" mapstructure:"nick"`                 // 昵称
	AvatarUrl string `gorm:"column:avatar_url;" json:"avatarUrl" mapstructure:"avatarUrl"` // 头像
	Mobile    string `gorm:"column:mobile;" json:"mobile" mapstructure:"mobile"`           // 手机号码
}

// / 拼单店铺信息
type PindanShopInfo struct {
	Id        int64   `gorm:"column:id;" json:"id" mapstructure:"id"`                       // 店铺ID
	Address   string  `gorm:"column:address;" json:"address" mapstructure:"address"`        // 店铺地址
	Name      string  `gorm:"column:name;" json:"name" mapstructure:"name"`                 // 店铺名称
	LinkPhone string  `gorm:"column:link_phone;" json:"linkPhone" mapstructure:"linkPhone"` // 店铺联系电话
	Latitude  float64 `gorm:"column:latitude;" json:"latitude" mapstructure:"latitude"`     // 店铺纬度
	Longitude float64 `gorm:"column:longitude;" json:"longitude" mapstructure:"longitude"`  // 店铺经度
}
