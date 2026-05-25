CREATE TABLE IF NOT EXISTS `jlg`.`suggestion` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '流水ID,自增序列号',
  `user_id` varchar(64) NOT NULL DEFAULT '' COMMENT '用户 ID',
  `user_nickname` varchar(64) NOT NULL DEFAULT '' COMMENT '用户昵称',
  `user_avatar_url` varchar(256) NOT NULL DEFAULT '' COMMENT '用户头像',
  `content` varchar(1024) NOT NULL DEFAULT '' COMMENT '建议内容',
  `contact` varchar(128) NOT NULL DEFAULT '' COMMENT '联系方式',
  `handle_status` int(8) unsigned NOT NULL DEFAULT '0' COMMENT '处理状态：0 未处理，1 已处理',
  `status` int(8) unsigned NOT NULL DEFAULT '1' COMMENT '记录状态：0 删除，1 正常',
  `create_time` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '记录创建时间',
  `update_time` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '记录更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_suggestion_status_handle_create` (`status`, `handle_status`, `create_time`),
  KEY `idx_suggestion_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC COMMENT='用户建议箱';
