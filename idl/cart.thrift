namespace go cart

include "model.thrift"

namespace go cart

/* struct AddGoodsIntoCartRequest 将商品添加到购物车，暴露
* @Param skuId skuID
* @Param count 数量
*/
struct AddGoodsIntoCartRequest{
    1: required i64 skuId,
    3: required i64 count,
}

struct AddGoodsIntoCartResponse{
    1: required model.BaseResp base,
}

/* struct ShowCartGoodsListRequest 查看购物车内容，暴露
* @Param pageNum 页码(一页默认15个商品)
*/
struct ShowCartGoodsListRequest{
    1: required i64 pageNum,
}

struct ShowCartGoodsListResponse{
    1: required model.BaseResp base,
    2: required list<model.Sku> goodsList,
    3: required i64 goodsCount,
}

/* struct UpdateCartGoodsRequest 更新购物车商品，暴露
* @Param skuId skuID
* @Param count 数量
*/
struct UpdateCartGoodsRequest{
    1: required i64 skuId,
    3: required i64 count,
}

struct UpdateCartGoodsResponse{
    1: required model.BaseResp base,
}

/* struct PayCartGoodsRequest 调用支付，暴露
* @Param skuIdList skuID列表
*/
struct PayCartGoodsRequest{
    1: required list<i64> skuIdList,
}

struct PayCartGoodsResponse{
    1: required model.BaseResp base,
}
