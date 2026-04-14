package service

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"gitee.com/stuinfer/bee-api/common"
	"gitee.com/stuinfer/bee-api/db"
	"gitee.com/stuinfer/bee-api/enum"
	"gitee.com/stuinfer/bee-api/kit"
	"gitee.com/stuinfer/bee-api/logger"
	"gitee.com/stuinfer/bee-api/model"
	"gitee.com/stuinfer/bee-api/proto"
	"gitee.com/stuinfer/bee-api/util"
	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type pindanSrv struct {
	BaseSrv
}

var pindanSrvOnce sync.Once
var pindanSrvInstance *pindanSrv

func GetPinDanServ() *pindanSrv {
	pindanSrvOnce.Do(func() {
		pindanSrvInstance = &pindanSrv{}
	})
	return pindanSrvInstance
}

func (s *pindanSrv) WxPay(ginCtx *gin.Context, pindanId int64, uid int64) (*proto.GetWxPayInfoRes, string, error) {

	var pindanRecord model.PinDanRecord
	var pindanRecordItems []model.BeePindanOrderItem
	type payAction struct {
		Type int32 `json:"type"`
		ID   int64 `json:"id"`
	}
	var wxPayConfig model.BeeWxPayConfig

	var payOrderId = util.GetRandInt64()
	if err := db.GetDB().Model(&model.PinDanRecord{BaseModel: common.BaseModel{Id: pindanId}}).First(&pindanRecord).Error; err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, payOrderId, err
	} else if user, err := GetUserSrv().GetUser(ginCtx); err != nil {
		return nil, payOrderId, errors.Wrap(err, "获取用户信息失败")
	} else if userAmount, err := GetUserSrv().Amount(ginCtx, user.Id); err != nil {
		return nil, payOrderId, errors.Wrap(err, "获取用户账户信息失败")
	} else if err := db.GetDB().Model(&model.BeePindanOrderItem{}).Where("pindan_id = ?", pindanId).Find(&pindanRecordItems).Error; err != nil {
		return nil, payOrderId, enum.NewBussErr(err, enum.ResCodeFail, "拼单项查询失败")
	} else if err := db.GetDB().Where("user_id = ? and is_deleted = 0", kit.GetUserId(ginCtx)).Take(&wxPayConfig).Error; err != nil {
		return nil, payOrderId, errors.Wrap(err, "获取微信配置失败！")
	} else if wxPayClient, err := GetPaySrv().GetWechatPayClient(ginCtx, &WxPayConfig{
		MchId:           wxPayConfig.MchId,
		Secret:          wxPayConfig.AppSecret,
		Token:           wxPayConfig.Token,
		ReturnUrl:       "",
		NotifyUrl:       GetPaySrv().GetWxPayNotifyUrl(ginCtx, &wxPayConfig),
		PrivateCertPath: wxPayConfig.PrivateCert, //"/Users/mac/project/bee-api/bee-api/lingwu_key.pem", //wxPayConfig.PrivateCert,
		Debug:           wxPayConfig.Debug,
	}); err != nil {
		return nil, payOrderId, errors.Wrap(err, "获取微信支付客户端失败！")
	} else if userOpenId, err := GetUserSrv().GetUserWxOpenId(ginCtx); err != nil {
		return nil, payOrderId, err
	} else {

		var totalOri = decimal.NewFromInt(0)
		var totalVip = decimal.NewFromInt(0)

		for _, v := range pindanRecordItems {
			totalVip = totalVip.Add(v.AmountVip.Mul(decimal.NewFromInt(v.GoodsNumber)))
			totalOri = totalOri.Add(v.Amount.Mul(decimal.NewFromInt(v.GoodsNumber)))
		}
		var total = decimal.NewFromInt(0)
		if user.VipLevel > 0 {
			total = totalVip
		} else {
			total = totalOri
		}

		if userAmount.Balance.GreaterThan(total) {
			//全部使用余额支付

			db.GetDB().Model(&model.BeeUserAmount{
				BaseModel: common.BaseModel{Id: userAmount.Id},
			}).Update("balance", gorm.Expr("balance - ?", total))
			total = decimal.NewFromInt(0)
			if err := db.GetDB().Transaction(func(tx *gorm.DB) error {
				payLog := &model.BeePayLog{
					BaseModel:   *kit.GetInsertBaseModel(ginCtx),
					Money:       total.Mul(decimal.NewFromInt(100)),
					NextAction:  fmt.Sprintf("{type: 16,id:%s}", payOrderId),
					OrderNo:     payOrderId,
					OrderNumber: strconv.FormatInt(int64(pindanId), 10),
					PayGate:     enum.PayGateWXAPP,
					Remark:      "拼单支付",
					Uid:         kit.GetUid(ginCtx),
					ShopId:      pindanRecord.ShopId,
					OrderType:   3,
					Status:      enum.PayLogStatusUnPaid,
				}
				if err := tx.Create(payLog).Error; err != nil {
					return err
				}
				//更新pindanRecord实付款

				if err := tx.Model(&model.PinDanRecord{BaseModel: common.BaseModel{Id: pindanRecord.Id}}).Updates(map[string]interface{}{
					"amount":      totalVip,
					"amount_real": total,
					"is_vip":      user.VipLevel > 0,
					"is_pay":      true,
				}).Error; err != nil {
					return err
				}
				return nil
			}); err != nil {
				return nil, payOrderId, errors.Wrap(err, "微信下单失败")
			} else {
				return nil, payOrderId, nil
			}
		} else {
			// 部分使用支付
			total = total.Sub(userAmount.Balance)
			if wxResp, err := wxPayClient.V3TransactionJsapi(ginCtx, gopay.BodyMap{
				"mchid":        wxPayConfig.MchId,
				"out_trade_no": payOrderId,
				"appid":        wxPayConfig.AppId,
				"description":  "拼单支付",
				"notify_url":   GetPaySrv().GetWxPayNotifyUrl(ginCtx, &wxPayConfig),
				"amount": map[string]interface{}{
					"total":    total.Mul(decimal.NewFromInt(100)).IntPart(),
					"currency": "CNY",
				},
				"time_expire": time.Now().Add(time.Hour * 1).Format(time.RFC3339), // 限制一小时内支付
				"payer": map[string]interface{}{
					"openid": userOpenId,
				},
			}); err != nil {
				return nil, payOrderId, err
			} else if wxResp.Code != 0 {
				return nil, payOrderId, errors.New("微信请求失败：" + wxResp.Error)
			} else if err := db.GetDB().Model(&model.BeeUserAmount{
				BaseModel: common.BaseModel{Id: userAmount.Id},
			}).Update("balance", 0).Error; err != nil {
				return nil, payOrderId, errors.Wrap(err, "更新用户余额失败")
			} else if jsapiSignInfo, err := wxPayClient.PaySignOfApplet(wxPayConfig.AppId, wxResp.Response.PrepayId); err != nil {
				return nil, payOrderId, errors.Wrap(err, "获取微信支付签名失败")
			} else if err := db.GetDB().Transaction(func(tx *gorm.DB) error {
				payLog := &model.BeePayLog{
					BaseModel:   *kit.GetInsertBaseModel(ginCtx),
					Money:       total.Mul(decimal.NewFromInt(100)),
					NextAction:  fmt.Sprintf("{type: 16,id:%s}", payOrderId),
					OrderNo:     payOrderId,
					OrderNumber: strconv.FormatInt(int64(pindanId), 10),
					PayGate:     enum.PayGateWXAPP,
					Remark:      "拼单支付",
					Uid:         kit.GetUid(ginCtx),
					ShopId:      pindanRecord.ShopId,
					OrderType:   3,
					Status:      enum.PayLogStatusUnPaid,
				}
				if err := tx.Create(payLog).Error; err != nil {
					return err
				}
				//更新pindanRecord实付款
				if err := tx.Model(&model.PinDanRecord{BaseModel: common.BaseModel{Id: pindanRecord.Id}}).Updates(map[string]interface{}{
					"amount":      totalVip,
					"amount_real": total,
					"is_vip":      user.VipLevel > 0,
					"is_pay":      false,
				}).Error; err != nil {
					return err
				}
				return nil
			}); err != nil {
				return nil, payOrderId, errors.Wrap(err, "微信下单失败")
			} else {

				return &proto.GetWxPayInfoRes{
					TimeStamp:  jsapiSignInfo.TimeStamp,
					OutTradeId: payOrderId,
					Package:    jsapiSignInfo.Package,
					PaySign:    jsapiSignInfo.PaySign,
					Appid:      wxPayConfig.AppId,
					Sign:       jsapiSignInfo.PaySign,
					SignType:   jsapiSignInfo.SignType,
					PrepayId:   payOrderId,
					PayAmount:  total,
					NonceStr:   jsapiSignInfo.NonceStr,
				}, payOrderId, nil
			}
		}

	}
}

