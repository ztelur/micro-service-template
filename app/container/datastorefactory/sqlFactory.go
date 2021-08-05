package datastorefactory

import (
	"database/sql"

	databaseConfig "github.com/jfeng45/gtransaction/config"
	"github.com/jfeng45/gtransaction/factory"
	"github.com/jfeng45/gtransaction/gdbc"

	"github.com/longjoy/micro-service/app/config"
	"github.com/longjoy/micro-service/app/container"
	"github.com/longjoy/micro-service/app/logger"
)

// sqlFactory is receiver for Build method
type sqlFactory struct{}

// implement Build method for SQL database
func (sf *sqlFactory) Build(c container.Container, dsc *config.DataStoreConfig) (DataStoreInterface, error) {
	logger.Log.Debug("sqlFactory")
	key := dsc.Code
	// Only non-transaction connection is cached
	if !dsc.Tx  {
		if value, found := c.Get(key); found {
			logger.Log.Debug("found db in container for key:", key)
			return value, nil
		}
	}
	tdbc := databaseConfig.DatabaseConfig{dsc.DriverName,dsc.UrlAddress, dsc.Tx}
	db, err := factory.BuildSqlDB(&tdbc)
	if err != nil {
		return nil, err
	}
	gdbc, err := buildGdbc(db, dsc.Tx)
	if err != nil {
		return nil, err
	}
	// Only non-transaction connection is cached
	if !dsc.Tx {
		c.Put(key, gdbc)
	}
	return gdbc, nil

}

func buildGdbc(sdb *sql.DB, tx bool) (gdbc.SqlGdbc, error) {
	var sdt gdbc.SqlGdbc
	if tx {
		tx, err := sdb.Begin()
		if err != nil {
			return nil, err
		}
		sdt = &gdbc.SqlConnTx{DB: tx}
		logger.Log.Debug("buildGdbc(), create TX:")
	} else {
		sdt = &gdbc.SqlDBTx{sdb}
		logger.Log.Debug("buildGdbc(), create DB:")
	}
	return sdt, nil
}
