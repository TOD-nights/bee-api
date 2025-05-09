package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/beeshop/api"
	"github.com/gin-gonic/gin"
)

func InitMemberCardRoute(privateRoute *gin.RouterGroup, publicRoute *gin.RouterGroup) {
	group := privateRoute.Group("/member-card")
	group.GET("/page", api.MemberCardController.Page)
	group.POST("/save", api.MemberCardController.Save)
	group.DELETE("/deleteOneById", api.MemberCardController.DeleteOneById)
	group.GET("/info", api.MemberCardController.Info)
	group.POST("/recoverOneById", api.MemberCardController.RecoverOneById)

}

func InitUserMemberCardRoute(privateRoute *gin.RouterGroup, publicRoute *gin.RouterGroup) {
	group := privateRoute.Group("/user-member-card")
	group.GET("/page", api.UserMemberCardController.Page)
	group.POST("/save", api.UserMemberCardController.Save)
	group.DELETE("/deleteOneById", api.UserMemberCardController.DeleteOneById)
	group.GET("/info", api.UserMemberCardController.Info)
}

func InitUserMemberCardUsedLogRoute(privateRoute *gin.RouterGroup, publicRoute *gin.RouterGroup) {
	group := privateRoute.Group("/user-member-card-used-log")
	group.GET("/page", api.UserMemberCardUsedLogController.Page)
	group.POST("/save", api.UserMemberCardUsedLogController.Save)
	group.DELETE("/deleteOneById", api.UserMemberCardUsedLogController.DeleteOneById)
}
