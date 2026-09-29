ALTER TABLE `plg_news_article`
  ADD COLUMN `collection_count` bigint unsigned NOT NULL DEFAULT '0' COMMENT '收藏数量' AFTER `share_count`;

CREATE TABLE IF NOT EXISTS `plg_news_article_collection` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `uid` varchar(32) COLLATE utf8mb4_general_ci NOT NULL COMMENT '收藏业务唯一标识',
  `article_uid` varchar(32) COLLATE utf8mb4_general_ci NOT NULL COMMENT '文章业务唯一标识',
  `member_sn` varchar(32) COLLATE utf8mb4_general_ci NOT NULL COMMENT '会员业务编号',
  `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态，1收藏，2取消',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_plg_news_article_collection_uid` (`uid`),
  UNIQUE KEY `uk_plg_news_article_collection_member` (`article_uid`,`member_sn`),
  KEY `idx_plg_news_article_collection_member_status` (`member_sn`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='新闻文章收藏记录表';
