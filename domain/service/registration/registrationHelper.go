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
package registration

import (
	"github.com/longjoy/micro-service/domain/model"
	"github.com/longjoy/micro-service/domain/repository"
	"github.com/pkg/errors"
	"strconv"
)

func modifyUser(ury repository.UserRepository, user *model.User) error {
	//loggera.Log.Debug("modifyUser")
	err := user.ValidatePersisted()
	if err != nil {
		return errors.Wrap(err, "user validation failed")
	}
	rowsAffected, err := ury.Update(user)
	if err != nil {
		return errors.Wrap(err, "")
	}
	if rowsAffected != 1 {
		return errors.New("Modify user failed. rows affected is " + strconv.Itoa(int(rowsAffected)))
	}
	return nil
}

func unregisterUser(ury repository.UserRepository, username string) error {
	affected, err := ury.Remove(username)
	if err != nil {
		return errors.Wrap(err, "")
	}
	if affected == 0 {
		errStr := "UnregisterUser failed. No such user " + username
		return errors.New(errStr)
	}

	if affected != 1 {
		errStr := "UnregisterUser failed. Number of users unregistered are  " + strconv.Itoa(int(affected))
		return errors.New(errStr)
	}
	return nil
}

// The business function will be wrapped inside a transaction or a non-transaction function
// It needs to be written in a way that every error will be returned so it can be caught by TxEnd() function,
// which will handle commit and rollback
func ModifyAndUnregister(ury repository.UserRepository, user *model.User) error {
	//loggera.Log.Debug("ModifyAndUnregister")
	err := modifyUser(ury, user)
	if err != nil {
		return errors.Wrap(err, "")
	}
	err = unregisterUser(ury, user.Name)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}
