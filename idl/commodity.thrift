namespace go commodity

include "model.thrift"

/*
* struct CreateCouponReq 创建优惠券信息
* @Param UserID 创建者ID
* @Param Name 名称
* @Param TypeInfo 类型
* @Param ConditionCost 用券门槛
* @Param DiscountAmount 满减金额
* @Param Discount 折扣
* @Param RangeType 对应范围
* @Param RangeID 对应范围的类型ID
* @Param Description 描述
* @Param ExpireTime 有效期
* @Param deadlineForGet 可领取优惠券的截止时间
*/
struct CreateCouponReq {
    1: required i64 userID;
    2: required string name;
    3: required i32 typeInfo;
    4: optional double conditionCost;
    5: optional double discountAmount;
    6: optional double discount;
    7: required i32 rangeType;
    8: required i64 rangeID;
    9: optional string description;
    10: required i64 expireTime;
    11: required i64 deadlineForGet;
}

struct CreateCouponResp {
    1: required model.BaseResp base;
    2: required i64 couponID;
}

/*
* struct DeleteCouponReq 删除优惠券信息
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

/*
* struct CreateUserCouponReq 用户领取优惠券
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

/*
* struct ViewCouponReq 查看优惠券信息或商家创建的优惠券
* @Param CouponID 优惠券信息ID
* @Param UserID 用户ID
* @Param PageNum 页数
* @Param PageSize 页面大小
*/
struct ViewCouponReq {
    1: required i64 couponID;
    2: optional i64 creatorID;
    3: optional i64 pageNum;
    4: optional i64 pageSize;
}

struct ViewCouponResp {
    1: required model.BaseResp base;
    2: required model.Coupon couponInfo;
}

/*
* struct ViewUserAllCouponReq 查看用户自己持有的优惠券
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
    2: required list<model.UserCoupon> coupons;
}

/*
* struct UseUserCouponReq 使用用户优惠券
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

/*
* struct CreateSpuReq 创建Spu请求
* @Param spuImages spu具体介绍图
* @Param name 名称
* @Param userID 创建者ID
* @Param description 描述
* @Param categoryID 类型ID
* @Param goodsHeadDrawing 款式头图
* @Param price 价格
* @Param forSale 是否出售
* @Param shipping 运费
*/
struct CreateSpuReq {
    1: optional list<binary> spuImages;
    2: required string name;
    3: required i64 userID;
    4: required string description;
    5: required i64 categoryID;
    6: required string goodsHeadDrawing;
    7: required double price;
    8: required i32 forSale;
    9: required double shipping;
}

struct CreateSpuResp {
    1: required model.BaseResp base;
    2: required i64 spuID;
}

/*
* struct UpdateSpuReq 创建spu请求
* @Param spuImages spu具体介绍图
* @Param name 名称
* @Param userID 创建者ID
* @Param description 描述
* @Param categoryID 类型ID
* @Param goodsHeadDrawing 款式头图
* @Param price 价钱
* @Param forSale 是否出售
* @Param shipping 运费
* @Param spuID spu的ID
*/
struct UpdateSpuReq {
    1: optional list<binary> spuImages;
    2: optional string name;
    3: required i64 userID;
    4: optional string description;
    5: optional i64 categoryID;
    6: optional string goodsHeadDrawing;
    7: optional double price;
    8: optional i32 forSale;
    9: optional double shipping;
    10: required i64 spuID;
}

struct UpdateSpuResp {
    1: required model.BaseResp base;
}

/* struct ViewSpuReq 查询商品(提供关键词，品牌，预算，是否免运费)
* @Param KeyWord 关键词
* @Param CategoryID 类型
* @Param minCost 最小花费
* @Param maxCost 最大花费
* @Param IsShipping 是否免运费
* @Param SpuID Spu对应ID
*/
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

/*
* struct DeleteSpuReq 删除spu请求
* @Param userID 创建者ID
* @Param spuID spuID
 */
