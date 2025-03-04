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
	"time"

	"github.com/bytedance/sonic"
	"golang.org/x/sync/errgroup"

	"github.com/west2-online/DomTok/app/commodity/domain/model"
	contextLogin "github.com/west2-online/DomTok/pkg/base/context"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/logger"
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
	go svc.SendCreateSpuMsg(ctx, spu)
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
			err := upyun.UploadImg(spu.GoodsHeadDrawing, spu.GoodsHeadDrawingUrl)
			if err != nil {
				return fmt.Errorf("service.UpdateSpu: upload spuImage failed: %w", err)
			}
			return nil
		})

		eg.Go(func() error {
			err := upyun.DeleteImg(originSpu.GoodsHeadDrawingUrl)
			if err != nil {
				return fmt.Errorf("service.UpdateSpu: delete spuImage failed: %w", err)
			}
			return nil
		})

		if err := eg.Wait(); err != nil {
			return fmt.Errorf("service.UpdateSpu: update spu failed: %w", err)
		}
	}
	go svc.SendUpdateSpuMsg(ctx, spu)
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
	go svc.SendDeleteSpuMsg(ctx, spuId)
	return nil
}

func (svc *CommodityService) DeleteAllSpuImages(ctx context.Context, spuId int64) error {
	var eg errgroup.Group

	ids, urls, err := svc.db.DeleteSpuImagesBySpuId(ctx, spuId)
	if err != nil {
		return fmt.Errorf("service.DeleteAllSpuImages: delete spuImages failed: %w", err)
	}

	if len(ids) < 1 {
		return nil
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

func (svc *CommodityService) SendCreateSpuMsg(ctx context.Context, spu *model.Spu) {
	err := svc.mq.SendCreateSpuInfo(ctx, spu)
	if err != nil {
		logger.Errorf("service.SendCreateSpuMsg failed: %v", err)
	}
}

func (svc *CommodityService) SendUpdateSpuMsg(ctx context.Context, spu *model.Spu) {
	err := svc.mq.SendCreateSpuInfo(ctx, spu)
	if err != nil {
		logger.Errorf("service.SendUpdateSpuMsg failed: %v", err)
	}
}

func (svc *CommodityService) SendDeleteSpuMsg(ctx context.Context, id int64) {
	err := svc.mq.SendDeleteSpuInfo(ctx, id)
	if err != nil {
		logger.Errorf("service.SendDeleteSpuMsg failed: %v", err)
	}
}

func (svc *CommodityService) ConsumeCreateSpuMsg(ctx context.Context) {
	msgCh := svc.mq.ConsumeCreateSpuInfo(ctx)
	go func() {
		for msg := range msgCh {
			req := new(model.Spu)
			err := sonic.Unmarshal(msg.V, req)
			if err != nil {
				logger.Errorf("service.ConsumeCreateSpuMsg Unmarshal failed: %v", err)
			}
			err = svc.es.AddItem(ctx, constants.SpuTableName, req)
			if err != nil {
				logger.Errorf("service.ConsumeCreateSpuMsg add item failed: %v", err)
			}
		}
	}()
}

func (svc *CommodityService) ConsumeUpdateSpuMsg(ctx context.Context) {
	msgCh := svc.mq.ConsumeUpdateSpuInfo(ctx)
	go func() {
		for msg := range msgCh {
			req := new(model.Spu)
			err := sonic.Unmarshal(msg.V, req)
			if err != nil {
				logger.Errorf("service.ConsumeUpdateSpuMsg Unmarshal failed: %v", err)
			}
			err = svc.es.UpdateItem(ctx, constants.SpuTableName, req)
			if err != nil {
				logger.Errorf("service.ConsumeUpdateSpuMsg update item failed: %v", err)
			}
		}
	}()
}

func (svc *CommodityService) ConsumeDeleteSpuMsg(ctx context.Context) {
	msgCh := svc.mq.ConsumeDeleteSpuInfo(ctx)
	go func() {
		for msg := range msgCh {
			idStr := string(msg.V)
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				logger.Errorf("CommodityService.ConsumeDeleteSpuMsg: invalid param, %v", errno.ParamVerifyError.WithMessage(err.Error()))
			}

			err = svc.es.RemoveItem(ctx, constants.SpuTableName, id)
			if err != nil {
				logger.Errorf("service.ConsumeDeleteSpuMsg delete item failed: %v", err)
			}
		}
	}()
}

func (svc *CommodityService) IsSpuMappingExist(ctx context.Context) error {
	var err error
	if !svc.es.IsExist(ctx, constants.SpuTableName) {
		err = svc.es.CreateIndex(ctx, constants.SpuTableName)
		if err != nil {
			return fmt.Errorf("service.IsSpuMappingExist CreateIndex failed: %w", err)
		}
	}
	return err
}

func (svc *CommodityService) CreateCategory(ctx context.Context, category *model.Category) error {
	category.Id = svc.nextID()
	if err := svc.db.CreateCategory(ctx, category); err != nil {
		return fmt.Errorf("create category failed: %w", err)
	}
	return nil
}

func (svc *CommodityService) DeleteCategory(ctx context.Context, category *model.Category) error {
	err := svc.db.DeleteCategory(ctx, category)
	if err != nil {
		return fmt.Errorf("delete category failed: %w", err)
	}
	return nil
}

func (svc *CommodityService) UpdateCategory(ctx context.Context, category *model.Category) error {
	err := svc.db.UpdateCategory(ctx, category)
	if err != nil {
		return fmt.Errorf("update category failed: %w", err)
	}
	return err
}

func (svc *CommodityService) Cached(ctx context.Context, infos []*model.SkuBuyInfo) bool {
	cached := true

	for _, info := range infos {
		key := svc.cache.GetLockStockKey(info.SkuID)
		stockKey := svc.cache.GetStockKey(info.SkuID)

		lockStockCached := true
		StockCached := true

		if !svc.cache.IsExist(ctx, key) {
			lockStockCached = false
			cached = false
		}
		if !svc.cache.IsExist(ctx, stockKey) {
			StockCached = false
			cached = false
		}

		if !lockStockCached || !StockCached {
			go func(id int64, stockKey, lockStockKey string, c context.Context) {
				data, err := svc.db.GetSkuById(c, id)
				if err != nil {
					return
				}
				if !lockStockCached {
					svc.cache.SetLockStockNum(c, lockStockKey, data.LockStock)
				}

				if !StockCached {
					svc.cache.SetLockStockNum(c, stockKey, data.Stock)
				}
			}(info.SkuID, stockKey, key, ctx)
		}
	}
	return cached
}

func (svc *CommodityService) IncrLockStockInNX(ctx context.Context, infos []*model.SkuBuyInfo) error {
	var err error
	if !svc.Cached(ctx, infos) {
		return errno.Errorf(errno.RedisKeyNotExist, "CommodityService.IncrLockStockInNX failed")
	}

	keys := make([]string, 0)

	for _, info := range infos {
		keys = append(keys, svc.cache.GetSkuKey(info.SkuID))
	}

	err = svc.cache.Lock(ctx, keys, constants.RedisNXExpireTime)
	if err != nil {
		return fmt.Errorf("CommodityService.IncrLockStockInNX failed: %w", err)
	}
	defer func() {
		err = svc.cache.UnLock(ctx, keys)
		if err != nil {
			logger.Errorf("CommodityService.IncrLockStockInNX failed: %v", err)
		}
	}()

	var eg errgroup.Group

	eg.Go(func() error {
		err = svc.cache.IncrLockStockNum(ctx, infos)
		if err != nil {
			return fmt.Errorf("CommodityService.IncrLockStock failed: %w", err)
		}
		return nil
	})

	eg.Go(func() error {
		err = svc.db.IncrLockStockInNX(ctx, infos)
		if err != nil {
			return fmt.Errorf("CommodityService.IncrLockStock failed: %w", err)
		}
		return nil
	})

	if err = eg.Wait(); err != nil {
		return err
	}
	return nil
}

func (svc *CommodityService) DecrLockStockInNX(ctx context.Context, infos []*model.SkuBuyInfo) error {
	var err error
	if !svc.Cached(ctx, infos) {
		return errno.Errorf(errno.RedisKeyNotExist, "CommodityService.DecrLockStockInNX failed")
	}

	keys := make([]string, 0)

	for _, info := range infos {
		keys = append(keys, svc.cache.GetSkuKey(info.SkuID))
	}

	err = svc.cache.Lock(ctx, keys, constants.RedisNXExpireTime)
	if err != nil {
		return fmt.Errorf("CommodityService.DecrStockInNX failed: %w", err)
	}
	defer func() {
		err = svc.cache.UnLock(ctx, keys)
		if err != nil {
			logger.Errorf("CommodityService.DecrStockInNX failed: %v", err)
		}
	}()

	var eg errgroup.Group

	eg.Go(func() error {
		err = svc.cache.DecrLockStockNum(ctx, infos)
		if err != nil {
			return fmt.Errorf("CommodityService.DecrLockStockInNX failed: %w", err)
		}
		return nil
	})

	eg.Go(func() error {
		err = svc.db.DecrLockStockInNX(ctx, infos)
		if err != nil {
			return fmt.Errorf("CommodityService.DecrLockStockInNX failed: %w", err)
		}
		return nil
	})

	if err = eg.Wait(); err != nil {
		return err
	}
	return nil
}

func (svc *CommodityService) DecrStockInNX(ctx context.Context, infos []*model.SkuBuyInfo) error {
	var err error
	if !svc.Cached(ctx, infos) {
		return errno.Errorf(errno.RedisKeyNotExist, "CommodityService.DecrStockInNX failed")
	}

	keys := make([]string, 0)

	for _, info := range infos {
		keys = append(keys, svc.cache.GetSkuKey(info.SkuID))
	}

	err = svc.cache.Lock(ctx, keys, constants.RedisNXExpireTime)
	if err != nil {
		return fmt.Errorf("CommodityService.DecrStockInNX failed: %w", err)
	}
	defer func() {
		err = svc.cache.UnLock(ctx, keys)
		if err != nil {
			logger.Errorf("CommodityService.DecrStockInNX failed: %v", err)
		}
	}()

	var eg errgroup.Group

	eg.Go(func() error {
		err = svc.cache.DecrStockNum(ctx, infos)
		if err != nil {
			return fmt.Errorf("CommodityService.DecrStockInNX failed: %w", err)
		}
		return nil
	})

	eg.Go(func() error {
		err = svc.db.DecrStockInNX(ctx, infos)
		if err != nil {
			return fmt.Errorf("CommodityService.DecrStockInNX failed: %w", err)
		}
		return nil
	})
	if err = eg.Wait(); err != nil {
		return err
	}
	return nil
}

func (svc *CommodityService) CreateSku(ctx context.Context, sku *model.Sku, ext string) (*model.Sku, error) {
	sku.SkuID = svc.nextID()
	sku.HistoryID = svc.nextID()
	sku.StyleHeadDrawingUrl = utils.GenerateFileName(constants.SkuDirDest, sku.SkuID) + ext
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
		return nil, err
	}
	s := &model.Sku{
		SkuID:     sku.SkuID,
		CreatorID: sku.CreatorID,
		HistoryID: sku.HistoryID,
	}
	return s, nil
}

