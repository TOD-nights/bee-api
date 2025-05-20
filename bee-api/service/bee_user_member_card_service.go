package service

import (
	"context"
	"gitee.com/stuinfer/bee-api/db"
	"gitee.com/stuinfer/bee-api/enum"
	"gitee.com/stuinfer/bee-api/kit"
	"gitee.com/stuinfer/bee-api/logger"
	"gitee.com/stuinfer/bee-api/model"
	"gitee.com/stuinfer/bee-api/proto"
	"gitee.com/stuinfer/bee-api/util"
	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"strconv"
	"sync"
	"time"
)

type userMemberCardService struct {
}

var UserMemberCardService = &userMemberCardService{}

// 余额支付
func (s *userMemberCardService) BalancePay(ginCtx *gin.Context, shopId int64, memberCardId int32) error {
	var memberCard model.BeeMemberCard
	if err := db.GetDB().Model(&model.BeeMemberCard{}).First(&memberCard, memberCardId).Error; err != nil {
		logger.GetLogger().Error(err.Error())
		return err
	} else if amount, err := GetBalanceSrv().GetAmount(ginCtx, kit.GetUid(ginCtx)); err != nil {
		logger.GetLogger().Error(err.Error())
		return err
	} else if decimal.NewFromFloat(memberCard.Amount).GreaterThan(amount.Balance) {
		return errors.New("余额不足")
	} else {

		err = db.GetDB().Transaction(func(tx *gorm.DB) error {

			amountBalance := amount.Balance.Sub(decimal.NewFromFloat(memberCard.Amount))
			// 更新支付状态
			// 扣除余额
			//amount.Balance = amountBalance
			//更新余额
			tx.Model(&model.BeeUserAmount{}).Where("uid", kit.GetUid(ginCtx)).Update("balance", amountBalance)

			// 添加user-memberCart

			payLog := &model.BeePayLog{
				BaseModel:    *kit.GetInsertBaseModel(ginCtx),
				Money:        decimal.NewFromFloat(memberCard.Amount),
				NextAction:   "",
				OrderNo:      util.GetRandInt64(),
				OrderNumber:  strconv.FormatInt(int64(memberCard.ID), 10),
				PayGate:      enum.PayGateBalance,
				Remark:       "购买会员卡",
				Uid:          kit.GetUid(ginCtx),
				ShopId:       shopId,
				OrderType:    1,
				Status:       enum.PayLogStatusPaid,
				MemberCardId: int64(memberCardId),
			}
			if err = tx.Create(payLog).Error; err != nil {
				return err
			}

			var targetTime = memberCard.CreateTime.AddDate(0, int(memberCard.ValidMonth), 0)
			var days = int(targetTime.Sub(memberCard.CreateTime).Hours() / 24)
			var userMemberCard = model.BeeUserMemberCard{
				UserID:     kit.GetUid(ginCtx), // 前端用户id
				CardID:     memberCardId,       // 会员卡id
				CreateTime: time.Now(),
				Amount:     memberCard.Amount,
				Name:       memberCard.Name, // 卡片名称
				ValidMonth: memberCard.ValidMonth,
				TotalCount: int32(days),
				LeftCount:  int32(days),
				OutTradeNo: payLog.OrderNo,
			}

			return db.GetDB().Save(&userMemberCard).Error
		})

		return err
	}

}

func (s *userMemberCardService) MyMemberCard(uid int64) ([]model.BeeUserMemberCard, error) {

	var list []model.BeeUserMemberCard
	var err = db.GetDB().Model(&model.BeeUserMemberCard{}).Where("user_id = ?", uid).Find(&list).Error
	return list, err

}