// / 查询我发起的拼单
func (s *pindanSrv) GetMyCreatedPindanRecord(page int64, status int64, uid int64) ([]proto.MyCreatePindanItem, error) {

	if page <= 0 {
		page = 1
	}
	tx := db.GetDB().Model(&model.PinDanRecord{}).Where("user_id = ?", uid)
	if status > -1 {
		tx.Where("status = ?", status)
	}
	tx.Offset((int(page) - 1) * 10).Limit(10)
	var records []model.PinDanRecord
	if err := tx.Find(&records).Error; err != nil {
		return nil, err
	}
	var res []proto.MyCreatePindanItem
	for _, item := range records {
		var target proto.MyCreatePindanItem

		if err := copier.Copy(&target, item); err == nil {
			target.ShopInfo = proto.PindanShopInfo{Name: "测试"}
			res = append(res, target)
		} else {
			fmt.Println(err)
		}
	}

	return res, nil
}
func (s *pindanSrv) JoinPindan(c context.Context, ip string, reqVo *proto.CreatePindanReq) (int64, error) {

	var pindanRecord model.PinDanRecord
	var user model.BeeUser
	var shopInfo model.BeeShopInfo
	var goods model.BeeShopGoods

	// 查询用户信息
	if err := db.GetDB().Model(&model.BeeUser{}).First(&user, kit.GetUid(c)).Error; err != nil {
		return 0, errors.New("用户信息不存在")
	} else if err := db.GetDB().Model(&model.PinDanRecord{}).First(&pindanRecord, reqVo.PindanId).Error; err != nil {
		return 0, errors.New("拼单信息不存在")
	} else if err := db.GetDB().Where("id = ? and is_deleted = 0", reqVo.GoodsId).First(&goods).Error; err != nil {
		return 0, errors.New("商品信息不存在")
	} else if err := db.GetDB().Where("id = ? and is_deleted = 0", pindanRecord.ShopId).First(&shopInfo).Error; err != nil {
		return 0, errors.New("门店信息不存在")
	} else {
		// 校验商品属性
		var goodsProps []model.BeeShopGoodsProp
		// 查询商品sku信息
		var skuSelected model.BeeShopGoodsSku
		var property_child_ids = ""
		for _, sku := range reqVo.Sku {
			property_child_ids += sku.Desc() + ","
		}
		var propertyTx = db.GetDB().Where("is_deleted = 0")
		// propertyTx.Where()
		var propertyW = make([]string, len(reqVo.Sku))
		var args = make([]interface{}, len(reqVo.Sku)*2)
		for index, sku := range reqVo.Sku {
			propertyW[index] = "(?,?)"
			args[index*2] = sku.OptionValueId
			args[index*2+1] = sku.OptionId
		}

		var condition = fmt.Sprintf("(id,property_id) in (%s)", strings.Join(propertyW, ","))
		propertyTx.Where(condition, args...)
		if err := propertyTx.Find(&goodsProps).Error; err != nil {
			return 0, errors.New("商品属性信息不存在")
		}
		if err := db.GetDB().Where("goods_id = ? and is_deleted = 0 and property_child_ids = ?", reqVo.GoodsId, property_child_ids).First(&skuSelected).Error; err != nil {
			return 0, errors.New("商品规格信息不存在")
		}
		if err := db.GetDB().Transaction(func(tx *gorm.DB) error {
			// 1. 更新拼单信息, 更新拼单商品数量,商品价格
			if err := tx.Model(&model.PinDanRecord{BaseModel: common.BaseModel{Id: pindanRecord.BaseModel.Id}}).Updates(map[string]interface{}{
				"amount":       pindanRecord.Amount.Add(skuSelected.Price.Mul(decimal.NewFromInt(reqVo.Number))),
				"goods_number": pindanRecord.GoodsNumber + reqVo.Number,
			}).Error; err != nil {
				return err
			}
			//检查是否已经参与了拼单
			var count int64
			if err := tx.Model(&model.BeePindanOrderItem{}).Where("user_id = ? and pindan_id = ?", kit.GetUid(c), pindanRecord.Id).Count(&count).Error; err != nil {
				return err
			}

			if count > 0 {
				return &enum.BussError{Code: 100000, Message: "已加入拼单,无需重复加入"}
			}
			// 2. 为发起人创建一个新的订单项
			pindanOrderItem := model.BeePindanOrderItem{
				BaseModel: common.BaseModel{
					UserId:     kit.GetUid(c),
					IsDeleted:  false,
					DateAdd:    common.JsonTime(time.Now()),
					DateUpdate: common.JsonTime(time.Now()),
				},

				Amount:           skuSelected.Price,
				AmountVip:        skuSelected.VipPrice,
				GoodsNumber:      reqVo.Number,
				Ip:               ip,
				IsPay:            false,
				OrderNumber:      util.GenerateOrderNo("PD", kit.GetUid(c)),
				OrderType:        enum.OrderTypePindan,
				ShopId:           goods.ShopId,
				ShopIdZt:         goods.ShopId,
				ShopNameZt:       shopInfo.Name,
				Status:           enum.OrderStatusUnPaid,
				PeisongType:      int8(enum.OrderPeisongTypeZiti),
				PindanId:         pindanRecord.Id,
				GoodsId:          goods.Id,
				GoodsPropertyIds: property_child_ids,
				GoodsPropertyNames: strings.Join(lo.Map(goodsProps, func(prop model.BeeShopGoodsProp, index int) string {
					return prop.Name
				}), ","),
			}

			if err := tx.Create(&pindanOrderItem).Error; err != nil {
				zap.Error(err)
				return errors.New("创建拼单订单项失败")
			}
			return nil

		}); err != nil {
			fmt.Println(err)
			return 0, err
		}
		return pindanRecord.Id, nil
	}
}

