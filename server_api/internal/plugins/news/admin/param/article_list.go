package param

// ArticleListReq 新闻文章列表请求。
type ArticleListReq struct {
	Keyword    string `form:"keyword" binding:"omitempty,max=200" comment:"标题或标识关键字"`
	CategoryID uint   `form:"category_id" comment:"分类ID"`
	Status     int    `form:"status" binding:"omitempty,oneof=1 2 3" comment:"状态，1草稿，2已发布，3已下线"`
	Page       int    `form:"page" comment:"页码"`
	PageSize   int    `form:"page_size" comment:"每页数量"`
}
