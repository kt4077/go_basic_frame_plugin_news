ALTER TABLE `plg_news_article`
  ADD COLUMN `uid` varchar(32) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '文章业务唯一标识' AFTER `id`,
  ADD COLUMN `like_count` bigint unsigned NOT NULL DEFAULT '0' COMMENT '点赞数量' AFTER `virtual_view_count`,
  ADD COLUMN `share_count` bigint unsigned NOT NULL DEFAULT '0' COMMENT '发起分享数量' AFTER `like_count`,
  ADD COLUMN `comment_count` bigint unsigned NOT NULL DEFAULT '0' COMMENT '已展示一级评论数量' AFTER `share_count`,
  ADD UNIQUE KEY `uk_plg_news_article_uid` (`uid`);

CREATE TABLE IF NOT EXISTS `plg_news_article_like` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `uid` varchar(32) COLLATE utf8mb4_general_ci NOT NULL COMMENT '点赞业务唯一标识',
  `article_uid` varchar(32) COLLATE utf8mb4_general_ci NOT NULL COMMENT '文章业务唯一标识',
  `member_sn` varchar(32) COLLATE utf8mb4_general_ci NOT NULL COMMENT '会员业务编号',
  `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态，1点赞，2取消',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_plg_news_article_like_uid` (`uid`),
  UNIQUE KEY `uk_plg_news_article_like_member` (`article_uid`,`member_sn`),
  KEY `idx_plg_news_article_like_member_status` (`member_sn`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='新闻文章点赞记录表';

CREATE TABLE IF NOT EXISTS `plg_news_comment` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `uid` varchar(32) COLLATE utf8mb4_general_ci NOT NULL COMMENT '评论业务唯一标识',
  `article_uid` varchar(32) COLLATE utf8mb4_general_ci NOT NULL COMMENT '文章业务唯一标识',
  `member_sn` varchar(32) COLLATE utf8mb4_general_ci NOT NULL COMMENT '评论会员业务编号',
  `request_uid` varchar(64) COLLATE utf8mb4_general_ci NOT NULL COMMENT '客户端幂等请求标识',
  `parent_uid` varchar(32) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '直接父评论业务标识，一级评论为空',
  `root_uid` varchar(32) COLLATE utf8mb4_general_ci NOT NULL COMMENT '一级评论业务标识',
  `reply_member_sn` varchar(32) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '被回复会员业务编号',
  `content` varchar(1000) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '评论内容',
  `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态，1待审核，2显示，3关闭',
  `like_count` bigint unsigned NOT NULL DEFAULT '0' COMMENT '点赞数量',
  `reply_count` bigint unsigned NOT NULL DEFAULT '0' COMMENT '回复数量',
  `platform` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '发布平台，复用用户注册来源枚举',
  `ip` varchar(64) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '发布IP，仅后台审计使用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_plg_news_comment_uid` (`uid`),
  UNIQUE KEY `uk_plg_news_comment_request` (`member_sn`,`request_uid`),
  KEY `idx_plg_news_comment_article_status_time` (`article_uid`,`status`,`created_at`),
  KEY `idx_plg_news_comment_root_status_time` (`root_uid`,`status`,`created_at`),
  KEY `idx_plg_news_comment_parent_uid` (`parent_uid`),
  KEY `idx_plg_news_comment_member_time` (`member_sn`,`created_at`),
  KEY `idx_plg_news_comment_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='新闻文章评论表';

CREATE TABLE IF NOT EXISTS `plg_news_comment_image` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `uid` varchar(32) COLLATE utf8mb4_general_ci NOT NULL COMMENT '评论图片业务唯一标识',
  `comment_uid` varchar(32) COLLATE utf8mb4_general_ci NOT NULL COMMENT '评论业务唯一标识',
  `path` varchar(1024) COLLATE utf8mb4_general_ci NOT NULL COMMENT '图片相对路径',
  `sort` int NOT NULL DEFAULT '0' COMMENT '排序值',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_plg_news_comment_image_uid` (`uid`),
  KEY `idx_plg_news_comment_image_comment_sort` (`comment_uid`,`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='新闻评论图片表';

CREATE TABLE IF NOT EXISTS `plg_news_comment_like` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `uid` varchar(32) COLLATE utf8mb4_general_ci NOT NULL COMMENT '评论点赞业务唯一标识',
  `comment_uid` varchar(32) COLLATE utf8mb4_general_ci NOT NULL COMMENT '评论业务唯一标识',
  `member_sn` varchar(32) COLLATE utf8mb4_general_ci NOT NULL COMMENT '会员业务编号',
  `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态，1点赞，2取消',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_plg_news_comment_like_uid` (`uid`),
  UNIQUE KEY `uk_plg_news_comment_like_member` (`comment_uid`,`member_sn`),
  KEY `idx_plg_news_comment_like_member_status` (`member_sn`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='新闻评论点赞记录表';

CREATE TABLE IF NOT EXISTS `plg_news_share_record` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `uid` varchar(32) COLLATE utf8mb4_general_ci NOT NULL COMMENT '分享记录业务唯一标识',
  `article_uid` varchar(32) COLLATE utf8mb4_general_ci NOT NULL COMMENT '文章业务唯一标识',
  `member_sn` varchar(32) COLLATE utf8mb4_general_ci NOT NULL COMMENT '会员业务编号',
  `request_uid` varchar(64) COLLATE utf8mb4_general_ci NOT NULL COMMENT '客户端请求唯一标识',
  `share_channel` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '分享渠道，1微信好友，2微信朋友圈，3系统分享，4复制链接',
  `platform` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '发起平台，复用用户注册来源枚举',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_plg_news_share_record_uid` (`uid`),
  UNIQUE KEY `uk_plg_news_share_request` (`member_sn`,`request_uid`),
  KEY `idx_plg_news_share_article_time` (`article_uid`,`created_at`),
  KEY `idx_plg_news_share_member_time` (`member_sn`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='新闻文章分享发起记录表';
