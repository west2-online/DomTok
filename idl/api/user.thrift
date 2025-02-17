namespace go api.user
include "../model.thrift"

//
struct RegisterRequest {
    1: required string name
    2: required string password
    3: required string email
}

struct RegisterResponse {
}

struct LoginRequest {
    1: required string name
    2: required string password
}

struct LoginResponse {
        1: model.BaseResp base,
        2: model.UserInfo user,
}


service UserService {
    RegisterResponse Register(1: RegisterRequest req)(api.get = "api/v1/user/register"),
    LoginResponse Login(1: LoginRequest req)(api.post = "api/v1/user/login")
}

// cart
struct AddGoodsIntoCartRequest{
    1: required i64 skuId,
    2: required i64 shop_id,
    3: required i64 count,
}

struct AddGoodsIntoCartResponse{
    1: required model.BaseResp base,
}

service CartService {
    AddGoodsIntoCartResponse AddGoodsIntoCart(1: AddGoodsIntoCartRequest req)(api.post = "api/cart/sku/add"),
}
