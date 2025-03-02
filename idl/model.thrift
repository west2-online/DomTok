namespace go model

struct BaseResp {
    1: i64 code,
    2: string msg,
}

struct UserInfo {
    1: i64 userId,
    2: string name,
}

struct LoginData {
    1: i64 userId,
}

struct CategoryInfo {
    1: required i64 categoryID;
    2: required string name;
}

struct Coupon {
    1: required i64 couponID;
    2: required i64 creatorID;
    3: required i64 deadlineForGet;
    4: required string name;
    5: required i32 typeInfo;
    6: required double conditionCost;
    7: optional double discountAmount;
    8: optional double discount;
    9: required i32 rangeType;
    10: required i64 rangeId;
    11: required i64 expireTime;
    12: required string description;
    13: required i64 createdAt;
    14: optional i64 updatedAt;
    15: optional i64 deletedAt;
}

struct AssignedCouponSpuInfo{
    1: required i64 spuId,
    2: required Coupon coupon,
    3: required double discount_price,
}

struct UserCoupon {
    1: required Coupon coupon,
    13: required i64 remainUserUseCount;
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
    6: required i64 updatedAt;
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
    3: required i64 creatorID;
    4: required string description;
    5: required i64 categoryID;
    6: required string goodsHeadDrawing;
    7: required double price;
    8: required i32 forSale;
    9: required double shipping;
    10: required i64 createdAt;
    11: required i64 updatedAt;
    12: optional i64 deletedAt;
}

struct Sku {
    1: required i64 skuID;
    2: required i64 creatorID;
    3: required double price;
    4: required string name;
    5: required string description;
    6: required i32 forSale;
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

struct SkuInfo {
    1: required i64 skuID;
    2: required i64 creatorID;
    3: required double price;
    4: required string name;
    5: required i32 forSale;
    6: required i64 lockStock;
    7: required string styleHeadDrawing;
    8: required i64 spuID;
    9: required i64 historyID;
}

struct SkuVersion {
    1: required i64 skuID;
    2: required i64 versionID;
}

struct PriceHistory {
    1: required i64 historyID;
    2: required i64 skuID;
    3: required i64 price;
    4: required i64 createdAt;
    5: optional i64 prevVersion;
}

/*
* struct SkuBuyInfo 实际扣除商品
* @Param skuID skuID
* @Param count 购买商品数
 */

struct SkuBuyInfo {
    1: required i64 skuID;
    2: required i64 count;
}

struct Order {
    1: required i64 id;
    2: required string status;
    3: required i64 uid;
    4: required double totalAmountOfGoods;
    5: required double totalAmountOfFreight;
    6: required double totalAmountOfDiscount;
    7: required double paymentAmount;
    8: required string paymentStatus;
    9: required i64 paymentAt;
    10: required string paymentStyle;
    11: required i64 orderedAt;
    12: required i64 deletedAt;
    13: required i64 deliveryAt
    14: required i64 addressID;
    15: required string addressInfo;
    16: i64 couponId; // 优惠券 ID
    17: string couponName; // 优惠券名称
}

struct BaseOrder {
    1: required i64 id;
    2: required string status;
    3: required double totalAmountOfGoods;
    4: required double paymentAmount;
    5: required string paymentStatus;
}

struct orderWithGoods {
    1: required Order order;
    2: required list<OrderGoods> goods;
}

struct baseOrderWithGoods {
    1: required BaseOrder order;
    2: required list<BaseOrderGoods> goods;
}

struct OrderGoods {
    1: required i64 merchantId; // 商家 ID
    2: required i64 goodsId; // 商品 ID
    3: required string goodsName; // 商品名字
    4: required i64 styleId; // 商品款式 ID
    5: required string styleName; // 款式名称
    6: required i64 goodsVersion; // 商品版本号
    7: required string styleHeadDrawing; // 款式头图
    8: required double originPrice; // 原价
    9: required double salePrice; // 售卖价
    10: required double singleFreightPrice;
    11: required i64 purchaseQuantity; // 购买数量
    12: required double totalAmount; // 本应总金额
    13: required double freightAmount; // 运费总金额
    14: required double discountAmount; // 折扣总金额
    15: required double paymentAmount; // 支付总金额
    16: required double singlePrice; // 最终单间
    17: i64 couponId;
    18: string couponName;
    19: required i64 orderId;
}

struct BaseOrderGoods {
    1: required i64 merchantID; // 商家 ID
    2: required i64 goodsID; // 商品 ID
    3: required i64 styleID; // 商品款式 ID
    4: required i64 purchaseQuantity; // 购买数量
    5: i64 couponID // 优惠券 ID
    6: required i64 goodsVersion; // 商品历史号
}

/*
 * struct CreditCardInfo 信用卡信息
 * @Param maskedCreditCardNumber 仅存储信用卡号掩码，如 **** **** **** 1234
 * @Param creditCardExpirationYear 信用卡到期年
 * @Param creditCardExpirationMonth 信用卡到期月
 * @Param creditCardCvv 信用卡
 */
struct CreditCardInfo {
    1: required string maskedCreditCardNumber
    2: required i64 creditCardExpirationYear
    3: required i64 creditCardExpirationMonth
    4: required i64 creditCardCvv
}

struct PaymentTokenInfo{
    1:required string paymentToken
    2:required i64 paymentTokenExpirationTime
}

struct RefundResponseInfo{
    1: required i64 refundID
    2: required i64 status
}

struct CartGoods {
    1: required i64 merchantId; // 商家 ID
    2: required i64 goodsId; // 商品 ID
    3: required string goodsName; // 商品名字
    4: required i64 skuId; // 商品款式 ID
    5: required string skuName; // 款式名称
    6: required i64 goodsVersion; // 商品版本号
    7: required string styleHeadDrawing; // 款式头图
    11: required i64 purchaseQuantity; // 购买数量
    12: required double totalAmount; // 本应总金额
    14: required double discountAmount; // 折扣总金额
}
