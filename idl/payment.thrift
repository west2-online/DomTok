namespace go payment

include "model.thrift"

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
    4: required CreditCardInfo creditCard
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

/*
 * service PaymentService 支付服务
 */
service PaymentService {
    PaymentResponse ProcessPayment(1: PaymentRequest request)
    RefundResponse ProcessRefund(1: RefundRequest request)
}
