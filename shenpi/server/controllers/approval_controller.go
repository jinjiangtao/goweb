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

type ApprovalController struct {
	approvalService *services.ApprovalService
}

func NewApprovalController() *ApprovalController {
	return &ApprovalController{
		approvalService: services.NewApprovalService(models.DB),
	}
}

func (c *ApprovalController) Submit(ctx *gin.Context) {
	var req struct {
		TemplateID uint                   `json:"template_id" binding:"required"`
		Title      string                 `json:"title" binding:"required"`
		FormData   map[string]interface{} `json:"form_data" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	userID := ctx.GetUint("user_id")

	request, err := c.approvalService.SubmitApproval(userID, req.TemplateID, req.Title, req.FormData)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "提交审批失败: "+err.Error())
		return
	}

	utils.Success(ctx, request)
}

func (c *ApprovalController) Get(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	request, err := c.approvalService.GetRequestByID(uint(id))
	if err != nil {
		utils.Error(ctx, http.StatusNotFound, "审批单不存在")
		return
	}

	var formData map[string]interface{}
	json.Unmarshal([]byte(request.FormData), &formData)

	var instance *models.ProcessInstance
	if request.ProcessInstance != "" {
		instance, _ = c.approvalService.GetProcessInstance(request)
	}

	utils.Success(ctx, gin.H{
		"request":    request,
		"form_data":  formData,
		"instance":   instance,
	})
}

func (c *ApprovalController) MyRequests(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	status := ctx.Query("status")

	userID := ctx.GetUint("user_id")

	requests, total, err := c.approvalService.GetMyRequests(userID, page, pageSize, status)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取我的申请失败: "+err.Error())
		return
	}

	utils.Success(ctx, gin.H{
		"list":     requests,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (c *ApprovalController) MyApprovals(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	status := ctx.Query("status")

	userID := ctx.GetUint("user_id")

	requests, total, err := c.approvalService.GetMyApprovals(userID, page, pageSize, status)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取我的审批失败: "+err.Error())
		return
	}

	utils.Success(ctx, gin.H{
		"list":     requests,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (c *ApprovalController) Records(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	records, err := c.approvalService.GetApprovalRecords(uint(id))
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取审批记录失败: "+err.Error())
		return
	}

	utils.Success(ctx, records)
}

func (c *ApprovalController) Approve(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	var req struct {
		Comment string `json:"comment"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	userID := ctx.GetUint("user_id")

	request, err := c.approvalService.ProcessApproval(uint(id), userID, models.ApprovalActionApprove, req.Comment)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "审批失败: "+err.Error())
		return
	}

	utils.Success(ctx, request)
}

func (c *ApprovalController) Reject(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	var req struct {
		Comment string `json:"comment" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	userID := ctx.GetUint("user_id")

	request, err := c.approvalService.ProcessApproval(uint(id), userID, models.ApprovalActionReject, req.Comment)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "驳回失败: "+err.Error())
		return
	}

	utils.Success(ctx, request)
}

func (c *ApprovalController) Revoke(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	var req struct {
		Comment string `json:"comment"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	userID := ctx.GetUint("user_id")

	request, err := c.approvalService.RevokeRequest(uint(id), userID, req.Comment)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "撤回失败: "+err.Error())
		return
	}

	utils.Success(ctx, request)
}

func (c *ApprovalController) AddSign(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	var req struct {
		SignUserID uint   `json:"sign_user_id" binding:"required"`
		Comment    string `json:"comment"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	userID := ctx.GetUint("user_id")

	request, err := c.approvalService.AddSign(uint(id), userID, req.SignUserID, req.Comment)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "加签失败: "+err.Error())
		return
	}

	utils.Success(ctx, request)
}

func (c *ApprovalController) Transfer(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	var req struct {
		TargetUserID uint   `json:"target_user_id" binding:"required"`
		Comment      string `json:"comment"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	userID := ctx.GetUint("user_id")

	request, err := c.approvalService.TransferApproval(uint(id), userID, req.TargetUserID, req.Comment)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "转交失败: "+err.Error())
		return
	}

	utils.Success(ctx, request)
}

func (c *ApprovalController) GetProgress(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	request, err := c.approvalService.GetRequestByID(uint(id))
	if err != nil {
		utils.Error(ctx, http.StatusNotFound, "审批单不存在")
		return
	}

	instance, err := c.approvalService.GetProcessInstance(request)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取流程实例失败: "+err.Error())
		return
	}

	records, err := c.approvalService.GetApprovalRecords(uint(id))
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取审批记录失败: "+err.Error())
		return
	}

	utils.Success(ctx, gin.H{
		"instance": instance,
		"records":  records,
		"status":   request.Status,
	})
}

func (c *ApprovalController) Stats(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	stats, err := c.approvalService.GetStats(userID)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取统计数据失败: "+err.Error())
		return
	}

	utils.Success(ctx, stats)
}
