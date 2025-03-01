namespace go api.commodity
include "../model.thrift"

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
    1: required i64 couponID;
}

struct DeleteCouponReq {
    1: required i64 couponID;
}

struct DeleteCouponResp {

}

struct CreateUserCouponReq {
    1: required i64 couponID;
    2: required i64 remaining_use,
}

struct CreateUserCouponResp {

}

struct ViewCouponReq {
    1: required i64 pageNum;
}

struct ViewCouponResp {
    1: required list<model.Coupon> couponInfo;
}

struct ViewUserAllCouponReq {
    1: required i64 isIncludeExpired;
    3: required i64 pageNum;
}

struct ViewUserAllCouponResp {
    1: required list<model.Coupon> coupons;
}

struct UseUserCouponReq {
    1: required i64 couponID;
}

struct UseUserCouponResp {

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
    1: required i64 spuID;
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
    1: required list<model.Spu> spus;
    2: required i64 total;
}

struct DeleteSpuReq {
    1: required i64 spuID;
}

struct DeleteSpuResp {
}

struct CreateSpuImageReq {
    1: required i64 spuID;
}

struct CreateSpuImageResp {
    1: required i64 imageID;
}

struct UpdateSpuImageReq {
    1: required i64 imageID;
}

struct UpdateSpuImageResp {
}

struct ViewSpuImageReq {
    1: required i64 spuID;
    2: optional i64 pageNum;
    3: optional i64 pageSize;
}

struct ViewSpuImageResp {
    1: required list<model.SpuImage> images;
    2: required i64 total;
}

struct DeleteSpuImageReq {
    1: required i64 spuImageID;
}

struct DeleteSpuImageResp {
}

struct CreateSkuReq {
    1: required string name;
    2: required i64 stock;
    3: required string description;
    5: required double price;
    6: required i32 forSale;
    7: required double shipping;
    8: required i64 spuID;

}

struct CreateSkuResp {
    1: required i64 skuID;
}

struct CreateSkuImageReq {
    1: required i64 skuID;
}

struct CreateSkuImageResp {
    1: required i64 imageID;
}

struct UpdateSkuReq {
    1: required i64 skuID;
    2: optional double shipping;
    3: optional string description;
    5: optional double price;
    6: optional i32 forSale;
    7: optional i64 Stock;

}

struct UpdateSkuResp {

}

struct UpdateSkuImageReq {
    1: required i64 imageID;
}

struct UpdateSkuImageResp {
}


struct DeleteSkuReq {
    1: required i64 skuID;
}

struct DeleteSkuResp {

}

struct DeleteSkuImageReq {
    1: required i64 skuImageID;
}

struct DeleteSkuImageResp {
}

struct ViewSkuImageReq {
    1: required i64 skuID;
    2: optional i64 pageNum;
    3: optional i64 pageSize;
}

struct ViewSkuImageResp {
    1: required list<model.SkuImage> images;
}

struct ViewSkuReq {
    1: optional i64 skuID;
    2: optional i64 spuID;
    3: optional i64 pageNum;
    4: optional i64 pageSize;
}

struct ViewSkuResp {
    1: required list<model.Sku> skus;
}

struct UploadSkuAttrReq {
    1: optional i64 skuID;
    2: required string saleAttr;
    3: required string saleValue;
}

struct UploadSkuAttrResp {

}

struct CreateCategoryReq {
    1: required string name;
}

struct CreateCategoryResp {
    1: required i64 categoryID;
}

struct DeleteCategoryReq {
    1: required i64 categoryID;
}

struct DeleteCategoryResp {
}

struct ViewCategoryReq {
    1: required i64 pageNum;
    2: required i64 pageSize;
}

struct ViewCategoryResp {
    1: list<model.CategoryInfo> categoryInfo;
}

 struct UpdateCategoryReq {
    1: required i64 categoryID;
    2: required string name;
 }

 struct UpdateCategoryResp {

 }


