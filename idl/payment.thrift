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
    4: required double refundAmount
    5: required string refundReason
}

struct RefundResponse {
    1: model.BaseResp base,
    2: required i64 refundID
    3: required i64 status
}


/*
 * service PaymentService 支付服务
 */
service PaymentService {
    PaymentResponse ProcessPayment(1: PaymentRequest request)
    RefundResponse ProcessRefund(1: RefundRequest request)
}
