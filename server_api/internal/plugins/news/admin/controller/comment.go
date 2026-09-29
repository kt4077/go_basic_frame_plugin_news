package controller

import (
	"github.com/gin-gonic/gin"

	"server_api/internal/plugins/news/admin/logic"
	"server_api/internal/plugins/news/admin/param"
	"server_api/pkg/response"
	requestvalidate "server_api/pkg/validate"
)

// CommentController 新闻评论管理控制器。
type CommentController struct {
	Logic *logic.CommentLogic
}

// List 查询一级评论。
func (h *CommentController) List(c *gin.Context) {
	var req param.CommentListReq
	if err := requestvalidate.Bind(c, &req); err != nil {
		response.Fail(c, response.CodeErrParams, err.Error())
		return
	}
	result, err := h.Logic.List(c, &req, false)
	if err != nil {
		response.Fail(c, response.CodeErrBusiness, err.Error())
		return
	}
	response.OK(c, result)
}

// Replies 查询一级评论下的回复。
func (h *CommentController) Replies(c *gin.Context) {
	var req param.CommentListReq
	if err := requestvalidate.Bind(c, &req); err != nil {
		response.Fail(c, response.CodeErrParams, err.Error())
		return
	}
	result, err := h.Logic.List(c, &req, true)
	if err != nil {
		response.Fail(c, response.CodeErrBusiness, err.Error())
		return
	}
	response.OK(c, result)
}

// UpdateStatus 修改评论显示或关闭状态。
func (h *CommentController) UpdateStatus(c *gin.Context) {
	var req param.CommentStatusReq
	if err := requestvalidate.Bind(c, &req); err != nil {
		response.Fail(c, response.CodeErrParams, err.Error())
		return
	}
	if err := h.Logic.UpdateStatus(c, &req); err != nil {
		response.Fail(c, response.CodeErrBusiness, err.Error())
		return
	}
	response.OK(c, nil)
}
