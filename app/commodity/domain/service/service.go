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

	"golang.org/x/sync/errgroup"

	"github.com/west2-online/DomTok/app/commodity/domain/model"
	contextLogin "github.com/west2-online/DomTok/pkg/base/context"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/upyun"
	"github.com/west2-online/DomTok/pkg/utils"
)

func (svc *CommodityService) nextID() int64 {
	id, _ := svc.sf.NextVal()
	return id
}

func (svc *CommodityService) CreateSpu(ctx context.Context, spu *model.Spu) (int64, error) {
	spu.SpuId = svc.nextID()
	spu.GoodsHeadDrawingUrl = utils.GenerateFileName(constants.SpuDirDest, spu.SpuId)
	var eg errgroup.Group

	eg.Go(func() error {
		if err := svc.db.CreateSpu(ctx, spu); err != nil {
			return fmt.Errorf("service.CreateSpu: create spu failed: %w", err)
		}
		return nil
	})

	eg.Go(func() error {
		if err := upyun.UploadImg(spu.GoodsHeadDrawing, spu.GoodsHeadDrawingUrl); err != nil {
			return fmt.Errorf("service.UploadImg: upload image failed: %w", err)
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
	spuImage.Url = utils.GenerateFileName(constants.SpuImageDirDest, spuImage.SpuID)

	var eg errgroup.Group

	eg.Go(func() error {
		if err := upyun.UploadImg(spuImage.Data, spuImage.Url); err != nil {
			return fmt.Errorf("service.CreateSpuImage: upload spuImage failed: %w", err)
		}
		return nil
	})

	eg.Go(func() error {
		if err := svc.db.CreateSpuImage(ctx, spuImage); err != nil {
			return fmt.Errorf("service.CreateSpuImage: create spuImage failed: %w", err)
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return 0, err
	}

	return spuImage.ImageID, nil
}

// TODO
func (svc *CommodityService) UpdateSpuImage(ctx context.Context, spuImage *model.SpuImage, originSpuImage *model.SpuImage) error {
	err := svc.db.UpdateSpuImage(ctx, spuImage)
	if err != nil {
		return fmt.Errorf("service.UpdateSpu: update spu failed: %w", err)
	}

	var eg errgroup.Group
	eg.Go(func() error {
		err = upyun.UploadImg(spuImage.Data, spuImage.Url)
		if err != nil {
			return fmt.Errorf("service.UpdateSpuImage: upload spuImage failed: %w", err)
		}
		return nil
	})

	eg.Go(func() error {
		err = upyun.DeleteImg(originSpuImage.Url)
		if err != nil {
			return fmt.Errorf("service.UpdateSpuImage: delete spuImage failed: %w", err)
		}
		return nil
	})

	if err = eg.Wait(); err != nil {
		return err
	}

	return nil
}

func (svc *CommodityService) UpdateSpu(ctx context.Context, spu *model.Spu, originSpu *model.Spu) error {
	err := svc.db.UpdateSpu(ctx, spu)
	if err != nil {
		return fmt.Errorf("service.UpdateSpu: update spu failed: %w", err)
	}

	if len(spu.GoodsHeadDrawing) > 0 {
		var eg errgroup.Group
		eg.Go(func() error {
			err = upyun.UploadImg(spu.GoodsHeadDrawing, spu.GoodsHeadDrawingUrl)
			if err != nil {
				return fmt.Errorf("service.UpdateSpu: upload spuImage failed: %w", err)
			}
			return nil
		})

		eg.Go(func() error {
			err = upyun.DeleteImg(originSpu.GoodsHeadDrawingUrl)
			if err != nil {
				return fmt.Errorf("service.UpdateSpu: delete spuImage failed: %w", err)
			}
			return nil
		})

		if err = eg.Wait(); err != nil {
			return fmt.Errorf("service.UpdateSpu: update spu failed: %w", err)
		}
	}
	return nil
}

func (svc *CommodityService) DeleteSpuImage(ctx context.Context, imageId int64, url string) error {
	var eg errgroup.Group
	eg.Go(func() error {
		if err := svc.db.DeleteSpuImage(ctx, imageId); err != nil {
			return fmt.Errorf("service.DeleteSpuImage: delete spuImage failed: %w", err)
		}
		return nil
	})

	eg.Go(func() error {
		if err := upyun.DeleteImg(upyun.GetImageUrl(url)); err != nil {
			return fmt.Errorf("service.DeleteSpuImage: delete spuImage failed: %w", err)
		}
		return nil
	})
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

func (svc *CommodityService) DeleteSpu(ctx context.Context, spuId int64, url string) error {
	var eg errgroup.Group
	eg.Go(func() error {
		if err := svc.db.DeleteSpu(ctx, spuId); err != nil {
			return fmt.Errorf("service.DeleteSpu: delete spu failed: %w", err)
		}
		return nil
	})

	eg.Go(func() error {
		if err := upyun.DeleteImg(upyun.GetImageUrl(url)); err != nil {
			return fmt.Errorf("service.DeleteSpu: delete spuImage failed: %w", err)
		}
		return nil
	})
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

func (svc *CommodityService) DeleteAllSpuImages(ctx context.Context, spuId int64) error {

	var eg errgroup.Group

	ids, urls, err := svc.db.DeleteSpuImagesBySpuId(ctx, spuId)
	if err != nil {
		return fmt.Errorf("service.DeleteAllSpuImages: delete spuImages failed: %w", err)
	}

	for i := 0; i < len(ids); i++ {
		eg.Go(func() error {
			if err = upyun.DeleteImg(upyun.GetImageUrl(urls[i])); err != nil {
				return fmt.Errorf("service.DeleteAllSpuImages: delete spuImages failed: %w", err)
			}
			return nil
		})
	}
	if err = eg.Wait(); err != nil {
		return err
	}
	return nil
}

func (svc *CommodityService) GetSpuFromImageId(ctx context.Context, imageId int64) (*model.Spu, *model.SpuImage, error) {
	img, err := svc.db.GetSpuImage(ctx, imageId)
	if err != nil {
		return nil, nil, fmt.Errorf("service.GetSpuFromImageId: get image info failed: %w", err)
	}

	ret, err := svc.db.GetSpuBySpuId(ctx, img.SpuID)
	if err != nil {
		return nil, nil, fmt.Errorf("service.GetSpuFromImageId: get spu info failed: %w", err)
	}
	return ret, img, nil
}

func (svc *CommodityService) IdentifyUser(ctx context.Context, uid int64) error {
	loginData, err := contextLogin.GetLoginData(ctx)
	if err != nil {
		return fmt.Errorf("usecase.DeleteSpu failed: %w", err)
	}

	if loginData != uid {
		return errno.AuthNoOperatePermission
	}
	return nil
}

func (svc *CommodityService) MatchDeleteSpuCondition(ctx context.Context, spuId int64) (*model.Spu, error) {
	exists, err := svc.db.IsExistSku(ctx, spuId)
	if err != nil {
		return nil, fmt.Errorf("usecase.DeleteSpu failed: %w", err)
	}
	if exists {
		return nil, errno.Errorf(errno.ServiceSkuExist, "usecase.DeleteSpu failed: spu-%d‘s sku already exists", spuId)
	}

	ret, err := svc.db.GetSpuBySpuId(ctx, spuId)
	if err != nil {
		return nil, fmt.Errorf("usecase.DeleteSpu failed: %w", err)
	}
	return ret, nil
}
