namespace go api.user
include "../model.thrift"

//
struct RegisterRequest {
    1: required string name
    2: required string password
    3: required string email
}

struct RegisterResponse {
    1: required i64 uid;
}

struct LoginRequest {
    1: required string name
    2: required string password
}

struct LoginResponse {
    2: model.UserInfo user,
}

struct GetAddressRequest {
    1: required i64 address_id
}

struct GetAddressResponse {
    1: model.AddressInfo address,
}

struct AddAddressRequest {
    1: required string province
    2: required string city
    3: required string detail
}

struct AddAddressResponse {
    1: required i64 address_id
}

struct BanUserReq {
    1: required i64 uid
}

struct BanUserResp {
}

struct LiftBanUserReq {
    1: required i64 uid
}

struct LiftBanUserResp {

}

struct LogoutReq {
}

struct LogoutResp {

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
    RegisterResponse Register(1: RegisterRequest req)(api.post = "api/v1/user/register"),
    LoginResponse Login(1: LoginRequest req)(api.post = "api/v1/user/login")
    GetAddressResponse GetAddress(1: GetAddressRequest req)(api.get = "api/v1/user/location"),
    AddAddressResponse AddAddress(1: AddAddressRequest req)(api.post = "api/v1/user/address"),
    BanUserResp BanUser(1: BanUserReq req) (api.post="api/v1/user/ban"),
    LiftBanUserResp LiftBandUser(1: LiftBanUserReq req) (api.post="api/v1/user/lift"),
    LogoutResp Logout(1: LogoutReq req) (api.post="api/v1/user/logout")
    SetAdministratorResp SetAdministrator(1:SetAdministratorReq req) (api.post="api/v1/user/administrator")
}

