package api

import (
	"fmt"

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