// /  创建拼单
func (s *pindanSrv) CreatePindan(c context.Context, ip string, reqVo *proto.CreatePindanReq) (int64, error) {

	var pindanRecord model.PinDanRecord
	var user model.BeeUser
	var shopInfo model.BeeShopInfo

	// 查询用户信息
	if err := db.GetDB().Model(&model.BeeUser{}).First(&user, kit.GetUid(c)).Error; err != nil {
		return 0, errors.New("用户信息不存在")
	}

	if reqVo.GoodsId > 0 {
		// 查询商品信息
		var goods model.BeeShopGoods
		if err := db.GetDB().Where("id = ? and is_deleted = 0", reqVo.GoodsId).First(&goods).Error; err != nil {
			return 0, errors.New("商品信息不存在")
		}

		//查询门店信息
		if err := db.GetDB().Where("id = ? and is_deleted = 0", goods.ShopId).First(&shopInfo).Error; err != nil {
			return 0, errors.New("门店信息不存在")
		}

		// //查询在当前门店未完成的拼单
		// if err := db.GetDB().Where("user_id = ? and shop_id = ? and status = 0 and is_deleted = 0 and is_del_user = 0 and is_end = 0 and is_pay = 0",
		// 	kit.GetUid(c), goods.ShopId).First(&pindanRecord).Error; err == nil {
		// 	return pindanRecord.Id, nil
		// }

		// 校验商品属性
		var goodsProps []model.BeeShopGoodsProp
		// 查询商品sku信息
		var skuSelected model.BeeShopGoodsSku
		var property_child_ids = ""
		for _, sku := range reqVo.Sku {
			property_child_ids += sku.Desc() + ","
		}
		var propertyTx = db.GetDB().Where("is_deleted = 0")

		var subQueryDb = db.GetDB()
		for _, sku := range reqVo.Sku {
			subQueryDb = subQueryDb.Or("id = ? and property_id = ?", sku.OptionValueId, sku.OptionId)
		}
		propertyTx.Where(subQueryDb)
		if err := propertyTx.Find(&goodsProps).Error; err != nil {
			return 0, errors.New("商品属性信息不存在")
		}
		// property_child_ids = strings.TrimRight(property_child_ids, ",")
		logger.GetLogger().Debug("property_child_ids:", zap.String("property_child_ids", property_child_ids))
		tx := db.GetDB().Where("goods_id = ? and is_deleted = 0", reqVo.GoodsId)
		for _, v := range strings.Split(property_child_ids, ",") {
			tx.Where("property_child_ids like ?", "%"+v+"%")
		}
		if err := tx.First(&skuSelected).Error; err != nil {
			return 0, errors.New("商品规格信息不存在")
		}

		err := db.GetDB().Transaction(func(tx *gorm.DB) error {

			// 1. 创建新的拼单记录
			pindanRecord = model.PinDanRecord{
				BaseModel: common.BaseModel{
					UserId:     kit.GetUid(c),
					IsDeleted:  false,
					DateAdd:    common.JsonTime(time.Now()),
					DateUpdate: common.JsonTime(time.Now()),
				},
				Amount:          skuSelected.Price,
				AmountCoupons:   decimal.NewFromInt(0),
				AmountLogistics: decimal.NewFromInt(0),
				AmountBalance:   decimal.NewFromInt(0),
				AmountReal:      decimal.NewFromInt(0),
				GoodsNumber:     reqVo.Number,
				Ip:              ip,
				IsDelUser:       false,
				IsEnd:           false,
				IsPay:           false,
				ShopId:          goods.ShopId,
				ShopIdZt:        goods.ShopId,
				ShopNameZt:      shopInfo.Name,
				Status:          enum.OrderStatusUnPaid,
				IsVip:           user.VipLevel > 0,
			}
			if err := tx.Create(&pindanRecord).Error; err != nil {
				zap.Error(err)
				return errors.New("创建拼单失败")
			}

			// 2. 为发起人创建一个新的订单项
			pindanOrderItem := model.BeePindanOrderItem{
				BaseModel: common.BaseModel{
					UserId:     kit.GetUid(c),
					IsDeleted:  false,
					DateAdd:    common.JsonTime(time.Now()),
					DateUpdate: common.JsonTime(time.Now()),
				},

				Amount:           skuSelected.Price,
				AmountVip:        skuSelected.VipPrice,
				GoodsNumber:      reqVo.Number,
				Ip:               ip,
				IsPay:            false,
				OrderNumber:      util.GenerateOrderNo("PD", kit.GetUid(c)),
				OrderType:        enum.OrderTypePindan,
				ShopId:           goods.ShopId,
				ShopIdZt:         goods.ShopId,
				ShopNameZt:       shopInfo.Name,
				Status:           enum.OrderStatusUnPaid,
				PeisongType:      int8(enum.OrderPeisongTypeZiti),
				PindanId:         pindanRecord.Id,
				GoodsId:          goods.Id,
				GoodsPropertyIds: property_child_ids,
				GoodsPropertyNames: strings.Join(lo.Map(goodsProps, func(prop model.BeeShopGoodsProp, index int) string {
					return prop.Name
				}), ","),
			}

			if err := tx.Create(&pindanOrderItem).Error; err != nil {
				zap.Error(err)
				return errors.New("创建拼单订单项失败")
			}

			return nil
		})
		return pindanRecord.Id, err
	}

	logger.GetLogger().Error("创建拼单", zap.Int64("ShopId", reqVo.ShopId))
	if reqVo.ShopId > 0 {
		//查询门店信息
		if err := db.GetDB().Where("id = ? and is_deleted = 0", reqVo.ShopId).First(&shopInfo).Error; err != nil {
			return 0, errors.New("门店信息不存在")
		}

		//查询在当前门店未完成的拼单
		if err := db.GetDB().Where("user_id = ? and shop_id = ? and status = 0 and is_deleted = 0 and is_del_user = 0 and is_end = 0 and is_pay = 0",
			kit.GetUid(c), reqVo.ShopId).First(&pindanRecord).Error; err == nil {
			return pindanRecord.Id, nil
		}

		// 添加新的拼单信息
		err := db.GetDB().Transaction(func(tx *gorm.DB) error {

			// 1. 创建新的拼单记录
			pindanRecord = model.PinDanRecord{
				BaseModel: common.BaseModel{
					UserId:     kit.GetUid(c),
					IsDeleted:  false,
					DateAdd:    common.JsonTime(time.Now()),
					DateUpdate: common.JsonTime(time.Now()),
				},
				AmountCoupons:   decimal.NewFromInt(0),
				AmountLogistics: decimal.NewFromInt(0),
				AmountBalance:   decimal.NewFromInt(0),
				AmountReal:      decimal.NewFromInt(0),
				GoodsNumber:     reqVo.Number,
				Ip:              ip,
				IsDelUser:       false,
				IsEnd:           false,
				IsPay:           false,
				ShopId:          reqVo.ShopId,
				ShopIdZt:        reqVo.ShopId,
				ShopNameZt:      shopInfo.Name,
				Status:          enum.OrderStatusUnPaid,
				IsVip:           user.VipLevel > 0,
			}
			if err := tx.Create(&pindanRecord).Error; err != nil {
				zap.Error(err)
				return errors.New("创建拼单失败")
			}

			return nil
		})
		return pindanRecord.Id, err
	}
	return pindanRecord.Id, errors.New("参数错误，商品ID和店铺ID不能同时为空")
}

