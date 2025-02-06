namespace go order

include "model.thrift"



struct CreateOrderReq {
    1: required i64 AddressID; // 地址信息 ID
    2: required string AddressInfo; // 简略地址信息
    3: required list<model.BaseOrderGoods> BaseOrderGoods; // 商品列表
}

struct CreateOrderResp {
    1: required model.BaseResp Base;
    2: required i64 OrderID; // 订单号
}

struct ViewOrderListReq {
    1: i32 Page;
    2: i32 Size;
}

struct ViewOrderListResp {
    1: required model.BaseResp Base;
    2: required i32 Total;
    3: required i64 OrderID;
    4: required list<model.OrderGoods> OrderGoods;
}

struct ViewOrderReq {
    1: required i64 OrderID;
}

struct ViewOrderResp {
    1: required model.BaseResp Base;
    2: required double TotalAmountOfGoods; // 商品总金额
    3: required double TotalAmountOfFreight; // 总运费
    4: required double TotalAmountOfDiscount; // 总优惠
    5: required double PaymentAmount; // 实际付款价
    6: required i64 AddressID; // 地址信息 ID
    7: required string AddressInfo; // 简略地址信息
    8: required string Status; // 订单状态
    9: required list<model.OrderGoods> OrderGoods; // 商品列表
}

struct CancelOrderReq {
    1: required i64 OrderID;
}

struct CancelOrderResp {
    1: required model.BaseResp Base;
}

struct ChangeDeliverAddressReq {
    1: required i64 AddressID;
    2: required string AddressInfo;
    3: required i64 OrderID;
}

struct ChangeDeliverAddressResp {
    1: required model.BaseResp Base;
}

struct DeleteOrderReq {
    1: required i64 OrderID;
}

struct DeleteOrderResp {
    1: required model.BaseResp Base;
}

service OrderService {
    CreateOrderResp CreateOrder(1:CreateOrderReq req) (api.post="/api/order/create")
    ViewOrderListResp ViewOrderList(1:ViewOrderListReq req) (api.get="/api/order/list")
    ViewOrderResp ViewOrder(1:ViewOrderReq req) (api.get="/api/order/view")
    CancelOrderResp CancelOrder(1:CancelOrderReq req) (api.delete="/api/order/cancel")
    ChangeDeliverAddressResp ChangeDeliverAddress(1:ChangeDeliverAddressReq req) (api.put="/api/order/change-address")
    DeleteOrderResp DeleteOrder(1:DeleteOrderReq req) (api.delete="/api/order/delete")
}
