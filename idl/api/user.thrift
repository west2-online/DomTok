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


service UserService {
    RegisterResponse Register(1: RegisterRequest req)(api.get = "api/v1/user/register"),
}
