package middleware

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"strconv"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
)

var casbinService = service.ServiceGroupApp.SystemServiceGroup.CasbinService

// CasbinHandler 拦截器
func CasbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		waitUse, _ := utils.GetClaims(c)
		//获取请求的PATH
		path := c.Request.URL.Path
		obj := strings.TrimPrefix(path, global.GVA_CONFIG.System.RouterPrefix)
		// 获取请求方法
		act := c.Request.Method
		// 获取用户的角色
		sub := strconv.Itoa(int(waitUse.AuthorityId))
		var sysUser system.SysUser
		if err := global.GVA_DB.Model(&system.SysUser{}).Preload("Authorities").Preload("Authorities.ShopInfos").First(&sysUser, waitUse.BaseClaims.ID).Error; err != nil {
			response.FailWithDetailed(gin.H{}, "权限不足", c)
			c.Abort()
		}
		var admin = false
		if len(sysUser.Authorities) == 0 {
			response.FailWithDetailed(gin.H{}, "权限不足", c)
			c.Abort()
		}
		for _, v := range sysUser.Authorities {
			if v.Admin == 1 {
				admin = true
			}
		}
		if admin {
			c.Set("admin", true)
			c.Next()
		} else {

			e := casbinService.Casbin() // 判断策略中是否存在
			success, _ := e.Enforce(sub, obj, act)
			if !success {
				response.FailWithDetailed(gin.H{}, "权限不足", c)
				c.Abort()
				return
			}
			var shopIds []int
			for _, v := range sysUser.Authorities {
				for _, shop := range v.ShopInfos {
					shopIds = append(shopIds, *shop.Id)
				}
			}
			c.Set("shopIds", shopIds)
			c.Next()
		}
	}
}
