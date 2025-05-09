package api

import (
	"github.com/gin-gonic/gin"
)

type beeUserMemberCardController struct{}

/*
*
会员卡分页
*/
// Page //member-card/page  //会员卡列表
func (c *beeUserMemberCardController) Page(ginCtx *gin.Context) {

	//var memberCardSearchInfo request.MemberCardSearchInfo
	//if err := ginCtx.ShouldBindQuery(&memberCardSearchInfo); err != nil {
	//	global.GVA_LOG.Error("参数异常!", zap.Error(err))
	//	response.FailWithMessage("参数异常", ginCtx)
	//} else {
	//	if result, total, err := bee.MemberCardService.Page(memberCardSearchInfo); err != nil {
	//		global.GVA_LOG.Error("查询异常!", zap.Error(err))
	//		response.FailWithMessage("查询异常", ginCtx)
	//	} else {
	//
	//		response.OkWithDetailed(map[string]interface{}{"list": result, "total": total}, "查询成功", ginCtx)
	//	}
	//}
}
func (c *beeUserMemberCardController) Save(context *gin.Context) {

}

func (c *beeUserMemberCardController) DeleteOneById(context *gin.Context) {

}

func (c *beeUserMemberCardController) Info(context *gin.Context) {

}

var UserMemberCardController = &beeUserMemberCardController{}
