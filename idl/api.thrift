namespace go api
include "model.thrift"

//
struct RegisterRequest {
    1: required string name
    2: required string password
    3: required string email
}

struct RegisterResponse {
}


service UserService {
    RegisterResponse Register(1: RegisterRequest req)(api.get = "api/v1/user/register"),
}

// cart
struct AddGoodsIntoCartRequest{
    1: required i64 skuId,
    3: required i64 count,
}

struct AddGoodsIntoCartResponse{
    1: required model.BaseResp base,
}

service CartService {
    AddGoodsIntoCartResponse AddGoodsIntoCart(1: AddGoodsIntoCartRequest req)(api.post = "api/cart/sku/add"),
}
