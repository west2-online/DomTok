namespace go user

include "model.thrift"

struct RegisterRequest {
    1: required string username,
    2: required string password,
    3: required string email,
    4: required string phone,
}

struct RegisterResponse {
    1: required model.BaseResp base,
    2: required i64 userID,
}

struct LoginRequest {
    1: string username,
    2: string password,
    3: string confirm_password,
}

struct LoginResponse {
    1: model.BaseResp base,
    2: model.UserInfo user,
}

service UserService {
    RegisterResponse Register(1: RegisterRequest req),
    LoginResponse Login(1: LoginRequest req),
}
