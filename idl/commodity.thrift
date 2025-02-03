namespace go commodity

include "model.thrift"

/* struct CreateCouponReq 创建优惠券信息
* @Param ActivityID 活动ID
* @Param Name 名称
* @Param TypeInfo 类型
* @Param ConditionCost 用券门槛
* @Param DiscountAmount 满减金额
* @Param Discount 折扣
* @Param RangeType 对应范围
* @Param RangeID 对应范围的类型ID
* @Param Description 描述
* @Param ExpireTime 过期历时
* @Param UserID 创建者ID
*/
struct CreateCouponReq {
    1: required i64 activityID;
    2: required string name;
    3: required string typeInfo;
    4: optional i64 conditionCost;
    5: optional i64 discountAmount;
    6: optional double discount;
    7: required string rangeType;
    8: required i64 rangeID;
    9: optional string description;
    10: required string expireTime;
    11: required i64 userID;
}

struct CreateCouponResp {
    1: required model.BaseResp base;
    2: required i64 couponID;
}

/* struct DeleteCouponReq 删除优惠券信息
* @Param CouponID 优惠券信息ID
* @Param UserID 用户ID
*/
struct DeleteCouponReq {
    1: required i64 couponID;
    2: required i64 userID;
}

struct DeleteCouponResp {
    1: required model.BaseResp base;
}

/* struct CreateUserCouponReq 用户领取优惠券
* @Param UserID 用户ID
* @Param CouponID 优惠券信息ID
*/
struct CreateUserCouponReq {
    1: required i64 userID;
    2: required i64 couponID;
}

struct CreateUserCouponResp {
    1: required model.BaseResp base;
}

/* struct ViewCouponReq 查看优惠券信息或用户持有的该优惠券ID
* @Param CouponID 优惠券信息ID
* @Param UserID 用户ID
*/
struct ViewConponReq {
    1: required i64 couponID;
    2: optional i64 userID;
}

struct ViewConponResp {
    1: required model.BaseResp base;
    2: required model.Coupon couponInfo;
}

/* struct ViewUserAllCouponReq 查看用户自己持有的优惠券
* @Param UserID 用户ID
* @Param IsIncludeExpired 是否包含过期的券
* @Param PageNum 页数
* @Param PageSize 页面大小
*/
struct ViewUserAllCouponReq {
    1: required i64 userID;
    2: required i64 isIncludeExpired;
    3: required i64 pageNum;
    4: required i64 pageSize;
}

struct ViewUserAllCouponResp {
    1: required model.BaseResp base;
    2: required list<model.Coupon> coupons;
}

/* struct UseUserCouponReq 使用用户优惠券
* @Param UserID 用户
* @Param CouponID 优惠券ID
*/
struct UseUserCouponReq {
    1: required i64 userID;
    2: required i64 couponID;
}

struct UseUserCouponResp {
    1: required model.BaseResp base;
}



struct CreateSpuReq {
    1: optional list<binary> spuImages;
    2: required string name;
    3: required i64 userID;
    4: required string description;
    5: optional i64 brandID;
    6: required i64 categoryID;
    7: required string goodsHeadDrawing;
    8: required double price;
    9: required string forSale;
    10: required double shipping;
}

struct CreateSpuResp {
    1: required model.BaseResp base;
    2: required i64 spuID;
}

struct UpdateSpuReq {
    1: optional list<binary> spuImages;
    2: optional string name;
    3: required i64 userID;
    4: optional string description;
    5: optional i64 brandID;
    6: optional i64 categoryID;
    7: optional string goodsHeadDrawing;
    8: optional double price;
    9: optional string forSale;
    10: optional double shipping;
    11: required i64 spuID;
}

struct UpdateSpuResp {
    1: required model.BaseResp base;
}

/* struct ViewSpuReq 查询商品(提供关键词，品牌，预算，是否免运费)
* @Param KeyWord 关键词
* @Param CategoryID 类型
* @Param BrandID 匹配
* @Param Budget 预算
* @Param IsShipping 是否免运费
* @Param SpuID Spu对应ID
*/
struct ViewSpuReq {
    1: optional string keyWord;
    2: optional i64 categoryID;
    3: optional i64 brandID;
    4: optional double minCost;
    5: optional double maxCost;
    6: optional bool isShipping;
    7: optional i64 spuID;
    8: optional i64 pageNum;
    9: optional i64 pageSize;
}

struct ViewSpuResp {
    1: required model.BaseResp base;
    2: required list<model.Spu> spus;
}

struct DeleteSpuReq {
    1: required i64 userID;
    2: required i64 spuID;
}

struct DeleteSpuResp {
    1: required model.BaseResp base;
}

struct ViewSpuImageReq {
    1: required i64 spuID;
    2: optional i64 pageNum;
    3: optional i64 pageSize;
}

struct ViewSpuImageResp {
    1: required model.BaseResp base;
    2: required list<model.SpuImage> images;
}


struct CreateSkuReq {
    1: optional list<binary> spuImages;
    2: required string name;
    3: required i64 userID;
    4: required string description;
    5: required string styleHeadDrawing;
    6: required double price;
    7: required string forSale;
    8: required double shipping;
    9: required i64 spuID;
    10: required i64 stock;
}

struct CreateSkuResp {
    1: required model.BaseResp base;
    2: required i64 skuID;
}

struct UpdateSkuReq {
    1: required i64 skuID;
    2: required i64 userID;
    3: optional list<binary> skuImages;
    4: optional string description;
    5: optional string styleHeadDrawing;
    6: optional double price;
    7: optional string forSale;
    8: optional i64 stock;
    9: optional double shipping;
}

struct UpdateSkuResp {
    1: required model.BaseResp base;
}

struct DeleteSkuReq {
    1: required i64 skuID;
    2: required i64 userID;
}

struct DeleteSkuResp {
    1: required model.BaseResp base;
}

struct ViewSkuImageReq {
    1: required i64 skuID;
    2: optional i64 pageNum;
    3: optional i64 pageSize;
}

struct ViewSkuImageResp {
    1: required model.BaseResp base;
    2: required list<model.SkuImage> images;
}

/* struct ViewSkuReq 查看sku信息
* @Param skuID 指定查看的skuID
* @Param spuID 指定查询该SPU下的所有sku
*/
struct ViewSkuReq {
    1: optional i64 skuID;
    2: optional i64 spuID;
    3: optional i64 pageNum;
    4: optional i64 pageSize;
}

struct ViewSkuResp {
    1: required model.BaseResp base;
    2: required list<model.Sku> skus;
}

struct UploadSkuAttrReq {
    1: optional i64 skuId;
    2: required string saleAttr;
    3: required string saleValue;
}

struct UploadSkuAttrResp {
    1: required model.BaseResp base;
}

struct CreateBrandReq {
    1: required string name;
    2: required string logoUrl;
    3: required i64 userID;
}

struct CreateBandResp {
    1: required model.BaseResp base;
    2: required i64 brandID;
}

struct DeleteBrandReq {
    1: required i64 userID;
    2: required i64 brandID;
}

struct DeleteBrandResp {
    1: required model.BaseResp base;
}
