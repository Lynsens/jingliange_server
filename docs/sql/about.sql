CREATE TABLE IF NOT EXISTS `jlg`.`desc` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '流水ID,自增序列号',
  `content` varchar(512) NOT NULL DEFAULT '' COMMENT '餐厅介绍',
  `status` int(8) unsigned NOT NULL DEFAULT '0' COMMENT '状态：0 删除，1 正常',
  `create_time` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '记录创建时间',
  `update_time` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '记录更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC COMMENT='餐厅介绍表';


CREATE TABLE IF NOT EXISTS `jlg`.`activity` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '流水ID,自增序列号',
  `title` varchar(64) NOT NULL DEFAULT '' COMMENT '标题',
  `content` varchar(512) NOT NULL DEFAULT '' COMMENT '活动介绍',
  `status` int(8) unsigned NOT NULL DEFAULT '0' COMMENT '状态：0 删除，1 正常',
  `create_time` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '记录创建时间',
  `update_time` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '记录更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC COMMENT='活动表';


CREATE TABLE IF NOT EXISTS `jlg`.`images` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '流水ID,自增序列号',
  `address` varchar(64) NOT NULL DEFAULT '' COMMENT '图片保存地址',
  `desc` varchar(512) NOT NULL DEFAULT '' COMMENT '图片简介',
  `status` int(8) unsigned NOT NULL DEFAULT '0' COMMENT '状态：0 删除，1 正常',
  `create_time` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '记录创建时间',
  `update_time` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '记录更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC COMMENT='图片表';