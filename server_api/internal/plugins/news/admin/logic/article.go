package logic

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"server_api/internal/common/app"
	commonupload "server_api/internal/common/upload"
	"server_api/internal/plugins/news/admin/param"
	"server_api/internal/plugins/news/admin/resp"
	newsEnums "server_api/internal/plugins/news/enums"
	"server_api/internal/plugins/news/model"
	"server_api/pkg/dberror"
	"server_api/pkg/pagination"
	"server_api/pkg/sn"
)

var articleSlugPattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*$`)

// ArticleLogic 新闻文章业务逻辑。
type ArticleLogic struct{ App *app.App }

// List 查询新闻文章分页列表并预加载分类。
func (l *ArticleLogic) List(c *gin.Context, req *param.ArticleListReq) (*resp.ArticleListRes, error) {
	var total int64
	var list []resp.ArticleItem
	db := l.App.DB.WithContext(c.Request.Context()).Model(&model.Article{})
	if keyword := strings.TrimSpace(req.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		db = db.Where("title LIKE ? OR slug LIKE ? OR author LIKE ? OR source_name LIKE ?", like, like, like, like)
	}
	if req.CategoryID != 0 {
		db = db.Where("category_id = ?", req.CategoryID)
	}
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}
	if req.EnableStatus != 0 {
		db = db.Where("enable_status = ?", req.EnableStatus)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, errors.New("查询文章失败")
	}
	page, size := pagination.Normalize(req.Page, req.PageSize)
	orderClause := "sort DESC, id DESC"
	if req.OrderBy != "" {
		orderField := req.OrderBy
		if orderField == "total_view_count" {
			orderField = "(view_count + virtual_view_count)"
		}
		if orderField == "comment_count" {
			orderField = "(SELECT COUNT(*) FROM plg_news_comment AS comment_stats WHERE comment_stats.article_uid = plg_news_article.uid AND comment_stats.status = 2 AND comment_stats.deleted_at IS NULL)"
		}
		orderDirection := "DESC"
		if strings.EqualFold(req.Order, "asc") {
			orderDirection = "ASC"
		}
		orderClause = orderField + " " + orderDirection + ", id DESC"
	}
	if err := db.Preload("Category").Order(orderClause).Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, errors.New("查询文章失败")
	}
	if err := completeCoverURLs(l.App, list); err != nil {
		return nil, errors.New("封面地址解析失败")
	}
	if err := l.completeArticleCommentCounts(c, list); err != nil {
		return nil, err
	}
	return &resp.ArticleListRes{List: list, Total: total}, nil
}

// completeArticleCommentCounts 批量统计文章全部层级的已展示评论数。
func (l *ArticleLogic) completeArticleCommentCounts(c *gin.Context, list []resp.ArticleItem) error {
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

// Detail 查询新闻文章详情。
func (l *ArticleLogic) Detail(c *gin.Context, req *param.ArticleIDReq) (*resp.ArticleItem, error) {
	var item resp.ArticleItem
	if err := l.App.DB.WithContext(c.Request.Context()).Preload("Category").First(&item, req.ID).Error; err != nil {
		return nil, errors.New("文章不存在")
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

// Save 新增或修改新闻文章，封面只保存相对路径。
func (l *ArticleLogic) Save(c *gin.Context, req *param.ArticleSaveReq) (*resp.ArticleItem, error) {
	if !articleSlugPattern.MatchString(strings.TrimSpace(req.Slug)) {
		return nil, errors.New("文章标识仅支持字母、数字、中划线和下划线")
	}
	sourceURL := strings.TrimSpace(req.SourceURL)
	if sourceURL != "" {
		parsedSourceURL, parseErr := url.ParseRequestURI(sourceURL)
		if parseErr != nil || parsedSourceURL.Host == "" || (parsedSourceURL.Scheme != "http" && parsedSourceURL.Scheme != "https") {
			return nil, errors.New("来源外部链接仅支持有效的HTTP或HTTPS地址")
		}
	}
	resolver, err := commonupload.NewURLResolver(l.App)
	if err != nil && strings.TrimSpace(req.Cover) != "" {
		return nil, errors.New("存储配置不可用")
	}
	cover := ""
	if resolver != nil {
		cover, err = resolver.Relative(req.Cover)
	}
	if err != nil {
		return nil, err
	}
	article := model.Article{CategoryID: req.CategoryID, Title: strings.TrimSpace(req.Title), Slug: strings.ToLower(strings.TrimSpace(req.Slug)), Summary: strings.TrimSpace(req.Summary), Author: strings.TrimSpace(req.Author), SourceName: strings.TrimSpace(req.SourceName), SourceURL: sourceURL, Cover: cover, Content: req.Content, Sort: req.Sort, Status: newsEnums.ArticleStatusDraft, EnableStatus: req.EnableStatus, VirtualViewCount: req.VirtualViewCount}
	if req.ID == 0 {
		article.UID, err = sn.Generate()
		if err != nil {
			return nil, errors.New("文章标识生成失败")
		}
	}
	err = l.App.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var categoryCount int64
		if err := tx.Model(&model.Category{}).Where("id = ?", req.CategoryID).Count(&categoryCount).Error; err != nil {
			return err
		}
		if categoryCount == 0 {
			return errors.New("请选择有效分类")
		}
		if req.ID == 0 {
			return tx.Create(&article).Error
		}
		var current model.Article
		if err := tx.First(&current, req.ID).Error; err != nil {
			return err
		}
		return tx.Model(&current).Updates(map[string]interface{}{"category_id": article.CategoryID, "title": article.Title, "summary": article.Summary, "author": article.Author, "source_name": article.SourceName, "source_url": article.SourceURL, "cover": article.Cover, "content": article.Content, "sort": article.Sort, "enable_status": article.EnableStatus, "virtual_view_count": article.VirtualViewCount}).Error
	})
	if err != nil {
		if dberror.IsDuplicateKey(err) {
			return nil, errors.New("文章标识已存在")
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在")
		}
		return nil, err
	}
	id := article.ID
	if req.ID != 0 {
		id = req.ID
	}
	return l.Detail(c, &param.ArticleIDReq{ID: id})
}

// UpdateStatus 发布或下线新闻文章。
func (l *ArticleLogic) UpdateStatus(c *gin.Context, req *param.ArticleStatusReq) error {
	updates := map[string]interface{}{"status": req.Status}
	if req.Status == newsEnums.ArticleStatusPublished {
		now := time.Now()
		updates["published_at"] = &now
	}
	result := l.App.DB.WithContext(c.Request.Context()).Model(&model.Article{}).Where("id = ?", req.ID).Updates(updates)
	if result.Error != nil {
		return errors.New("更新文章状态失败")
	}
	if result.RowsAffected == 0 {
		return errors.New("文章不存在")
	}
	return nil
}

// Delete 软删除新闻文章。
func (l *ArticleLogic) Delete(c *gin.Context, req *param.ArticleIDReq) error {
	result := l.App.DB.WithContext(c.Request.Context()).Delete(&model.Article{}, req.ID)
	if result.Error != nil {
		return errors.New("删除文章失败")
	}
	if result.RowsAffected == 0 {
		return errors.New("文章不存在")
	}
	return nil
}

func completeCoverURLs(application *app.App, list []resp.ArticleItem) error {
	for index := range list {
		if err := completeCoverURL(application, &list[index]); err != nil {
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
