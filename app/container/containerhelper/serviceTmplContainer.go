package containerhelper

import (
	"github.com/longjoy/micro-service/app/config"
	"github.com/longjoy/micro-service/app/container"
	"github.com/longjoy/micro-service/domain/service"
	"github.com/pkg/errors"
)

func GetListUserUseCase(c container.Container) (service.ListUserUseCaseInterface, error) {
	key := config.LIST_USER
	value, err := c.BuildUseCase(key)
	if err != nil {
		//logger.Log.Errorf("%+v\n", err)
		return nil, errors.Wrap(err, "")
	}
	return value.(service.ListUserUseCaseInterface), nil
}

func GetRegistrationUseCase(c container.Container) (service.RegistrationUseCaseInterface, error) {
	key := config.REGISTRATION
	value, err := c.BuildUseCase(key)
	if err != nil {
		//logger.Log.Errorf("%+v\n", err)
		return nil, errors.Wrap(err, "")
	}
	return value.(service.RegistrationUseCaseInterface), nil

}

func GetRegistrationTxUseCase(c container.Container) (service.RegistrationTxUseCaseInterface, error) {
	key := config.REGISTRATION_TX
	value, err := c.BuildUseCase(key)
	if err != nil {
		//logger.Log.Errorf("%+v\n", err)
		return nil, errors.Wrap(err, "")
	}
	return value.(service.RegistrationTxUseCaseInterface), nil

}
