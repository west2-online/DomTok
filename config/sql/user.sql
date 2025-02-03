-- 用户服务 --

-- 用户表
CREATE TABLE `user` (
                        `id` BIGINT PRIMARY KEY COMMENT '使用 Snowflake 算法生成',
                        `username` VARCHAR(30) NOT NULL COMMENT '用户名最多 10 个中文字符或等长英文字符',
                        `password` CHAR(16) NOT NULL COMMENT '数字+字母组合，总长度上限 16',
                        `email` VARCHAR(50) NOT NULL COMMENT '邮箱'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;