func (svc *CommodityService) UpdateSku(ctx context.Context, sku *model.Sku, originSpu *model.Sku) error {
	sku.HistoryID = svc.nextID()
	ret, err := svc.db.GetSkuById(ctx, sku.SkuID)
	if err != nil {
		return fmt.Errorf("service.UpdateSku: get sku failed: %w", err)
	}
	if len(sku.StyleHeadDrawing) > 0 {
		var eg errgroup.Group
		eg.Go(func() error {
			err := upyun.DeleteImg(originSpu.StyleHeadDrawingUrl)
			if err != nil {
				return errno.UpYunFileError.WithMessage(err.Error())
			}
			return nil
		})

		eg.Go(func() error {
			err := upyun.UploadImg(sku.StyleHeadDrawing, sku.StyleHeadDrawingUrl)
			if err != nil {
				return errno.UpYunFileError.WithMessage(err.Error())
			}
			return nil
		})

		if err := eg.Wait(); err != nil {
			return err
		}
	}

	if err := svc.db.UpdateSku(ctx, sku, ret); err != nil {
		return fmt.Errorf("service.UpdateSku: update sku failed: %w", err)
	}

	return nil
}

func (svc *CommodityService) DeleteSku(ctx context.Context, sku *model.Sku) error {
	err := svc.db.DeleteSku(ctx, sku)
	if err != nil {
		return fmt.Errorf("usecase.DeleteSku failed: %w", err)
	}

	err = upyun.DeleteImg(sku.StyleHeadDrawingUrl)
	if err != nil {
		return errno.UpYunFileError.WithMessage(err.Error())
	}

	return nil
}