// 查询商品sku信息
// 根据拼单ID查询拼单信息
func (s *pindanSrv) GetPindanInfoByUserIdAndId(c context.Context, userId int64, id int64) (proto.PindanResp, error) {
	var resp proto.PindanResp
	var pindanRecord model.PinDanRecord
	if err := db.GetDB().Where("id = ? and status = 0 and is_deleted = 0 and is_del_user = 0 and is_end = 0 and is_pay = 0",
		id).First(&pindanRecord).Error; err != nil {
		return resp, &enum.BussError{Code: 500001, Message: "拼单信息不存在"}
	}

	copier.Copy(&resp, &pindanRecord)

	// 校验门店信息
	var shopInfo model.BeeShopInfo
	if err := db.GetDB().Where("id = ? and is_deleted = 0", resp.ShopId).First(&shopInfo).Error; err != nil {
		return resp, errors.New("门店信息不存在")
	}
	copier.Copy(&resp.PindanShopInfo, &shopInfo)

	// 查询拼单项列表
	var pindanOrderItems []*model.BeePindanOrderItem
	if err := db.GetDB().Where("pindan_id = ? and is_deleted = 0", pindanRecord.Id).Find(&pindanOrderItems).Error; err != nil {
		return resp, errors.New("拼单项信息查询失败")
	}
	if len(pindanOrderItems) > 0 {
		copier.Copy(&resp.Items, &pindanOrderItems)

		// 查询商品信息
		var goods []model.BeeShopGoods
		var users []model.BeeUser
		var goodInfos []proto.PindanItemGoods
		var goodIds = make([]int64, 0)
		var userIds = make([]int64, 0)
		for _, item := range pindanOrderItems {
			if slices.Index(goodIds, item.GoodsId) == -1 {
				goodIds = append(goodIds, item.GoodsId)
			}
			if slices.Index(userIds, item.UserId) == -1 {
				userIds = append(userIds, item.UserId)
			}
		}
		if err := db.GetDB().Where("id in ? and is_deleted = 0", goodIds).Find(&goods).Error; err != nil {
			return resp, errors.New("商品信息不存在")
		}

		if err := db.GetDB().Where("id in ? and is_deleted = 0", userIds).Find(&users).Error; err != nil {
			return resp, errors.New("用户信息不存在")
		}
		for _, g := range goods {
			var goodInfo proto.PindanItemGoods
			copier.Copy(&goodInfo, &g)
			goodInfos = append(goodInfos, goodInfo)
		}

		var pindanUsers = lo.Map(users, func(item model.BeeUser, index int) proto.PinDanUserInfo {
			var userInfo proto.PinDanUserInfo
			copier.Copy(&userInfo, &item)
			return userInfo
		})

		for _, item := range resp.Items {
			if one, found := lo.Find(goodInfos, func(g proto.PindanItemGoods) bool {
				return g.Id == item.GoodsId
			}); found {
				item.GoodsInfo = one
			}

			if one, found := lo.Find(pindanUsers, func(g proto.PinDanUserInfo) bool {
				return g.Id == item.UserId
			}); found {
				item.UserInfo = one
			}
		}
		fmt.Println("resp.Items:", resp.Items)
	}
	return resp, nil
}

