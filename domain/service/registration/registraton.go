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
)

// RegistrationUseCase implements RegistrationUseCaseInterface.
// It has UserDataInterface, which can be used to access persistence layer
// TxDataInterface is needed to support transaction
type RegistrationUseCase struct {
	UserRepository repository.UserRepository
}

func (ruc *RegistrationUseCase) RegisterUser(user *model.User) (*model.User, error) {
	err := user.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "user validation failed")
	}
	isDup, err := ruc.isDuplicate(user.Name)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	if isDup {
		return nil, errors.New("duplicate user for " + user.Name)
	}
	resultUser, err := ruc.UserRepository.Insert(user)

	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	return resultUser, nil
}

func (ruc *RegistrationUseCase) ModifyUser(user *model.User) error {
	return modifyUser(ruc.UserRepository, user)
}

func (ruc *RegistrationUseCase) isDuplicate(name string) (bool, error) {
	user, err := ruc.UserRepository.FindByName(name)
	//logger.Log.Debug("isDuplicate() user:", user)
	if err != nil {
		return false, errors.Wrap(err, "")
	}
	if user != nil {
		return true, nil
	}
	return false, nil
}

func (ruc *RegistrationUseCase) UnregisterUser(username string) error {
	return unregisterUser(ruc.UserRepository, username)
}

// The use case of ModifyAndUnregister without transaction
func (ruc *RegistrationUseCase) ModifyAndUnregister(user *model.User) error {
	return ModifyAndUnregister(ruc.UserRepository, user)
}
