package api

import (
	"gitee.com/stuinfer/bee-api/enum"
	"gitee.com/stuinfer/bee-api/service"
	"github.com/gin-gonic/gin"
)

type memberCardController struct {
	BaseApi
}

var MemberCardController = &memberCardController{}

func (c *memberCardController) ListAll(ginCtx *gin.Context) {

	if list, err := service.MemberCardService.ListAll(); err != nil {
		c.Fail(ginCtx, enum.ResCodeFail, err.Error())
	} else {
		c.Success(ginCtx, list)
	}
}
