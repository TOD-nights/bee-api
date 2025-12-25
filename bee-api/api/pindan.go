package api

import (
	"gitee.com/stuinfer/bee-api/proto"
	"gitee.com/stuinfer/bee-api/service"
	"github.com/gin-gonic/gin"
)

// / 拼单接口
type PindanApi struct {
	BaseApi
}

func (api PindanApi) Create(c *gin.Context) {

	var req proto.CreatePindanReq
	if err := c.Bind(&req); err != nil {
		api.Res(c, nil, err)
		return
	}
	resp, err := service.GetOrderSrv().CreateOrder(c, c.ClientIP(), &req)
	api.Res(c, resp, err)
}
