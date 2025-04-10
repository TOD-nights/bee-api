package bee

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/beeshop/dto"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/beeshop/model/bee"
	beeReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/beeshop/model/bee/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/beeshop/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BeeUserBalanceLogService struct{}

// CreateBeeUserBalanceLog 创建用户消费记录记录
// Author [piexlmax](https://github.com/piexlmax)
func (beeUserBalanceLogService *BeeUserBalanceLogService) CreateBeeUserBalanceLog(beeUserBalanceLog *bee.BeeUserBalanceLog) (err error) {
	beeUserBalanceLog.DateAdd = utils.NowPtr()
	beeUserBalanceLog.DateUpdate = utils.NowPtr()
	err = GetBeeDB().Create(beeUserBalanceLog).Error
	return err
}

// DeleteBeeUserBalanceLog 删除用户消费记录记录
// Author [piexlmax](https://github.com/piexlmax)
func (beeUserBalanceLogService *BeeUserBalanceLogService) DeleteBeeUserBalanceLog(id string, shopUserId int) (err error) {
	err = GetBeeDB().Model(&bee.BeeUserBalanceLog{}).Where("id = ?", id).Where("user_id = ?", shopUserId).
		Updates(map[string]interface{}{
			"is_deleted":  1,
			"date_delete": utils.NowPtr(),
		}).Error
	return err
}

// DeleteBeeUserBalanceLogByIds 批量删除用户消费记录记录
// Author [piexlmax](https://github.com/piexlmax)
func (beeUserBalanceLogService *BeeUserBalanceLogService) DeleteBeeUserBalanceLogByIds(ids []string, shopUserId int) (err error) {
	err = GetBeeDB().Model(&bee.BeeUserBalanceLog{}).Where("id = ?", ids).Where("user_id = ?", shopUserId).
		Updates(map[string]interface{}{
			"is_deleted":  1,
			"date_delete": utils.NowPtr(),
		}).Error
	return err
}

// UpdateBeeUserBalanceLog 更新用户消费记录记录
// Author [piexlmax](https://github.com/piexlmax)
func (beeUserBalanceLogService *BeeUserBalanceLogService) UpdateBeeUserBalanceLog(beeUserBalanceLog bee.BeeUserBalanceLog, shopUserId int) (err error) {
	beeUserBalanceLog.DateUpdate = utils.NowPtr()
	err = GetBeeDB().Model(&bee.BeeUserBalanceLog{}).Where("id = ? and user_id = ?", beeUserBalanceLog.Id, shopUserId).Updates(&beeUserBalanceLog).Error
	return err
}

// GetBeeUserBalanceLog 根据id获取用户消费记录记录
// Author [piexlmax](https://github.com/piexlmax)
func (beeUserBalanceLogService *BeeUserBalanceLogService) GetBeeUserBalanceLog(id string, shopUserId int) (beeUserBalanceLog bee.BeeUserBalanceLog, err error) {
	err = GetBeeDB().Where("id = ? and user_id = ?", id, shopUserId).First(&beeUserBalanceLog).Error
	return
}