struct DeleteSpuReq {
    1: required i64 userID;
    2: required i64 spuID;
}

struct DeleteSpuResp {
    1: required model.BaseResp base;
}

/*
* struct ViewSpuImageReq 查看spu具体介绍图片
* @Param spuID spuID
* @Param pageNum 页数
* @Param pageSize 页尺寸
 */
struct ViewSpuImageReq {
    1: required i64 spuID;
    2: optional i64 pageNum;
    3: optional i64 pageSize;
}

struct ViewSpuImageResp {
    1: required model.BaseResp base;
    2: required list<model.SpuImage> images;
}

/*
* struct CreateSkuReq 创建Sku请求
* @Param skuImages sku图片
* @Param name 名称
* @Param userID 创建者ID
* @Param description 描述
* @Param styleHeadDrawing 款式头图
* @Param price 价钱
* @Param forSale 是否出售
* @Param shipping 运费
* @Param spuID spuID
* @Param stock 库存
*/
struct CreateSkuReq {
    1: optional list<binary> skuImages;
    2: required string name;
    3: required i64 creatorID;
    4: required string description;
    5: required string styleHeadDrawing;
    6: required double price;
    7: required i32 forSale;
    8: required double shipping;
    9: required i64 spuID;
    10: required i64 stock;
}

struct CreateSkuResp {
    1: required model.BaseResp base;
    2: required i64 skuID;
}

/*
* struct UpdateSkuReq 更新sku请求
* @Param skuID skuID
* @Param creatorID 创建者ID
* @Param skuImages sku图片
* @Param description 描述
* @Param styleHeadDrawing 款式头图
* @Param price 价钱
* @Param forSale 是否出售
* @Param Stock 库存量
* @Param shipping 运费
*/
struct UpdateSkuReq {
    1: required i64 skuID;
    2: required i64 creatorID;
    3: optional list<binary> skuImages;
    4: optional string description;
    5: optional string styleHeadDrawing;
    6: optional double price;
    7: optional i32 forSale;
    8: optional i64 Stock;
    9: optional double shipping;
}

struct UpdateSkuResp {
    1: required model.BaseResp base;
}

/*
* struct DeleteSkuReq 删除sku请求
* @Param skuID skuID
* @Param userID 创建者ID
 */
struct DeleteSkuReq {
    1: required i64 skuID;
    2: required i64 userID;
}

struct DeleteSkuResp {
    1: required model.BaseResp base;
}

/*
* struct ViewSkuImageReq 查看sku展示图片请求
* @Param skuID skuID
* @Param pageNum 页数
* @Param pageSize 页尺寸
 */
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

/*
* struct UploadSkuAttrReq 上传属性请求
* @Param skuID skuID
* @Param saleAttr 属性名
* @Param saleValue 属性值
 */
struct UploadSkuAttrReq {
    1: optional i64 skuID;
    2: required string saleAttr;
    3: required string saleValue;
}

struct UploadSkuAttrResp {
    1: required model.BaseResp base;
}

/*
* struct CreateCategoryReq 创建种类
* @Param name 名称
* @Param creatorID 创建者ID
 */

struct CreateCategoryReq {
    1: required string name;
    2: required i64 creatorID;
}

struct CreateCategoryResp {
    1: required model.BaseResp base;
    2: required i64 categoryID;
}

/*
* struct DeleteCategoryReq 删除种类请求
* @Param userID 用户ID
* @Param categoryID 种类ID
 */
struct DeleteCategoryReq {
    1: required i64 userID;
    2: required i64 categoryID;
}

struct DeleteCategoryResp {
    1: required model.BaseResp base;
}

/*
* struct ViewCategoryReq 查看类型请求
* @Param pageNum 页数
* @Param pageSize 页尺寸
 */
struct ViewCategoryReq {
    1: required i64 pageNum;
    2: required i64 pageSize;
}

