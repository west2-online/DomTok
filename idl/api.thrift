namespace go api
include "model.thrift"

//
struct RegisterRequest {
    1: required string name
    2: required string password
    3: required string email
}

struct RegisterResponse {
}


service UserService {
    RegisterResponse Register(1: RegisterRequest req)(api.get = "api/v1/user/register"),
}


struct CreateCouponReq {
    1: required i64 deadlineForGet;
    2: required string name;
    3: required i32 typeInfo;
    4: optional double conditionCost;
    5: optional double discountAmount;
    6: optional double discount;
    7: required i32 rangeType;
    8: required i64 rangeID;
    9: optional string description;
    10: required i64 expireTime;
}

struct CreateCouponResp {
    1: required model.BaseResp base;
    2: required i64 couponID;
}

struct DeleteCouponReq {
    1: required i64 couponID;
}

struct DeleteCouponResp {
    1: required model.BaseResp base;
}

struct CreateUserCouponReq {
    1: required i64 couponID;
}

struct CreateUserCouponResp {
    1: required model.BaseResp base;
}

struct ViewCouponReq {
    1: required i64 couponID;
    2: optional i64 pageNum;
    3: optional i64 pageSize;
}

struct ViewCouponResp {
    1: required model.BaseResp base;
    2: required model.Coupon couponInfo;
}

struct ViewUserAllCouponReq {
    1: required i64 isIncludeExpired;
    2: required i64 pageNum;
    3: required i64 pageSize;
}

struct ViewUserAllCouponResp {
    1: required model.BaseResp base;
    2: required list<model.UserCoupon> coupons;
}

struct UseUserCouponReq {
    1: required i64 couponID;
}

struct UseUserCouponResp {
    1: required model.BaseResp base;
}

// 文件表单直接formFile获取即可
struct CreateSpuReq {

    1: required string name;
    2: required string description;
    3: required i64 categoryID;
    4: required double price;
    5: required i32 forSale;
    6: required double shipping;
}

struct CreateSpuResp {
    1: required model.BaseResp base;
    2: required i64 spuID;
}

// 文件表单直接formFile获取即可
struct UpdateSpuReq {

    1: optional string name;
    2: required i64 spuID;
    3: optional string description;
    4: optional i64 categoryID;
    5: optional double price;
    6: optional i32 forSale;
    7: optional double shipping;
    8: optional i64 spuImageId;
}

struct UpdateSpuResp {
    1: required model.BaseResp base;
}

struct ViewSpuReq {
    1: optional string keyWord;
    2: optional i64 categoryID;
    3: optional i64 spuID;
    4: optional double minCost;
    5: optional double maxCost;
    6: optional bool isShipping;
    7: optional i64 pageNum;
    8: optional i64 pageSize;
}

struct ViewSpuResp {
    1: required model.BaseResp base;
    2: required list<model.Spu> spus;
}

struct DeleteSpuReq {
    1: required i64 spuID;
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
    1: optional list<binary> skuImages;
    2: required string name;
    3: required i64 stock;
    4: required string description;
    5: required string styleHeadDrawing;
    6: required double price;
    7: required i32 forSale;
    8: required double shipping;
    9: required i64 spuID;

}

struct CreateSkuResp {
    1: required model.BaseResp base;
    2: required i64 skuID;
}

struct UpdateSkuReq {
    1: required i64 skuID;
    2: optional double shipping;
    3: optional list<binary> skuImages;
    4: optional string description;
    5: optional string styleHeadDrawing;
    6: optional double price;
    7: optional i32 forSale;
    8: optional i64 Stock;

}

struct UpdateSkuResp {
    1: required model.BaseResp base;
}


struct DeleteSkuReq {
    1: required i64 skuID;
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
    1: optional i64 skuID;
    2: required string saleAttr;
    3: required string saleValue;
}

struct UploadSkuAttrResp {
    1: required model.BaseResp base;
}

struct CreateCategoryReq {
    1: required string name;
}

struct CreateCategoryResp {
    1: required model.BaseResp base;
    2: required i64 categoryID;
}

struct DeleteCategoryReq {
    1: required i64 categoryID;
}