func (svc *CommodityService) ViewSku(ctx context.Context, skuIds []*int64, pageNum, pageSize int) ([]*model.Sku, int64, error) {
	skus, total, err := svc.db.ViewSku(ctx, skuIds, pageNum, pageSize)
	if err != nil {
		return nil, -1, fmt.Errorf("usecase.ViewSku failed: %w", err)
	}

	for _, sku := range skus {
		stockKey := svc.cache.GetStockKey(sku.SkuID)
		lockStockKey := svc.cache.GetLockStockKey(sku.SkuID)

		if svc.cache.IsExist(ctx, stockKey) {
			sku.Stock, err = svc.cache.GetLockStockNum(ctx, stockKey)
			if err != nil {
				return nil, -1, fmt.Errorf("service.ViewSku failed: %w", err)
			}
		} else {
			svc.cache.SetLockStockNum(ctx, stockKey, sku.Stock)
		}

		if svc.cache.IsExist(ctx, lockStockKey) {
			sku.LockStock, err = svc.cache.GetLockStockNum(ctx, lockStockKey)
			if err != nil {
				return nil, -1, fmt.Errorf("service.ViewSku failed: %w", err)
			}
		} else {
			svc.cache.SetLockStockNum(ctx, lockStockKey, sku.LockStock)
		}
	}

	return skus, total, nil
}

