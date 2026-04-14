package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/beeshop/api"
	"github.com/gin-gonic/gin"
)

func InitPindanRoute(privateRoute *gin.RouterGroup, publicRoute *gin.RouterGroup) {
	group := privateRoute.Group("/pindan")
	group.GET("/page", api.PindanController.Page)

}
