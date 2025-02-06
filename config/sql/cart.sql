-- 购物车服务

-- 购物车表
CREATE TABLE `cart` (
                           `user_id` BIGINT NOT NULL PRIMARY KEY COMMENT '用户ID',
                           `sku_json` TEXT NOT NULL COMMENT '商品json',
                           `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                           `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                           `deleted_at` TIMESTAMP DEFAULT NULL COMMENT '删除时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
