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


service UserService {
    RegisterResponse Register(1: RegisterRequest req)(api.post = "api/v1/user/register"),
    LoginResponse Login(1: LoginRequest req)(api.post = "api/v1/user/login")
    GetAddressResponse GetAddress(1: GetAddressRequest req)(api.get = "api/v1/user/location"),
    AddAddressResponse AddAddress(1: AddAddressRequest req)(api.post = "api/v1/user/address")
}
