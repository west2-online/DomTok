-- 支付服务

-- 支付订单表
CREATE TABLE `payment_orders` (
                                  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '支付订单的唯一标识',
                                  `order_id` BIGINT NOT NULL COMMENT '商户订单号',
                                  `user_id` BIGINT NOT NULL COMMENT '用户的唯一标识',
                                  `amount` DECIMAL(15,4) NOT NULL COMMENT '订单总金额',
                                  `status` TINYINT NOT NULL DEFAULT 0 COMMENT '支付状态：0-待支付，1-处理中，2-成功支付 3-支付失败',
                                  `masked_credit_card_number` VARCHAR(19) COMMENT '信用卡号 国际信用卡号的最大长度为19 (仅存储掩码，如 **** **** **** 1234)',
                                  `credit_card_expiration_year` INT COMMENT '信用卡到期年',
                                  `credit_card_expiration_month` INT COMMENT '信用卡到期月',
                                  `description` VARCHAR(255) COMMENT '订单描述信息',
                                  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '订单创建时间',
                                  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '订单最后更新时间',
                                  `deleted_at` TIMESTAMP NULL COMMENT '订单删除时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 退款表
CREATE TABLE `payment_refunds` (
                                   `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '支付退款的唯一标识',
                                   `order_id` VARCHAR(64) NOT NULL COMMENT '关联的商户订单号',
                                   `user_id` BIGINT NOT NULL COMMENT '用户的唯一标识',
                                   `refund_amount` DECIMAL(15,4) NOT NULL COMMENT '退款金额，单位为元',
                                   `refund_reason` VARCHAR(255) COMMENT '退款原因',
                                   `status` TINYINT NOT NULL DEFAULT 0 COMMENT '退款状态：0-待处理，1-处理中，2-成功退款 3-退款失败',
                                   `masked_credit_card_number` VARCHAR(19) COMMENT '信用卡号 国际信用卡号的最大长度为19 (仅存储掩码，如 **** **** **** 1234)',
                                   `credit_card_expiration_year` INT COMMENT '信用卡到期年',
                                   `credit_card_expiration_month` INT COMMENT '信用卡到期月',
                                   `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '退款申请时间',
                                   `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '退款最后更新时间',
                                   `deleted_at` TIMESTAMP NULL COMMENT '退款记录删除时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 流水信息表
CREATE TABLE `payment_ledger` (
                                  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '流水ID',
                                  `reference_id` BIGINT NOT NULL COMMENT '关联的支付订单或退款订单ID',
                                  `user_id` BIGINT NOT NULL COMMENT '用户ID',
                                  `amount` DECIMAL(15,4) NOT NULL COMMENT '交易金额（正数表示收入，负数表示支出）',
                                  `transaction_type` TINYINT NOT NULL COMMENT '交易类型：1-支付，2-退款，3-手续费，4-调整',
                                  `status` TINYINT NOT NULL DEFAULT 0 COMMENT '交易状态：0-待处理，1-成功，2-失败',
                                  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '交易创建时间',
                                  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '交易更新时间',
                                  `deleted_at` TIMESTAMP NULL COMMENT '交易记录删除时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf
