package service

import (
	"context"
	"fmt"
	"slices"
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

		//查询在当前门店未完成的拼单
		if err := db.GetDB().Where("user_id = ? and shop_id = ? and status = 0 and is_deleted = 0 and is_del_user = 0 and is_end = 0 and is_pay = 0",
			kit.GetUid(c), goods.ShopId).First(&pindanRecord).Error; err == nil {
			return pindanRecord.Id, nil
		}

		// 校验商品属性
		var goodsProps []model.BeeShopGoodsProp
		// 查询商品sku信息
		var skuSelected model.BeeShopGoodsSku
		var property_child_ids = ""
		for _, sku := range reqVo.Sku {
			property_child_ids += sku.Desc() + ","
		}
		var propertyTx = db.GetDB().Where("is_deleted = 0")

		propertyTx.Scopes(func(db *gorm.DB) *gorm.DB {
			for _, sku := range reqVo.Sku {
				db.Or("id = ? and property_id = ?", sku.OptionValueId, sku.OptionId)
			}
			return db
		})
		if err := propertyTx.Find(&goodsProps).Error; err != nil {
			return 0, errors.New("商品属性信息不存在")
		}
		// property_child_ids = strings.TrimRight(property_child_ids, ",")
		logger.GetLogger().Debug("property_child_ids:", zap.String("property_child_ids", property_child_ids))
		if err := db.GetDB().Where("goods_id = ? and is_deleted = 0 and property_child_ids = ?", reqVo.GoodsId, property_child_ids).First(&skuSelected).Error; err != nil {
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

				Amount:           goods.OriginalPrice,
				AmountVip:        goods.VipPrice,
				GoodsNumber:      reqVo.Number,
				Ip:               ip,
				IsPay:            false,
				OrderNumber:      util.GenerateOrderNo("PD", kit.GetUid(c)),
				OrderType:        enum.OrderTypePindan,
				ShopId:           goods.ShopId,
				ShopIdZt:         goods.ShopId,
				ShopNameZt:       goods.Name,
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