struct ViewHistoryPriceReq {
    1: required i64 historyID;
    2: required i64 skuID;
    3: required i64 pageSize;
    4: required i64 pageNum;
}

struct ViewHistoryPriceResp {
    1: required list<model.PriceHistory> records;
}


service CommodityService {
    // 优惠券
    CreateCouponResp CreateCoupon(1: CreateCouponReq req) (api.post="/api/v1/commodity/coupon/create");
    DeleteCouponResp DeleteCoupon(1: DeleteCouponReq req) (api.delete="/api/v1/commodity/coupon/delete");
    CreateUserCouponResp CreateUserCoupon(1: CreateUserCouponReq req) (api.post="/api/v1/commodity/coupon/receive");
    ViewCouponResp ViewCoupon(1: ViewCouponReq req) (api.get="/api/v1/commodity/coupon/search");
    ViewUserAllCouponResp ViewUserAllCoupon(1: ViewUserAllCouponReq req) (api.get="/api/v1/commodity/coupon/all");

    // SPU
    CreateSpuResp CreateSpu(1: CreateSpuReq req) (api.post="/api/v1/commodity/spu/create");
    UpdateSpuResp UpdateSpu(1: UpdateSpuReq req) (api.post="/api/v1/commodity/spu/update");
    ViewSpuResp ViewSpu(1: ViewSpuReq req) (api.get="/api/v1/commodity/spu/search");
    DeleteSpuResp DeleteSpu(1: DeleteSpuReq req) (api.delete="/api/v1/commodity/spu/delete");
    ViewSpuImageResp ViewSpuImage(1: ViewSpuImageReq req) (api.get="/api/v1/commodity/spu/image/search");
    CreateSpuImageResp CreateSpuImage(1: CreateSpuImageReq req) (api.post = "/api/v1/commodity/spu/image/create");
    UpdateSpuImageResp UpdateSpuImage(1: UpdateSpuImageReq req) (api.post = "/api/v1/commodity/spu/image/update");
    DeleteSpuImageResp DeleteSpuImage(1: DeleteSpuImageReq req) (api.delete="/api/v1/commodity/spu/image/delete");


    //SKU
    CreateSkuResp CreateSku(1: CreateSkuReq req) (api.post="/api/v1/commodity/sku/create");
    UpdateSkuResp UpdateSku(1: UpdateSkuReq req) (api.post="/api/v1/commodity/sku/upadte");
    DeleteSkuResp DeleteSku(1: DeleteSkuReq req) (api.delete="/api/v1/commodity/sku/delete");
    ViewSkuImageResp ViewSkuImage(1: ViewSkuImageReq req) (api.get="/api/v1/commodity/sku/image");
    ViewSkuResp ViewSku(1: ViewSkuReq req) (api.get="/api/v1/commodity/sku/search");
    UploadSkuAttrResp UploadSkuAttr(1: UploadSkuAttrReq req) (api.post="/api/v1/commodity/sku/attr");
    CreateSkuImageResp CreateSkuImage(1: CreateSkuImageReq req) (api.post = "/api/v1/commodity/sku/image/create");
    UpdateSkuImageResp UpdateSkuImage(1: UpdateSkuImageReq req) (api.post = "/api/v1/commodity/sku/image/update");
    DeleteSkuImageResp DeleteSkuImage(1: DeleteSkuImageReq req) (api.delete="/api/v1/commodity/sku/image/delete");
    ViewHistoryPriceResp ViewHistory(1: ViewHistoryPriceReq req) (api.get="/api/v1/commodity/price/history")

    //category
    CreateCategoryResp CreateCategory(1: CreateCategoryReq req) (api.post = "/api/v1/commodity/category/create");
    DeleteCategoryResp DeleteCategory(1: DeleteCategoryReq req) (api.delete="/api/v1/commodity/category/delete");
    ViewCategoryResp ViewCategory(1: ViewCategoryReq req) (api.get="/api/v1/commodity/category/search");
    UpdateCategoryResp UpdateCategory(1: UpdateCategoryReq req) (api.post="/api/v1/commodity/category/update");
}
