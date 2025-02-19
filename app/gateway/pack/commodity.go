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

package pack

import (
	"io/ioutil"
	"mime/multipart"

	"github.com/west2-online/DomTok/pkg/errno"
)

func BuildFileDataBytes(file *multipart.FileHeader) ([]byte, error) {
	src, err := file.Open()
	if err != nil {
		return nil, errno.OSOperationError.WithError(err)
	}
	defer src.Close()

	data, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, errno.IOOperationError.WithError(err)
	}
	return data, err
}
