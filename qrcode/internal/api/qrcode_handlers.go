package api

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"qrcode/internal/model"
	"qrcode/internal/qrcode"
)

type generateReq struct {
	Content      string            `json:"content"`
	StyleConfig  model.StyleConfig `json:"styleConfig"`
	SaveTemplate *bool             `json:"saveTemplate"`
	TemplateName string            `json:"templateName"`
	TemplateID   *int64            `json:"templateId"`
}

type generateResp struct {
	ID         int64         `json:"id"`
	PreviewURL string        `json:"previewUrl"`
	Record     *model.Record `json:"record,omitempty"`
}

func (svc *Service) preview(c *gin.Context) {
	var req generateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Err("参数错误: "+err.Error()))
		return
	}
	png, err := qrcode.Generate(req.Content, req.StyleConfig.WithDefaults())
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Err("生成失败: "+err.Error()))
		return
	}
	svc.imagePNG(c, png)
}

func (svc *Service) generate(c *gin.Context) {
	var req generateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Err("参数错误: "+err.Error()))
		return
	}
	style := req.StyleConfig.WithDefaults()
	png, err := qrcode.Generate(req.Content, style)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Err("生成失败: "+err.Error()))
		return
	}
	path, err := svc.savePNG(png)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Err("保存失败: "+err.Error()))
		return
	}
	rec := &model.Record{
		Content:     req.Content,
		ContentType: qrcode.DetectType(req.Content),
		Source:      "single",
		StyleConfig: style,
		FilePath:    path,
		Status:      "active",
		TemplateID:  req.TemplateID,
	}
	if err := svc.Store.CreateRecord(rec); err != nil {
		c.JSON(http.StatusInternalServerError, model.Err("归档失败: "+err.Error()))
		return
	}

	if req.SaveTemplate != nil && *req.SaveTemplate {
		name := strings.TrimSpace(req.TemplateName)
		if name == "" {
			name = "模板 " + rec.CreatedAt
		}
		_ = svc.Store.CreateTemplate(&model.Template{Name: name, StyleConfig: style})
	}

	c.JSON(http.StatusOK, model.OK(generateResp{
		ID:         rec.ID,
		PreviewURL: "/api/records/" + strconv.FormatInt(rec.ID, 10) + "/image",
		Record:     rec,
	}))
}

type batchItem struct {
	Content     string             `json:"content"`
	StyleConfig *model.StyleConfig `json:"styleConfig"`
}

type batchReq struct {
	Items       []batchItem       `json:"items"`
	StyleConfig model.StyleConfig `json:"styleConfig"`
}

type batchResp struct {
	BatchID string         `json:"batchId"`
	Count   int            `json:"count"`
	Records []model.Record `json:"records"`
}

func (svc *Service) batch(c *gin.Context) {
	var req batchReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Err("参数错误: "+err.Error()))
		return
	}
	if len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, model.Err("批量内容为空"))
		return
	}
	if len(req.Items) > 500 {
		c.JSON(http.StatusBadRequest, model.Err("单次批量上限 500 条"))
		return
	}
	baseStyle := req.StyleConfig.WithDefaults()
	batchID := randID() + "-" + time.Now().Format("20060102")
	records := make([]model.Record, 0, len(req.Items))
	for _, it := range req.Items {
		style := baseStyle
		if it.StyleConfig != nil {
			style = it.StyleConfig.WithDefaults()
		}
		content := strings.TrimSpace(it.Content)
		if content == "" {
			continue
		}
		png, err := qrcode.Generate(content, style)
		if err != nil {
			continue
		}
		path, err := svc.savePNG(png)
		if err != nil {
			continue
		}
		rec := &model.Record{
			Content:     content,
			ContentType: qrcode.DetectType(content),
			Source:      "batch",
			StyleConfig: style,
			FilePath:    path,
			Status:      "active",
			BatchID:     batchID,
		}
		if err := svc.Store.CreateRecord(rec); err == nil {
			records = append(records, *rec)
		}
	}
	c.JSON(http.StatusOK, model.OK(batchResp{
		BatchID: batchID,
		Count:   len(records),
		Records: records,
	}))
}
