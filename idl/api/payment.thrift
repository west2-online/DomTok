namespace go api.payment

include "../model.thrift"

/*
 * struct PaymentTokenRequest 支付令牌请求
 * @Param orderID 商户订单号
 * @Param userID 用户ID
 */
struct PaymentTokenRequest {
    1: required i64 orderID
}

/*
 * struct PaymentTokenResponse 支付令牌响应
 * @Param token 支付令牌
 * @Param expirationTime 令牌过期时间
 * @Param status 请求状态
 */
struct PaymentTokenResponse {
    1: model.BaseResp base,
    2: model.PaymentTokenInfo tokenInfo,
}
/*
 * struct PaymentRequest 订单结算请求
 * @Param orderID 商户订单号
 * @Param userID 用户ID
 * @Param creditCard 信用卡信息
 * @Param description 订单描述
 */
struct PaymentRequest {
    1: required i64 orderID
    2: required i64 userID
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
}

/*
 * struct RefundTokenResponse 退款令牌响应
 * @Param refundToken 退款令牌
 * @Param expirationTime 令牌过期时间
 * @Param status 请求状态
 */
struct RefundTokenResponse {
    1: model.BaseResp base,
    2: required i64 refundID
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
    2: required string refundReason
}

/*
 * struct RefundResponse 退款响应
 * @Param refundID 退款ID
 * @Param status 请求状态
 */
struct RefundResponse {
    1: model.BaseResp base,
    2: model.RefundResponseInfo refundInfo
}

/*
 * struct RefundReviewRequest 退款审核请求
 * @Param orderID 关联的商户订单号
 * @Param passed 审核是否通过
 */
struct RefundReviewRequest {
    1: required i64 orderID (api.body="order_id")
    2: required bool passed
}

/*
 * struct RefundReviewResponse 退款审核响应
 */
struct RefundReviewResponse {
    1: model.BaseResp base
}

/*
 * struct PaymentCheckoutRequest 订单结算请求
 * @Param orderID 订单号
 * @Param token 支付令牌
 */
struct PaymentCheckoutRequest {
    1: required i64 orderID (api.body="order_id")
    2: required string token
}

/*
 * struct PaymentCheckoutResponse 订单支付响应
 */
struct PaymentCheckoutResponse {
    1: model.BaseResp base
}

/*
 * service PaymentService 支付服务
 * @Method RequestPaymentToken 请求支付令牌
 * @Method ProcessPayment 处理支付
 * @Method PaymentCheckout 订单支付
 * @Method RequestRefundToken 请求退款令牌
 * @Method RefundReview 退款审核
 */
service PaymentService {
    PaymentResponse ProcessPayment(1: PaymentRequest request) (api.post="/api/payment/process")
    PaymentTokenResponse RequestPaymentToken(1: PaymentTokenRequest request) (api.get="/api/payment/token")
    PaymentCheckoutResponse RequestPaymentCheckout(1: PaymentCheckoutRequest request) (api.post="/api/payment/checkout")
    RefundReviewResponse RefundReview(1: RefundReviewRequest request) (api.post="/api/payment/refund/review")
    RefundResponse RequestRefund(1: RefundRequest request) (api.get="/api/payment/refund")
}

