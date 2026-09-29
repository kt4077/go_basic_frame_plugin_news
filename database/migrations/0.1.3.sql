ALTER TABLE `plg_news_advertisement_config`
  ADD COLUMN `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '启用状态：1启用，2停用' AFTER `position`,
  ADD KEY `idx_news_ad_status` (`status`);
