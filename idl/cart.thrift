namespace go cart

include "model.thrift"

namespace go cart

/* struct AddGoodsIntoCartRequest 将商品添加到购物车，暴露
* @Param skuId skuID
* @Param count 数量
*/
struct AddGoodsIntoCartRequest{
    1: required i64 skuId,
    2: required i64 shop_id,
    3: required i64 version_id
    4: required i64 count,
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
    2: required list<model.CartGoods> goodsList,
    3: required i64 goodsCount,
}

/* struct UpdateCartGoodsRequest 更新购物车商品，暴露
* @Param skuId skuID
* @Param count 数量
*/
struct UpdateCartGoodsRequest{
    1: required i64 skuId,
    2: required i64 shop_id,
    3: required i64 count,
}

struct UpdateCartGoodsResponse{
    1: required model.BaseResp base,
}

/* struct PurChaseCartGoodsRequest 购买购物车商品，暴露
* @Param sku_id skuID
*/
struct PurChaseCartGoodsRequest{
    1: required list<model.CartGoods> cartGoods,
}

struct PurChaseCartGoodsResponse{
    1: required model.BaseResp base,
    2: required i64 order_id,
}
/* struct DeleteAllCartGoodsRequest 清空购物车，暴露
*/
struct DeleteAllCartGoodsRequest{
}

struct DeleteAllCartGoodsResponse{
    1: required model.BaseResp base,
}

service CartService {
    AddGoodsIntoCartResponse AddGoodsIntoCart(1: AddGoodsIntoCartRequest req),
    ShowCartGoodsListResponse ShowCartGoodsList(1: ShowCartGoodsListRequest req),
    UpdateCartGoodsResponse UpdateCartGoods(1: UpdateCartGoodsRequest req),
    PurChaseCartGoodsResponse PurChaseCartGoods(1: PurChaseCartGoodsRequest req)
    DeleteAllCartGoodsResponse DeleteAllCartGoods(1:DeleteAllCartGoodsRequest req),
}
