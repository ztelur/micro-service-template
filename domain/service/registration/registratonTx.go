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
)

// RegistrationTxUseCase implements RegistrationTxUseCaseInterface.
// It has UserDataInterface, which can be used to access persistence layer
type RegistrationTxUseCase struct {
	UserRepository repository.UserRepository
}

// The use case of ModifyAndUnregister with transaction
func (rtuc *RegistrationTxUseCase) ModifyAndUnregisterWithTx(user *model.User) error {

	udi := rtuc.UserRepository
	return udi.EnableTx(func() error {
		// wrap the business function inside the TxEnd function
		return ModifyAndUnregister(udi, user)
	})
}
