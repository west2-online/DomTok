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
	"net/http"
	"strings"
	"time"

	"github.com/west2-online/DomTok/config"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/logger"
)

// UploadImg 又拍云上传文件
func UploadImg(file []byte, url string) error {
	body := bytes.NewReader(file)
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return errno.UpYunFileError.WithMessage(err.Error())
	}
	req.SetBasicAuth(config.Upyun.Operator, config.Upyun.Password)
	req.Header.Add("Date", time.Now().UTC().Format(http.TimeFormat))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errno.UpYunFileError.WithMessage(err.Error())
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			logger.Errorf("uploadimg close request meet error: %v", err)
		}
	}()
	if res.StatusCode != http.StatusOK {
		return errno.UpYunFileError
	}
	return nil
}

// DeleteImg 又拍云删除文件
func DeleteImg(url string) error {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return errno.UpYunFileError.WithMessage(err.Error())
	}
	req.SetBasicAuth(config.Upyun.Operator, config.Upyun.Password)
	req.Header.Add("Date", time.Now().UTC().Format(http.TimeFormat))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errno.UpYunFileError.WithMessage(err.Error())
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			logger.Errorf("deleteImg close request meet error: %v", err)
		}
	}()
	if res.StatusCode != http.StatusOK {
		return errno.UpYunFileError
	}
	return nil
}

// 获取图片URL
func GetImageUrl(url string) string {
	return strings.Join([]string{config.Upyun.DownloadDomain, strings.TrimPrefix(url, config.Upyun.UssDomain)}, "")
}
