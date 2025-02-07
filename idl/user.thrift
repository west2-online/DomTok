namespace go user

include "model.thrift"

/*
 * struct RegisterRequest 用户注册请求
 * @Param username 用户名
 * @Param password 密码
 * @Param email 用户邮箱
 */
struct RegisterRequest {
    1: string username     /* 用户名 */
    2: string password     /* 密码 */
    3: string email        /* 用户邮箱 */
}

struct RegisterResponse {
    1: model.BaseResp base
    2: i64 userID
}

/*
 * struct LoginRequest 用户登录请求
 * @Param username 用户名
 * @Param password 密码
 */
struct LoginRequest {
    1: string username
    2: string password
}

struct LoginResponse {
    1: model.BaseResp base
    2: model.UserInfo user
}


/* 用户服务 */
service UserService {
    RegisterResponse register(1: RegisterRequest request)
    LoginResponse login(1: LoginRequest request)
}
