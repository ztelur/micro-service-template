package sqldb

import (
	"database/sql"
	"github.com/longjoy/micro-service/app/bootstrap"
	"github.com/pkg/errors"
)

func NewSqlDB(dc *bootstrap.MySQLConf) (*sql.DB, error) {
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
