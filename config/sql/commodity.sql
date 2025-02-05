-- 商品服务

-- 类别表
CREATE TABLE `category` (
                            `id` BIGINT NOT NULL PRIMARY KEY COMMENT '分类ID',
                            `name` VARCHAR(255) NOT NULL COMMENT '类别名',
                            `creator_id` BIGINT NOT NULL COMMENT '创建者ID',
                            `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                            `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                            `deleted_at` TIMESTAMP COMMENT '删除时间',
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- 优惠券信息表
CREATE TABLE `coupon_info` (
                               `id` BIGINT NOT NULL PRIMARY KEY COMMENT '优惠券ID',
                               `uid` BIGINT NOT NULL COMMENT '用户ID',
                               `name` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '优惠券名称',
                               `type_info` TINYINT NOT NULL COMMENT '1：满减券，2：满减折扣',
                               `condition_cost` DECIMAL(15,4) DEFAULT 0 COMMENT '用券门槛',
                               `discount_amount` DECIMAL(15,4) DEFAULT 0.0 COMMENT '满减金额',
                               `discount` DECIMAL(2,1) DEFAULT 0.0 COMMENT '折扣，例如0.8表示八折',
                               `range_type` TINYINT NOT NULL COMMENT '优惠券的范围 1-商品(spu_id)，2-商品类型，3-任意类型',
                               `range_id` BIGINT NOT NULL COMMENT '优惠券的范围对应类型ID',
                               `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                               `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                               `deleted_at` TIMESTAMP COMMENT '删除时间',
                               `expire_time` TIMESTAMP NOT NULL COMMENT '有效期',
                                `deadline_for_get` TIMESTAMP NOT NULL COMMENT '可以领取该券的截止时间',
                               `description` VARCHAR(255) DEFAULT '' COMMENT '描述'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 优惠券用户关系表
CREATE TABLE `user_coupon` (
                               `coupon_id` BIGINT NOT NULL COMMENT '优惠券ID',
                               `uid` BIGINT NOT NULL COMMENT '用户ID',
                               `remaining_uses` TINYINT DEFAULT 1 COMMENT '优惠券剩余的可使用次数',
                               `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                               `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                               `deleted_at` TIMESTAMP COMMENT '删除时间',
                               PRIMARY KEY (`coupon_id`, `uid`) -- 复合主键存疑
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- stock keeping unit 库存量单位表
CREATE TABLE `sku_info` (
                            `id` BIGINT NOT NULL PRIMARY KEY COMMENT 'SKU ID',
                            `creator_id` BIGINT NOT NULL COMMENT '创建者ID',
                            `price` DECIMAL(11,4) NOT NULL COMMENT '价格',
                            `name` VARCHAR(255) DEFAULT '' COMMENT '商品名称',
                            `description` VARCHAR(255) DEFAULT '' COMMENT '商品规格描述',
                            `for_sale` TINYINT NOT NULL COMMENT '是否出售 1-是, 0-否',
                            `stock` BIGINT NOT NULL COMMENT '库存',
                            `lock_stock` BIGINT NOT NULL COMMENT '预留库存',
                            `history_version_id` bigint not null comment '历史版本号',
                            `style_head_drawing` VARCHAR(512) NOT NULL COMMENT '款式头图 URL',
                            `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                            `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                            `deleted_at` TIMESTAMP COMMENT '删除时间',
                            INDEX `idx_user_delete_forSale` (`uid`, `deleted_at`, `for_sale`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- sku 属性表
CREATE TABLE `sku_sale_attr` (
                                 `id` BIGINT NOT NULL PRIMARY KEY COMMENT '属性ID',
                                 `sku_id` BIGINT NOT NULL COMMENT 'SKU ID',
                                 `history_version_id` bigint not null comment 'SKU 历史版本号',
                                 `sale_attr` VARCHAR(255) DEFAULT NULL COMMENT 'SKU属性（商品属性）',
                                 `sale_value` VARCHAR(255) DEFAULT NULL COMMENT 'SKU属性值',
                                 `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                 `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                 `deleted_at` TIMESTAMP COMMENT '删除时间',
                                 INDEX `idx_skuId` (`sku_id`, `deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- sku spu 关系表
CREATE TABLE `spu_to_sku` (
                              `sku_id` BIGINT NOT NULL COMMENT 'SKU ID',
                              `spu_id` BIGINT NOT NULL COMMENT 'SPU ID',
                              `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                              `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                              `deleted_at` TIMESTAMP COMMENT '删除时间',
                              PRIMARY KEY (`sku_id`, `spu_id`),
                              INDEX `idx_deleted_created` (`deleted_at`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- spu 信息表
CREATE TABLE `spu_info` (
                            `id` BIGINT NOT NULL PRIMARY KEY COMMENT 'SPU ID',
                            `name` VARCHAR(255) NOT NULL COMMENT 'SPU名称',
                            `creator_id` BIGINT NOT NULL COMMENT '创建者ID',
                            `description` VARCHAR(255) DEFAULT '' COMMENT '描述',
                            `category_id` BIGINT NOT NULL COMMENT '类别ID',
                            `goods_head_drawing` VARCHAR(512) NOT NULL COMMENT '商品头图 URL',
                            `price` DECIMAL(11,4) NOT NULL COMMENT '价格',
                            `for_sale` TINYINT NOT NULL COMMENT '是否出售 1-是, 0-否',
                            `shipping` DECIMAL(11,4) NOT NULL DEFAULT 0.0 COMMENT '运费',
                            `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                            `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                            `deleted_at` TIMESTAMP COMMENT '删除时间',
                            INDEX `idx_user_delete_forSale` (`uid`, `deleted_at`, `for_sale`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- sku 的轮播图表
CREATE TABLE `sku_image` (
                             `id` BIGINT NOT NULL PRIMARY KEY COMMENT '图片ID',
                             `url` VARCHAR(512) NOT NULL COMMENT '图片URL',
                             `sku_id` BIGINT NOT NULL COMMENT 'SKU ID',
                             `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                             `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                             `deleted_at` TIMESTAMP COMMENT '删除时间',
                             INDEX `idx_skuId_delete_created` (`sku_id`, `deleted_at`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- spu 的轮播图表
CREATE TABLE `spu_image` (
                             `id` BIGINT NOT NULL PRIMARY KEY COMMENT '图片ID',
                             `url` VARCHAR(255) NOT NULL COMMENT '图片URL',
                             `spu_id` BIGINT NOT NULL COMMENT 'SPU ID',
                             `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                             `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                             `deleted_at` TIMESTAMP COMMENT '删除时间',
                             INDEX `idx_spuId_delete_created` (`spu_id`, `deleted_at`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- sku 历史快照
CREATE TABLE `sku_price_history` (
                                     `id` BIGINT NOT NULL PRIMARY KEY COMMENT '版本ID',
                                     `sku_id` BIGINT NOT NULL COMMENT 'SKU ID',
                                     `mark_price` DECIMAL(11,4) NOT NULL COMMENT '该版本对应的价格',
                                     `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                     `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                     `deleted_at` TIMESTAMP COMMENT '删除时间',
                                     `prev_version` BIGINT COMMENT '上个版本的ID',
                                     INDEX `idx_skuId_created` (`sku_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
