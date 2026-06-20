package controllers

import (
	"bytes"
	"fmt"
	"jianli-server/models"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/signintech/gopdf"
)

type PDFController struct{}

func NewPDFController() *PDFController {
	return &PDFController{}
}

var zhFontRegular = "zh"
var zhFontBold = "zh-bold"
var zhFontAvailable = false

func (pc *PDFController) ExportPDF(c *gin.Context) {
	var req struct {
		Content     models.JSON `json:"content" binding:"required"`
		StyleConfig models.JSON `json:"style_config" binding:"required"`
		Filename    string      `json:"filename"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	zhFontAvailable = false
	if setupFonts(&pdf) {
		zhFontAvailable = true
	}

	primaryColor := getStringFromStyle(req.StyleConfig, "primaryColor", "#1a1a2e")
	if primaryColor == "#1a1a2e" {
		if v, ok := req.StyleConfig["primary_color"].(string); ok && v != "" {
			primaryColor = v
		}
	}
	accentColor := getStringFromStyle(req.StyleConfig, "accentColor", "#0f3460")
	if accentColor == "#0f3460" {
		if v, ok := req.StyleConfig["accent_color"].(string); ok && v != "" {
			accentColor = v
		}
	}
	fontSize := getFloatFromStyle(req.StyleConfig, "fontSize", 12)
	if fontSize == 12 {
		if v := getFloatFromStyle(req.StyleConfig, "font_size", 0); v > 0 {
			fontSize = v
		}
	}
	lineSpacing := getFloatFromStyle(req.StyleConfig, "lineSpacing", 1.5)
	if lineSpacing == 1.5 {
		if v := getFloatFromStyle(req.StyleConfig, "line_spacing", 0); v > 0 {
			lineSpacing = v
		}
	}

	marginLeft := getFloatFromStyle(req.StyleConfig, "marginLeft", 25)
	if marginLeft == 25 {
		if v := getFloatFromStyle(req.StyleConfig, "margin_left", 0); v > 0 {
			marginLeft = v
		}
	}
	marginTop := getFloatFromStyle(req.StyleConfig, "marginTop", 20)
	if marginTop == 20 {
		if v := getFloatFromStyle(req.StyleConfig, "margin_top", 0); v > 0 {
			marginTop = v
		}
	}
	marginRight := getFloatFromStyle(req.StyleConfig, "marginRight", 25)
	if marginRight == 25 {
		if v := getFloatFromStyle(req.StyleConfig, "margin_right", 0); v > 0 {
			marginRight = v
		}
	}
	marginBottom := getFloatFromStyle(req.StyleConfig, "marginBottom", 20)
	if marginBottom == 20 {
		if v := getFloatFromStyle(req.StyleConfig, "margin_bottom", 0); v > 0 {
			marginBottom = v
		}
	}

	fontReg := zhFontRegular
	fontBd := zhFontBold
	if !zhFontAvailable {
		pdf.AddTTFFont("helvetica", findHelvetica())
		fontReg = "helvetica"
		fontBd = "helvetica"
	}

	pdf.AddPage()
	pdf.SetMargins(marginLeft, marginTop, marginRight, marginBottom)

	ctx := &pdfRenderContext{
		pdf:         &pdf,
		fontReg:     fontReg,
		fontBd:      fontBd,
		pageW:       gopdf.PageSizeA4.W,
		leftMargin:  marginLeft,
		rightMargin: marginRight,
	}

	renderResumeContent(ctx, req.Content, primaryColor, accentColor, fontSize, lineSpacing)

	var buf bytes.Buffer
	if err := pdf.Write(&buf); err != nil {
		log.Printf("[PDF ERROR] Write failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成PDF失败",
			"error":   err.Error(),
		})
		return
	}

	filename := req.Filename
	if filename == "" {
		filename = fmt.Sprintf("简历_%s.pdf", time.Now().Format("20060102150405"))
	}
	if !strings.HasSuffix(strings.ToLower(filename), ".pdf") {
		filename += ".pdf"
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(http.StatusOK, "application/pdf", buf.Bytes())
}

type pdfRenderContext struct {
	pdf         *gopdf.GoPdf
	fontReg     string
	fontBd      string
	pageW       float64
	leftMargin  float64
	rightMargin float64
}

func (c *pdfRenderContext) contentW() float64 {
	return c.pageW - c.leftMargin - c.rightMargin
}

func (c *pdfRenderContext) setFont(style string, size float64) {
	font := c.fontReg
	if style == "B" {
		font = c.fontBd
	}
	c.pdf.SetFont(font, "", size)
}

func (c *pdfRenderContext) setTextColor(hex string) {
	r, g, b := hexToRGB(hex)
	c.pdf.SetTextColor(uint8(r), uint8(g), uint8(b))
}

func (c *pdfRenderContext) setRGBTextColor(r, g, b int) {
	c.pdf.SetTextColor(uint8(r), uint8(g), uint8(b))
}

func (c *pdfRenderContext) setDrawColor(hex string) {
	r, g, b := hexToRGB(hex)
	c.pdf.SetStrokeColor(uint8(r), uint8(g), uint8(b))
}

func (c *pdfRenderContext) setFillColor(hex string) {
	r, g, b := hexToRGB(hex)
	c.pdf.SetFillColor(uint8(r), uint8(g), uint8(b))
}

func (c *pdfRenderContext) cell(w, h float64, text string) {
	opt := gopdf.CellOption{Align: gopdf.Left}
	c.pdf.CellWithOption(&gopdf.Rect{W: w, H: h}, text, opt)
}

func (c *pdfRenderContext) br(h float64) {
	c.pdf.Br(h)
}

func (c *pdfRenderContext) multiCell(w, h float64, text string) {
	if text == "" {
		return
	}
	opt := gopdf.CellOption{Align: gopdf.Left}
	c.pdf.MultiCellWithOption(&gopdf.Rect{W: w, H: h}, text, opt)
}

func (c *pdfRenderContext) rect(x, y, w, h float64) {
	c.pdf.RectFromUpperLeftWithStyle(x, y, w, h, "F")
}

func (c *pdfRenderContext) line(x1, y1, x2, y2 float64) {
	c.pdf.Line(x1, y1, x2, y2)
}

func (c *pdfRenderContext) getY() float64 {
	return c.pdf.GetY()
}

func (c *pdfRenderContext) setY(y float64) {
	c.pdf.SetY(y)
}

func findHelvetica() string {
	candidates := []string{
		"C:\\Windows\\Fonts\\arial.ttf",
		"C:\\Windows\\Fonts\\ARIAL.TTF",
		"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
		"/usr/share/fonts/TTF/DejaVuSans.ttf",
		"/System/Library/Fonts/Helvetica.ttc",
	}
	for _, p := range candidates {
		if fileExists(p) && fileSize(p) > 1000 {
			return p
		}
	}
	return ""
}

func setupFonts(pdf *gopdf.GoPdf) bool {
	fontDir, _ := filepath.Abs("./fonts")
	if _, err := os.Stat(fontDir); os.IsNotExist(err) {
		fontDir, _ = filepath.Abs("../server/fonts")
	}

	candidates := []struct {
		regular string
		bold    string
	}{
		{filepath.Join(fontDir, "NotoSansSC-Regular.ttf"), filepath.Join(fontDir, "NotoSansSC-Bold.ttf")},
		{filepath.Join(fontDir, "NotoSansSC-Regular.otf"), filepath.Join(fontDir, "NotoSansSC-Bold.otf")},
		{filepath.Join(fontDir, "SimHei.ttf"), filepath.Join(fontDir, "SimHei.ttf")},
		{filepath.Join(fontDir, "SimKai.ttf"), filepath.Join(fontDir, "SimKai.ttf")},
		{filepath.Join(fontDir, "MicrosoftYaHei.ttc"), filepath.Join(fontDir, "MicrosoftYaHei.ttc")},
		{"C:\\Windows\\Fonts\\simhei.ttf", "C:\\Windows\\Fonts\\simhei.ttf"},
		{"C:\\Windows\\Fonts\\msyh.ttc", "C:\\Windows\\Fonts\\msyh.ttc"},
		{"C:\\Windows\\Fonts\\simsun.ttc", "C:\\Windows\\Fonts\\simsun.ttc"},
		{"C:\\Windows\\Fonts\\simkai.ttf", "C:\\Windows\\Fonts\\simkai.ttf"},
	}

	for _, c := range candidates {
		if fileExists(c.regular) && fileSize(c.regular) > 1000 {
			boldPath := c.regular
			if fileExists(c.bold) && fileSize(c.bold) > 1000 {
				boldPath = c.bold
			}
			log.Printf("[PDF INFO] Trying gopdf font: regular=%s bold=%s", c.regular, boldPath)
			if err := pdf.AddTTFFont(zhFontRegular, c.regular); err != nil {
				log.Printf("[PDF WARN] AddTTFFont regular failed for %s: %v", c.regular, err)
				continue
			}
			if err := pdf.AddTTFFont(zhFontBold, boldPath); err != nil {
				log.Printf("[PDF WARN] AddTTFFont bold failed for %s: %v", boldPath, err)
				continue
			}
			log.Printf("[PDF INFO] gopdf fonts loaded successfully")
			return true
		}
	}
	log.Printf("[PDF WARN] No Chinese font found for gopdf")
	return false
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func fileSize(filename string) int64 {
	if info, err := os.Stat(filename); err == nil {
		return info.Size()
	}
	return 0
}

func getFloatFromStyle(style models.JSON, key string, defaultValue float64) float64 {
	if val, ok := style[key]; ok {
		switch v := val.(type) {
		case float64:
			return v
		case float32:
			return float64(v)
		case int:
			return float64(v)
		case int64:
			return float64(v)
		}
	}
	return defaultValue
}

func getStringFromStyle(style models.JSON, key string, defaultValue string) string {
	if val, ok := style[key]; ok {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return defaultValue
}

func hexToRGB(hex string) (int, int, int) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) == 3 {
		hex = string([]byte{hex[0], hex[0], hex[1], hex[1], hex[2], hex[2]})
	}
	r, _ := strconv.ParseInt(hex[0:2], 16, 0)
	g, _ := strconv.ParseInt(hex[2:4], 16, 0)
	b, _ := strconv.ParseInt(hex[4:6], 16, 0)
	return int(r), int(g), int(b)
}

func getModuleType(mod map[string]interface{}) string {
	if t, ok := mod["type"].(string); ok && t != "" {
		return t
	}
	if id, ok := mod["id"].(string); ok {
		known := map[string]bool{
			"basic": true, "summary": true, "evaluation": true,
			"education": true, "experience": true,
			"skills": true, "projects": true,
		}
		if known[id] {
			return id
		}
	}
	return ""
}

func getBasicField(basic map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if v, ok := basic[k].(string); ok && v != "" {
			return v
		}
	}
	return ""
}

func renderResumeContent(c *pdfRenderContext, content models.JSON, primaryColor, accentColor string, fontSize, lineSpacing float64) {
	modules, ok := content["modules"].([]interface{})
	if !ok {
		return
	}

	moduleList := make([]map[string]interface{}, 0)
	for _, m := range modules {
		if mod, ok := m.(map[string]interface{}); ok {
			visible, _ := mod["visible"].(bool)
			if visible || mod["visible"] == nil {
				moduleList = append(moduleList, mod)
			}
		}
	}

	sort.Slice(moduleList, func(i, j int) bool {
		orderI, _ := moduleList[i]["order"].(float64)
		orderJ, _ := moduleList[j]["order"].(float64)
		return orderI < orderJ
	})

	for _, mod := range moduleList {
		moduleName, _ := mod["name"].(string)
		if moduleName == "" {
			moduleName = getModuleType(mod)
		}
		moduleType := getModuleType(mod)
		if moduleType == "" {
			continue
		}

		renderModuleHeader(c, moduleName, accentColor, fontSize)

		switch moduleType {
		case "basic":
			renderBasicInfo(c, content, fontSize, lineSpacing, primaryColor)
		case "summary":
			renderSummary(c, content, fontSize, lineSpacing)
		case "evaluation":
			renderEvaluation(c, content, fontSize, lineSpacing)
		case "education":
			renderEducation(c, content, fontSize, lineSpacing, accentColor)
		case "experience":
			renderExperience(c, content, fontSize, lineSpacing, accentColor)
		case "skills":
			renderSkills(c, content, fontSize, lineSpacing)
		case "projects":
			renderProjects(c, content, fontSize, lineSpacing, accentColor)
		}

		c.br(5)
	}
}

func renderModuleHeader(c *pdfRenderContext, name string, accentColor string, fontSize float64) {
	c.setDrawColor(accentColor)
	c.setFillColor(accentColor)
	x := c.leftMargin
	y := c.getY()
	c.rect(x, y, 3, fontSize+2)

	c.setFont("B", fontSize+2)
	c.setTextColor(accentColor)
	c.br(0)
	x2 := c.leftMargin + 5
	c.pdf.SetX(x2)
	c.cell(c.contentW()-5, fontSize+2, name)
	c.br(fontSize + 6)

	c.setDrawColor("#C8C8C8")
	y2 := c.getY()
	c.line(x, y2, c.pageW-c.rightMargin, y2)
	c.br(4)
}

func renderBasicInfo(c *pdfRenderContext, content models.JSON, fontSize, lineSpacing float64, primaryColor string) {
	basic, ok := content["basic"].(map[string]interface{})
	if !ok {
		return
	}

	name := getBasicField(basic, "name")
	lineH := fontSize * lineSpacing
	if name != "" {
		c.setFont("B", fontSize+6)
		c.setTextColor(primaryColor)
		c.cell(c.contentW(), fontSize+4, name)
		c.br(fontSize + 6)
	}

	jobIntention := getBasicField(basic, "job_intention", "title")
	if jobIntention != "" {
		c.setFont("", fontSize)
		c.setRGBTextColor(100, 100, 100)
		c.cell(c.contentW(), lineH, jobIntention)
		c.br(lineH + 2)
	}

	c.setFont("", fontSize)
	c.setRGBTextColor(50, 50, 50)

	infoRows := [][]string{
		{"phone", "电话"},
		{"email", "邮箱"},
		{"city", "城市"},
		{"location", "城市"},
		{"address", "地址"},
		{"github", "GitHub"},
		{"website", "网站"},
	}
	written := map[string]bool{}
	for _, row := range infoRows {
		key := row[0]
		label := row[1]
		if written[label] {
			continue
		}
		val := getBasicField(basic, key)
		if val != "" {
			c.setFont("B", fontSize)
			c.setRGBTextColor(80, 80, 80)
			c.pdf.SetX(c.leftMargin)
			c.cell(18, lineH, label+":")
			c.setFont("", fontSize)
			c.setRGBTextColor(50, 50, 50)
			c.cell(c.contentW()-18, lineH, val)
			c.br(lineH)
			written[label] = true
		}
	}
}

func renderSummary(c *pdfRenderContext, content models.JSON, fontSize, lineSpacing float64) {
	summary, ok := content["summary"].(string)
	if !ok || summary == "" {
		return
	}

	c.setFont("", fontSize)
	c.setRGBTextColor(50, 50, 50)
	lineH := fontSize * lineSpacing
	c.multiCell(c.contentW(), lineH, summary)
}

func renderEvaluation(c *pdfRenderContext, content models.JSON, fontSize, lineSpacing float64) {
	evaluation, ok := content["evaluation"].(string)
	if !ok || evaluation == "" {
		return
	}

	c.setFont("", fontSize)
	c.setRGBTextColor(50, 50, 50)
	lineH := fontSize * lineSpacing
	c.multiCell(c.contentW(), lineH, evaluation)
}

func renderEducation(c *pdfRenderContext, content models.JSON, fontSize, lineSpacing float64, accentColor string) {
	education, ok := content["education"].([]interface{})
	if !ok || len(education) == 0 {
		return
	}

	lineH := fontSize * lineSpacing

	for _, item := range education {
		edu, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		school, _ := edu["school"].(string)
		degree, _ := edu["degree"].(string)
		major, _ := edu["major"].(string)
		startDate, _ := edu["start_date"].(string)
		endDate, _ := edu["end_date"].(string)
		description, _ := edu["description"].(string)

		c.setFont("B", fontSize)
		c.setTextColor(accentColor)
		c.cell(c.contentW(), lineH, school)
		c.br(lineH)

		c.setFont("", fontSize)
		c.setRGBTextColor(80, 80, 80)
		degreeInfo := fmt.Sprintf("%s  %s", degree, major)
		if startDate != "" || endDate != "" {
			degreeInfo += fmt.Sprintf("  |  %s - %s", startDate, endDate)
		}
		c.cell(c.contentW(), lineH, degreeInfo)
		c.br(lineH)

		if description != "" {
			c.setRGBTextColor(50, 50, 50)
			c.multiCell(c.contentW(), lineH, description)
		}

		c.br(3)
	}
}

func renderExperience(c *pdfRenderContext, content models.JSON, fontSize, lineSpacing float64, accentColor string) {
	experience, ok := content["experience"].([]interface{})
	if !ok || len(experience) == 0 {
		return
	}

	lineH := fontSize * lineSpacing

	for _, item := range experience {
		exp, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		company, _ := exp["company"].(string)
		position, _ := exp["position"].(string)
		startDate, _ := exp["start_date"].(string)
		endDate, _ := exp["end_date"].(string)
		description, _ := exp["description"].(string)

		c.setFont("B", fontSize)
		c.setTextColor(accentColor)
		title := company
		if position != "" {
			title += "  -  " + position
		}
		c.cell(c.contentW(), lineH, title)
		c.br(lineH)

		c.setFont("", fontSize)
		c.setRGBTextColor(80, 80, 80)
		dateStr := startDate
		if endDate != "" {
			dateStr += " - " + endDate
		} else if startDate != "" {
			dateStr += " - 至今"
		}
		if dateStr != "" {
			c.cell(c.contentW(), lineH, dateStr)
			c.br(lineH)
		}

		if description != "" {
			c.setRGBTextColor(50, 50, 50)
			c.multiCell(c.contentW(), lineH, description)
		}

		if highlights, ok := exp["highlights"].([]interface{}); ok && len(highlights) > 0 {
			c.br(2)
			for _, h := range highlights {
				if hl, ok := h.(string); ok && hl != "" {
					c.setRGBTextColor(50, 50, 50)
					c.pdf.SetX(c.leftMargin)
					c.cell(5, lineH, "•")
					c.multiCell(c.contentW()-5, lineH, " "+hl)
				}
			}
		}

		c.br(3)
	}
}

func renderSkills(c *pdfRenderContext, content models.JSON, fontSize, lineSpacing float64) {
	skills, ok := content["skills"].([]interface{})
	if !ok || len(skills) == 0 {
		return
	}

	lineH := fontSize * lineSpacing
	c.setFont("", fontSize)
	c.setRGBTextColor(50, 50, 50)

	for _, item := range skills {
		skill, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		category, _ := skill["category"].(string)
		items, _ := skill["items"].(string)

		if category != "" {
			c.setFont("B", fontSize)
			c.pdf.SetX(c.leftMargin)
			c.cell(25, lineH, category+":")
		}
		c.setFont("", fontSize)
		c.multiCell(c.contentW()-25, lineH, items)
		c.br(2)
	}
}

func renderProjects(c *pdfRenderContext, content models.JSON, fontSize, lineSpacing float64, accentColor string) {
	projects, ok := content["projects"].([]interface{})
	if !ok || len(projects) == 0 {
		return
	}

	lineH := fontSize * lineSpacing

	for _, item := range projects {
		proj, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		name, _ := proj["name"].(string)
		role, _ := proj["role"].(string)
		startDate, _ := proj["start_date"].(string)
		endDate, _ := proj["end_date"].(string)
		description, _ := proj["description"].(string)

		c.setFont("B", fontSize)
		c.setTextColor(accentColor)
		title := name
		if role != "" {
			title += "  -  " + role
		}
		c.cell(c.contentW(), lineH, title)
		c.br(lineH)

		c.setFont("", fontSize)
		c.setRGBTextColor(80, 80, 80)
		dateStr := startDate
		if endDate != "" {
			dateStr += " - " + endDate
		}
		if dateStr != "" {
			c.cell(c.contentW(), lineH, dateStr)
			c.br(lineH)
		}

		if description != "" {
			c.setRGBTextColor(50, 50, 50)
			c.multiCell(c.contentW(), lineH, description)
		}

		if highlights, ok := proj["highlights"].([]interface{}); ok && len(highlights) > 0 {
			c.br(2)
			for _, h := range highlights {
				if hl, ok := h.(string); ok && hl != "" {
					c.setRGBTextColor(50, 50, 50)
					c.pdf.SetX(c.leftMargin)
					c.cell(5, lineH, "•")
					c.multiCell(c.contentW()-5, lineH, " "+hl)
				}
			}
		}

		c.br(3)
	}
}
