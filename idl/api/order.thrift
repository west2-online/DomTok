namespace go api.order

include "../model.thrift"

struct CreateOrderReq {
    1: required i64 addressID; // 地址信息 ID
    2: required string addressInfo; // 简略地址信息
    3: required list<model.BaseOrderGoods> baseOrderGoods; // 商品列表
}

struct CreateOrderResp {
    1: required model.BaseResp base;
    2: required i64 orderID; // 订单号
}

struct ViewOrderListReq {
    1: i32 page;
    2: i32 size;
}

struct ViewOrderListResp {
    1: required model.BaseResp base;
    2: required i32 total;
    3: required i64 orderID;
    4: required list<model.OrderGoods> orderGoods;
}

struct ViewOrderReq {
    1: required i64 orderID;
}

struct ViewOrderResp {
    1: required model.BaseResp base;
    2: required double totalAmountOfGoods; // 商品总金额
    3: required double totalAmountOfFreight; // 总运费
    4: required double totalAmountOfDiscount; // 总优惠
    5: required double paymentAmount; // 实际付款价
    6: required i64 addressID; // 地址信息 ID
    7: required string addressInfo; // 简略地址信息
    8: required string status; // 订单状态
    9: required list<model.OrderGoods> orderGoods; // 商品列表
}

struct CancelOrderReq {
    1: required i64 orderID;
}

struct CancelOrderResp {
    1: required model.BaseResp base;
}

struct ChangeDeliverAddressReq {
    1: required i64 addressID;
    2: required string addressInfo;
    3: required i64 orderID;
}

struct ChangeDeliverAddressResp {
    1: required model.BaseResp base;
}

struct DeleteOrderReq {
    1: required i64 orderID;
}

struct DeleteOrderResp {
    1: required model.BaseResp base;
}

service OrderService {
    CreateOrderResp CreateOrder(1:CreateOrderReq req) (api.post="/api/v1/order/create")
    ViewOrderListResp ViewOrderList(1:ViewOrderListReq req) (api.get="/api/v1/order/list")
    ViewOrderResp ViewOrder(1:ViewOrderReq req) (api.get="/api/v1/order/view")
    CancelOrderResp CancelOrder(1:CancelOrderReq req) (api.delete="/api/v1/order/cancel")
    ChangeDeliverAddressResp ChangeDeliverAddress(1:ChangeDeliverAddressReq req) (api.put="/api/v1/order/change-address")
    DeleteOrderResp DeleteOrder(1:DeleteOrderReq req) (api.delete="/api/v1/order/delete")
}
