package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/beeshop/model/bee/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/beeshop/service/bee"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type beePindanController struct{}

// Page //member-card/page  //会员卡列表
func (c *beePindanController) Page(ginCtx *gin.Context) {

	var reqVo request.BeePindanItemReq
	if err := ginCtx.ShouldBindQuery(&reqVo); err != nil {
		global.GVA_LOG.Error("参数异常!", zap.Error(err))
		response.FailWithMessage("参数异常", ginCtx)
	} else {
		if result, total, err := bee.PindanService.Page(reqVo); err != nil {
			global.GVA_LOG.Error("查询异常!", zap.Error(err))
			response.FailWithMessage("查询异常", ginCtx)
		} else {

			response.OkWithDetailed(map[string]interface{}{"list": result, "total": total}, "查询成功", ginCtx)
		}
	}
}

var PindanController = &beePindanController{}
