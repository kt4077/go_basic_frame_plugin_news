package controller

import (
	"github.com/gin-gonic/gin"
	"server_api/internal/plugins/news/api/logic"
	"server_api/internal/plugins/news/api/param"
	"server_api/pkg/response"
	requestvalidate "server_api/pkg/validate"
)

// NewsController 用户端新闻控制器。
type NewsController struct{ Logic *logic.NewsLogic }

func (h *NewsController) Advertisements(c *gin.Context) {
	var req param.AdvertisementListReq
	if err := requestvalidate.Bind(c, &req); err != nil {
		response.Fail(c, response.CodeErrParams, err.Error())
		return
	}
	result, err := h.Logic.Advertisements(c, &req)
	if err != nil {
		response.Fail(c, response.CodeErrBusiness, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *NewsController) Categories(c *gin.Context) {
	result, err := h.Logic.Categories(c)
	if err != nil {
		response.Fail(c, response.CodeErrBusiness, err.Error())
		return
	}
	response.OK(c, result)
}
func (h *NewsController) Articles(c *gin.Context) {
	var req param.ArticleListReq
	if err := requestvalidate.Bind(c, &req); err != nil {
		response.Fail(c, response.CodeErrParams, err.Error())
		return
	}
	result, err := h.Logic.Articles(c, &req)
	if err != nil {
		response.Fail(c, response.CodeErrBusiness, err.Error())
		return
	}
	response.OK(c, result)
}
func (h *NewsController) Article(c *gin.Context) {
	var req param.ArticleDetailReq
	if err := requestvalidate.Bind(c, &req); err != nil {
		response.Fail(c, response.CodeErrParams, err.Error())
		return
	}
	result, err := h.Logic.Article(c, &req)
	if err != nil {
		response.Fail(c, response.CodeErrBusiness, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *NewsController) Comments(c *gin.Context) {
	var req struct {
		param.ArticleUIDReq
		param.CommentListReq
	}
	if err := requestvalidate.Bind(c, &req); err != nil {
		response.Fail(c, response.CodeErrParams, err.Error())
		return
	}
	result, err := h.Logic.Comments(c, req.ArticleUID, &req.CommentListReq)
	if err != nil {
		response.Fail(c, response.CodeErrBusiness, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *NewsController) Replies(c *gin.Context) {
	var req struct {
		param.CommentUIDReq
		param.CommentListReq
	}
	if err := requestvalidate.Bind(c, &req); err != nil {
		response.Fail(c, response.CodeErrParams, err.Error())
		return
	}
	result, err := h.Logic.Replies(c, req.CommentUID, &req.CommentListReq)
	if err != nil {
		response.Fail(c, response.CodeErrBusiness, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *NewsController) CreateComment(c *gin.Context) {
	var req param.ArticleCommentCreateReq
	if err := requestvalidate.Bind(c, &req); err != nil {
		response.Fail(c, response.CodeErrParams, err.Error())
		return
	}
	result, err := h.Logic.CreateComment(c, req.ArticleUID, "", &req.CommentCreateReq)
	if err != nil {
		response.Fail(c, response.CodeErrBusiness, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *NewsController) ReplyComment(c *gin.Context) {
	var req param.CommentReplyCreateReq
	if err := requestvalidate.Bind(c, &req); err != nil {
		response.Fail(c, response.CodeErrParams, err.Error())
		return
	}
	result, err := h.Logic.ReplyComment(c, req.CommentUID, &req.CommentCreateReq)
	if err != nil {
		response.Fail(c, response.CodeErrBusiness, err.Error())
		return
	}
	response.OK(c, result)
}

// ToggleArticleLike 原子切换当前用户的文章点赞状态。
func (h *NewsController) ToggleArticleLike(c *gin.Context) {
	var req param.ArticleUIDReq
	if err := requestvalidate.Bind(c, &req); err != nil {
		response.Fail(c, response.CodeErrParams, err.Error())
		return
	}
	result, err := h.Logic.ToggleArticleLike(c, req.ArticleUID)
	if err != nil {
		response.Fail(c, response.CodeErrBusiness, err.Error())
		return
	}
	response.OK(c, result)
}

// ToggleArticleCollection 原子切换当前用户的文章收藏状态。
func (h *NewsController) ToggleArticleCollection(c *gin.Context) {
	var req param.ArticleUIDReq
	if err := requestvalidate.Bind(c, &req); err != nil {
		response.Fail(c, response.CodeErrParams, err.Error())
		return
	}
	result, err := h.Logic.ToggleArticleCollection(c, req.ArticleUID)
	if err != nil {
		response.Fail(c, response.CodeErrBusiness, err.Error())
		return
	}
	response.OK(c, result)
}

// ToggleCommentLike 原子切换当前用户的评论点赞状态。
func (h *NewsController) ToggleCommentLike(c *gin.Context) {
	var req param.CommentUIDReq
	if err := requestvalidate.Bind(c, &req); err != nil {
		response.Fail(c, response.CodeErrParams, err.Error())
		return
	}
	result, err := h.Logic.ToggleCommentLike(c, req.CommentUID)
	if err != nil {
		response.Fail(c, response.CodeErrBusiness, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *NewsController) Share(c *gin.Context) {
	var req param.ShareReq
	if err := requestvalidate.Bind(c, &req); err != nil {
		response.Fail(c, response.CodeErrParams, err.Error())
		return
	}
	if err := h.Logic.RecordShare(c, req.ArticleUID, &req); err != nil {
		response.Fail(c, response.CodeErrBusiness, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *NewsController) InteractionState(c *gin.Context) {
	var req param.ArticleUIDReq
	if err := requestvalidate.Bind(c, &req); err != nil {
		response.Fail(c, response.CodeErrParams, err.Error())
		return
	}
	result, err := h.Logic.InteractionState(c, req.ArticleUID)
	if err != nil {
		response.Fail(c, response.CodeErrBusiness, err.Error())
		return
	}
	response.OK(c, result)
}
