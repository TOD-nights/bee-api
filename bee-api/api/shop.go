package api

import (
	"gitee.com/stuinfer/bee-api/enum"
	"gitee.com/stuinfer/bee-api/proto"
	"gitee.com/stuinfer/bee-api/service"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type ShopApi struct {
	BaseApi
}

// / 获取门店详细
func (api ShopApi) GetShopInfo(ctx *gin.Context) {
	id := cast.ToInt64(ctx.Query("id"))
	lat := cast.ToFloat64(ctx.Query("lat"))
	lng := cast.ToFloat64(ctx.Query("lng"))
	info, err := service.GetShopSrv().GetShopInfo(ctx, id, lat, lng)
	if err != nil {
		api.Fail(ctx, enum.ResCodeFail, "系统错误")
		return
	}
	resp := &proto.ShopInfoResp{
		Info:    info,
		ExtJson: make(map[string]interface{}), //@todo
	}
	api.Success(ctx, resp)
}
func (api ShopApi) SubShopDetail(c *gin.Context) {
	id := cast.ToInt64(c.Query("id"))
	info, err := service.GetShopSrv().GetShopInfo(c, id, 39.9042, 116.4074)
	if err != nil {
		api.Fail(c, enum.ResCodeFail, "系统错误")
		return
	}
	resp := &proto.ShopInfoResp{
		Info:    info,
		ExtJson: make(map[string]interface{}), //@todo
	}
	api.Success(c, resp)
}

func (api ShopApi) SubShopList(c *gin.Context) {
	var req proto.ListShopReq
	if err := c.Bind(&req); err != nil {
		api.Res(c, nil, err)
		return
	}
	rsp, err := service.GetShopSrv().List(c, &req)
	api.Res(c, rsp, err)
}