// / 根据用户ID和店铺ID获取拼单信息
func (s *pindanSrv) GetPindanInfoByUserIdAndShopId(c context.Context, shopId int64, userId int64) (proto.PindanResp, error) {
	var resp proto.PindanResp
	var pindanRecord model.PinDanRecord
	if err := db.GetDB().Where("user_id = ? and shop_id = ? and status = 0 and is_deleted = 0 and is_del_user = 0 and is_end = 0 and is_pay = 0",
		userId, shopId).First(&pindanRecord).Error; err != nil {
		return resp, &enum.BussError{Code: 500001, Message: "拼单信息不存在"}
	}

	copier.Copy(&resp, &pindanRecord)

	// 查询拼单项列表
	var pindanOrderItems []*model.BeePindanOrderItem
	if err := db.GetDB().Where("pindan_id = ? and is_deleted = 0", pindanRecord.Id).Find(&pindanOrderItems).Error; err != nil {
		return resp, errors.New("拼单项信息查询失败")
	}
	if len(pindanOrderItems) > 0 {
		copier.Copy(&resp.Items, &pindanOrderItems)

		// 查询商品信息
		var goods []model.BeeShopGoods
		var users []model.BeeUser
		var goodInfos []proto.PindanItemGoods
		var goodIds = make([]int64, 0)
		var userIds = make([]int64, 0)
		for _, item := range pindanOrderItems {
			if slices.Index(goodIds, item.GoodsId) == -1 {
				goodIds = append(goodIds, item.GoodsId)
			}
			if slices.Index(userIds, item.UserId) == -1 {
				userIds = append(userIds, item.UserId)
			}
		}
		if err := db.GetDB().Where("id in ? and is_deleted = 0", goodIds).Find(&goods).Error; err != nil {
			return resp, errors.New("商品信息不存在")
		}

		if err := db.GetDB().Where("id in ? and is_deleted = 0", userIds).Find(&users).Error; err != nil {
			return resp, errors.New("用户信息不存在")
		}
		for _, g := range goods {
			var goodInfo proto.PindanItemGoods
			copier.Copy(&goodInfo, &g)
			goodInfos = append(goodInfos, goodInfo)
		}

		var pindanUsers = lo.Map(users, func(item model.BeeUser, index int) proto.PinDanUserInfo {
			var userInfo proto.PinDanUserInfo
			copier.Copy(&userInfo, &item)
			return userInfo
		})

		for _, item := range resp.Items {
			if one, found := lo.Find(goodInfos, func(g proto.PindanItemGoods) bool {
				return g.Id == item.GoodsId
			}); found {
				item.GoodsInfo = one
			}

			if one, found := lo.Find(pindanUsers, func(g proto.PinDanUserInfo) bool {
				return g.Id == item.UserId
			}); found {
				item.UserInfo = one
			}
		}

		fmt.Println("resp.Items:", resp.Items)
	}
	return resp, nil
}

