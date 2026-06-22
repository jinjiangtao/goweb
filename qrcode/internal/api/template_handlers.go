package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"qrcode/internal/model"
	"qrcode/internal/qrcode"
)

type templateReq struct {
	Name        string            `json:"name"`
	StyleConfig model.StyleConfig `json:"styleConfig"`
}

func (svc *Service) listTemplates(c *gin.Context) {
	keyword := c.Query("keyword")
	list, err := svc.Store.ListTemplates(keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Err("查询失败: "+err.Error()))
		return
	}
	if list == nil {
		list = []model.Template{}
	}
	c.JSON(http.StatusOK, model.OK(list))
}

func (svc *Service) createTemplate(c *gin.Context) {
	var req templateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Err("参数错误: "+err.Error()))
		return
	}
	t := &model.Template{Name: req.Name, StyleConfig: req.StyleConfig.WithDefaults()}
	if err := svc.Store.CreateTemplate(t); err != nil {
		c.JSON(http.StatusInternalServerError, model.Err("保存失败: "+err.Error()))
		return
	}
	c.JSON(http.StatusOK, model.OK(t))
}

func (svc *Service) updateTemplate(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Err("无效 ID"))
		return
	}
	var req templateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Err("参数错误: "+err.Error()))
		return
	}
	if err := svc.Store.UpdateTemplate(id, req.Name, req.StyleConfig.WithDefaults()); err != nil {
		c.JSON(http.StatusInternalServerError, model.Err("更新失败: "+err.Error()))
		return
	}
	t, _ := svc.Store.GetTemplate(id)
	c.JSON(http.StatusOK, model.OK(t))
}

func (svc *Service) deleteTemplate(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Err("无效 ID"))
		return
	}
	if err := svc.Store.DeleteTemplate(id); err != nil {
		c.JSON(http.StatusInternalServerError, model.Err("删除失败: "+err.Error()))
		return
	}
	c.JSON(http.StatusOK, model.OK(gin.H{"ok": true}))
}

func (svc *Service) templatePreview(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Err("无效 ID"))
		return
	}
	t, err := svc.Store.GetTemplate(id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.Err("模板不存在"))
		return
	}
	png, err := qrcode.Generate("https://qrforge.dev/template", t.StyleConfig.WithDefaults())
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Err("生成预览失败: "+err.Error()))
		return
	}
	svc.imagePNG(c, png)
}
