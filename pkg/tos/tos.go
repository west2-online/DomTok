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

package tos

import (
	"bytes"
	"context"
	"sync"

	"github.com/volcengine/ve-tos-golang-sdk/v2/tos"
	"github.com/volcengine/ve-tos-golang-sdk/v2/tos/enum"

	"github.com/west2-online/DomTok/config"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/logger"
)

var (
	tosClient *tos.ClientV2
	once      sync.Once

	initErr error
)

func newClient() (*tos.ClientV2, error) {
	once.Do(func() {
		tosClient, initErr = tos.NewClientV2(config.Tos.Endpoint,
			tos.WithRegion(config.Tos.Region),
			tos.WithCredentials(tos.NewStaticCredentials(config.Tos.AccessKey, config.Tos.SecretKey)),
		)
	})
	return tosClient, initErr
}

func UploadImage(ctx context.Context, image []byte, filePath string) error {
	client, err := newClient()
	if err != nil {
		logger.Errorf("Failed to create tos client: %v", err)
		return errno.Errorf(errno.InternalServiceErrorCode, "Failed to create tos client: %v", err)
	}

	resp, err := client.PutObjectV2(ctx, &tos.PutObjectV2Input{
		PutObjectBasicInput: tos.PutObjectBasicInput{
			Bucket: config.Tos.Bucket,
			Key:    filePath,
		},
		Content: bytes.NewReader(image),
	})
	if err != nil {
		logger.Errorf("Failed to upload image: %v", err)
		return errno.Errorf(errno.InternalServiceErrorCode, "Failed to upload image: %v", err)
	}

	logger.Infof("Upload image to tos success, object key: %s, etag: %s", filePath, resp.ETag)
	return nil
}

func DeleteImage(ctx context.Context, filePath string) error {
	client, err := newClient()
	if err != nil {
		logger.Errorf("Failed to create tos client: %v", err)
		return errno.Errorf(errno.InternalServiceErrorCode, "Failed to create tos client: %v", err)
	}

	_, err = client.DeleteObjectV2(ctx, &tos.DeleteObjectV2Input{
		Bucket: config.Tos.Bucket,
		Key:    filePath,
	})
	if err != nil {
		logger.Errorf("Failed to delete image: %v", err)
		return errno.Errorf(errno.InternalServiceErrorCode, "Failed to delete image: %v", err)
	}

	logger.Infof("Delete image from tos success, object key: %s", filePath)
	return nil
}

// SignedURL 会返回一个签名后的 URL
func SignedURL(fileName string) (string, error) {
	client, err := newClient()
	if err != nil {
		logger.Errorf("Failed to create tos client: %v", err)
		return "", errno.Errorf(errno.InternalServiceErrorCode, "Failed to create tos client: %v", err)
	}

	resp, err := client.PreSignedURL(&tos.PreSignedURLInput{
		HTTPMethod: enum.HttpMethodGet,
		Bucket:     config.Tos.Bucket,
		Key:        fileName,
	})
	if err != nil {
		logger.Errorf("Failed to sign URL: %v", err)
		return "", errno.Errorf(errno.InternalServiceErrorCode, "Failed to sign URL: %v", err)
	}
	logger.Infof("Sign URL to tos success, object key: %s", fileName)
	return resp.SignedUrl, nil
}

func MustSignedURL(fileName string) string {
	url, err := SignedURL(fileName)
	if err != nil {
		logger.Errorf("Failed to sign URL: %v", err)
	}
	return url
}
