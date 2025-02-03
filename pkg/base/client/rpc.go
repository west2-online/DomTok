/*
Copyright 2024 The west2-online Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package client

import (
	"errors"
	"fmt"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/streamclient"
	etcd "github.com/kitex-contrib/registry-etcd"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/kitex_gen/academic/academicservice"
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom/classroomservice"
	"github.com/west2-online/fzuhelper-server/kitex_gen/common/commonservice"
	"github.com/west2-online/fzuhelper-server/kitex_gen/course/courseservice"
	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen/launchscreenservice"
	"github.com/west2-online/fzuhelper-server/kitex_gen/paper/paperservice"
	"github.com/west2-online/fzuhelper-server/kitex_gen/user/userservice"
	"github.com/west2-online/fzuhelper-server/kitex_gen/version/versionservice"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
)

// 通用的RPC客户端初始化函数
func initRPCClient[T any](serviceName string, newClientFunc func(string, ...client.Option) (T, error)) (*T, error) {
	if config.Etcd == nil || config.Etcd.Addr == "" {
		return nil, errors.New("config.Etcd.Addr is nil")
	}
	// 初始化Etcd Resolver
	r, err := etcd.NewEtcdResolver([]string{config.Etcd.Addr})
	if err != nil {
		return nil, fmt.Errorf("initRPCClient etcd.NewEtcdResolver failed: %w", err)
	}
	// 初始化具体的RPC客户端
	client, err := newClientFunc(serviceName, client.WithResolver(r), client.WithMuxConnection(constants.MuxConnection))
	if err != nil {
		return nil, fmt.Errorf("initRPCClient NewClient failed: %w", err)
	}
	return &client, nil
}

func InitUserRPC() (*userservice.Client, error) {
	return initRPCClient(constants.UserServiceName, userservice.NewClient)
}

func InitClassroomRPC() (*classroomservice.Client, error) {
	return initRPCClient(constants.ClassroomServiceName, classroomservice.NewClient)
}

func InitCourseRPC() (*courseservice.Client, error) {
	return initRPCClient(constants.CourseServiceName, courseservice.NewClient)
}

func InitLaunchScreenRPC() (*launchscreenservice.Client, error) {
	if config.Etcd == nil || config.Etcd.Addr == "" {
		return nil, errors.New("config.Etcd.Addr is nil")
	}
	r, err := etcd.NewEtcdResolver([]string{config.Etcd.Addr})
	if err != nil {
		return nil, fmt.Errorf("InitLaunchScreenRPC etcd.NewEtcdResolver failed: %w", err)
	}
	client, err := launchscreenservice.NewClient(constants.LaunchScreenServiceName, client.WithResolver(r), client.WithMuxConnection(constants.MuxConnection))
	if err != nil {
		return nil, fmt.Errorf("InitLaunchScreenRPC NewClient failed: %w", err)
	}
	return &client, nil
}

func InitLaunchScreenStreamRPC() (*launchscreenservice.StreamClient, error) {
	if config.Etcd == nil || config.Etcd.Addr == "" {
		return nil, errors.New("config.Etcd.Addr is nil")
	}
	r, err := etcd.NewEtcdResolver([]string{config.Etcd.Addr})
	if err != nil {
		return nil, fmt.Errorf("InitLaunchScreenStreamRPC etcd.NewEtcdResolver failed: %w", err)
	}
	streamClient := launchscreenservice.MustNewStreamClient(constants.LaunchScreenServiceName, streamclient.WithResolver(r))

	return &streamClient, nil
}

func InitPaperRPC() (*paperservice.Client, error) {
	return initRPCClient(constants.PaperServiceName, paperservice.NewClient)
}

func InitAcademicRPC() (*academicservice.Client, error) {
	return initRPCClient(constants.AcademicServiceName, academicservice.NewClient)
}

func InitVersionRPC() (*versionservice.Client, error) {
	return initRPCClient(constants.VersionServiceName, versionservice.NewClient)
}

func InitCommonRPC() (*commonservice.Client, error) {
	return initRPCClient(constants.CommonServiceName, commonservice.NewClient)
}