// wxPay //微信支付
func (s *userMemberCardService) WxPay(ginCtx *gin.Context, shopId int64, memberCardId int32) (*proto.GetWxPayInfoRes, error) {
	var memberCard model.BeeMemberCard
	if err := db.GetDB().Model(&model.BeeMemberCard{ID: memberCardId}).First(&memberCard).Error; err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	} else {
		type payAction struct {
			Type int32 `json:"type"`
			ID   int64 `json:"id"`
		}
		var wxPayConfig model.BeeWxPayConfig

		var payOrderId = util.GetRandInt64()
		if err := db.GetDB().Where("user_id = ? and is_deleted = 0", kit.GetUserId(ginCtx)).Take(&wxPayConfig).Error; err != nil {
			return nil, errors.Wrap(err, "获取微信配置失败！")
		}
		wxPayClient, err := GetPaySrv().GetWechatPayClient(ginCtx, &WxPayConfig{
			MchId:           wxPayConfig.MchId,
			Secret:          wxPayConfig.AppSecret,
			Token:           wxPayConfig.Token,
			ReturnUrl:       "",
			NotifyUrl:       GetPaySrv().GetWxPayNotifyUrl(ginCtx, &wxPayConfig),
			PrivateCertPath: wxPayConfig.PrivateCert,
			Debug:           wxPayConfig.Debug,
		})
		if err != nil {
			return nil, errors.Wrap(err, "获取微信支付客户端失败！")
		}
		userOpenId, err := GetUserSrv().GetUserWxOpenId(ginCtx)
		if err != nil {
			return nil, err
		}

		wxResp, err := wxPayClient.V3TransactionJsapi(ginCtx, gopay.BodyMap{
			"mchid":        wxPayConfig.MchId,
			"out_trade_no": payOrderId,
			"appid":        wxPayConfig.AppId,
			"description":  "会与卡购买",
			"notify_url":   GetPaySrv().GetWxPayNotifyUrl(ginCtx, &wxPayConfig),
			"amount": map[string]interface{}{
				"total":    int64(memberCard.Amount * 100.0),
				"currency": "CNY",
			},
			"time_expire": time.Now().Add(time.Hour * 1).Format(time.RFC3339), // 限制一小时内支付
			"payer": map[string]interface{}{
				"openid": userOpenId,
			},
		})
		if err != nil {
			return nil, err
		}
		if wxResp.Code != 0 {
			return nil, errors.New("微信请求失败：" + wxResp.Error)
		}

		if err := db.GetDB().Transaction(func(tx *gorm.DB) error {
			payLog := &model.BeePayLog{
				BaseModel:   *kit.GetInsertBaseModel(ginCtx),
				Money:       decimal.NewFromFloat(memberCard.Amount),
				NextAction:  "",
				OrderNo:     util.GetRandInt64(),
				OrderNumber: strconv.FormatInt(int64(memberCard.ID), 10),
				PayGate:     enum.PayGateWXAPP,
				Remark:      "购买会员卡",
				Uid:         kit.GetUid(ginCtx),
				ShopId:      shopId,
				OrderType:   1,
				Status:      enum.PayLogStatusUnPaid,
			}
			if err = tx.Create(payLog).Error; err != nil {
				return err
			}

			return nil
		}); err != nil {
			return nil, errors.Wrap(err, "微信下单失败")
		}
		jsapiSignInfo, err := wxPayClient.PaySignOfApplet(wxPayConfig.AppId, wxResp.Response.PrepayId)
		if err != nil {
			return nil, errors.Wrap(err, "获取微信支付签名失败")
		}
		return &proto.GetWxPayInfoRes{
			TimeStamp:  jsapiSignInfo.TimeStamp,
			OutTradeId: payOrderId,
			Package:    jsapiSignInfo.Package,
			PaySign:    jsapiSignInfo.PaySign,
			Appid:      wxPayConfig.AppId,
			Sign:       jsapiSignInfo.PaySign,
			SignType:   jsapiSignInfo.SignType,
			PrepayId:   payOrderId,
			NonceStr:   jsapiSignInfo.NonceStr,
		}, nil
	}
}

// 查询指定会员卡详情
func (s *userMemberCardService) GetMemberCardHxInfo(id int64, uid int64) (proto.UserMemberCardRes, error) {

	var res proto.UserMemberCardRes
	if err := db.GetDB().Model(&model.BeeUserMemberCard{}).First(&res.BeeUserMemberCard, id).Error; err != nil {
		return res, err
	} else if err := db.GetDB().Model(&model.BeeMemberCard{}).First(&res.MemberCard, res.CardID).Error; err != nil {
		return res, err
	} else {
		db.GetDB().Model(&model.BeeUserMemberCardUseLog{}).Where("user_card_id = ? and date_format(use_time,'%Y-%m-%d') = ?", id, time.Now().Format(time.DateOnly)).
			First(&res.UseLog)
		return res, nil
	}
}

// 自动扣除已过期的次数
func (s *userMemberCardService) AutoSubInValidCount(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				now := time.Now()
				if now.Hour() == 0 {
					s.autoSubInValidCountHandler()
				}
			}
		}
	}()
}

// 每天凌晨0点执行,检查前一天是否有没有领取的,没有领取的,可用次数自动减1
func (s *userMemberCardService) autoSubInValidCountHandler() {

	var now = time.Now().Add(-1 * time.Hour).Format(time.DateOnly)
	var list []model.BeeUserMemberCard
	if err := db.GetDB().Model(&model.BeeUserMemberCard{}).Find(&list).Error; err != nil {
		logger.GetLogger().Error(err.Error())
	}
	for _, item := range list {
		var count int64 = 0
		db.GetDB().Model(&model.BeeUserMemberCardUseLog{}).Where("user_card_id = ? and date_format(use_time,'%Y-%m-%d') = ?", item.ID, now).
			Count(&count)
		if count > 0 {
			continue
		} else {
			db.GetDB().Model(&model.BeeUserMemberCard{ID: item.ID}).Update("left_count", item.LeftCount-1)
		}
	}
}
