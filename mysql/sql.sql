--- 创建压测表
CREATE TABLE `stress_test` (
                               `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
                               `uuid` char(36) NOT NULL DEFAULT '' COMMENT 'UUID字符串',
                               `user_id` int unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
                               `amount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '金额',
                               `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                               `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态',
                               PRIMARY KEY (`id`),
                               KEY `idx_user_id` (`user_id`),
                               KEY `idx_create_time` (`create_time`),
                               KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='压测表';

