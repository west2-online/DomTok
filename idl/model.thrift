namespace go model

struct BaseResp {
    1: i64 code,
    2: string msg,
}

struct UserInfo {
    1: required i64 uid,
    2: required string name,
    3: required string email,
    4: required string phone,
}

struct Coupon {
    1: required i64 couponID;
    2: required i64 userID;
    3: required i64 activityID;
    4: required string name;
    5: required string typeInfo;
    6: required i64 conditionCost;
    7: optional double discountAmount;
    8: optional double discount;
    9: required i64 rangeType;
    10: required i64 rangeId;
    11: required i64 expireAt;
    12: required string description;
    13: required i64 createdAt;
    14: optional i64 updatedAt;
    15: optional i64 deletedAt;
    16: required i64 remainUserUseCount;
    17: required i64 expireTime
}

struct AttrValue {
    1: required string saleAttr;
    2: required string saleValue;
}

struct SpuImage {
    1: required i64 imageID;
    2: required i64 spuID;
    3: required string url;
    4: required i64 createdAt;
    5: optional i64 deletedAt;
}

struct SkuImage {
     1: required i64 imageID;
    2: required i64 skuID;
    3: required string url;
    4: required i64 createdAt;
    5: optional i64 deletedAt;
}

struct Spu {
    1: required i64 spuID;
    2: required string name;
    3: required i64 userID;
    4: required string description;
    5: optional i64 brandID;
    6: required i64 categoryID;
    7: required string goodsHeadDrawing;
    8: required double price;
    9: required string forSale;
    10: required double shipping;
    11: required i64 createdAt;
    12: required i64 updatedAt;
    13: optional i64 deletedAt;
}

struct Sku {
    1: required i64 skuID;
    2: required i64 userID;
    3: required double price;
    4: required string name;
    5: required string description;
    6: required string forSale;
    7: required i64 stock;
    8: required string styleHeadDrawing;
    9: required i64 createdAt;
    10: required i64 updatedAt;
    11: optional i64 deletedAt;
    12: required i64 spuID;
    13: optional list<AttrValue> saleAttr;
    14: required i64 historyID;
    15: required i64 lockStock;
}

struct OrderGoods {
    1: required i64 MerchantID; // 商家 ID
    2: required i64 GoodsID; // 商品 ID
    3: required string GoodsName; // 商品名字
    4: required string GoodsHeadDrawing; // 商品头图链接
    5: required i64 StyleID; // 商品款式 ID
    6: required string StyleName; // 款式名称
    7: required string StyleHeadDrawing; // 款式头图
    8: required double OriginCast; // 原价
    9: required double SaleCast; // 售卖价
    10: required double PurchaseQuantity; // 购买数量
    11: required double PaymentAmount; // 支付金额
    12: required double FreightAmount; // 运费金额
    13: required double SettlementAmount; // 结算金额
    14: required double DiscountAmount; // 优惠金额
    15: required double SingleCast; // 下单单价
    16: i64 CouponID // 优惠券 ID
}

struct BaseOrderGoods {
    1: required i64 MerchantID; // 商家 ID
    2: required i64 GoodsID; // 商品 ID
    3: required i64 StyleID; // 商品款式 ID
    4: required double PurchaseQuantity; // 购买数量
    5: i64 CouponID // 优惠券 ID
}