// / 更新拼单的配送方式
func (s *pindanSrv) UpdatePindanPeisongTypeById(id int64) error {
	var pindanRecord model.PinDanRecord
	if err := db.GetDB().Model(&model.PinDanRecord{}).First(&pindanRecord, id).Error; err != nil {
		return err
	} else if err := db.GetDB().Model(&model.PinDanRecord{BaseModel: common.BaseModel{Id: id}}).Updates(map[string]interface{}{
		"peisong_type": lo.Ternary(pindanRecord.PeisongType == 1, 2, 1),
		"date_update":  common.JsonTime(time.Now()),
	}).Error; err != nil {
		return err
	}

	return nil
}

// / 我参与的拼单列表
func (s pindanSrv) GetMyJoinedPindanRecord(page int64, status int64, uid int64) ([]proto.MyJoinedPindanItem, error) {
	var result []proto.MyJoinedPindanItem
	var tx = db.GetDB().Model(&model.BeePindanOrderItem{}).Where("bee_pindan_order_item.user_id = ?", uid)

	tx = tx.Joins("left join bee_shop_info a on a.id = bee_pindan_order_item.shop_id  ").
		Joins("left join bee_shop_goods b on b.id = bee_pindan_order_item.goods_id").
		Select("bee_pindan_order_item.*,a.id as shopinfo_id,a.name as shopinfo_name,b.id as goodsinfo_id,b.name as goodsinfo_name,b.pic as goodsinfo_pic")

	if status > -1 {
		tx = tx.Where("bee_pindan_order_item.status = ?", status)
	}
	//分页
	if page <= 0 {
		page = 1
	}
	err := tx.Offset((int(page) - 1) * 10).Limit(10).Scan(&result).Error
	return result, err
}

