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

package upyun

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/upyun/go-sdk/v3/upyun"

	"github.com/west2-online/DomTok/config"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/logger"
	"github.com/west2-online/DomTok/pkg/utils"
)

var UpYun *upyun.UpYun

func NewUpYun() {
	UpYun = upyun.NewUpYun(
		&upyun.UpYunConfig{
			Bucket:   config.Upyun.Bucket,
			Operator: config.Upyun.Operator,
			Password: config.Upyun.Password,
		},
	)
	go uploadToUpYun(constants.TempSpuStorage, constants.SpuDirDEST)
	go uploadToUpYun(constants.TempSpuImageStorage, constants.SpuImageDirDest)
}

func uploadFile(src, dest string) error {
	return UpYun.Put(&upyun.PutObjectConfig{
		Path:      dest,
		UseMD5:    true,
		LocalPath: src,
	})
}

func GetImageUrl(uri string) (string, error) {
	etime := strconv.FormatInt(time.Now().Unix()+config.Upyun.TokenTimeout, 10)
	sign := utils.MD5(strings.Join([]string{config.Upyun.TokenSecret, etime, uri}, "&"))
	url := fmt.Sprintf("%s%s?_upt=%s%s", config.Upyun.UssDomain, utils.UriEncode(uri), sign[12:20], etime)
	return url, nil
}

func SaveFile(data []byte, tmpFile, destDir, filename string) error {
	err := os.MkdirAll(tmpFile, os.ModePerm)
	if err != nil {
		return errno.OSOperationError.WithError(err)
	}

	out, err := os.Create(tmpFile + filename)
	if err != nil {
		return errno.OSOperationError.WithError(err)
	}
	defer out.Close()
	_, err = io.Copy(out, bytes.NewReader(data))
	if err != nil {
		return errno.IOOperationError.WithError(err)
	}
	// go uploadToUpYun(tmpFile, destDir)
	return nil
}

func uploadToUpYun(tempDir string, destDir string) {
	ticker := time.NewTimer(constants.TickerTimer)
	defer ticker.Stop()

	for range ticker.C {
		files, err := os.ReadDir(tempDir)
		if err != nil {
			logger.Errorf("upyun ReadDir failed, err: %v", err)
			return
		}

		for _, file := range files {
			filename := file.Name()
			src := filepath.Join(tempDir, filename)
			dest := filepath.Join(destDir, filename)

			err = uploadFile(src, dest)
			if err != nil {
				logger.Errorf("upyun upload file failed, err: %v", err)
			}

			err = os.Remove(src)
			if err != nil {
				logger.Errorf("upyun remove file failed, err: %v", err)
			}
		}
		ticker.Reset(constants.TickerTimer)
	}
}
