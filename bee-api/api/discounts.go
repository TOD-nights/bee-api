package api

import (
	"gitee.com/stuinfer/bee-api/enum"
	"gitee.com/stuinfer/bee-api/service"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"strconv"
	"strings"
)

type DiscountsApi struct {
	BaseApi
}

func (api DiscountsApi) Statistics(c *gin.Context) {
	userInfo := api.GetUserInfo(c)
	data, err := service.GetCouponSrv().GetMyCouponStatistics(userInfo.Id)
	if err != nil {
		api.Fail(c, enum.ResCodeFail, err.Error())
		return
	}
	api.Success(c, data)
}

func (api DiscountsApi) My(c *gin.Context) {
	userInfo := api.GetUserInfo(c)
	status := c.Query("status")
	var statusEnum []enum.CouponStatus
	for _, s := range strings.Split(status, ",") {
		atoi, err := strconv.Atoi(s)
		if err != nil {
			api.Fail(c, enum.ResCodeFail, err.Error())
			return
		}
		statusEnum = append(statusEnum, enum.CouponStatus(atoi))
	}
	data, err := service.GetCouponSrv().GetMyCouponListByStatus(userInfo.Id, statusEnum)
	if err != nil {
		api.Fail(c, enum.ResCodeFail, err.Error())
		return
	}
	api.Success(c, data)
}

func (api DiscountsApi) Coupons(c *gin.Context) {
	data, err := service.GetCouponSrv().GetCoupons(1)
	api.Res(c, data, err)
}

func (api DiscountsApi) Fetch(c *gin.Context) {
	//@todo 领取优惠券
	userInfo := api.GetUserInfo(c)
	id_, _ := c.GetPostForm("id")
	id := cast.ToInt64(id_)
	err := service.GetCouponSrv().FetchCoupon(c, userInfo, id)
	api.Res(c, "success", err)
	//api.Fail(c, 30001, "user score is not enough") //积分不足
}

func (api DiscountsApi) CouponDetail(c *gin.Context) {
	id := cast.ToInt64(c.Query("id"))
	resp, err := service.GetCouponSrv().CouponDetail(c, id)
	api.Res(c, resp, err)
}

func (api DiscountsApi) Exchange(c *gin.Context) {
	api.Res(c, nil, enum.ErrNotImplement)
}
