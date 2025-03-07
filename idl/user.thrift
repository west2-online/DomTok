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

struct GetAddressRequest {
    1: required i64 address_id
}

struct GetAddressResponse {
    1: required model.BaseResp base,
    2: optional model.AddressInfo address,
}

struct AddAddressRequest {
    1: required string province,
    2: required string city,
    3: required string detail,
}

struct AddAddressResponse {
    1: required model.BaseResp base,
    2: required i64 addressID,
}

struct BanUserReq {
    1: required i64 uid
}

struct BanUserResp {
    1: required model.BaseResp base,
}

struct LiftBanUserReq {
    1: required i64 uid
}

struct LiftBanUserResp {
    1: required model.BaseResp base,
}

struct LogoutReq {
}

struct LogoutResp {
    1: required model.BaseResp base,
}

struct SetAdministratorReq {
    1: required i64 uid
    2: required string password
    3: required i16 action
}

struct SetAdministratorResp {
    1: required model.BaseResp base,
}

service UserService {
    RegisterResponse Register(1: RegisterRequest req),
    LoginResponse Login(1: LoginRequest req),
    GetAddressResponse GetAddress(1: GetAddressRequest req),
    AddAddressResponse AddAddress(1: AddAddressRequest req),
    BanUserResp BanUser(1: BanUserReq req),
    LiftBanUserResp LiftBandUser(1: LiftBanUserReq req) ,
    LogoutResp Logout(1: LogoutReq req),
    SetAdministratorResp SetAdministrator(1:SetAdministratorReq req)
}
