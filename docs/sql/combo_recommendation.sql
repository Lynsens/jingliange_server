CREATE TABLE IF NOT EXISTS `jlg`.`combo_recommendation` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '流水ID,自增序列号',
  `title` varchar(64) NOT NULL DEFAULT '' COMMENT '套餐标题',
  `description` varchar(512) NOT NULL DEFAULT '' COMMENT '套餐说明',
  `is_active` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否首页启用：0 否，1 是',
  `status` int(8) unsigned NOT NULL DEFAULT '1' COMMENT '状态：0 删除，1 正常',
  `create_time` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '记录创建时间',
  `update_time` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '记录更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_combo_recommendation_status_active` (`status`, `is_active`, `update_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC COMMENT='套餐推荐';

CREATE TABLE IF NOT EXISTS `jlg`.`combo_recommendation_item` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '流水ID,自增序列号',
  `combo_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '套餐推荐 ID',
  `menu_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '菜品 ID',
  `sort_order` int(8) unsigned NOT NULL DEFAULT '0' COMMENT '排序',
  `status` int(8) unsigned NOT NULL DEFAULT '1' COMMENT '状态：0 删除，1 正常',
  `create_time` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '记录创建时间',
  `update_time` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '记录更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_combo_item_combo_status_sort` (`combo_id`, `status`, `sort_order`),
  KEY `idx_combo_item_menu` (`menu_id`),
  CONSTRAINT `fk_combo_item_combo` FOREIGN KEY (`combo_id`) REFERENCES `combo_recommendation` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_combo_item_menu` FOREIGN KEY (`menu_id`) REFERENCES `menu` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC COMMENT='套餐推荐菜品';
