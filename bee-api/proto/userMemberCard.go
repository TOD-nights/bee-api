package proto

import "gitee.com/stuinfer/bee-api/model"

type UserMemberCardRes struct {
	model.BeeUserMemberCard
	MemberCard model.BeeMemberCard           `json:"member_card"`
	UseLog     model.BeeUserMemberCardUseLog `json:"useLog"`
}
