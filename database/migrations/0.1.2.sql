CREATE TABLE IF NOT EXISTS `plg_news_advertisement_config` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `advertisement_id` bigint unsigned NOT NULL COMMENT '核心广告配置ID',
  `position` tinyint unsigned NOT NULL COMMENT '展示位置：1新闻列表，2新闻详情',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_news_ad_position` (`advertisement_id`,`position`),
  KEY `idx_news_ad_position` (`position`),
  KEY `idx_news_ad_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='新闻插件广告展示配置';
