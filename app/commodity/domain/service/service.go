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

package service

import (
	"context"
	"fmt"
	"strconv"

	"golang.org/x/sync/errgroup"

	"github.com/west2-online/DomTok/app/commodity/domain/model"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/upyun"
)

func (svc *CommodityService) nextID() int64 {
	id, _ := svc.sf.NextVal()
	return id
}

func (svc *CommodityService) CreateSpu(ctx context.Context, spu *model.Spu) (int64, error) {
	spu.SpuId = svc.nextID()
	spu.GoodsHeadDrawingName = strconv.FormatInt(spu.SpuId, 10) + "-" + spu.GoodsHeadDrawingName
	var eg errgroup.Group

	eg.Go(func() error {
		if err := upyun.SaveFile(spu.GoodsHeadDrawing, constants.TempSpuStorage, spu.GoodsHeadDrawingName); err != nil {
			return fmt.Errorf("service.CreateSpu: save file failed: %w", err)
		}
		return nil
	})

	eg.Go(func() error {
		if err := svc.db.CreateSpu(ctx, spu); err != nil {
			return fmt.Errorf("service.CreateSpu: create spu failed: %w", err)
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return 0, err
	}
	return spu.SpuId, nil
}

func (svc *CommodityService) CreateSpuImage(ctx context.Context, spuImage *model.SpuImage) (int64, error) {
	spuImage.ImageID = svc.nextID()
	spuImage.Url = strconv.FormatInt(spuImage.ImageID, 10) + "-" + spuImage.Url
	var eg errgroup.Group

	eg.Go(func() error {
		if err := upyun.SaveFile(spuImage.Data, constants.TempSpuImageStorage, spuImage.Url); err != nil {
			return fmt.Errorf("service.CreateSpuImage: create spuImage failed: %w", err)
		}
		return nil
	})

	eg.Go(func() error {
		if err := svc.db.CreateSpuImage(ctx, spuImage); err != nil { // TODO 后续可引入mq执行
			return fmt.Errorf("service.CreateSpuImage: create spuImage failed: %w", err)
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return 0, err
	}
	return spuImage.ImageID, nil
}
