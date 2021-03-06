package service_impl

import (
	"github.com/tientruongcao51/oauth2-sever/models"
	"sync"
)

type UserService interface {
	Put(username string, user models.OauthUser) (id string, err error)
	GetByUsername(username string) (user *models.OauthUser, err error)
}

var UserServiceIns UserService

func init() {
	// sync
	syncOne := sync.Once{}
	syncOne.Do(func() {
		UserServiceIns = UserNewService()
	})
}
