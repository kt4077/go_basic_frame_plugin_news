package controller

import (
	"github.com/gin-gonic/gin"
	newslogic "server_api/internal/plugins/news/admin/logic"
	newsparam "server_api/internal/plugins/news/admin/param"
	"server_api/pkg/response"
	requestvalidate "server_api/pkg/validate"
)

type AdvertisementController struct{ Logic *newslogic.AdvertisementLogic }

func (h *AdvertisementController) List(c *gin.Context) {
	result, err := h.Logic.List(c)
	if err != nil {
		response.Fail(c, response.CodeErrBusiness, err.Error())
		return
	}
	response.OK(c, result)
}
func (h *AdvertisementController) Options(c *gin.Context) {
	result, err := h.Logic.Options(c)
	if err != nil {
		response.Fail(c, response.CodeErrBusiness, err.Error())
		return
	}
	response.OK(c, result)
}
func (h *AdvertisementController) Save(c *gin.Context) {
	var req newsparam.AdvertisementSaveReq
	if err := requestvalidate.Bind(c, &req); err != nil {
		response.Fail(c, response.CodeErrParams, err.Error())
		return
	}
	if err := h.Logic.Save(c, &req); err != nil {
		response.Fail(c, response.CodeErrBusiness, err.Error())
		return
	}
	response.OK(c, nil)
}
func (h *AdvertisementController) Delete(c *gin.Context) {
	var req newsparam.IDReq
	if err := requestvalidate.Bind(c, &req); err != nil {
		response.Fail(c, response.CodeErrParams, err.Error())
		return
	}
	if err := h.Logic.Delete(c, req.ID); err != nil {
		response.Fail(c, response.CodeErrBusiness, err.Error())
		return
	}
	response.OK(c, nil)
}
