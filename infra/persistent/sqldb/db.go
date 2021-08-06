package sqldb

import (
	"database/sql"
	"github.com/jfeng45/gtransaction/config"
	"github.com/pkg/errors"
)

func NewSqlDB(dc *config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open(dc.DriverName, dc.DataSourceName)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	// check the connection
	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	return db, nil

}
