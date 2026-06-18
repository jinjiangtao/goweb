package services

import (
	"encoding/json"
	"shenpi/models"

	"gorm.io/gorm"
)

type TemplateService struct {
	db *gorm.DB
}

func NewTemplateService(db *gorm.DB) *TemplateService {
	return &TemplateService{db: db}
}

func (s *TemplateService) CreateTemplate(template *models.ProcessTemplate) error {
	return s.db.Create(template).Error
}

func (s *TemplateService) GetTemplateByID(id uint) (*models.ProcessTemplate, error) {
	var template models.ProcessTemplate
	err := s.db.Preload("Creator").First(&template, id).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (s *TemplateService) GetTemplateList(page, pageSize int, keyword string, templateType string, creatorID uint) ([]models.ProcessTemplate, int64, error) {
	var templates []models.ProcessTemplate
	var total int64

	query := s.db.Model(&models.ProcessTemplate{}).Preload("Creator")

	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	if templateType != "" {
		query = query.Where("type = ?", templateType)
	}
	if creatorID > 0 {
		query = query.Where("creator_id = ?", creatorID)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&templates).Error
	if err != nil {
		return nil, 0, err
	}

	return templates, total, nil
}

func (s *TemplateService) UpdateTemplate(template *models.ProcessTemplate) error {
	return s.db.Save(template).Error
}

func (s *TemplateService) DeleteTemplate(id uint) error {
	return s.db.Delete(&models.ProcessTemplate{}, id).Error
}

func (s *TemplateService) GetTemplateNodes(template *models.ProcessTemplate) ([]models.Node, error) {
	engine := NewProcessEngine()
	return engine.ParseNodes(template)
}

func (s *TemplateService) GetTemplateEdges(template *models.ProcessTemplate) ([]models.Edge, error) {
	engine := NewProcessEngine()
	return engine.ParseEdges(template)
}

func (s *TemplateService) GetTemplateFormFields(template *models.ProcessTemplate) ([]models.FormField, error) {
	engine := NewProcessEngine()
	return engine.ParseFormConfig(template)
}

func (s *TemplateService) BuildDefaultTemplate(name string, templateType string, creatorID uint) *models.ProcessTemplate {
	nodes := []models.Node{
		{ID: "start", Type: "start", Label: "开始", X: 100, Y: 200, ApproverType: "all"},
		{ID: "approval1", Type: "approval", Label: "部门主管审批", X: 300, Y: 200, ApproverType: "specified", Approvers: []uint{2, 3}},
		{ID: "approval2", Type: "approval", Label: "财务审批", X: 500, Y: 200, ApproverType: "specified", Approvers: []uint{4}},
		{ID: "end", Type: "end", Label: "结束", X: 700, Y: 200, ApproverType: "all"},
	}

	edges := []models.Edge{
		{ID: "e1", Source: "start", Target: "approval1"},
		{ID: "e2", Source: "approval1", Target: "approval2"},
		{ID: "e3", Source: "approval2", Target: "end"},
	}

	var formConfig []models.FormField
	if templateType == "leave" {
		formConfig = s.buildLeaveForm()
	} else if templateType == "reimburse" {
		formConfig = s.buildReimburseForm()
	} else {
		formConfig = s.buildGeneralForm()
	}

	nodesJSON, _ := json.Marshal(nodes)
	edgesJSON, _ := json.Marshal(edges)
	formJSON, _ := json.Marshal(formConfig)

	return &models.ProcessTemplate{
		Name:        name,
		Description: "系统默认模板",
		Type:        templateType,
		FormConfig:  string(formJSON),
		Nodes:       string(nodesJSON),
		Edges:       string(edgesJSON),
		Status:      "published",
		CreatorID:   creatorID,
	}
}

func (s *TemplateService) buildLeaveForm() []models.FormField {
	return []models.FormField{
		{Name: "leave_type", Label: "请假类型", Type: "select", Required: true, Options: []map[string]string{
			{"label": "年假", "value": "annual"},
			{"label": "事假", "value": "personal"},
			{"label": "病假", "value": "sick"},
			{"label": "婚假", "value": "marriage"},
			{"label": "产假", "value": "maternity"},
		}},
		{Name: "start_date", Label: "开始日期", Type: "date", Required: true},
		{Name: "end_date", Label: "结束日期", Type: "date", Required: true},
		{Name: "days", Label: "请假天数", Type: "number", Required: true, Validation: map[string]interface{}{"min": 0.5}},
		{Name: "reason", Label: "请假事由", Type: "textarea", Required: true, Placeholder: "请详细说明请假原因"},
	}
}

func (s *TemplateService) buildReimburseForm() []models.FormField {
	return []models.FormField{
		{Name: "reimburse_type", Label: "报销类型", Type: "select", Required: true, Options: []map[string]string{
			{"label": "差旅费", "value": "travel"},
			{"label": "办公费", "value": "office"},
			{"label": "招待费", "value": "entertainment"},
			{"label": "交通费", "value": "transport"},
			{"label": "其他", "value": "other"},
		}},
		{Name: "amount", Label: "报销金额", Type: "number", Required: true, Validation: map[string]interface{}{"min": 0}},
		{Name: "description", Label: "费用说明", Type: "textarea", Required: true, Placeholder: "请详细说明费用情况"},
		{Name: "invoice_no", Label: "发票号码", Type: "input", Required: false},
		{Name: "payment_method", Label: "支付方式", Type: "select", Required: true, Options: []map[string]string{
			{"label": "个人垫付", "value": "personal"},
			{"label": "公司账户", "value": "company"},
		}},
	}
}

func (s *TemplateService) buildGeneralForm() []models.FormField {
	return []models.FormField{
		{Name: "title", Label: "申请标题", Type: "input", Required: true, Placeholder: "请输入申请标题"},
		{Name: "content", Label: "申请内容", Type: "textarea", Required: true, Placeholder: "请详细说明申请内容"},
		{Name: "remark", Label: "备注", Type: "textarea", Required: false},
	}
}
