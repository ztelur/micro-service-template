package dataservicefactory

import (
	"github.com/longjoy/micro-service/app/config"
	"github.com/longjoy/micro-service/app/container"
	"github.com/longjoy/micro-service/app/container/datastorefactory"
	"github.com/longjoy/micro-service/app/logger"
	"github.com/longjoy/micro-service/assembler"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// cacheDataServiceFactory is a empty receiver for Build method
type cacheDataServiceFactory struct{}

func (cdsf *cacheDataServiceFactory) Build(c container.Container, dataConfig *config.DataConfig) (DataServiceInterface, error) {
	logger.Log.Debug("cacheDataServiceFactory")
	dsc := dataConfig.DataStoreConfig
	dsi, err := datastorefactory.GetDataStoreFb(dsc.Code).Build(c, &dsc)
	grpcConn := dsi.(*grpc.ClientConn)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	cdg := assembler.CacheDataGrpc{grpcConn}
	//logger.Log.Debug("udm:", udm.DB)

	return &cdg, nil
}