struct DeleteCategoryResp {
    1: required model.BaseResp base;
}

struct ViewCategoryReq {
    1: required i64 pageNum;
    2: required i64 pageSize;
}

struct ViewCategoryResp {
    1: required model.BaseResp base;
    2: list<model.CategoryInfo> categoryInfo;
}

 struct UpdateCategoryReq {
    1: required i64 categoryID;
    2: required string name;
 }

 struct UpdateCategoryResp {
    1: required model.BaseResp base;
 }

struct ListSkuInfoReq {
    1: required list<i64> skuIDs;
    2: required i64 pageNum;
    3: required i64 pageSize;
}

struct ListSkuInfoResp {
    1: required model.BaseResp base;
    2: required list<model.SkuInfo> skuInfos;
}

struct ViewHistoryPriceReq {
    1: required i64 historyID;
    2: required i64 skuID;
    3: required i64 pageSize;
    4: required i64 pageNum;
}

struct ViewHistoryPriceResp {
    1: required model.BaseResp base;
    2: required list<model.PriceHistory> records;
}

 service CommodityService {
     // 优惠券
     CreateCouponResp CreateCoupon(1: CreateCouponReq req) (api.post="/api/commodity/coupon/create");
     DeleteCouponResp DeleteCoupon(1: DeleteCouponReq req) (api.delete="/api/commodity/coupon/delete");
     CreateUserCouponResp CreateUserCoupon(1: CreateUserCouponReq req) (api.post="/api/commodity/coupon/receive");
     ViewCouponResp ViewCoupon(1: ViewCouponReq req) (api.get="/api/commodity/coupon/search");
     ViewUserAllCouponResp ViewUserAllCoupon(1: ViewUserAllCouponReq req) (api.get="/api/commodity/coupon/all");
     UseUserCouponResp UseUserCoupon(1: UseUserCouponReq req) (api.post="/api/commodity/coupon/use");

     // SPU
     CreateSpuResp CreateSpu(1: CreateSpuReq req) (api.post="/api/commodity/spu/create");
     UpdateSpuResp UpdateSpu(1: UpdateSpuReq req) (api.post="/api/commodity/spu/update");
     ViewSpuResp ViewSpu(1: ViewSpuReq req) (api.get="/api/commodity/spu/search");
     DeleteSpuResp DeleteSpu(1: DeleteSpuReq req) (api.delete="/api/commodity/spu/delete");
     ViewSpuImageResp ViewSpuImage(1: ViewSpuImageReq req) (api.get="/api/commodity/spu/image");

     //SKU
     CreateSkuResp CreateSku(1: CreateSkuReq req) (api.post="/api/commodity/sku/create");
     UpdateSkuResp UpdateSku(1: UpdateSkuReq req) (api.post="/api/commodity/sku/upadte");
     DeleteSkuResp DeleteSku(1: DeleteSkuReq req) (api.delete="/api/commodity/sku/delete");
     ViewSkuImageResp ViewSkuImage(1: ViewSkuImageReq req) (api.get="/api/commodity/sku/image");
     ViewSkuResp ViewSku(1: ViewSkuReq req) (api.get="/api/commodity/sku/search");
     UploadSkuAttrResp UploadSkuAttr(1: UploadSkuAttrReq req) (api.post="/api/commodity/sku/attr");
     ListSkuInfoResp ListSkuInfo(1: ListSkuInfoReq req) (api.get="/api/commodity/sku/list");
     ViewHistoryPriceResp ViewHistory(1: ViewHistoryPriceReq req) (api.get="/api/commodity/price/history")

     //category
     CreateCategoryResp CreateCategory(1: CreateCategoryReq req) (api.post = "/api/commodity/category/create");
     DeleteCategoryResp DeleteCategory(1: DeleteCategoryReq req) (api.delete="/api/commodity/category/delete");
     ViewCategoryResp ViewCategory(1: ViewCategoryReq req) (api.get="/api/commodity/category/search");
     UpdateCategoryResp UpdateCategory(1: UpdateCategoryReq req) (api.post="/api/commodity/category/update");
 }
