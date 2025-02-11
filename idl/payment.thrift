namespace go payment

include "model.thrift"

/*
 * struct PaymentTokenRequest 支付令牌请求
 * @Param orderID 商户订单号
 * @Param userID 用户ID
 */
struct PaymentTokenRequest {
    1: required i64 orderID
    2: required i64 userID
}

/*
 * struct PaymentTokenResponse 支付令牌响应
 * @Param token 支付令牌
 * @Param expirationTime 令牌过期时间
 * @Param status 请求状态
 */
struct PaymentTokenResponse {
    1: model.BaseResp base,
    2: required string paymentToken
    3: required i64 expirationTime
}

/*
 * struct PaymentRequest 订单结算请求
 * @Param orderID 商户订单号
 * @Param userID 用户ID
 * @Param paymentToken 支付令牌
 * @Param creditCard 信用卡信息
 * @Param description 订单描述
 */
struct PaymentRequest {
    1: required i64 orderID
    2: required i64 userID
    3: required string paymentToken
    4: required model.CreditCardInfo creditCard
    5: optional string description
}

struct PaymentResponse {
    1: model.BaseResp base,
    2: required i64 paymentID
    3: required i64 status
}

/*
 * struct RefundRequest 退款请求
 * @Param orderID 关联的商户订单号
 * @Param userID 用户ID
 * @Param refundAmount 退款金额
 * @Param refundReason 退款原因
 */
struct RefundRequest {
    1: required i64 orderID
    2: required i64 userID
    3: required double refundAmount
    4: required string refundReason
}

struct RefundResponse {
    1: model.BaseResp base,
    2: required i64 refundID
    3: required i64 status
}

/*
 * struct PaymentOrder 支付订单
 * @Param id 支付订单唯一标识
 * @Param orderID 商户订单号
 * @Param userID 用户ID(用户的唯一标识)
 * @Param amount 订单总金额
 * @Param status 支付状态：0-待支付，1-处理中，2-成功支付 3-支付失败
 * @Param maskedCreditCardNumber 信用卡号国际信用卡号的最大长度为19 (仅存储掩码，如 **** **** **** 1234)
 * @Param creditCardExpirationYear 信用卡到期年
 * @Param creditCardExpirationMonth 信用卡到期月
 * @Param description 订单描述
 * @Param createdAt 订单创建时间
 * @Param updatedAt 订单更新时间
 * @Param deletedAt 订单删除时间
 */
struct PaymentOrder {
    1: required i64 id
    2: required i64 orderID
    3: required i64 userID
    4: required double amount
    5: required i32 status
    6: optional string maskedCreditCardNumber
    7: optional i32 creditCardExpirationYear
    8: optional i32 creditCardExpirationMonth
    9: optional string description
    10: optional string createdAt
    11: optional string updatedAt
    12: optional string deletedAt
}

/*
 * struct RefundOrder 退款订单
 * @Param id 支付退款唯一标识
 * @Param orderID 关联的商户订单号
 * @Param userID 用户ID(用户的唯一标识)
 * @Param refundAmount 退款金额,单位为元
 * @Param refundReason 退款原因
 * @Param status 退款状态：0-待处理，1-处理中，2-成功退款 3-退款失败
 * @Param maskedCreditCardNumber 信用卡号国际信用卡号的最大长度为19(仅存储掩码，如 **** **** **** 1234)
 * @Param creditCardExpirationYear 信用卡到期年
 * @Param creditCardExpirationMonth 信用卡到期月
 * @Param createdAt 退款申请时间
 * @Param updatedAt 退款最后更新时间
 * @Param deletedAt 退款记录删除时间
 */
struct RefundOrder {
    1: required i64 id
    2: required string orderID
    3: required i64 userID
    4: required double refundAmount
    5: optional string refundReason
    6: required i32 status
    7: optional string maskedCreditCardNumber
    8: optional i32 creditCardExpirationYear
    9: optional i32 creditCardExpirationMonth
    10: optional string createdAt
    11: optional string updatedAt
    12: optional string deletedAt
}

/*
 * struct PaymentLedger 支付流水
 * @Param id 流水ID
 * @Param referenceID 关联的支付订单或退款订单ID
 * @Param userID 用户ID
 * @Param amount 交易金额(正数表示收入，负数表示支出)
 * @Param transactionType 交易类型：1-支付，2-退款，3-手续费，4-调整'
 * @Param status 交易状态：0-待处理，1-成功，2-失败'
 * @Param createdAt 交易创建时间
 * @Param updatedAt 交易更新时间
 * @Param deletedAt 交易记录删除时间
 */
struct PaymentLedger {
    1: required i64 id
    2: required i64 referenceID
    3: required i64 userID
    4: required double amount
    5: required i32 transactionType
    6: required i32 status
    7: optional string createdAt
    8: optional string updatedAt
    9: optional string deletedAt
}

/*
 * service PaymentService 支付服务
 */
service PaymentService {
    PaymentResponse ProcessPayment(1: PaymentRequest request)
    RefundResponse ProcessRefund(1: RefundRequest request)
}
