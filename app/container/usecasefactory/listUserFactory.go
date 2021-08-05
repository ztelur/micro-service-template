package usecasefactory

import (
	"github.com/longjoy/micro-service/app/config"
	"github.com/longjoy/micro-service/app/container"
	"github.com/longjoy/micro-service/domain/service/listuser"
	"github.com/pkg/errors"
)

type ListUserFactory struct{}

func (luf *ListUserFactory) Build(c container.Container, appConfig *config.AppConfig, key string) (UseCaseInterface, error) {
	uc := appConfig.UseCaseConfig.ListUser

	udi, err := buildUserData(c, &uc.UserDataConfig)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	luuc := listuser.ListUserUseCase{udi}
	return &luuc, nil
}
