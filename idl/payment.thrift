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

/*
 * struct PaymentResponse 订单支付响应
 * @Param paymentID 支付ID
 * @Param status 请求状态
 */
struct PaymentResponse {
    1: model.BaseResp base,
    2: required i64 paymentID
    3: required i64 status
}

/*
 * struct RefundTokenRequest 退款令牌请求
 * @Param orderID 商户订单号
 * @Param userID 管理员ID
 */
struct RefundTokenRequest {
    1: required i64 orderID
    2: required i64 userID
}

/*
 * struct RefundTokenResponse 退款令牌响应
 * @Param refundToken 退款令牌
 * @Param expirationTime 令牌过期时间
 * @Param status 请求状态
 */
struct RefundTokenResponse {
    1: model.BaseResp base,
    2: required string refundToken,
    3: required i64 expirationTime
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

/*
 * struct RefundResponse 退款响应
 * @Param refundID 退款ID
 * @Param status 请求状态
 */
struct RefundResponse {
    1: model.BaseResp base,
    2: required i64 refundID
    3: required i64 status
}


/*
 * service PaymentService 支付服务
 * @Method RequestPaymentToken 请求支付令牌
 * @Method ProcessPayment 处理支付
 * @Method RequestRefundToken 请求退款令牌
 * @Method ProcessRefund 处理退款
 */
service PaymentService {
    PaymentResponse ProcessPayment(1: PaymentRequest request) (api.post="/api/payment/process")
    PaymentTokenResponse RequestPaymentToken(1: PaymentTokenRequest request) (api.get="/api/payment/token")
    RefundResponse ProcessRefund(1: RefundRequest request) (api.post="/api/payment/refund")
    RefundTokenResponse RequestRefundToken(1: RefundTokenRequest request) (api.get="/api/payment/refund-token")
}
