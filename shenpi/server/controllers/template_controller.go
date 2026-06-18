package controllers

import (
	"encoding/json"
	"net/http"
	"shenpi/models"
	"shenpi/services"
	"shenpi/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TemplateController struct {
	templateService *services.TemplateService
}

func NewTemplateController() *TemplateController {
	return &TemplateController{
		templateService: services.NewTemplateService(models.DB),
	}
}

func (c *TemplateController) Create(ctx *gin.Context) {
	var req struct {
		Name        string          `json:"name" binding:"required"`
		Description string          `json:"description"`
		Type        string          `json:"type"`
		FormConfig  json.RawMessage `json:"form_config"`
		Nodes       json.RawMessage `json:"nodes"`
		Edges       json.RawMessage `json:"edges"`
		Status      string          `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	userID := ctx.GetUint("user_id")

	template := &models.ProcessTemplate{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		FormConfig:  string(req.FormConfig),
		Nodes:       string(req.Nodes),
		Edges:       string(req.Edges),
		Status:      req.Status,
		CreatorID:   userID,
	}

	if err := c.templateService.CreateTemplate(template); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "创建模板失败: "+err.Error())
		return
	}

	utils.Success(ctx, template)
}

func (c *TemplateController) Get(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	template, err := c.templateService.GetTemplateByID(uint(id))
	if err != nil {
		utils.Error(ctx, http.StatusNotFound, "模板不存在")
		return
	}

	var nodes interface{}
	if template.Nodes != "" {
		json.Unmarshal([]byte(template.Nodes), &nodes)
	}

	var edges interface{}
	if template.Edges != "" {
		json.Unmarshal([]byte(template.Edges), &edges)
	}

	var formConfig interface{}
	if template.FormConfig != "" {
		json.Unmarshal([]byte(template.FormConfig), &formConfig)
	}

	utils.Success(ctx, gin.H{
		"id":          template.ID,
		"name":        template.Name,
		"description": template.Description,
		"type":        template.Type,
		"form_config": formConfig,
		"nodes":       nodes,
		"edges":       edges,
		"status":      template.Status,
		"creator_id":  template.CreatorID,
		"creator":     template.Creator,
		"created_at":  template.CreatedAt,
		"updated_at":  template.UpdatedAt,
	})
}

func (c *TemplateController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	keyword := ctx.Query("keyword")
	templateType := ctx.Query("type")
	creatorID, _ := strconv.ParseUint(ctx.Query("creator_id"), 10, 32)

	templates, total, err := c.templateService.GetTemplateList(page, pageSize, keyword, templateType, uint(creatorID))
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取模板列表失败: "+err.Error())
		return
	}

	result := make([]gin.H, 0, len(templates))
	for _, t := range templates {
		var nodes interface{}
		if t.Nodes != "" {
			json.Unmarshal([]byte(t.Nodes), &nodes)
		}
		var edges interface{}
		if t.Edges != "" {
			json.Unmarshal([]byte(t.Edges), &edges)
		}
		var formConfig interface{}
		if t.FormConfig != "" {
			json.Unmarshal([]byte(t.FormConfig), &formConfig)
		}
		result = append(result, gin.H{
			"id":          t.ID,
			"name":        t.Name,
			"description": t.Description,
			"type":        t.Type,
			"form_config": formConfig,
			"nodes":       nodes,
			"edges":       edges,
			"status":      t.Status,
			"creator_id":  t.CreatorID,
			"creator":     t.Creator,
			"created_at":  t.CreatedAt,
			"updated_at":  t.UpdatedAt,
		})
	}

	utils.Success(ctx, gin.H{
		"list":     result,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (c *TemplateController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	template, err := c.templateService.GetTemplateByID(uint(id))
	if err != nil {
		utils.Error(ctx, http.StatusNotFound, "模板不存在")
		return
	}

	var req struct {
		Name        string          `json:"name"`
		Description string          `json:"description"`
		Type        string          `json:"type"`
		FormConfig  json.RawMessage `json:"form_config"`
		Nodes       json.RawMessage `json:"nodes"`
		Edges       json.RawMessage `json:"edges"`
		Status      string          `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if req.Name != "" {
		template.Name = req.Name
	}
	if req.Description != "" {
		template.Description = req.Description
	}
	if req.Type != "" {
		template.Type = req.Type
	}
	if string(req.FormConfig) != "" {
		template.FormConfig = string(req.FormConfig)
	}
	if string(req.Nodes) != "" {
		template.Nodes = string(req.Nodes)
	}
	if string(req.Edges) != "" {
		template.Edges = string(req.Edges)
	}
	if req.Status != "" {
		template.Status = req.Status
	}

	if err := c.templateService.UpdateTemplate(template); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "更新模板失败: "+err.Error())
		return
	}

	utils.Success(ctx, template)
}

func (c *TemplateController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	if err := c.templateService.DeleteTemplate(uint(id)); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "删除模板失败: "+err.Error())
		return
	}

	utils.Success(ctx, nil)
}

func (c *TemplateController) GetFormFields(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	template, err := c.templateService.GetTemplateByID(uint(id))
	if err != nil {
		utils.Error(ctx, http.StatusNotFound, "模板不存在")
		return
	}

	fields, err := c.templateService.GetTemplateFormFields(template)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取表单字段失败: "+err.Error())
		return
	}

	utils.Success(ctx, fields)
}

func (c *TemplateController) CreateDefault(ctx *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
		Type string `json:"type" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	userID := ctx.GetUint("user_id")
	template := c.templateService.BuildDefaultTemplate(req.Name, req.Type, userID)

	if err := c.templateService.CreateTemplate(template); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "创建默认模板失败: "+err.Error())
		return
	}

	utils.Success(ctx, template)
}
