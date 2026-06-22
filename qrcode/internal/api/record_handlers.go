package api

import (
	"archive/zip"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	"qrcode/internal/model"
	"qrcode/internal/qrcode"
)

type statusReq struct {
	Status string `json:"status"`
}

func (svc *Service) listRecords(c *gin.Context) {
	status := c.Query("status")
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, total, err := svc.Store.ListRecords(status, keyword, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Err("查询失败: "+err.Error()))
		return
	}
	if list == nil {
		list = []model.Record{}
	}
	c.JSON(http.StatusOK, model.OK(gin.H{"total": total, "list": list}))
}

func (svc *Service) getRecord(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Err("无效 ID"))
		return
	}
	rec, err := svc.Store.GetRecord(id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.Err("记录不存在"))
		return
	}
	c.JSON(http.StatusOK, model.OK(rec))
}

func (svc *Service) recordImage(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Err("无效 ID"))
		return
	}
	rec, err := svc.Store.GetRecord(id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.Err("记录不存在"))
		return
	}
	data, err := os.ReadFile(svc.absPath(rec.FilePath))
	if err != nil {
		c.JSON(http.StatusNotFound, model.Err("图片文件缺失"))
		return
	}
	svc.imagePNG(c, data)
}

func (svc *Service) updateStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Err("无效 ID"))
		return
	}
	var req statusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Err("参数错误: "+err.Error()))
		return
	}
	status := req.Status
	if status != "active" && status != "invalid" {
		c.JSON(http.StatusBadRequest, model.Err("状态取值无效"))
		return
	}
	if err := svc.Store.UpdateRecordStatus(id, status); err != nil {
		c.JSON(http.StatusInternalServerError, model.Err("更新失败: "+err.Error()))
		return
	}
	c.JSON(http.StatusOK, model.OK(gin.H{"ok": true}))
}

type contentReq struct {
	Content string `json:"content"`
}

func (svc *Service) updateContent(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Err("无效 ID"))
		return
	}
	var req contentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Err("参数错误: "+err.Error()))
		return
	}
	rec, err := svc.Store.GetRecord(id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.Err("记录不存在"))
		return
	}
	png, err := qrcode.Generate(req.Content, rec.StyleConfig.WithDefaults())
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Err("重新生成失败: "+err.Error()))
		return
	}
	if err := svc.overwritePNG(rec.FilePath, png); err != nil {
		c.JSON(http.StatusInternalServerError, model.Err("写入失败: "+err.Error()))
		return
	}
	typ := qrcode.DetectType(req.Content)
	if err := svc.Store.UpdateRecordContent(id, req.Content, typ, rec.FilePath); err != nil {
		c.JSON(http.StatusInternalServerError, model.Err("更新失败: "+err.Error()))
		return
	}
	rec, _ = svc.Store.GetRecord(id)
	c.JSON(http.StatusOK, model.OK(rec))
}

func (svc *Service) deleteRecord(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Err("无效 ID"))
		return
	}
	rec, _ := svc.Store.GetRecord(id)
	_ = svc.Store.DeleteRecord(id)
	if rec.FilePath != "" {
		_ = os.Remove(svc.absPath(rec.FilePath))
	}
	c.JSON(http.StatusOK, model.OK(gin.H{"ok": true}))
}

type exportReq struct {
	IDs []int64 `json:"ids"`
}

func (svc *Service) exportRecords(c *gin.Context) {
	var req exportReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Err("参数错误: "+err.Error()))
		return
	}
	if len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, model.Err("未选择记录"))
		return
	}
	records, err := svc.Store.GetRecordsByIDs(req.IDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Err("查询失败: "+err.Error()))
		return
	}

	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", `attachment; filename="qrcode-export.zip"`)
	zw := zip.NewWriter(c.Writer)
	used := map[string]int{}
	for _, r := range records {
		data, err := os.ReadFile(svc.absPath(r.FilePath))
		if err != nil {
			continue
		}
		base := sanitizeFilename(r.Content, 24)
		name := base
		if n, ok := used[name]; ok {
			used[name] = n + 1
			name = base + "_" + strconv.Itoa(n+1)
		} else {
			used[name] = 1
		}
		name = name + "_" + strconv.FormatInt(r.ID, 10) + ".png"
		w, err := zw.Create(name)
		if err != nil {
			continue
		}
		_, _ = w.Write(data)
	}
	_ = zw.Close()
}
