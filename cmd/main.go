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
package main

import (
	"fmt"
	"github.com/longjoy/micro-service/app"
	"github.com/longjoy/micro-service/app/container"
	"github.com/longjoy/micro-service/app/container/containerhelper"
	"github.com/longjoy/micro-service/app/logger"
	"github.com/longjoy/micro-service/domain/model"
	"github.com/longjoy/micro-service/tool/timea"
	"github.com/pkg/errors"
	"time"
)

const (
	DEV_CONFIG  string = "../app/config/appConfigDev.yaml"
	PROD_CONFIG string = "../app/config/appConfigProd.yaml"
)

func main() {
	testMySql()
	//testCouchDB()
}

func testMySql() {

	filename := DEV_CONFIG
	container, err := buildContainer(filename)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	testListUser(container)
	testFindById(container)
	testRegisterUser(container)
	testModifyUser(container)
	testUnregister(container)
	testModifyAndUnregister(container)
	testModifyAndUnregisterWithTx(container)

}
func testCouchDB() {
	filename := PROD_CONFIG
	container, err := buildContainer(filename)
	if err != nil {
		fmt.Printf("%+v\n", err)
		//logger.Log.Errorf("%+v\n", err)
		return
	}
	testFindById(container)
}

func testUnregister(container container.Container) {

	ruci, err := containerhelper.GetRegistrationUseCase(container)
	if err != nil {
		logger.Log.Fatal("registration interface build failed:%+v\n", err)
	}
	username := "Aditi"
	err = ruci.UnregisterUser(username)
	if err != nil {
		logger.Log.Fatalf("testUnregister failed:%+v\n", err)
	}
	logger.Log.Infof("testUnregister successully")
}

func testRegisterUser(container container.Container) {
	ruci, err := containerhelper.GetRegistrationUseCase(container)
	if err != nil {
		logger.Log.Fatal("registration interface build failed:%+v\n", err)
	}
	created, err := time.Parse(timea.FORMAT_ISO8601_DATE, "2018-12-09")
	if err != nil {
		logger.Log.Errorf("date format err:%+v\n", err)
	}

	user := model.User{Name: "Brian", Department: "Marketing", Created: created}

	resultUser, err := ruci.RegisterUser(&user)
	if err != nil {
		logger.Log.Errorf("user registration failed:%+v\n", err)
	} else {
		logger.Log.Info("new user registered:", resultUser)
	}
}

func testModifyUser(container container.Container) {
	ruci, err := containerhelper.GetRegistrationUseCase(container)
	if err != nil {
		logger.Log.Fatal("registration interface build failed:%+v\n", err)
	}
	created, err := time.Parse(timea.FORMAT_ISO8601_DATE, "2019-12-01")
	if err != nil {
		logger.Log.Errorf("date format err:%+v\n", err)
	}
	user := model.User{Id: 29, Name: "Aditi", Department: "HR", Created: created}
	err = ruci.ModifyUser(&user)
	if err != nil {
		logger.Log.Infof("Modify user failed:%+v\n", err)
	} else {
		logger.Log.Info("user modified succeed:", user)
	}
}

func testListUser(container container.Container) {
	rluf, err := containerhelper.GetListUserUseCase(container)
	if err != nil {
		logger.Log.Fatal("RetrieveListUser interface build failed:", err)
	}
	users, err := rluf.ListUser()
	if err != nil {
		logger.Log.Errorf("user list failed:%+v\n", err)
	}
	logger.Log.Info("user list:", users)
}

func testModifyAndUnregister(container container.Container) {
	ruci, err := containerhelper.GetRegistrationUseCase(container)
	if err != nil {
		logger.Log.Fatal("RegisterRegistration interface build failed:%+v\n", err)
	}
	created, err := time.Parse(timea.FORMAT_ISO8601_DATE, "2018-12-09")
	if err != nil {
		logger.Log.Errorf("date format err:%+v\n", err)
	}
	user := model.User{Id: 31, Name: "Richard", Department: "Sales", Created: created}
	err = ruci.ModifyAndUnregister(&user)
	if err != nil {
		logger.Log.Errorf("ModifyAndUnregister failed:%+v\n", err)
	} else {
		logger.Log.Infof("ModifyAndUnregister succeed")
	}
}

func testModifyAndUnregisterWithTx(container container.Container) {
	rtuci, err := containerhelper.GetRegistrationTxUseCase(container)
	if err != nil {
		logger.Log.Fatal("RegisterRegistration interface build failed:%+v\n", err)
	}
	created, err := time.Parse(timea.FORMAT_ISO8601_DATE, "2018-12-09")
	if err != nil {
		logger.Log.Errorf("date format err:%+v\n", err)
	}
	user := model.User{Id: 32, Name: "Richard", Department: "Finance", Created: created}
	err = rtuci.ModifyAndUnregisterWithTx(&user)
	if err != nil {
		logger.Log.Errorf("ModifyAndUnregisterWithTx failed:%+v\n", err)
	} else {
		logger.Log.Infof("ModifyAndUnregisterWithTx succeed")
	}
}

func testFindById(container container.Container) {
	//It is uid in database. Make sure you have it in database, otherwise it won't find it.
	id := 10
	rluf, err := containerhelper.GetListUserUseCase(container)
	if err != nil {
		logger.Log.Fatalf("RetrieveListUser interface build failed:%+v\n", err)
	}
	user, err := rluf.Find(id)
	if err != nil {
		logger.Log.Errorf("fin user failed failed:%+v\n", err)
	}
	logger.Log.Info("find user:", user)
}

func buildContainer(filename string) (container.Container, error) {
	container, err := app.InitApp(filename)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	return container, nil
}
