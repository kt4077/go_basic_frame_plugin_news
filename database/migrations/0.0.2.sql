ALTER TABLE `plg_news_article`
  ADD COLUMN `author` varchar(100) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '文章作者' AFTER `summary`,
  ADD COLUMN `source_name` varchar(100) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '来源平台名称' AFTER `author`,
  ADD COLUMN `source_url` varchar(1024) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '来源外部链接' AFTER `source_name`;
