package controllers

import (
	"bytes"
	"fmt"
	"jianli-server/models"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf/v2"
)

type PDFController struct{}

func NewPDFController() *PDFController {
	return &PDFController{}
}

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

	pdf := gofpdf.New("P", "mm", "A4", "")

	fontDir := "./fonts"
	addChineseFont(pdf, fontDir)

	primaryColor := getStringFromStyle(req.StyleConfig, "primary_color", "#1a1a2e")
	accentColor := getStringFromStyle(req.StyleConfig, "accent_color", "#0f3460")
	fontSize := getFloatFromStyle(req.StyleConfig, "font_size", 12)
	lineSpacing := getFloatFromStyle(req.StyleConfig, "line_spacing", 1.5)

	pdf.SetFont("zh", "", fontSize)

	pdf.SetMargins(
		getFloatFromStyle(req.StyleConfig, "margin_left", 25),
		getFloatFromStyle(req.StyleConfig, "margin_top", 20),
		getFloatFromStyle(req.StyleConfig, "margin_right", 25),
	)
	pdf.SetAutoPageBreak(true, getFloatFromStyle(req.StyleConfig, "margin_bottom", 20))

	pdf.AddPage()

	renderResumeContent(pdf, req.Content, req.StyleConfig, primaryColor, accentColor, fontSize, lineSpacing)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
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

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(http.StatusOK, "application/pdf", buf.Bytes())
}

