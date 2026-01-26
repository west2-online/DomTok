# User Address Management Enhancements

## 1. Context
To support user delivery address management, we have extended the User Service with capabilities to list and delete addresses. This complements the existing `GetAddress` and `AddAddress` functionalities, providing a complete CRUD experience for user addresses.

## 2. Interface Design

### 2.1 Thrift Definition Updates (`idl/user.thrift`)

The following interfaces were implemented:

#### ListAddress
Retrieves a paginated list of addresses for the current logged-in user.

```thrift
struct ListAddressRequest {
    1: required i64 pageNum,
    2: required i64 pageSize,
}

struct ListAddressResponse {
    1: required model.BaseResp base,
    2: required list<model.AddressInfo> addresses,
}
```

#### DeleteAddress
Deletes a specific address by ID. Includes ownership verification.

```thrift
struct DeleteAddressRequest {
    1: required i64 address_id
}

struct DeleteAddressResponse {
    1: required model.BaseResp base,
}
```

## 3. Database Changes

### 3.1 Schema Update (`config/sql/user.sql`)
A new `address` table has been introduced to store address information.

```sql
CREATE TABLE `address` (
    `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '地址ID',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `province` VARCHAR(255) NOT NULL COMMENT '省份',
    `city` VARCHAR(255) NOT NULL COMMENT '城市',
    `detail` VARCHAR(255) NOT NULL COMMENT '详细地址',
    INDEX `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户地址表';
```

## 4. Implementation Details

### 4.1 Domain Layer
- **Model**: Updated `Address` model to include `Uid` (User ID) for ownership verification.
- **Repository**: Added `ListAddress` and `DeleteAddress` methods to `UserDB` interface.
- **Service**: 
    - `ListAddress`: Retrieves addresses filtered by the current user's ID.
    - `DeleteAddress`: Fetches the address first to verify that the `user_id` matches the current user, ensuring data security.

### 4.2 Infrastructure Layer
- **MySQL**: Implemented GORM-based methods for `ListAddress` (with pagination) and `DeleteAddress`.

### 4.3 Controller Layer
- **RPC Handler**: Implemented `ListAddress` and `DeleteAddress` handlers, mapping Thrift requests to UseCase calls and converting domain models back to Thrift responses.

## 5. Testing Plan

### 5.1 Unit Tests
- Verify `DeleteAddress` fails if the user does not own the address.
- Verify `ListAddress` returns correct pagination results.
- Mock DB interactions to ensure service logic is correct.

### 5.2 Integration Tests
- Run against a real MySQL instance to verify SQL correctness.
- Test full flow from RPC handler to DB.
