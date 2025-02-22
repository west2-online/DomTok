namespace go order

include "model.thrift"



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
    3: required list<model.baseOrderWithGoods> orderList;
}

struct ViewOrderReq {
    1: required i64 orderID;
}

struct ViewOrderResp {
    1: required model.BaseResp base;
    2: required model.orderWithGoods data;
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
    CreateOrderResp CreateOrder(1:CreateOrderReq req) (api.post="/api/order/create")
    ViewOrderListResp ViewOrderList(1:ViewOrderListReq req) (api.get="/api/order/list")
    ViewOrderResp ViewOrder(1:ViewOrderReq req) (api.get="/api/order/view")
    CancelOrderResp CancelOrder(1:CancelOrderReq req) (api.delete="/api/order/cancel")
    ChangeDeliverAddressResp ChangeDeliverAddress(1:ChangeDeliverAddressReq req) (api.put="/api/order/change-address")
    DeleteOrderResp DeleteOrder(1:DeleteOrderReq req) (api.delete="/api/order/delete")
}
