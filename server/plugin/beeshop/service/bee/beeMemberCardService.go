package bee

import "sync"

type beeMemberCardService struct {
}

var once sync.Once
var instance *beeMemberCardService

func GetMemberCardServiceInstance() *beeMemberCardService {
	once.Do(func() {
		instance = &beeMemberCardService{}
	})
	return instance
}
