package request

type MemberCardSearchInfo struct {
	Name       string `json:"name" form:"name" `
	ValidMonth int    `json:"validMonth" form:"validMonth" `
	DeleteFlag int    `json:"deleteFlag" form:"deleteFlag" `
	Page       int    `json:"page" form:"page" `
	PageSize   int    `json:"pageSize" form:"pageSize" `
}
