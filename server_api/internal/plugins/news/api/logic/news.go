package logic

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"server_api/internal/common/app"
	commonmiddleware "server_api/internal/common/middleware"
	commonupload "server_api/internal/common/upload"
	"server_api/internal/plugins/news/api/param"
	"server_api/internal/plugins/news/api/resp"
	newsEnums "server_api/internal/plugins/news/enums"
	"server_api/internal/plugins/news/model"
	"server_api/pkg/pagination"
)

// NewsLogic 用户端新闻业务逻辑。
type NewsLogic struct{ App *app.App }

// Advertisements 返回当前用户端平台、当前新闻展示位置可用的启用广告。
func (l *NewsLogic) Advertisements(c *gin.Context, req *param.AdvertisementListReq) ([]resp.AdvertisementItem, error) {
	platform := strconv.Itoa(commonmiddleware.CtxPlatformSource(c))
	var list []resp.AdvertisementItem
	err := l.App.DB.WithContext(c.Request.Context()).Table("plg_news_advertisement_config c").
		Select("a.id,a.name,a.ad_id,a.format,a.description").
		Joins("JOIN sys_advertisement a ON a.id=c.advertisement_id AND a.deleted_at IS NULL").
		Joins("JOIN sys_advertisement_plugin ap ON ap.advertisement_id=a.id").
		Joins("JOIN sys_plugin p ON p.id=ap.plugin_id AND p.plugin_id=? AND p.status=1 AND p.deleted_at IS NULL", "news").
		Where("c.position=? AND c.status=1 AND a.status=1 AND TRIM(a.ad_id)<>'' AND FIND_IN_SET(?,a.platforms)", req.Position, platform).
		Order("c.id ASC").Scan(&list).Error
	if err != nil {
		return nil, errors.New("查询新闻广告失败")
	}
	return list, nil
}

// Categories 查询启用的新闻分类。
func (l *NewsLogic) Categories(c *gin.Context) ([]resp.CategoryItem, error) {
	var list []resp.CategoryItem
	if err := l.App.DB.WithContext(c.Request.Context()).Model(&model.Category{}).Where("status = ?", newsEnums.CategoryStatusEnabled).Order("sort DESC,id ASC").Find(&list).Error; err != nil {
		return nil, errors.New("查询分类失败")
	}
	return list, nil
}

// Articles 查询已发布新闻，分类条件使用稳定标识。
func (l *NewsLogic) Articles(c *gin.Context, req *param.ArticleListReq) (*resp.ArticleListRes, error) {
	var total int64
	var list []resp.ArticleItem
	db := l.App.DB.WithContext(c.Request.Context()).Model(&model.Article{}).
		Where("plg_news_article.status = ? AND plg_news_article.enable_status = ?", newsEnums.ArticleStatusPublished, newsEnums.ArticleEnableStatusEnabled)
	categorySlug := normalizeOptionalQuery(req.CategorySlug)
	if categorySlug != "" {
		db = db.Joins("JOIN plg_news_category ON plg_news_category.id = plg_news_article.category_id AND plg_news_category.deleted_at IS NULL").Where("plg_news_category.slug = ? AND plg_news_category.status = ?", categorySlug, newsEnums.CategoryStatusEnabled)
	}
	if keyword := normalizeOptionalQuery(req.Keyword); keyword != "" {
		db = db.Where("plg_news_article.title LIKE ?", "%"+keyword+"%")
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, errors.New("查询新闻失败")
	}
	page, size := pagination.Normalize(req.Page, req.PageSize)
	if err := db.Preload("Category").Select("plg_news_article.*").Omit("content", "source_url").Order("sort DESC,published_at DESC,id DESC").Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, errors.New("查询新闻失败")
	}
	if err := completeCoverURLs(l.App, list); err != nil {
		return nil, errors.New("封面地址解析失败")
	}
	if err := l.completeArticleCommentCounts(c, list); err != nil {
		return nil, err
	}
	return &resp.ArticleListRes{List: list, Total: total}, nil
}

func normalizeOptionalQuery(value string) string {
	value = strings.TrimSpace(value)
	if strings.EqualFold(value, "undefined") || strings.EqualFold(value, "null") {
		return ""
	}
	return value
}

// Article 查询已发布文章详情，并以原子更新累加浏览量。
func (l *NewsLogic) Article(c *gin.Context, req *param.ArticleDetailReq) (*resp.ArticleItem, error) {
	var item resp.ArticleItem
	err := l.App.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Preload("Category").Where("slug = ? AND status = ? AND enable_status = ?", req.Slug, newsEnums.ArticleStatusPublished, newsEnums.ArticleEnableStatusEnabled).First(&item).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.Article{}).Where("uid = ?", item.UID).UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error; err != nil {
			return err
		}
		item.ViewCount++
		item.TotalViewCount = item.ViewCount + item.VirtualViewCount
		return nil
	})
	if err != nil {
		return nil, errors.New("新闻不存在")
	}
	if err := completeCoverURL(l.App, &item); err != nil {
		return nil, errors.New("封面地址解析失败")
	}
	items := []resp.ArticleItem{item}
	if err := l.completeArticleCommentCounts(c, items); err != nil {
		return nil, err
	}
	item.CommentCount = items[0].CommentCount
	return &item, nil
}

// completeArticleCommentCounts 批量统计文章全部层级的已展示评论数。
func (l *NewsLogic) completeArticleCommentCounts(c *gin.Context, list []resp.ArticleItem) error {
	if len(list) == 0 {
		return nil
	}
	uids := make([]string, 0, len(list))
	for _, item := range list {
		uids = append(uids, item.UID)
	}
	type articleCount struct {
		ArticleUID string `gorm:"column:article_uid"`
		Total      uint64 `gorm:"column:total"`
	}
	var counts []articleCount
	if err := l.App.DB.WithContext(c.Request.Context()).Model(&model.Comment{}).
		Select("article_uid, COUNT(*) AS total").
		Where("article_uid IN ? AND status = ?", uids, newsEnums.CommentStatusVisible).
		Group("article_uid").Find(&counts).Error; err != nil {
		return errors.New("查询文章评论数量失败")
	}
	countMap := make(map[string]uint64, len(counts))
	for _, count := range counts {
		countMap[count.ArticleUID] = count.Total
	}
	for index := range list {
		list[index].CommentCount = countMap[list[index].UID]
	}
	return nil
}

func completeCoverURLs(application *app.App, list []resp.ArticleItem) error {
	for i := range list {
		if err := completeCoverURL(application, &list[i]); err != nil {
			return err
		}
	}
	return nil
}
func completeCoverURL(application *app.App, item *resp.ArticleItem) error {
	url, err := commonupload.FileURL(application, item.Cover)
	if err != nil {
		return err
	}
	item.CoverURL = url
	item.TotalViewCount = item.ViewCount + item.VirtualViewCount
	return nil
}
