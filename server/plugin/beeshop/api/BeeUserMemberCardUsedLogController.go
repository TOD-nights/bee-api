package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/beeshop/model/bee/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/beeshop/service/bee"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type beeUserMemberCardUsedLogController struct{}

/*
*
会员卡分页
*/
func (c *beeUserMemberCardUsedLogController) Page(ginCtx *gin.Context) {
	var memberCardSearchInfo request.MemberCardSearchInfo
	if err := ginCtx.ShouldBindQuery(&memberCardSearchInfo); err != nil {
		global.GVA_LOG.Error("参数异常!", zap.Error(err))
		response.FailWithMessage("参数异常", ginCtx)
	} else {
		if result, total, err := bee.UserMemberCardLogService.Page(memberCardSearchInfo); err != nil {
			global.GVA_LOG.Error("查询异常!", zap.Error(err))
			response.FailWithMessage("查询异常", ginCtx)
		} else {

			response.OkWithDetailed(map[string]interface{}{"list": result, "total": total}, "查询成功", ginCtx)
		}
	}
}
func (c *beeUserMemberCardUsedLogController) Save(context *gin.Context) {

}

func (c *beeUserMemberCardUsedLogController) DeleteOneById(context *gin.Context) {

}

var UserMemberCardUsedLogController = &beeUserMemberCardUsedLogController{}
