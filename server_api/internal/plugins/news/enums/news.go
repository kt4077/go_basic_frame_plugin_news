// Package enums 定义新闻插件业务枚举，所有枚举值从1开始。
package enums

const (
	// CategoryStatusEnabled 表示分类启用。
	CategoryStatusEnabled = 1
	// CategoryStatusDisabled 表示分类停用。
	CategoryStatusDisabled = 2
)

const (
	// ArticleStatusDraft 表示文章草稿。
	ArticleStatusDraft = 1
	// ArticleStatusPublished 表示文章已发布。
	ArticleStatusPublished = 2
	// ArticleStatusOffline 表示文章已下线。
	ArticleStatusOffline = 3
)
