namespace go user

include "model.thrift"

struct RegisterRequest {
    1: string username,
    2: string password,
    3: string email,
}

struct RegisterResponse {
    1: model.BaseResp base,
    2: i64 userID,
}

struct LoginRequest {
    1: string username,
    2: string password,
}

struct LoginResponse {
    1: model.BaseResp base,
    2: model.UserInfo user,
}