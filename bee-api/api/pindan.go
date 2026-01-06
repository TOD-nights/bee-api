package api

import (
	"errors"
	"fmt"
	"strconv"

	"gitee.com/stuinfer/bee-api/enum"
	"gitee.com/stuinfer/bee-api/kit"
	"gitee.com/stuinfer/bee-api/proto"
	"gitee.com/stuinfer/bee-api/service"
	"github.com/gin-gonic/gin"
)

// / 拼单接口
type PindanApi struct {
	BaseApi
}

// / 发起拼单
func (api PindanApi) Create(c *gin.Context) {

	var req proto.CreatePindanReq
	if err := c.Bind(&req); err != nil {
		fmt.Println("PindanApi Create bind err:", err.Error())
		api.Res(c, nil, err)
		return
	}
	fmt.Println("req:", req)
	resp, err := service.GetPinDanServ().CreatePindan(c, c.ClientIP(), &req)
	api.Res(c, resp, err)
}

// / 加入拼单
func (api PindanApi) Join(c *gin.Context) {

	var req proto.CreatePindanReq
	if err := c.Bind(&req); err != nil {
		fmt.Println("PindanApi Create bind err:", err.Error())
		api.Res(c, nil, err)
		return
	}
	fmt.Println("req:", req)
	resp, err := service.GetPinDanServ().JoinPindan(c, c.ClientIP(), &req)
	api.Res(c, resp, err)
}

// 查询拼单详情
func (api PindanApi) GetPinDanInfo(c *gin.Context) {

	var req proto.GetPindanInfoReq
	if err := c.Bind(&req); err != nil {
		api.Res(c, nil, err)
		return
	}
	resp, err := service.GetPinDanServ().GetPindanInfoByUserIdAndShopId(c, req.ShopID, kit.GetUid(c))
	api.Res(c, resp, err)
}

// 查询拼单详情
func (api PindanApi) GetPinDanInfoById(c *gin.Context) {

	var req proto.GetPindanInfoReq
	if err := c.Bind(&req); err != nil {
		api.Res(c, nil, err)
		return
	}
	resp, err := service.GetPinDanServ().GetPindanInfoByUserIdAndId(c, kit.GetUid(c), req.ID)
	api.Res(c, resp, err)
}

// / 更新指定拼单的配送方式
func (api PindanApi) UpdatePindanPeisongTypeById(ctx *gin.Context) {
	if id, err := strconv.ParseInt(ctx.Param("id"), 10, 64); err != nil {
		api.Res(ctx, nil, errors.New("拼单标识不合法"))
	} else if err := service.GetPinDanServ().UpdatePindanPeisongTypeById(id); err != nil {
		api.Res(ctx, nil, errors.New("拼单配送方式更新失败"))
	} else {
		api.Res(ctx, true, nil)
	}
}

// / 我创建的拼单列表, 查询条件; 状态,页数, status page
func (api PindanApi) GetMyCreatedPindanRecord(ctx *gin.Context) {
	var reqVo orderQueryParam
	if err := ctx.ShouldBindQuery(&reqVo); err != nil {
		api.Res(ctx, nil, errors.New("请求参数不合法"))
	} else if res, err := service.GetPinDanServ().GetMyCreatedPindanRecord(reqVo.Page, reqVo.Status, kit.GetUid(ctx)); err != nil {
		api.Res(ctx, nil, &enum.BussError{Code: 100000, Message: "查询失败", Err: err})
	} else {
		api.Res(ctx, res, nil)
	}
}

// / 拼单 微信支付
func (api PindanApi) PindanWxPay(ctx *gin.Context) {

	pindanId, _ := strconv.ParseInt(ctx.PostForm("pindanId"), 10, 64)
	if res, err := service.GetPinDanServ().WxPay(ctx, pindanId, kit.GetUid(ctx)); err != nil {
		fmt.Println(err)
		api.Fail(ctx, enum.ResCodeFail, "微信支付失败")
	} else {
		api.Success(ctx, res)
	}
}

// / 我创建的拼单列表, 查询条件; 状态,页数,
func (api PindanApi) GetMyJoinedPindanRecord(ctx *gin.Context) {
	var reqVo orderQueryParam
	if err := ctx.ShouldBindQuery(&reqVo); err != nil {
		api.Res(ctx, nil, errors.New("请求参数不合法"))
	} else if res, err := service.GetPinDanServ().GetMyJoinedPindanRecord(reqVo.Page, reqVo.Status, kit.GetUid(ctx)); err != nil {
		api.Res(ctx, nil, &enum.BussError{Code: 100000, Message: "查询失败", Err: err})
	} else {
		api.Res(ctx, res, nil)
	}
}

type orderQueryParam struct {
	Page   int64 `json:"page" form:"page"`
	Status int64 `json:"status" form:"status"`
}
