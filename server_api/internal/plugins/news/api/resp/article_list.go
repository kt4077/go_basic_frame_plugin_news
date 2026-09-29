package resp

import "time"

// ArticleCategoryItem 文章关联分类响应，关联只在Resp定义。
type ArticleCategoryItem struct {
	ID   uint   `json:"-" gorm:"column:id" comment:"内部分类主键ID"`
	Name string `json:"name" gorm:"column:name" comment:"分类名称"`
	Slug string `json:"slug" gorm:"column:slug" comment:"分类标识"`
}

func (ArticleCategoryItem) TableName() string { return "plg_news_category" }

// ArticleItem 用户端新闻文章响应。
type ArticleItem struct {
	ID               uint                `json:"-" gorm:"column:id" comment:"内部文章主键ID"`
	UID              string              `json:"uid" gorm:"column:uid" comment:"文章业务唯一标识"`
	CategoryID       uint                `json:"-" gorm:"column:category_id" comment:"内部分类主键ID"`
	Category         ArticleCategoryItem `json:"category" gorm:"foreignKey:CategoryID;references:ID" comment:"关联分类"`
	Title            string              `json:"title" gorm:"column:title" comment:"文章标题"`
	Slug             string              `json:"slug" gorm:"column:slug" comment:"文章标识"`
	Summary          string              `json:"summary" gorm:"column:summary" comment:"文章摘要"`
	Author           string              `json:"author" gorm:"column:author" comment:"文章作者"`
	SourceName       string              `json:"source_name" gorm:"column:source_name" comment:"来源平台名称"`
	SourceURL        string              `json:"source_url,omitempty" gorm:"column:source_url" comment:"来源外部链接，列表省略"`
	Cover            string              `json:"cover" gorm:"column:cover" comment:"封面相对路径"`
	CoverURL         string              `json:"cover_url" gorm:"-" comment:"封面完整访问地址"`
	Content          string              `json:"content,omitempty" gorm:"column:content" comment:"文章正文"`
	ViewCount        uint64              `json:"view_count" gorm:"column:view_count" comment:"真实浏览次数"`
	VirtualViewCount uint64              `json:"virtual_view_count" gorm:"column:virtual_view_count" comment:"虚拟浏览次数"`
	TotalViewCount   uint64              `json:"total_view_count" gorm:"-" comment:"展示浏览次数，真实与虚拟之和"`
	LikeCount        uint64              `json:"like_count" gorm:"column:like_count" comment:"点赞数量"`
	ShareCount       uint64              `json:"share_count" gorm:"column:share_count" comment:"发起分享数量"`
	CollectionCount  uint64              `json:"collection_count" gorm:"column:collection_count" comment:"收藏数量"`
	CommentCount     uint64              `json:"comment_count" gorm:"column:comment_count" comment:"已展示一级评论数量"`
	PublishedAt      *time.Time          `json:"published_at" gorm:"column:published_at" comment:"发布时间"`
}

func (ArticleItem) TableName() string { return "plg_news_article" }

// ArticleListRes 用户端新闻文章分页响应。
type ArticleListRes struct {
	List  []ArticleItem `json:"list" comment:"文章列表"`
	Total int64         `json:"total" comment:"数据总数"`
}
