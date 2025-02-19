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

struct UserCoupon {
    1: required i64 couponID;
    2: required i64 creatorID;
    3: required i64 userID;
    4: required string name;
    5: required i32 typeInfo;
    6: required i64 conditionCost;
    7: optional double discountAmount;
    8: optional double discount;
    9: required i64 rangeType;
    10: required i64 rangeId;
    11: required i64 expireAt;
    12: required string description;
    13: required i64 remainUserUseCount;
    14: required i64 createdAt;
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

struct PriceHistory {
    1: required i64 historyID;
    2: required i64 skuID;
    3: required i64 price;
    4: required i64 createdAt;
    5: optional i64 prevVersion;
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
    1: required i64 merchantID; // 商家 ID
    2: required i64 goodsID; // 商品 ID
    3: required string goodsName; // 商品名字
    4: required string goodsHeadDrawing; // 商品头图链接
    5: required i64 styleID; // 商品款式 ID
    6: required string styleName; // 款式名称
    7: required string styleHeadDrawing; // 款式头图
    8: required double originCast; // 原价
    9: required double saleCast; // 售卖价
    10: required i64 purchaseQuantity; // 购买数量
    11: required double paymentAmount; // 支付金额
    12: required double freightAmount; // 运费金额
    13: required double settlementAmount; // 结算金额
    14: required double discountAmount; // 优惠金额
    15: required double singleCast; // 下单单价
    16: i64 couponID // 优惠券 ID
    17: required i64 orderID;
    18: required i32 goodsVersion;
}

struct BaseOrderGoods {
    1: required string merchantName; // 商家 ID
    2: required i64 goodsName; // 商品 ID
    3: required i64 styleName; // 商品款式 ID
    4: required i64 purchaseQuantity; // 购买数量
    5: required string styleHeadDrawing; // 款式头图链接
    6: i64 couponID // 优惠券 ID
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