func addChineseFont(pdf *gofpdf.Fpdf, fontDir string) {
	regularFont := fontDir + "/NotoSansSC-Regular.ttf"
	boldFont := fontDir + "/NotoSansSC-Bold.ttf"

	if fileExists(regularFont) && fileExists(boldFont) {
		pdf.AddUTF8Font("zh", "", regularFont)
		pdf.AddUTF8Font("zh", "B", boldFont)
	} else {
		pdf.AddFont("Arial", "", "")
		pdf.AddFont("Arial", "B", "")
		pdf.SetFont("Arial", "", 12)
		pdf.AliasNbPages("{nb}")
	}
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func getFloatFromStyle(style models.JSON, key string, defaultValue float64) float64 {
	if val, ok := style[key]; ok {
		switch v := val.(type) {
		case float64:
			return v
		case int:
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

func renderResumeContent(pdf *gofpdf.Fpdf, content models.JSON, style models.JSON, primaryColor, accentColor string, fontSize, lineSpacing float64) {
	modules, ok := content["modules"].([]interface{})
	if !ok {
		return
	}

	moduleList := make([]map[string]interface{}, 0)
	for _, m := range modules {
		if mod, ok := m.(map[string]interface{}); ok {
			if visible, _ := mod["visible"].(bool); visible {
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
		moduleID, _ := mod["id"].(string)
		moduleName, _ := mod["name"].(string)

		renderModuleHeader(pdf, moduleName, accentColor, fontSize)

		switch moduleID {
		case "basic":
			renderBasicInfo(pdf, content, fontSize, lineSpacing, primaryColor)
		case "summary":
			renderSummary(pdf, content, fontSize, lineSpacing)
		case "education":
			renderEducation(pdf, content, fontSize, lineSpacing, accentColor)
		case "experience":
			renderExperience(pdf, content, fontSize, lineSpacing, accentColor)
		case "skills":
			renderSkills(pdf, content, fontSize, lineSpacing)
		case "projects":
			renderProjects(pdf, content, fontSize, lineSpacing, accentColor)
		}

		pdf.Ln(5)
	}
}

func renderModuleHeader(pdf *gofpdf.Fpdf, name string, accentColor string, fontSize float64) {
	r, g, b := hexToRGB(accentColor)
	pdf.SetDrawColor(r, g, b)
	pdf.SetFillColor(r, g, b)
	pdf.Rect(10, pdf.GetY(), 3, fontSize+2, "F")

	pdf.SetFont("zh", "B", fontSize+2)
	pdf.SetTextColor(r, g, b)
	pdf.Cell(5, fontSize+2, "")
	pdf.Cell(0, fontSize+2, name)
	pdf.Ln(fontSize + 6)

	pdf.SetDrawColor(200, 200, 200)
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(4)
}

func renderBasicInfo(pdf *gofpdf.Fpdf, content models.JSON, fontSize, lineSpacing float64, primaryColor string) {
	basic, ok := content["basic"].(map[string]interface{})
	if !ok {
		return
	}

	name, _ := basic["name"].(string)
	if name != "" {
		r, g, b := hexToRGB(primaryColor)
		pdf.SetFont("zh", "B", fontSize+6)
		pdf.SetTextColor(r, g, b)
		pdf.Cell(0, fontSize+4, name)
		pdf.Ln(fontSize + 6)
	}

	pdf.SetFont("zh", "", fontSize)
	pdf.SetTextColor(50, 50, 50)

	infoItems := []string{"title", "phone", "email", "github", "website", "address", "location"}
	lineH := fontSize * lineSpacing

	for _, key := range infoItems {
		if val, ok := basic[key].(string); ok && val != "" {
			label := getLabelForKey(key)
			pdf.Cell(20, lineH, label+":")
			pdf.Cell(0, lineH, val)
			pdf.Ln(lineH)
		}
	}
}

func getLabelForKey(key string) string {
	labels := map[string]string{
		"title":    "职位",
		"phone":    "电话",
		"email":    "邮箱",
		"github":   "GitHub",
		"website":  "网站",
		"address":  "地址",
		"location": "城市",
		"name":     "姓名",
	}
	if label, ok := labels[key]; ok {
		return label
	}
	return key
}

func renderSummary(pdf *gofpdf.Fpdf, content models.JSON, fontSize, lineSpacing float64) {
	summary, ok := content["summary"].(string)
	if !ok || summary == "" {
		return
	}

	pdf.SetFont("zh", "", fontSize)
	pdf.SetTextColor(50, 50, 50)
	lineH := fontSize * lineSpacing
	pdf.MultiCell(0, lineH, summary, "", "", false)
}

func renderEducation(pdf *gofpdf.Fpdf, content models.JSON, fontSize, lineSpacing float64, accentColor string) {
	education, ok := content["education"].([]interface{})
	if !ok || len(education) == 0 {
		return
	}

	lineH := fontSize * lineSpacing
	r, g, b := hexToRGB(accentColor)

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

		pdf.SetFont("zh", "B", fontSize)
		pdf.SetTextColor(r, g, b)
		pdf.Cell(0, lineH, school)
		pdf.Ln(lineH)

		pdf.SetFont("zh", "", fontSize)
		pdf.SetTextColor(80, 80, 80)
		degreeInfo := fmt.Sprintf("%s  %s", degree, major)
		if startDate != "" || endDate != "" {
			degreeInfo += fmt.Sprintf("  |  %s - %s", startDate, endDate)
		}
		pdf.Cell(0, lineH, degreeInfo)
		pdf.Ln(lineH)

		if description != "" {
			pdf.SetTextColor(50, 50, 50)
			pdf.MultiCell(0, lineH, description, "", "", false)
		}

		pdf.Ln(3)
	}
}

func renderExperience(pdf *gofpdf.Fpdf, content models.JSON, fontSize, lineSpacing float64, accentColor string) {
	experience, ok := content["experience"].([]interface{})
	if !ok || len(experience) == 0 {
		return
	}

	lineH := fontSize * lineSpacing
	r, g, b := hexToRGB(accentColor)

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

		pdf.SetFont("zh", "B", fontSize)
		pdf.SetTextColor(r, g, b)
		pdf.Cell(0, lineH, fmt.Sprintf("%s  -  %s", company, position))
		pdf.Ln(lineH)

		pdf.SetFont("zh", "", fontSize)
		pdf.SetTextColor(80, 80, 80)
		dateStr := startDate
		if endDate != "" {
			dateStr += " - " + endDate
		} else {
			dateStr += " - 至今"
		}
		pdf.Cell(0, lineH, dateStr)
		pdf.Ln(lineH)

		if description != "" {
			pdf.SetTextColor(50, 50, 50)
			pdf.MultiCell(0, lineH, description, "", "", false)
		}

		if highlights, ok := exp["highlights"].([]interface{}); ok && len(highlights) > 0 {
			pdf.Ln(2)
			for _, h := range highlights {
				if hl, ok := h.(string); ok && hl != "" {
					pdf.SetTextColor(50, 50, 50)
					pdf.Cell(5, lineH, "•")
					pdf.MultiCell(0, lineH, " "+hl, "", "", false)
				}
			}
		}

		pdf.Ln(3)
	}
}

func renderSkills(pdf *gofpdf.Fpdf, content models.JSON, fontSize, lineSpacing float64) {
	skills, ok := content["skills"].([]interface{})
	if !ok || len(skills) == 0 {
		return
	}

	lineH := fontSize * lineSpacing
	pdf.SetFont("zh", "", fontSize)
	pdf.SetTextColor(50, 50, 50)

	for _, item := range skills {
		skill, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		category, _ := skill["category"].(string)
		items, _ := skill["items"].(string)

		if category != "" {
			pdf.SetFont("zh", "B", fontSize)
			pdf.Cell(25, lineH, category+":")
		}
		pdf.SetFont("zh", "", fontSize)
		pdf.MultiCell(0, lineH, items, "", "", false)
		pdf.Ln(2)
	}
}

func renderProjects(pdf *gofpdf.Fpdf, content models.JSON, fontSize, lineSpacing float64, accentColor string) {
	projects, ok := content["projects"].([]interface{})
	if !ok || len(projects) == 0 {
		return
	}

	lineH := fontSize * lineSpacing
	r, g, b := hexToRGB(accentColor)

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

		pdf.SetFont("zh", "B", fontSize)
		pdf.SetTextColor(r, g, b)
		title := name
		if role != "" {
			title += "  -  " + role
		}
		pdf.Cell(0, lineH, title)
		pdf.Ln(lineH)

		pdf.SetFont("zh", "", fontSize)
		pdf.SetTextColor(80, 80, 80)
		dateStr := startDate
		if endDate != "" {
			dateStr += " - " + endDate
		}
		if dateStr != "" {
			pdf.Cell(0, lineH, dateStr)
			pdf.Ln(lineH)
		}

		if description != "" {
			pdf.SetTextColor(50, 50, 50)
			pdf.MultiCell(0, lineH, description, "", "", false)
		}

		if highlights, ok := proj["highlights"].([]interface{}); ok && len(highlights) > 0 {
			pdf.Ln(2)
			for _, h := range highlights {
				if hl, ok := h.(string); ok && hl != "" {
					pdf.SetTextColor(50, 50, 50)
					pdf.Cell(5, lineH, "•")
					pdf.MultiCell(0, lineH, " "+hl, "", "", false)
				}
			}
		}

		pdf.Ln(3)
	}
}
