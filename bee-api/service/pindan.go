package service

import (
	"context"
	"sync"
	"time"

	"gitee.com/stuinfer/bee-api/db"
	"gitee.com/stuinfer/bee-api/enum"
	"gitee.com/stuinfer/bee-api/kit"
	"gitee.com/stuinfer/bee-api/model"
	"gitee.com/stuinfer/bee-api/proto"
	"github.com/samber/lo"
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
func (s *pindanSrv) CreatePindan(c context.Context, reqVo proto.CreatePindanReq) error {

	return db.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.BeeOrder{}).Where("id = ? and uid = ?", orderId, kit.GetUid(c)).Updates(map[string]interface{}{
			"status":      enum.OrderStatusClose,
			"date_update": time.Now(),
			"remark":      remark,
		}).Error; err != nil {
			return err
		}
		orderLog := &model.BeeOrderLog{
			BaseModel: *kit.GetInsertBaseModel(c),
			OrderId:   orderId,
			Type:      enum.OrderLogTypeClose,
		}
		if err := tx.Create(orderLog).Error; err != nil {
			return err
		}
		//@todo 退换优惠券、积分
		if orderInfo.Status == enum.OrderStatusUnPaid {
			couponIds := lo.Map(orderCoupons, func(item *model.BeeOrderCoupon, _ int) int64 {
				return item.CouponId
			})
			if len(couponIds) > 0 {
				if err := tx.Model(&model.BeeOrderCoupon{}).Where("coupon_id in ? and uid = ? and is_deleted = 0", couponIds, kit.GetUid(c)).Updates(map[string]interface{}{
					"is_deleted":  1,
					"date_delete": time.Now(),
				}).Error; err != nil {
					return err
				}
				if err := tx.Model(&model.BeeUserCoupon{}).Where("id in ? and is_deleted = 0", couponIds).Updates(map[string]interface{}{
					"status":      enum.CouponStatusNormal,
					"date_update": time.Now(),
				}).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}
