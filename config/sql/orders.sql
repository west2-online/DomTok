-- 订单服务

-- 订单表
CREATE TABLE `order` (
                          `id` BIGINT UNSIGNED NOT NULL PRIMARY KEY COMMENT '订单ID',
                          `status` TINYINT NOT NULL COMMENT '订单状态',
                          `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
                          `total_amount_of_goods` DECIMAL(15,4) NOT NULL COMMENT '商品总金额',
                          `total_amount_of_freight` DECIMAL(15,4) NOT NULL COMMENT '商品总运费',
                          `total_amount_of_discount` DECIMAL(15,4) NOT NULL COMMENT '商品总优惠',
                          `payment_amount` DECIMAL(15,4) NOT NULL COMMENT '支付金额',
                          `payment_status` TINYINT NOT NULL COMMENT '支付状态',
                          `payment_at` BIGINT DEFAULT 0 COMMENT '支付时间(毫秒级时间戳)',
                          `payment_style` VARCHAR(32) NOT NULL COMMENT '支付类型',
                          `ordered_at` BIGINT UNSIGNED NOT NULL COMMENT '下单时间(毫秒级时间戳)',
                          `deleted_at` BIGINT UNSIGNED DEFAULT NULL COMMENT '订单删除时间(毫秒级时间戳)',
                          `delivery_at` BIGINT UNSIGNED DEFAULT NULL COMMENT '发货时间(毫秒级时间戳)',
                          `address_id` BIGINT UNSIGNED NOT NULL COMMENT '地址信息ID',
                          `address_info` VARCHAR(255) NOT NULL COMMENT '简略地址信息',
                          `coupon_id` BIGINT NOT NULL DEFAULT 0 COMMENT 'coupon id 用于全局优惠',
                          `coupon_name` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '优惠券名称',
                          INDEX `idx_order_deleted_at` (`deleted_at`),
                          INDEX `idx_order_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `order_goods` (
                               `order_id` BIGINT UNSIGNED NOT NULL COMMENT '订单ID',
                               `merchant_id` BIGINT UNSIGNED NOT NULL COMMENT '商家ID 对应 user 中的 uid',
                               `goods_id` BIGINT UNSIGNED NOT NULL COMMENT '商品ID 对应 commodity 中的 spu_id',
                               `goods_name` VARCHAR(128) NOT NULL COMMENT '商品名称',
                               `style_id` TINYINT NOT NULL COMMENT '款式ID 对应 commodity 中的 sku_id',
                               `style_name` VARCHAR(128) NOT NULL COMMENT '款式名称',
                               `goods_version` INT UNSIGNED NOT NULL COMMENT '快照版本号',
                               `style_head_drawing` VARCHAR(512) NOT NULL COMMENT '款式头图',
                               `origin_price` DECIMAL(11,4) NOT NULL COMMENT '原价',
                               `sale_price` DECIMAL(11,4) NOT NULL COMMENT '售卖价',
                               `single_freight_price` DECIMAL(11,4) NOT NULL COMMENT '单个运费',
                               `purchase_quantity` SMALLINT UNSIGNED NOT NULL COMMENT '购买数量',
                               `freight_amount` DECIMAL(15,4) NOT NULL COMMENT '运费金额',
                               `discount_amount` DECIMAL(15,4) NOT NULL COMMENT '优惠金额',
                               `payment_amount` DECIMAL(15,4) NOT NULL COMMENT '应付金额',
                               `total_amount` DECIMAL(15,4) NOT NULL COMMENT '结算金额 = 优惠前 - 优惠金额',
                               `single_price` DECIMAL(11,4) NOT NULL COMMENT '下单单价',
                               `coupon_id` BIGINT NOT NULL DEFAULT 0 COMMENT '优惠券ID',
                               `coupon_name` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '优惠券名称',
                               PRIMARY KEY (`order_id`, `goods_id`, `style_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


