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

func (svc *CommodityService) UpdateSpuImage(ctx context.Context, spuImage *model.SpuImage, originSpuImage *model.SpuImage) error {
	var eg errgroup.Group
	var err error
	eg.Go(func() error {
		err = svc.db.UpdateSpuImage(ctx, spuImage)
		if err != nil {
			return fmt.Errorf("service.UpdateSpu: update spu failed: %w", err)
		}
		return nil
	})

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
		if err := upyun.DeleteImg(url); err != nil {
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
		if err := upyun.DeleteImg(url); err != nil {
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
			if err = upyun.DeleteImg(urls[i]); err != nil {
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
		return fmt.Errorf("service.IdentifyUser failed: %w", err)
	}

	if loginData != uid {
		return errno.AuthNoOperatePermission
	}
	return nil
}

func (svc *CommodityService) IdentifyUserInStreamCtx(ctx context.Context, uid int64) error {
	loginData, err := contextLogin.GetStreamLoginData(ctx)
	if err != nil {
		return fmt.Errorf("service.IdentifyUserInStreamCtx failed: %w", err)
	}
	if loginData != uid {
		return errno.AuthNoOperatePermission
	}
	return nil
}

func (svc *CommodityService) MatchDeleteSpuCondition(ctx context.Context, spuId int64) (*model.Spu, error) {
	exists, err := svc.db.IsExistSku(ctx, spuId)
	if err != nil {
		return nil, fmt.Errorf("service.MatchDeleteSpuCondition failed: %w", err)
	}
	if exists {
		return nil, errno.Errorf(errno.ServiceSkuExist, "service.MatchDeleteSpuCondition failed: spu-%d‘s sku already exists", spuId)
	}

	ret, err := svc.db.GetSpuBySpuId(ctx, spuId)
	if err != nil {
		return nil, fmt.Errorf("service.MatchDeleteSpuConditionu failed: %w", err)
	}
	return ret, nil
}

func (svc *CommodityService) GetSpuImages(ctx context.Context, spuId int64, offset, limit int) ([]*model.SpuImage, int64, error) {
	key := fmt.Sprintf("spuImgs:%d:%d", spuId, offset)
	if svc.cache.IsExist(ctx, key) {
		ret, err := svc.cache.GetSpuImages(ctx, key)
		if err != nil {
			return nil, 0, fmt.Errorf("service.GetSpuImages failed: %w", err)
		}
		return ret.Images, ret.Total, nil
	}
	_, err := svc.db.GetSpuBySpuId(ctx, spuId)
	if err != nil {
		return nil, 0, fmt.Errorf("usecase.ViewSpuImages failed: %w", err)
	}

	imgs, total, err := svc.db.GetImagesBySpuId(ctx, spuId, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("usecase.ViewSpuImages get images failed: %w", err)
	}
	go svc.cache.SetSpuImages(ctx, key, &model.SpuImages{Images: imgs, Total: total})
	return imgs, total, nil
}

func (svc *CommodityService) CreateSku(ctx context.Context, sku *model.Sku) (int64, error) {
	sku.SkuID = svc.nextID()
	sku.StyleHeadDrawingUrl = utils.GenerateFileName(constants.SkuDirDest, sku.SkuID)

	var eg errgroup.Group
	eg.Go(func() error {
		if err := svc.db.CreateSku(ctx, sku); err != nil {
			return fmt.Errorf("service.CreateSku: create sku failed: %w", err)
		}
		return nil
	})

	eg.Go(func() error {
		if err := upyun.UploadImg(sku.StyleHeadDrawing, sku.StyleHeadDrawingUrl); err != nil {
			return fmt.Errorf("service.UploadImg: upload image failed: %w", err)
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return 0, err
	}

	return sku.SkuID, nil
}

func (svc *CommodityService) UpdateSku(ctx context.Context, sku *model.Sku, originSpu *model.Sku) error {
	if err := svc.db.UpdateSku(ctx, sku); err != nil {
		return fmt.Errorf("service.UpdateSku: update sku failed: %w", err)
	}

	if len(sku.StyleHeadDrawing) > 0 {
		var eg errgroup.Group
		eg.Go(func() error {
			err := upyun.UploadImg(sku.StyleHeadDrawing, sku.StyleHeadDrawingUrl)
			if err != nil {
				return errno.UpYunFileError.WithMessage(err.Error())
			}
			return nil
		})

		eg.Go(func() error {
			err := upyun.DeleteImg(originSpu.StyleHeadDrawingUrl)
			if err != nil {
				return errno.UpYunFileError.WithMessage(err.Error())
			}
			return nil
		})

		if err := eg.Wait(); err != nil {
			return err
		}
	}

	return nil
}

func (svc *CommodityService) GetSkuIdBySpuID(ctx context.Context, spuID int64, pageNum int, pageSize int) ([]*int64, error) {
	skuIds, err := svc.db.GetSkuIdBySpuID(ctx, spuID, pageNum, pageSize)
	if err != nil {
		return nil, fmt.Errorf("service.GetSkuIdBySpuID: get sku id by spu id failed: %w", err)
	}
	return skuIds, nil
}

func (svc *CommodityService) SetCreatorID(ctx context.Context, sku *model.Sku) error {
	uid, err := contextLogin.GetLoginData(ctx)
	if err != nil {
		return errno.Errorf(errno.AuthNoTokenCode, "no token find")
	}
	sku.CreatorID = uid
	return nil
}

func (svc *CommodityService) NormalizePagination(pageNum, pageSize *int64) (int, int) {
	var (
		pNum  int64
		pSize int64
	)

	if pageNum != nil && *pageNum > 0 {
		pNum = *pageNum
	}
	if pageSize != nil && *pageSize > 0 {
		pSize = *pageSize
	}

	return int(pNum), int(pSize)
}