func (svc *CommodityService) UploadSkuAttr(ctx context.Context, attr *model.AttrValue, sku *model.Sku) error {
	id := svc.nextID()

	if err := svc.db.UploadSkuAttr(ctx, sku, attr, id); err != nil {
		return fmt.Errorf("usecase.UploadSkuAttr failed: %w", err)
	}

	return nil
}

func (svc *CommodityService) ListSkuInfo(ctx context.Context, skuInfo []*model.SkuVersion, pageNum, pageSize int64) ([]*model.Sku, int64, error) {
	var (
		skuInfos []*model.Sku
		total    int64
	)

	skuInfos, err := svc.db.ListSkuInfo(ctx, skuInfo, int(pageNum), int(pageSize))
	if err != nil {
		return nil, -1, fmt.Errorf("usecase.ListSkuInfo failed: %w", err)
	}

	return skuInfos, total, nil
}

func (svc *CommodityService) ViewSkuPriceHistory(ctx context.Context, s *model.SkuPriceHistory, pNum int64, pSize int64) ([]*model.SkuPriceHistory, error) {
	histories, err := svc.db.ViewSkuPriceHistory(ctx, s, int(pNum), int(pSize))
	if err != nil {
		return nil, fmt.Errorf("usecase.ViewSkuPriceHistory failed: %w", err)
	}
	return histories, nil
}