struct ViewCategoryResp {
    1: required model.BaseResp base;
    2: list<model.CategoryInfo> categoryInfo;
}

/*
* struct UpdateCategoryReq 更新种类
* @Param categoryID 种类ID
* @Param name 名称
 */
 struct UpdateCategoryReq {
    1: required i64 categoryID;
    2: required string name;
 }

 struct UpdateCategoryResp {
    1: required model.BaseResp base;
 }


/*
* struct ListSkuInfoReq 列出sku信息
* @Param skuIDs skuID列表
* @Param pageNum 页数
* @Param pageSize 页尺寸
 */
struct ListSkuInfoReq {
    1: required list<i64> skuIDs;
    2: required i64 pageNum;
    3: required i64 pageSize;
}

struct ListSkuInfoResp {
    1: required model.BaseResp base;
    2: required list<model.SkuInfo> skuInfos;
}

/*
* struct DescSkuLockStockReq 预扣商品
* @Param skuID skuID
* @Param count 购买商品数
 */
struct DescSkuLockStockReq {
    1: required i64 skuID;
    2: required i64 count
}

struct DescSkuLockStockResp {
    1: required model.BaseResp base;
}

/*
* struct IncrSkuLockStockReq 回滚商品数
* @Param skuID skuID
* @Param count 原购买商品数
 */
struct IncrSkuLockStockReq {
    1: required i64 skuID;
    2: required i64 count;
}

struct IncrSkuLockStockResp {
    1: required model.BaseResp base;
}

/*
* struct DescSkuStockReq 实际扣除商品
* @Param skuID skuID
* @Param count 购买商品数
 */
struct DescSkuStockReq {
    1: required i64 skuID;
    2: required i64 count;
}

struct DescSkuStockResp {
    1: required model.BaseResp base;
}


service CommodityService {
    // 优惠券
    CreateCouponResp CreateCoupon(1: CreateCouponReq req);
    DeleteCouponResp DeleteCoupon(1: DeleteCouponReq req);
    CreateUserCouponResp CreateUserCoupon(1: CreateCouponReq req);
    ViewCouponResp ViewCoupon(1: ViewCouponReq req);
    ViewUserAllCouponResp ViewUserAllCoupon(1: ViewCouponReq req);
    UseUserCouponResp UseUserCoupon(1: UseUserCouponReq req);

    // SPU
    CreateSpuResp CreateSpu(1: CreateSpuReq req);
    UpdateSpuResp UpdateSpu(1: UpdateSkuReq req);
    ViewSpuResp ViewSpu(1: ViewSpuReq req);
    DeleteSpuResp DeleteSpu(1: DeleteSpuReq req);
    ViewSpuImageResp ViewSpuImage(1: ViewSpuImageReq req);

    //SKU
    CreateSkuResp CreateSku(1: CreateSkuReq req);
    UpdateSkuResp UpdateSku(1: UpdateSkuReq req);
    DeleteSkuResp DeleteSku(1: DeleteSkuReq req);
    ViewSkuImageResp ViewSkuImage(1: ViewSkuImageReq req);
    ViewSkuResp ViewSku(1: ViewSkuReq req);
    UploadSkuAttrResp UploadSkuAttr(1: UploadSkuAttrReq req);
    ListSkuInfoResp ListSkuInfo(1: ListSkuInfoReq req);

    //供订单服务调用
    DescSkuLockStockResp DescSkuLockStock(1: DescSkuLockStockReq req);
    IncrSkuLockStockResp IncrSkuLockStock(1: IncrSkuLockStockReq req);
    DescSkuStockResp DescSkuStock(1: DescSkuStockReq req);

    //category
    CreateCategoryResp CreateCategory(1: CreateCategoryReq req);
    DeleteCategoryResp DeleteCategory(1: DeleteCategoryReq req);
    ViewCategoryResp ViewCategory(1: ViewCategoryReq req);
    UpdateCategoryResp UpdateCategory(1: UpdateCategoryReq req);
}
