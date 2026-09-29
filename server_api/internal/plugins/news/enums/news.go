// Package enums 定义新闻插件业务枚举，所有枚举值从1开始。
package enums

const (
	AdvertisementPositionList   = 1
	AdvertisementPositionDetail = 2
)

func IsValidAdvertisementPosition(value int) bool {
	return value == AdvertisementPositionList || value == AdvertisementPositionDetail
}

const (
	// CategoryStatusEnabled 表示分类启用。
	CategoryStatusEnabled = 1
	// CategoryStatusDisabled 表示分类停用。
	CategoryStatusDisabled = 2
)

const (
	InteractionStatusActive   = 1
	InteractionStatusCanceled = 2
)

const (
	CommentStatusPending = 1
	CommentStatusVisible = 2
	CommentStatusClosed  = 3
)

const (
	ShareChannelWechatFriend   = 1
	ShareChannelWechatTimeline = 2
	ShareChannelSystem         = 3
	ShareChannelCopyLink       = 4
)

const (
	// ArticleStatusDraft 表示文章草稿。
	ArticleStatusDraft = 1
	// ArticleStatusPublished 表示文章已发布。
	ArticleStatusPublished = 2
	// ArticleStatusOffline 表示文章已下线。
	ArticleStatusOffline = 3
)

const (
	// ArticleEnableStatusEnabled 表示文章允许在用户端展示。
	ArticleEnableStatusEnabled = 1
	// ArticleEnableStatusDisabled 表示文章不允许在用户端展示。
	ArticleEnableStatusDisabled = 2
)