// GetBeeUserBalanceLogInfoList 分页获取用户消费记录记录
// Author [piexlmax](https://github.com/piexlmax)
func (beeUserBalanceLogService *BeeUserBalanceLogService) GetBeeUserBalanceLogInfoList(info beeReq.BeeUserBalanceLogSearch, shopUserId int, loginUserId uint) (list []dto.BeeUserBalanceLogDto, total int64, err error) {
	var userInfo system.SysUser
	GetBeeDB().Model(&userInfo).Preload("Authorities").First(&userInfo, loginUserId)
	var adminFlag = false
	roleIds := []uint{}
	for _, role := range userInfo.Authorities {
		roleIds = append(roleIds, role.AuthorityId)
		if role.Admin == 1 {
			adminFlag = true
		}
	}

	var roles []system.SysAuthority
	if err := global.GVA_DB.Model(&system.SysAuthority{}).Preload("ShopInfos").Find(&roles, roleIds).Error; err != nil {
		return list, 0, err
	}

	var shopIds = []int{}
	for _, role := range roles {
		if len(role.ShopInfos) > 0 {
			for _, shop := range role.ShopInfos {
				shopIds = append(shopIds, int(*shop.Id))
			}
		}
	}
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	// 修改查询逻辑，使用 LEFT JOIN 而不是 INNER JOIN
	db := GetBeeDB().Debug().Table(bee.BeeUserBalanceLog{}.TableName() + " as log ").
		Joins("left join bee_order a on concat('pay',a.order_number) = log.order_id").
		Joins("left join bee_shop_info b on a.shop_id = b.id")

	// 基础条件
	db = db.Where("log.user_id = ?", shopUserId)

	// 根据 type 参数过滤记录
	if info.Type == "payment" {
		db = db.Where("log.mark = ?", "订单支付")
	} else if info.Type == "recharge" {
		db = db.Where("log.mark = ?", "充值")
	}
	if !adminFlag {
		if len(shopIds) > 0 {
			db = db.Where("a.shop_id in ?", shopIds)
		} else {
			db = db.Where("a.shop_id = -1")
		}
	}
	var beeUserBalanceLogs []dto.BeeUserBalanceLogDto
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.Uid != nil {
		db = db.Where("log.uid = ?", info.Uid)
	}
	if info.StartDateAdd != nil && info.EndDateAdd != nil {
		db = db.Where("log.date_add BETWEEN ? AND ? ", info.StartDateAdd, info.EndDateAdd)
	}

	if info.ShopId > 0 {
		db = db.Where("a.shop_id = ?", info.ShopId)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	var OrderStr string
	orderMap := make(map[string]bool)
	orderMap["id"] = true
	orderMap["date_add"] = true
	if orderMap[info.Sort] {
		OrderStr = info.Sort
		if info.Order == "descending" {
			OrderStr = OrderStr + " desc"
		}
		db = db.Order(OrderStr)
	}
	if limit != 0 {
		db = db.Limit(limit).Offset(offset).Select("log.*,b.name as shopName,COALESCE(a.amount, log.num) as amount")
	}

	err = db.Find(&beeUserBalanceLogs).Error
	return beeUserBalanceLogs, total, err
}

func (s *BeeUserBalanceLogService) GetBeeUserBalanceLogInfoCount(c *gin.Context, orderType string, info beeReq.BeeOrderSearch) (float64, float64) {

	var orderPrefix = ""
	if orderType == "payment" {
		orderPrefix = "pay_"
	} else {
		orderPrefix = "recharge_"
	}
	// 创建db
	// 修改查询逻辑，使用 LEFT JOIN 而不是 INNER JOIN
	db := global.GVA_DB.Session(&gorm.Session{}).Debug().Table(bee.BeePayLog{}.TableName() + " as log ").
		Joins("inner join bee_user_balance_log a on concat('" + orderPrefix + "',log.order_no) = a.order_id").
		Joins("left join bee_shop_info b on log.shop_id = b.id")

	// 基础条件
	//db = db.Where("log.user_id = ?", shopUserId)

	// 根据 type 参数过滤记录
	if orderType == "payment" {
		db = db.Where("a.order_id like ?", "pay_%")
	} else if orderType == "recharge" {
		db = db.Where("a.order_id like ?", "recharge_%")
	}
	if _, exist := c.Get("admin"); !exist {
		if shopIds, exist := c.Get("shopIds"); exist {
			db = db.Where("log.shop_id IN ?", shopIds)
		}
	}

	if info.StartDateAdd != nil && info.EndDateAdd != nil {
		db = db.Where("log.date_add BETWEEN ? AND ? ", info.StartDateAdd, info.EndDateAdd)
	}
	var sum float64
	var sumSelected float64
	if info.ShopId != nil && *info.ShopId > 0 {
		tx := db.Session(&gorm.Session{}).Where("log.shop_id = ?", info.ShopId)
		tx.Select("ifnull(sum(money),0)").Scan(&sumSelected)
	}

	db.Select("ifnull(sum(money),0)").Scan(&sum)
	return sum, sumSelected
}
