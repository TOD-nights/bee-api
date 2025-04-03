package bee

import (
	"gitee.com/stuinfer/bee-api/model"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/beeshop/model/bee"
	beeReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/beeshop/model/bee/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/beeshop/utils"
)

type BeeShopInfoService struct{}

// CreateBeeShopInfo 创建商店信息记录
// Author [piexlmax](https://github.com/piexlmax)
func (beeShopInfoService *BeeShopInfoService) CreateBeeShopInfo(beeShopInfo *bee.BeeShopInfo) (err error) {
	beeShopInfo.DateAdd = utils.NowPtr()
	beeShopInfo.DateUpdate = utils.NowPtr()
	err = global.GVA_DB.Create(beeShopInfo).Error
	return err
}

// DeleteBeeShopInfo 删除商店信息记录
// Author [piexlmax](https://github.com/piexlmax)
func (beeShopInfoService *BeeShopInfoService) DeleteBeeShopInfo(id string, shopUserId int) (err error) {
	err = global.GVA_DB.Model(&bee.BeeShopInfo{}).Where("id = ?", id).Where("user_id = ?", shopUserId).
		Updates(map[string]interface{}{
			"is_deleted":  1,
			"date_delete": utils.NowPtr(),
		}).Error
	return err
}

// DeleteBeeShopInfoByIds 批量删除商店信息记录
// Author [piexlmax](https://github.com/piexlmax)
func (beeShopInfoService *BeeShopInfoService) DeleteBeeShopInfoByIds(ids []string, shopUserId int) (err error) {
	err = global.GVA_DB.Model(&bee.BeeShopInfo{}).Where("id = ?", ids).Where("user_id = ?", shopUserId).
		Updates(map[string]interface{}{
			"is_deleted":  1,
			"date_delete": utils.NowPtr(),
		}).Error
	return err
}

// UpdateBeeShopInfo 更新商店信息记录
// Author [piexlmax](https://github.com/piexlmax)
func (beeShopInfoService *BeeShopInfoService) UpdateBeeShopInfo(beeShopInfo bee.BeeShopInfo, shopUserId int) (err error) {
	beeShopInfo.DateUpdate = utils.NowPtr()
	err = global.GVA_DB.Model(&bee.BeeShopInfo{}).Where("id = ? and user_id = ?", beeShopInfo.Id, shopUserId).Updates(&beeShopInfo).Error
	return err
}

// GetBeeShopInfo 根据id获取商店信息记录
// Author [piexlmax](https://github.com/piexlmax)
func (beeShopInfoService *BeeShopInfoService) GetBeeShopInfo(id string, shopUserId int) (beeShopInfo bee.BeeShopInfo, err error) {
	err = global.GVA_DB.Where("id = ? and user_id = ?", id, shopUserId).First(&beeShopInfo).Error
	return
}

// GetBeeShopInfoInfoList 分页获取商店信息记录
// Author [piexlmax](https://github.com/piexlmax)
func (beeShopInfoService *BeeShopInfoService) GetBeeShopInfoInfoList(info beeReq.BeeShopInfoSearch, shopUserId int, loginUserId uint) (list []bee.BeeShopInfo, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	var shopIds = []int{}
	if loginUserId > 0 {
		var userInfo system.SysUser
		global.GVA_DB.Model(&userInfo).Preload("Authorities").First(&userInfo, loginUserId)
		roleIds := []uint{}
		for _, role := range userInfo.Authorities {
			roleIds = append(roleIds, role.AuthorityId)
		}

		var roles []system.SysAuthority
		if err := global.GVA_DB.Model(&system.SysAuthority{}).Preload("ShopInfos").Find(&roles, roleIds).Error; err != nil {
			return list, 0, err
		}

		for _, role := range roles {
			if len(role.ShopInfos) > 0 {
				for _, shop := range role.ShopInfos {
					shopIds = append(shopIds, int(*shop.Id))
				}
			}
		}
	}

	// 创建db
	db := global.GVA_DB.Model(&bee.BeeShopInfo{})
	// db = db.Where("user_id = ?", shopUserId)
	if len(shopIds) > 0 {
		db = db.Where("id in ?", shopIds)
	}
	var beeShopInfos []bee.BeeShopInfo
	// 如果有条件搜索 下方会自动创建搜索语句
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&beeShopInfos).Error
	return beeShopInfos, total, err
}

func (service *BeeShopInfoService) GetAllBeeShopInfos() ([]model.BeeShopInfo, error) {
	var list []model.BeeShopInfo
	if err := global.GVA_DB.Debug().Model(&model.BeeShopInfo{}).Order("id desc").Find(&list).Error; err != nil {
		return list, err
	} else {
		return list, nil
	}
}
