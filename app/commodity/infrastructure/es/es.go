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

package es

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/olivere/elastic/v7"
	"github.com/west2-online/DomTok/app/commodity/domain/model"
	"github.com/west2-online/DomTok/kitex_gen/commodity"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/logger"
)

func (es *CommodityElastic) IsExist(ctx context.Context, indexName string) bool {
	res, err := es.client.IndexExists(indexName).Do(ctx)
	if err != nil {
		logger.Errorf("CommodityElastic.IsExist Error checking if index exists: %v", err)
		return false
	}
	return res
}

func (es *CommodityElastic) CreateIndex(ctx context.Context, indexName string) error {
	_, err := es.client.CreateIndex(indexName).BodyString(mapping).Do(ctx)
	if err != nil {
		return errno.Errorf(errno.InternalESErrorCode, err.Error())
	}
	return nil
}

func (es *CommodityElastic) AddItem(ctx context.Context, indexName string, spu *model.Spu) error {
	spuEs := &model.SpuES{
		SpuId:      spu.SpuId,
		Name:       spu.Name,
		CategoryId: spu.CategoryId,
		Price:      spu.Price,
		Shipping:   spu.Shipping > 0,
	}

	_, err := es.client.Index().Index(indexName).
		Id(fmt.Sprintf("%d", spu.SpuId)).
		BodyJson(spuEs).Do(ctx)
	if err != nil {
		return errno.Errorf(errno.InternalESErrorCode, err.Error())
	}

	return nil
}

func (es *CommodityElastic) RemoveItem(ctx context.Context, indexName string, id int64) error {
	_, err := es.client.Delete().Index(indexName).Id(fmt.Sprintf("%d", id)).Do(ctx)
	if err != nil {
		return errno.Errorf(errno.InternalESErrorCode, err.Error())
	}

	return nil
}

func structToMapUsingJSON(obj interface{}) map[string]interface{} {
	data, _ := sonic.Marshal(obj)
	var result map[string]interface{}
	_ = sonic.Unmarshal(data, &result)
	return result
}

func (es *CommodityElastic) UpdateItem(ctx context.Context, indexName string, spu *model.Spu) error {
	spuEs := &model.SpuES{
		SpuId:      spu.SpuId,
		Name:       spu.Name,
		CategoryId: spu.CategoryId,
		Price:      spu.Price,
		Shipping:   spu.Shipping > 0,
	}
	_, err := es.client.Update().Index(indexName).
		Id(fmt.Sprintf("%d", spu.SpuId)).Doc(structToMapUsingJSON(spuEs)).
		Do(ctx)
	if err != nil {
		return errno.Errorf(errno.InternalESErrorCode, err.Error())
	}

	return nil
}

func (es *CommodityElastic) SearchItems(ctx context.Context, indexName string, query *commodity.ViewSpuReq) ([]int64, int64, error) {
	q := es.BuildQuery(query)
	pageSize := int(query.GetPageSize())
	pageNum := int(query.GetPageNum())

	result, err := es.client.Search().Index(indexName).
		Query(q).
		From(pageNum * pageSize).Size(pageSize).
		Do(ctx)
	if err != nil {
		return nil, 0, errno.Errorf(errno.InternalESErrorCode, err.Error())
	}

	rets := make([]int64, 0)
	for _, hit := range result.Hits.Hits {
		var spuEs model.SpuES
		data, err := hit.Source.MarshalJSON()
		if err != nil {
			return nil, 0, errno.Errorf(errno.InternalServiceErrorCode, err.Error())
		}
		err = sonic.Unmarshal(data, &spuEs)
		if err != nil {
			return nil, 0, errno.Errorf(errno.InternalServiceErrorCode, err.Error())
		}
		spu := &model.Spu{
			SpuId: spuEs.SpuId,
		}
		rets = append(rets, spu.SpuId)
	}
	return rets, result.TotalHits(), nil
}

func (es *CommodityElastic) BuildQuery(req *commodity.ViewSpuReq) *elastic.BoolQuery {
	query := elastic.NewBoolQuery()
	hasCondition := false
	// 处理关键词
	if req.KeyWord != nil && req.GetKeyWord() != "" {
		query = query.Must(elastic.NewMatchQuery("Name", req.GetKeyWord()))
		hasCondition = true
	}

	//处理分类 ID
	if req.CategoryID != nil && req.GetCategoryID() != 0 {
		query = query.Must(elastic.NewMatchQuery("CateGoryId", req.GetCategoryID()))
		hasCondition = true
	}

	// 处理 Spu ID
	if req.SpuID != nil && req.GetSpuID() != 0 {
		query = query.Must(elastic.NewMatchQuery("SpuId", req.GetSpuID()))
		hasCondition = true
	}
	if req.MinCost != nil || req.MaxCost != nil {
		// 价格范围查询
		rangeQuery := elastic.NewRangeQuery("Price")
		minCost := constants.CommodityDefaultMinCost
		maxCost := constants.CommodityDefaultMaxCost

		if req.MinCost != nil {
			minCost = req.GetMinCost()
			hasCondition = true
		}
		if req.MaxCost != nil {
			maxCost = req.GetMaxCost()
			hasCondition = true
		}

		// 确保 minCost <= maxCost，避免查询无效
		if minCost > maxCost {
			minCost, maxCost = maxCost, minCost
		}
		rangeQuery.Gte(minCost).Lte(maxCost)
		query = query.Must(rangeQuery)
	}

	if req.IsShipping != nil {
		query = query.Must(elastic.NewMatchQuery("Shipping", req.IsShipping))
	}

	if !hasCondition {
		query = query.Must(elastic.NewMatchAllQuery())
	}

	return query
}
