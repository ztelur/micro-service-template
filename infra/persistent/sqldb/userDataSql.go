/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package sqldb

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jfeng45/gtransaction/gdbc"
	"github.com/longjoy/micro-service/app/logger"
	"github.com/longjoy/micro-service/domain/model"
	"github.com/longjoy/micro-service/tool/timea"
	"github.com/pkg/errors"
	"time"
)

const (
	// test rollback
	//DELETE_USER        string = "delete from userinf where username=?"
	DELETE_USER        string = "delete from userinfo where username=?"
	QUERY_USER_BY_ID   string = "SELECT * FROM userinfo where uid =?"
	QUERY_USER_BY_NAME        = "SELECT * FROM userinfo where username =?"
	QUERY_USER                = "SELECT * FROM userinfo "
	UPDATE_USER               = "update userinfo set username=?, department=?, created=? where uid=?"
	INSERT_USER               = "INSERT userinfo SET username=?,department=?,created=?"
)

// UserDataSql is the SQL implementation of UserRepository
type UserDataSql struct {
	DB gdbc.SqlGdbc
}

func (uds *UserDataSql) Remove(username string) (int64, error) {

	stmt, err := uds.DB.Prepare(DELETE_USER)
	if err != nil {
		return 0, errors.Wrap(err, "")
	}
	defer stmt.Close()

	res, err := stmt.Exec(username)
	if err != nil {
		return 0, errors.Wrap(err, "")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "")
	}

	logger.Log.Debug("remove:row affected ", rowsAffected)
	return rowsAffected, nil
}

func (uds *UserDataSql) Find(id int) (*model.User, error) {
	rows, err := uds.DB.Query(QUERY_USER_BY_ID, id)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	defer rows.Close()
	return retrieveUser(rows)
}
func retrieveUser(rows *sql.Rows) (*model.User, error) {
	if rows.Next() {
		return rowsToUser(rows)
	}
	return nil, nil
}
func rowsToUser(rows *sql.Rows) (*model.User, error) {
	var ds string
	user := &model.User{}
	err := rows.Scan(&user.Id, &user.Name, &user.Department, &ds)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	created, err := time.Parse(timea.FORMAT_ISO8601_DATE, ds)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	user.Created = created

	logger.Log.Debug("rows to User:", user)
	return user, nil
}
func (uds *UserDataSql) FindByName(name string) (*model.User, error) {
	//logger.Log.Debug("call FindByName() and name is:", name)
	rows, err := uds.DB.Query(QUERY_USER_BY_NAME, name)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	defer rows.Close()
	return retrieveUser(rows)
}

func (uds *UserDataSql) FindAll() ([]model.User, error) {

	rows, err := uds.DB.Query(QUERY_USER)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	defer rows.Close()
	users := []model.User{}

	//var ds string
	for rows.Next() {
		user, err := rowsToUser(rows)
		if err != nil {
			return users, errors.Wrap(err, "")
		}
		users = append(users, *user)

	}
	//need to check error for rows.Next()
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "")
	}
	logger.Log.Debug("find user list:", users)
	return users, nil
}

func (uds *UserDataSql) Update(user *model.User) (int64, error) {

	stmt, err := uds.DB.Prepare(UPDATE_USER)

	if err != nil {
		return 0, errors.Wrap(err, "")
	}
	defer stmt.Close()
	res, err := stmt.Exec(user.Name, user.Department, user.Created, user.Id)
	if err != nil {
		return 0, errors.Wrap(err, "")
	}
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return 0, errors.Wrap(err, "")
	}
	logger.Log.Debug("update: rows affected: ", rowsAffected)

	return rowsAffected, nil
}

func (uds *UserDataSql) Insert(user *model.User) (*model.User, error) {

	stmt, err := uds.DB.Prepare(INSERT_USER)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	defer stmt.Close()
	res, err := stmt.Exec(user.Name, user.Department, user.Created)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	user.Id = int(id)
	logger.Log.Debug("user inserted:", user)
	return user, nil
}

func (uds *UserDataSql) EnableTx(txFunc func() error) error {
	return uds.DB.TxEnd(txFunc)
}