// 取单
func (s pindanSrv) Qudan(id int64, userId int64) error {

	var pindanItem model.BeePindanOrderItem
	if err := db.GetDB().Model(&model.BeePindanOrderItem{}).First(&pindanItem, id).Error; err != nil {
		return &enum.BussError{Code: 10000, Message: "拼单不存在", Err: err}
	} else if pindanItem.UserId != userId {
		return &enum.BussError{Code: 100001, Message: "只能取自己的拼单"}
	} else if pindanItem.Status != enum.OrderStatusPaid {
		return &enum.BussError{Code: 100002, Message: "拼单状态不合法"}
	} else {
		db.GetDB().Model(&model.BeePindanOrderItem{BaseModel: common.BaseModel{Id: id}}).Updates(map[string]interface{}{
			"status":     enum.OrderStatusQudan,
			"qudan_time": time.Now(),
		})

		// 如果所有拼单项都已取单,更新拼单状态
		var unQudanCount int64
		if err := db.GetDB().Model(&model.BeePindanOrderItem{}).Where("status !=  ?", enum.OrderStatusQudan).Count(&unQudanCount).Error; err == nil {
			if unQudanCount == 0 {
				// 都已取单
				db.GetDB().Model(&model.PinDanRecord{BaseModel: common.BaseModel{Id: pindanItem.PindanId}}).Update("status", enum.OrderStatusQudan)
			}
		}

		return nil
	}
}
