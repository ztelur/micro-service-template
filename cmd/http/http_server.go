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
package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/longjoy/micro-service/app"
	_ "github.com/longjoy/micro-service/app/bootstrap"
	"github.com/longjoy/micro-service/app/container/servicecontainer"
	"github.com/longjoy/micro-service/infra/logger"
	"github.com/pkg/errors"
)

const (
	DEV_CONFIG  string = "../../app/config/appConfigDev.yaml"
	PROD_CONFIG string = "../../app/config/appConfigProd.yaml"
)

func main() {

	filename := DEV_CONFIG
	//filename := PROD_CONFIG
	container, err := buildContainer(filename)
	if err != nil {
		fmt.Printf("%+v\n", err)
		//logger.Log.Errorf("%+v\n", err)
		panic(err)
	}
	if err := runServer(container); err != nil {
		logger.Log.Errorf("Failed to run user server: %+v\n", err)
		panic(err)
	} else {
		logger.Log.Info("server started")
	}

	r := SetupRouter()
	err = r.Run()
	if err != nil {
		logger.Log.Errorf(err.Error())
	}
}

func buildContainer(filename string) (*servicecontainer.ServiceContainer, error) {

	container, err := app.InitApp(filename)
	sc := container.(*servicecontainer.ServiceContainer)
	if err != nil {
		//logger.Log.Errorf("%+v\n", err)
		return nil, errors.Wrap(err, "")
	}
	return sc, nil
}

func runServer(sc *servicecontainer.ServiceContainer) error {
	logger.Log.Debug("start runserver")
	r := SetupRouter()
	ugc := sc.AppConfig.UserGrpcConfig

	err := r.Run(ugc.UrlAddress)

	if err != nil {
		logger.Log.Errorf(err.Error())
	}
	return nil
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/config/api/base", GetBaseInfo)
}