func (svc *CommodityService) CreateSkuImage(ctx context.Context, skuImage *model.SkuImage, data []byte) (int64, error) {
	skuImage.ImageID = svc.nextID()
	skuImage.Url = utils.GenerateFileName(constants.SkuImageDirDest, skuImage.ImageID)
	var eg errgroup.Group

	eg.Go(func() error {
		if err := upyun.UploadImg(data, skuImage.Url); err != nil {
			return fmt.Errorf("service.CreateSkuImage: upload skuImage failed: %w", err)
		}
		return nil
	})

	eg.Go(func() error {
		if err := svc.db.CreateSkuImage(ctx, skuImage); err != nil {
			return fmt.Errorf("service.CreateSkuImage: create skuImage failed: %w", err)
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return 0, err
	}

	return skuImage.ImageID, nil
}

func (svc *CommodityService) UpdateSkuImage(ctx context.Context, skuImage *model.SkuImage, originSkuImages *model.SkuImage, data []byte) error {
	var eg errgroup.Group

	eg.Go(func() error {
		err := upyun.DeleteImg(originSkuImages.Url)
		if err != nil {
			return errno.UpYunFileError.WithMessage(err.Error())
		}
		return nil
	})

	eg.Go(func() error {
		err := upyun.UploadImg(data, skuImage.Url)
		if err != nil {
			return errno.UpYunFileError.WithMessage(err.Error())
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	if err := svc.db.UpdateSkuImage(ctx, skuImage); err != nil {
		return fmt.Errorf("service.UpdateSkuImage: update skuImage failed: %w", err)
	}

	return nil
}

func (svc *CommodityService) ViewSkuImages(ctx context.Context, sku *model.Sku, pageNum, pageSize int) ([]*model.SkuImage, int64, error) {
	offset := pageNum * pageSize
	key := fmt.Sprintf("skuImgs:%d:%d", sku.SkuID, offset)
	if svc.cache.IsExist(ctx, key) {
		ret, err := svc.cache.GetSkuImages(ctx, key)
		if err != nil {
			return nil, -1, fmt.Errorf("service.GetSkuImages failed: %w", err)
		}
		return ret, -1, nil
	}

	images, total, err := svc.db.ViewSkuImage(ctx, sku, pageNum, pageSize)
	if err != nil {
		return nil, -1, fmt.Errorf("usecase.ViewSkuImage failed: %w", err)
	}

	svc.cache.SetSkuImages(ctx, key, images)

	return images, total, nil
}

func (svc *CommodityService) DeleteSkuImage(ctx context.Context, imageId int64, url string) error {
	if err := svc.db.DeleteSkuImage(ctx, imageId); err != nil {
		return fmt.Errorf("service.DeleteSkuImage: delete skuImage failed: %w", err)
	}

	if err := upyun.DeleteImg(url); err != nil {
		return fmt.Errorf("service.DeleteSkuImage: delete skuImage failed: %w", err)
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

func (svc *CommodityService) GetSkuFromImageId(ctx context.Context, imageId int64) (*model.Sku, *model.SkuImage, error) {
	img, err := svc.db.GetSkuImageByImageId(ctx, imageId)
	if err != nil {
		return nil, nil, fmt.Errorf("service.GetSkuFromImageId: get image info failed: %w", err)
	}

	ret, err := svc.db.GetSkuBySkuId(ctx, img.SkuID)
	if err != nil {
		return nil, nil, fmt.Errorf("service.GetSkuFromImageId: get sku info failed: %w", err)
	}
	return ret, img, nil
}

func (svc *CommodityService) CheckoutRedisHealth() {
	for {
		err := svc.cache.IsHealthy(context.Background())
		if err != nil {
			RedisAvailable.Store(constants.RedisUnHealthy)
		} else {
			RedisAvailable.Store(constants.RedisHealthy)
		}
		time.Sleep(constants.RedisCheckoutInterval)
	}
}

func (svc *CommodityService) IsHealthy() bool {
	return RedisAvailable.Load()
}
