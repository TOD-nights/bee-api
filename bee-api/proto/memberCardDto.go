package proto

import "gitee.com/stuinfer/bee-api/model"

type MemberCardDto struct {
	model.BeeMemberCard
	TypeName string `json:"typeName" gorm:"column:dictLabel"`
}
