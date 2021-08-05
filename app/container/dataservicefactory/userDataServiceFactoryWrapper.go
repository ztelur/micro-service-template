package dataservicefactory

import (
	"github.com/longjoy/micro-service/app/config"
	"github.com/longjoy/micro-service/app/container"
	"github.com/longjoy/micro-service/app/container/dataservicefactory/userdataservicefactory"
	"github.com/longjoy/micro-service/app/logger"
	"github.com/pkg/errors"
)

// userDataServiceFactory is a empty receiver for Build method
type userDataServiceFactoryWrapper struct{}

func (udsfw *userDataServiceFactoryWrapper) Build(c container.Container, dataConfig *config.DataConfig) (DataServiceInterface, error) {
	logger.Log.Debug("UserDataServiceFactory")
	key := dataConfig.DataStoreConfig.Code
	udsi, err := userdataservicefactory.GetUserDataServiceFb(key).Build(c, dataConfig)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	return udsi, nil
}

