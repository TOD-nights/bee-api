package bee

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/beeshop/model/bee"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/beeshop/model/bee/request"
)

type pindanService struct {
}

var PindanService = &pindanService{}

func (s *pindanService) Page(info request.BeePindanItemReq) ([]request.BeePindanItemResp, int64, error) {

	tx := global.GVA_DB.Model(&bee.BeePindanOrderItem{})
	tx.Joins("left join bee_shop_goods a on a.id = bee_pindan_order_item.goods_id left join bee_shop_info b on b.id = bee_pindan_order_item.shop_id")
	var count int64 = 0
	var list []request.BeePindanItemResp
	tx.Count(&count)

	if info.Page > 0 {
		tx = tx.Offset((info.Page - 1) * info.PageSize).Limit(info.PageSize).Select("bee_pindan_order_item.*, a.id as goodsInfo_id,a.name as goodsInfo_name,a.pic as goodsInfo_pic,b.id as shopInfo_id,b.address as shopInfo_address,b.name as shopInfo_name").Find(&list)
	}

	return list, count, tx.Error
}
