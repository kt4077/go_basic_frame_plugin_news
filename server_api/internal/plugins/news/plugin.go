// Package news 提供新闻分类、文章发布和用户端新闻查询能力。
package news

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"server_api/internal/common/app"
	commonplugin "server_api/internal/common/plugin"
	adminController "server_api/internal/plugins/news/admin/controller"
	adminLogic "server_api/internal/plugins/news/admin/logic"
	apiController "server_api/internal/plugins/news/api/controller"
	apiLogic "server_api/internal/plugins/news/api/logic"
	"server_api/internal/plugins/news/model"
	"server_api/pkg/dberror"
	"server_api/pkg/sn"
)

const (
	migrationV001Checksum = "da37268cc9545f2811aac01aedd79dd5b878df6f061bb528d902475b0870ffd3"
	migrationV002Checksum = "cb01d837c81e556821e894e855954a83d398240df459907120744d7879b8d9ae"
	migrationV010Checksum = "74834583dd99da27ca555859872b876215638ae54ef2a29c60ace707893edc85"
	migrationV011Checksum = "f238f98bd859c00f96ae712c52bf88cbc701f331d16c10f8e26bc17af80f069c"
	migrationV012Checksum = "fe840845f712730ac9caf4b276dc462bb5b38d3ec2faa65a3a27b570de3361d4"
	migrationV013Checksum = "19091a09a63ac207ca1796c1e9585747274cf3365ee5a00d94105f401b5c4174"
)

// Plugin 新闻资讯插件。
type Plugin struct {
	application *app.App
}

// New 创建新闻资讯插件。
func New() *Plugin { return &Plugin{} }

// Manifest 返回插件元信息。
func (p *Plugin) Manifest() commonplugin.Manifest {
	return commonplugin.Manifest{PluginID: "news", Name: "新闻资讯", Version: "0.1.4", Logo: "", Author: "kt4077", Homepage: "https://github.com/kt4077/go_basic_frame_plugin_news", Description: "提供新闻分类、文章发布、评论、点赞、收藏、分享及多端内容展示能力。", CoreVersion: commonplugin.CoreVersion, Dependencies: []commonplugin.Dependency{}}
}

// Migrations 返回插件迁移声明。
func (p *Plugin) Migrations() []commonplugin.Migration {
	return []commonplugin.Migration{
		{Version: "0.0.1", Description: "创建新闻分类与文章完整基线结构", Checksum: migrationV001Checksum},
		{Version: "0.0.2", Description: "增加文章作者、来源平台与来源外链字段", Checksum: migrationV002Checksum},
		{Version: "0.1.0", Description: "增加文章业务标识、评论、点赞与分享数据结构", Checksum: migrationV010Checksum},
		{Version: "0.1.1", Description: "增加文章收藏与互动统计排序能力", Checksum: migrationV011Checksum},
		{Version: "0.1.2", Description: "增加新闻列表与详情广告展示配置", Checksum: migrationV012Checksum},
		{Version: "0.1.3", Description: "增加新闻广告独立投放状态", Checksum: migrationV013Checksum},
	}
}

// RegisterAdminRoutes 注册受鉴权和权限校验保护的管理端路由。
func (p *Plugin) RegisterAdminRoutes(ctx *commonplugin.Context, groups commonplugin.AdminRouteGroups) error {
	p.application = ctx.App
	category := &adminController.CategoryController{Logic: &adminLogic.CategoryLogic{App: ctx.App}}
	article := &adminController.ArticleController{Logic: &adminLogic.ArticleLogic{App: ctx.App}}
	comment := &adminController.CommentController{Logic: &adminLogic.CommentLogic{App: ctx.App}}
	advertisement := &adminController.AdvertisementController{Logic: &adminLogic.AdvertisementLogic{App: ctx.App}}
	group := groups.Permission.Group("/plugin/news")
	group.GET("/category/list", category.List)
	group.POST("/category/save", category.Save)
	group.POST("/category/delete", category.Delete)
	group.GET("/article/list", article.List)
	group.GET("/article/detail", article.Detail)
	group.POST("/article/save", article.Save)
	group.POST("/article/status", article.UpdateStatus)
	group.POST("/article/delete", article.Delete)
	group.GET("/article/comment/list", comment.List)
	group.GET("/article/comment/replies", comment.Replies)
	group.POST("/article/comment/status", comment.UpdateStatus)
	group.GET("/advertisement/list", advertisement.List)
	group.GET("/advertisement/options", advertisement.Options)
	group.POST("/advertisement/save", advertisement.Save)
	group.POST("/advertisement/delete", advertisement.Delete)
	return nil
}

// RegisterAPIRoutes 注册用户端公开新闻路由。
func (p *Plugin) RegisterAPIRoutes(ctx *commonplugin.Context, groups commonplugin.APIRouteGroups) error {
	p.application = ctx.App
	controller := &apiController.NewsController{Logic: &apiLogic.NewsLogic{App: ctx.App}}
	publicGroup := groups.Public.Group("/plugin/news")
	publicGroup.GET("/categories", controller.Categories)
	publicGroup.GET("/articles", controller.Articles)
	publicGroup.GET("/article/detail", controller.Article)
	publicGroup.GET("/article/comments", controller.Comments)
	publicGroup.GET("/advertisements", controller.Advertisements)
	publicGroup.GET("/comment/replies", controller.Replies)
	authGroup := groups.Auth.Group("/plugin/news")
	authGroup.POST("/article/comment", controller.CreateComment)
	authGroup.POST("/comment/reply", controller.ReplyComment)
	authGroup.POST("/article/like", controller.ToggleArticleLike)
	authGroup.POST("/article/collection", controller.ToggleArticleCollection)
	authGroup.POST("/comment/like", controller.ToggleCommentLike)
	authGroup.POST("/article/share", controller.Share)
	authGroup.GET("/article/interaction", controller.InteractionState)
	return nil
}

// Start 启动插件并为历史文章补齐业务唯一标识。
func (p *Plugin) Start(ctx context.Context, service commonplugin.ServiceType) error {
	if p.application == nil {
		return errors.New("新闻插件运行时上下文未初始化")
	}
	for {
		var articles []model.Article
		err := p.application.DB.WithContext(ctx).Where("uid IS NULL OR uid = ''").Order("id ASC").Limit(100).Find(&articles).Error
		if err != nil {
			return errors.New("查询待补齐文章业务标识失败")
		}
		if len(articles) == 0 {
			return nil
		}
		for _, article := range articles {
			if err := backfillArticleUID(ctx, p.application, article.ID); err != nil {
				return err
			}
		}
	}
}

func backfillArticleUID(ctx context.Context, application *app.App, articleID uint) error {
	for attempt := 0; attempt < 5; attempt++ {
		uid, err := sn.Generate()
		if err != nil {
			return errors.New("生成文章业务标识失败")
		}
		result := application.DB.WithContext(ctx).Model(&model.Article{}).
			Where("id = ? AND (uid IS NULL OR uid = '')", articleID).Update("uid", uid)
		if result.Error == nil {
			return nil
		}
		if !dberror.IsDuplicateKey(result.Error) {
			return errors.New("补齐文章业务标识失败")
		}
	}
	return gorm.ErrInvalidData
}

// Stop 停止插件，本版本没有后台任务。
func (p *Plugin) Stop(ctx context.Context) error { return nil }
