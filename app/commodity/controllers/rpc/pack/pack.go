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
	"github.com/west2-online/DomTok/app/commodity/domain/model"
	modelKitex "github.com/west2-online/DomTok/kitex_gen/model"
	"github.com/west2-online/DomTok/pkg/base"
)

func BuildImage(img *model.SpuImage) *modelKitex.SpuImage {
	return &modelKitex.SpuImage{
		ImageID:   img.ImageID,
		SpuID:     img.SpuID,
		Url:       img.Url,
		CreatedAt: img.CreatedAt,
		UpdatedAt: img.UpdatedAt,
	}
}

func BuildImages(imgs []*model.SpuImage) []*modelKitex.SpuImage {
	return base.BuildTypeList(imgs, BuildImage)
}
