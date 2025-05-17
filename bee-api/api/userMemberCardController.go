package api

import (
	"gitee.com/stuinfer/bee-api/enum"
	"gitee.com/stuinfer/bee-api/kit"
	"gitee.com/stuinfer/bee-api/model"
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
func (c *userMemberCardController) GetMemberCardHxInfo(ginCtx *gin.Context) {
	var uid = kit.GetUid(ginCtx)
	if id, err := strconv.ParseInt(ginCtx.Query("id"), 10, 64); err != nil {
		c.Fail(ginCtx, enum.ResCodeFail, err.Error())
	} else if res, err := service.UserMemberCardService.GetMemberCardHxInfo(id, uid); err != nil {
		c.Fail(ginCtx, enum.ResCodeFail, err.Error())
	} else {
		c.Success(ginCtx, res)
	}
}

// 认领
func (c *userMemberCardController) LineQu(context *gin.Context) {

	var useLog linequReqBody
	if err := context.Bind(&useLog); err != nil {
		c.Fail(context, enum.ResCodeFail, err.Error())
	} else if err := service.MemberCardService.SaveUserMemberCardLog(useLog.BeeUserMemberCardUseLog); err != nil {
		c.Fail(context, enum.ResCodeFail, err.Error())
	} else {
		c.Success(context, "领取成功")
	}
}

type balancePayBody struct {
	MemberCardId int32 `json:"memberCardId"`
	ShopId       int64 `json:"shopId"`
}

type linequReqBody struct {
	model.BeeUserMemberCardUseLog
	Token string
}
