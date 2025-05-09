package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	bee2 "github.com/flipped-aurora/gin-vue-admin/server/plugin/beeshop/model/bee"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/beeshop/model/bee/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/beeshop/service/bee"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

type beeMemberCardController struct{}

// Page //member-card/page  //会员卡列表
func (c *beeMemberCardController) Page(ginCtx *gin.Context) {

	var memberCardSearchInfo request.MemberCardSearchInfo
	if err := ginCtx.ShouldBindQuery(&memberCardSearchInfo); err != nil {
		global.GVA_LOG.Error("参数异常!", zap.Error(err))
		response.FailWithMessage("参数异常", ginCtx)
	} else {
		if result, total, err := bee.MemberCardService.Page(memberCardSearchInfo); err != nil {
			global.GVA_LOG.Error("查询异常!", zap.Error(err))
			response.FailWithMessage("查询异常", ginCtx)
		} else {

			response.OkWithDetailed(map[string]interface{}{"list": result, "total": total}, "查询成功", ginCtx)
		}
	}
}
func (c *beeMemberCardController) Save(context *gin.Context) {
	var memberCard bee2.BeeMemberCard
	if err := context.ShouldBindJSON(&memberCard); err != nil {
		global.GVA_LOG.Error("参数异常!", zap.Error(err))
		response.FailWithMessage("参数异常", context)
	} else {
		if err := bee.MemberCardService.SaveOrUpdate(memberCard); err != nil {
			global.GVA_LOG.Error("查询异常!", zap.Error(err))
			response.FailWithMessage("查询异常", context)
		} else {

			response.OkWithMessage("提交成功", context)
		}
	}
}

func (c *beeMemberCardController) DeleteOneById(context *gin.Context) {
	if idStr, existed := context.GetQuery("id"); existed {
		if id, err := strconv.Atoi(idStr); err != nil || id <= 0 {
			response.FailWithMessage("参数不合法", context)
		} else {
			bee.MemberCardService.DeleteOneById(int32(id))
			response.OkWithMessage("删除成功", context)
		}
	} else {
		response.FailWithMessage("参数不存在", context)
	}
}

func (c *beeMemberCardController) Info(context *gin.Context) {
	if idStr, existed := context.GetQuery("id"); existed {
		if id, err := strconv.Atoi(idStr); err != nil || id <= 0 {
			response.FailWithMessage("参数不合法", context)
		} else if memberCard, err := bee.MemberCardService.Info(id); err != nil {
			response.FailWithMessage("查询失败", context)
		} else {

			response.OkWithDetailed(memberCard, "删除成功", context)
		}
	} else {
		response.FailWithMessage("参数不存在", context)
	}
}

func (c *beeMemberCardController) RecoverOneById(context *gin.Context) {
	if idStr, existed := context.GetQuery("id"); existed {
		if id, err := strconv.Atoi(idStr); err != nil || id <= 0 {
			response.FailWithMessage("参数不合法", context)
		} else {
			bee.MemberCardService.RecoverOneById(int32(id))
			response.OkWithMessage("删除成功", context)
		}
	} else {
		response.FailWithMessage("参数不存在", context)
	}
}

var MemberCardController = &beeMemberCardController{}
