package usecasefactory

import (
	"github.com/longjoy/micro-service/app/config"
	"github.com/longjoy/micro-service/app/container"
	"github.com/longjoy/micro-service/app/container/dataservicefactory"
	"github.com/longjoy/micro-service/domain/repository"
	"github.com/pkg/errors"
)

func buildUserData(c container.Container, dc *config.DataConfig) (repository.UserRepository, error) {
	dsi, err := dataservicefactory.GetDataServiceFb(dc.Code).Build(c, dc)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	udi := dsi.(repository.UserRepository)
	return udi, nil
}

