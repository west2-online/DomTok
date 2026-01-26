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

package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"os"
	"strings"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
	"github.com/shopspring/decimal"

	"github.com/west2-online/DomTok/config"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/logger"
)

// GetMysqlDSN 会拼接 mysql 的 DSN
func GetMysqlDSN() (string, error) {
	if config.Mysql == nil {
		return "", errors.New("config not found")
	}

	dsn := strings.Join([]string{
		config.Mysql.Username, ":", config.Mysql.Password,
		"@tcp(", config.Mysql.Addr, ")/",
		config.Mysql.Database, "?charset=" + config.Mysql.Charset + "&parseTime=true",
	}, "")

	return dsn, nil
}

// GetEsHost 会获取 ElasticSearch 的客户端
func GetEsHost() (string, error) {
	if config.Elasticsearch == nil {
		return "", errors.New("elasticsearch not found")
	}

	return config.Elasticsearch.Host, nil
}

// AddrCheck 会检查当前的监听地址是否已被占用
func AddrCheck(addr string) bool {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}
	defer func() {
		if err := l.Close(); err != nil {
			logger.Errorf("utils.AddrCheck: failed to close listener: %v", err.Error())
		}
	}()
	return true
}

// GetAvailablePort 会尝试获取可用的监听地址
func GetAvailablePort() (string, error) {
	if config.Service.AddrList == nil {
		return "", errors.New("utils.GetAvailablePort: config.Service.AddrList is nil")
	}
	for _, addr := range config.Service.AddrList {
		if ok := AddrCheck(addr); ok {
			return addr, nil
		}
	}
	return "", errors.New("utils.GetAvailablePort: not available port from config")
}

// CheckImageFileType 检查文件格式是否合规
func CheckImageFileType(header *multipart.FileHeader) (string, bool) {
	file, err := header.Open()
	if err != nil {
		return "", false
	}
	defer func() {
		// 捕获并处理关闭文件时可能发生的错误
		if err := file.Close(); err != nil {
			logger.Errorf("utils.CheckImageFileType: failed to close file: %v", err.Error())
		}
	}()

	buffer := make([]byte, constants.CheckFileTypeBufferSize)
	_, err = file.Read(buffer)
	if err != nil {
		return "", false
	}

	kind, _ := filetype.Match(buffer)

	// 检查是否为jpg、png
	switch kind {
	case types.Get("jpg"):
		return "jpg", true
	case types.Get("png"):
		return "png", true
	default:
		return "", false
	}
}

// GetImageFileType 获得图片格式
func GetImageFileType(fileBytes *[]byte) (string, error) {
	buffer := (*fileBytes)[:constants.CheckFileTypeBufferSize]

	kind, _ := filetype.Match(buffer)

	// 检查是否为jpg、png
	switch kind {
	case types.Get("jpg"):
		return "jpg", nil
	case types.Get("png"):
		return "png", nil
	default:
		return "", errno.InternalServiceError
	}
}

func GetFloat64(d *decimal.Decimal) float64 {
	v, _ := d.Float64()
	return v
}

func FileToBytes(file *multipart.FileHeader) (ret [][]byte, err error) {
	if file == nil {
		return nil, errno.ParamMissingError
	}

	fileOpen, err := file.Open()
	if err != nil {
		return nil, errno.OSOperationError.WithMessage(err.Error())
	}
	defer fileOpen.Close()

	for {
		buf := make([]byte, constants.FileStreamBufferSize)
		_, err := fileOpen.Read(buf)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, errno.InternalServiceError.WithMessage(err.Error())
		}
		ret = append(ret, buf)
	}
	return ret, nil
}

// GenerateTosFilePath 会生成tos的文件路径, 格式为: path/id
func GenerateTosFilePath(path string, id int64) string {
	return fmt.Sprintf("%s/%d", path, id)
}

func DecimalFloat64(d *decimal.Decimal) float64 {
	v, _ := d.Float64()
	return v
}

func EnvironmentEnable() bool {
	return os.Getenv(constants.EnvironmentStartEnv) == constants.EnvironmentStartFlag &&
		os.Getenv(constants.EtcdEnv) != ""
}
