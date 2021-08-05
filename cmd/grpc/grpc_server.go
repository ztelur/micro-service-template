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
	"context"
	"fmt"
	"github.com/longjoy/micro-service/assembler"
	"github.com/longjoy/micro-service/pb"
	"net"

	"github.com/longjoy/micro-service/app"
	"github.com/longjoy/micro-service/app/container/containerhelper"
	"github.com/longjoy/micro-service/app/container/servicecontainer"
	"github.com/longjoy/micro-service/app/logger"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

const (
	DEV_CONFIG  string = "../../app/config/appConfigDev.yaml"
	PROD_CONFIG string = "../../app/config/appConfigProd.yaml"
)

type UserService struct {
	//container container.Container
	container *servicecontainer.ServiceContainer
}

func catchPanic() {
	if p := recover(); p != nil {
		logger.Log.Errorf("%+v\n", p)
	}
}

func (uss *UserService) RegisterUser(ctx context.Context, req *pb.RegisterUserReq) (*pb.RegisterUserResp, error) {
	defer catchPanic()
	logger.Log.Debug("RegisterUser called")

	ruci, err := containerhelper.GetRegistrationUseCase(uss.container)
	if err != nil {
		logger.Log.Errorf("%+v\n", err)
		return nil, errors.Wrap(err, "")
	}
	mu, err := assembler.GrpcToUser(req.User)

	if err != nil {
		logger.Log.Errorf("%+v\n", err)
		return nil, errors.Wrap(err, "")
	}
	logger.Log.Debug("mu:", mu)
	resultUser, err := ruci.RegisterUser(mu)
	if err != nil {
		logger.Log.Errorf("%+v\n", err)
		return nil, errors.Wrap(err, "")
	}
	logger.Log.Debug("resultUser:", resultUser)
	gu, err := assembler.UserToGrpc(resultUser)
	if err != nil {
		logger.Log.Errorf("%+v\n", err)
		return nil, errors.Wrap(err, "")
	}

	logger.Log.Debug("user registered: ", gu)

	return &pb.RegisterUserResp{User: gu}, nil

}

func (uss *UserService) ListUser(ctx context.Context, in *pb.ListUserReq) (*pb.ListUserResp, error) {
	defer catchPanic()
	logger.Log.Debug("ListUser called")

	luci, err := containerhelper.GetListUserUseCase(uss.container)
	if err != nil {
		logger.Log.Errorf("%+v\n", err)
		return nil, errors.Wrap(err, "")
	}

	lu, err := luci.ListUser()
	if err != nil {
		logger.Log.Errorf("%+v\n", err)
		return nil, errors.Wrap(err, "")
	}
	gu, err := assembler.UserListToGrpc(lu)
	if err != nil {
		logger.Log.Errorf("%+v\n", err)
		return nil, errors.Wrap(err, "")
	}

	logger.Log.Debug("user list: ", gu)

	return &pb.ListUserResp{User: gu}, nil

}
func runServer(sc *servicecontainer.ServiceContainer) error {
	logger.Log.Debug("start runserver")

	srv := grpc.NewServer()

	cs := &UserService{sc}
	pb.RegisterUserServiceServer(srv, cs)
	//l, err:=net.Listen(GRPC_NETWORK, GRPC_ADDRESS)
	ugc := sc.AppConfig.UserGrpcConfig
	logger.Log.Debugf("userGrpcConfig: %+v\n", ugc)
	l, err := net.Listen(ugc.DriverName, ugc.UrlAddress)
	if err != nil {
		return errors.Wrap(err, "")
	} else {
		logger.Log.Debug("server listening")
	}
	return srv.Serve(l)
}

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
