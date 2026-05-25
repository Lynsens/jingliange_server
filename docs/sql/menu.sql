CREATE TABLE IF NOT EXISTS `jlg`.`menu` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '流水ID,自增序列号',
  `image_url` varchar(256) NOT NULL DEFAULT '' COMMENT '图片保存地址',
  `desc` varchar(512) NOT NULL DEFAULT '' COMMENT '菜品介绍',
  `nutrition` varchar(512) NOT NULL DEFAULT '' COMMENT '营养价值表，json 格式',
  `status` int(8) unsigned NOT NULL DEFAULT '0' COMMENT '状态：0 删除，1 正常',
  `create_time` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '记录创建时间',
  `update_time` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '记录更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC COMMENT='菜单';

ALTER TABLE `jlg`.`menu`
ADD COLUMN `name` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '菜品名称' AFTER `id`;

ALTER TABLE `jlg`.`menu`
ADD COLUMN `ingredients` JSON NOT NULL COMMENT '成分表，JSON 格式' AFTER `nutrition`;

ALTER TABLE `jlg`.`menu`
ADD COLUMN `is_recommended` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否今日推荐：0 否，1 是' AFTER `ingredients`;

ALTER TABLE `jlg`.`menu`
ADD COLUMN `is_archived` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否下架：0 上架，1 下架' AFTER `is_recommended`;

ALTER TABLE `jlg`.`menu`
ADD COLUMN `archive_time` DATETIME(3) NULL DEFAULT NULL COMMENT '下架时间' AFTER `is_archived`;

CREATE TABLE IF NOT EXISTS `jlg`.`menu_feedback` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '流水ID,自增序列号',
  `menu_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '菜单 ID',
  `user_id` varchar(64) NOT NULL DEFAULT '' COMMENT '用户 ID',
  `preference` int(8) unsigned NOT NULL DEFAULT '0' COMMENT '状态：0 默认，1 喜欢，2 不喜欢',
  `comment` varchar(128) NOT NULL DEFAULT '' COMMENT '评论',
  `status` int(8) unsigned NOT NULL DEFAULT '0' COMMENT '状态：0 删除，1 正常',
  `create_time` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '记录创建时间',
  `update_time` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '记录更新时间',
  PRIMARY KEY (`id`),
  FOREIGN KEY (`menu_id`) REFERENCES `menu`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC COMMENT='菜单反馈';

ALTER TABLE `jlg`.`menu_feedback`
ADD COLUMN `user_nickname` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '评论展示昵称' AFTER `user_id`;

ALTER TABLE `jlg`.`menu_feedback`
ADD COLUMN `user_avatar_url` VARCHAR(256) NOT NULL DEFAULT '' COMMENT '评论展示头像 URL' AFTER `user_nickname`;

ALTER TABLE `jlg`.`menu_feedback`
ADD UNIQUE KEY `uk_menu_feedback_menu_user` (`menu_id`, `user_id`);

CREATE TABLE IF NOT EXISTS `jlg`.`menu_comment_like` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '流水ID,自增序列号',
  `comment_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '评论反馈 ID',
  `user_id` varchar(64) NOT NULL DEFAULT '' COMMENT '用户 ID',
  `status` int(8) unsigned NOT NULL DEFAULT '1' COMMENT '状态：0 取消，1 已点赞',
  `create_time` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '记录创建时间',
  `update_time` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '记录更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_menu_comment_like_comment_user` (`comment_id`, `user_id`),
  KEY `idx_menu_comment_like_comment_status` (`comment_id`, `status`),
  FOREIGN KEY (`comment_id`) REFERENCES `menu_feedback`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC COMMENT='菜单评论点赞';
