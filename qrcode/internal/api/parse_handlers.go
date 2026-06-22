package api

import (
	"archive/zip"
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"

	"qrcode/internal/model"
	"qrcode/internal/qrcode"
)

type parseTextReq struct {
	Texts []string `json:"texts"`
}

type parseTextResp struct {
	Results []model.ParseResult `json:"results"`
}

func (svc *Service) parseText(c *gin.Context) {
	var req parseTextReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Err("参数错误: "+err.Error()))
		return
	}
	results := make([]model.ParseResult, 0, len(req.Texts))
	for _, t := range req.Texts {
		raw := strings.TrimSpace(t)
		if raw == "" {
			continue
		}
		typ := qrcode.DetectType(raw)
		results = append(results, model.ParseResult{Raw: raw, Type: typ, Content: raw})
		svc.Store.CreateParseLog("text", raw, typ)
	}
	c.JSON(http.StatusOK, model.OK(parseTextResp{Results: results}))
}

type parseImageResp struct {
	Results []model.ParseResult `json:"results"`
}

func (svc *Service) parseImage(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Err("未接收到文件: "+err.Error()))
		return
	}
	files := form.File["files"]
	results := make([]model.ParseResult, 0)
	for _, fh := range files {
		f, err := fh.Open()
		if err != nil {
			continue
		}
		buf, _ := io.ReadAll(f)
		f.Close()

		lower := strings.ToLower(fh.Filename)
		if strings.HasSuffix(lower, ".zip") {
			results = append(results, parseZip(buf)...)
			continue
		}
		if r := decodeAndParse(buf, fh.Filename); r != nil {
			results = append(results, *r)
		}
	}
	c.JSON(http.StatusOK, model.OK(parseImageResp{Results: results}))
}

func parseZip(buf []byte) []model.ParseResult {
	zr, err := zip.NewReader(bytes.NewReader(buf), int64(len(buf)))
	if err != nil {
		return nil
	}
	out := make([]model.ParseResult, 0)
	for _, f := range zr.File {
		if f.FileInfo().IsDir() {
			continue
		}
		ext := strings.ToLower(path.Ext(f.Name))
		if ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".gif" && ext != ".bmp" {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			continue
		}
		b, _ := io.ReadAll(rc)
		rc.Close()
		if r := decodeAndParse(b, f.Name); r != nil {
			out = append(out, *r)
		}
	}
	return out
}

func decodeAndParse(b []byte, name string) *model.ParseResult {
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return &model.ParseResult{Raw: name, Type: "error", Content: "无法识别图片"}
	}
	text, err := qrcode.Parse(img)
	if err != nil {
		return &model.ParseResult{Raw: name, Type: "error", Content: "未检测到二维码"}
	}
	typ := qrcode.DetectType(text)
	return &model.ParseResult{Raw: name, Type: typ, Content: text}
}
