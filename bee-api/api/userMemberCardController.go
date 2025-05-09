package api

import (
	"gitee.com/stuinfer/bee-api/enum"
	"gitee.com/stuinfer/bee-api/kit"
	"gitee.com/stuinfer/bee-api/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

// 用户-会员卡
type userMemberCardController struct {
	BaseApi
}

var UserMemberCardController = &userMemberCardController{}

// BalancePay //余额支付
func (c *userMemberCardController) BalancePay(ginCtx *gin.Context) {
	memberCardId, _ := strconv.ParseInt(ginCtx.PostForm("memberCardId"), 10, 32)
	shopId, _ := strconv.ParseInt(ginCtx.PostForm("shopId"), 10, 64)
	ginCtx.PostForm("shopId")
	if err := service.UserMemberCardService.BalancePay(ginCtx, shopId, int32(memberCardId)); err != nil {
		c.Fail(ginCtx, enum.ResCodeFail, err.Error())
		return
	}
	c.Success(ginCtx, "购买成功")
}

// WxPay //微信支付
func (c *userMemberCardController) WxPay(ginCtx *gin.Context) {

	memberCardId, _ := strconv.ParseInt(ginCtx.PostForm("memberCardId"), 10, 32)
	shopId, _ := strconv.ParseInt(ginCtx.PostForm("shopId"), 10, 64)
	ginCtx.PostForm("shopId")
	if res, err := service.UserMemberCardService.WxPay(ginCtx, shopId, int32(memberCardId)); err != nil {
		c.Fail(ginCtx, enum.ResCodeFail, "微信支付失败")
	} else {
		c.Success(ginCtx, res)
	}
}

// MyMemberCard //我的会员卡列表
func (c *userMemberCardController) MyMemberCard(ginCtx *gin.Context) {
	var uid = kit.GetUid(ginCtx)
	if list, err := service.UserMemberCardService.MyMemberCard(uid); err != nil {
		c.Fail(ginCtx, enum.ResCodeFail, err.Error())
	} else {
		c.Success(ginCtx, list)
	}
}

type balancePayBody struct {
	MemberCardId int32 `json:"memberCardId"`
	ShopId       int64 `json:"shopId"`
}
