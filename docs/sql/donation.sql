-- 创建用户表
CREATE TABLE IF NOT EXISTS `user` (
    `id` VARCHAR(64) NOT NULL COMMENT '用户ID（来自微信小程序）',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 创建捐款记录表
CREATE TABLE IF NOT EXISTS `donation` (
    `id` INT(11) NOT NULL AUTO_INCREMENT COMMENT '捐款记录ID',
    `user_id` VARCHAR(64) NOT NULL COMMENT '用户ID',
    `donor_name` VARCHAR(32) NOT NULL COMMENT '捐款人昵称',
    `amount` DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '捐款金额',
    `donate_time` DATETIME NOT NULL COMMENT '捐款时间',
    `is_visible` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否显示在榜单：0隐藏，1显示',
    `message` VARCHAR(256) NOT NULL DEFAULT '' COMMENT '留言',
    `remarks` VARCHAR(256) NOT NULL DEFAULT '' COMMENT '备注（管理员使用）',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_donate_time` (`donate_time`),
    KEY `idx_amount` (`amount`),
    KEY `idx_is_visible` (`is_visible`),
    CONSTRAINT `fk_donation_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='捐款记录表';

-- 插入一些测试数据
INSERT INTO `user` (`id`) VALUES 
('user123'),
('user456'),
('user789');

INSERT INTO `donation` (`user_id`, `donor_name`, `amount`, `donate_time`, `is_visible`, `message`) VALUES
('user123', '善心人士', 100.00, '2025-07-08 10:00:00', 1, '祝愿净莲阁越来越好'),
('user456', '匿名捐款者', 50.00, '2025-07-07 15:30:00', 1, '随喜功德'),
('user789', '佛心居士', 200.00, '2025-07-06 09:15:00', 1, '阿弥陀佛'),
('user123', '善心人士', 20.00, '2025-01-15 14:20:00', 1, '新年快乐'),
('user456', '匿名捐款者', 88.00, '2025-01-10 16:45:00', 1, '心诚则灵');
