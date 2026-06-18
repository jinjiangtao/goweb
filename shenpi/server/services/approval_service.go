package services

import (
	"encoding/json"
	"fmt"
	"shenpi/models"
	"time"

	"gorm.io/gorm"
)

type ApprovalService struct {
	db             *gorm.DB
	processEngine  *ProcessEngine
	templateService *TemplateService
}

func NewApprovalService(db *gorm.DB) *ApprovalService {
	return &ApprovalService{
		db:             db,
		processEngine:  NewProcessEngine(),
		templateService: NewTemplateService(db),
	}
}

func (s *ApprovalService) GenerateRequestNo() string {
	now := time.Now()
	return fmt.Sprintf("REQ%s%06d", now.Format("20060102150405"), time.Now().UnixNano()%1000000)
}

func (s *ApprovalService) SubmitApproval(userID uint, templateID uint, title string, formData map[string]interface{}) (*models.ApprovalRequest, error) {
	template, err := s.templateService.GetTemplateByID(templateID)
	if err != nil {
		return nil, err
	}

	instance, err := s.processEngine.InitProcessInstance(template)
	if err != nil {
		return nil, err
	}

	formDataJSON, _ := json.Marshal(formData)
	instanceJSON, _ := json.Marshal(instance)

	request := &models.ApprovalRequest{
		RequestNo:       s.GenerateRequestNo(),
		Title:           title,
		Type:            template.Type,
		TemplateID:      templateID,
		FormData:        string(formDataJSON),
		CurrentNodeID:   instance.CurrentNodeID,
		Status:          models.ApprovalStatusPending,
		ApplicantID:     userID,
		ProcessInstance: string(instanceJSON),
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(request).Error; err != nil {
			return err
		}

		startNode := s.processEngine.FindStartNode(s.getNodesFromInstance(instance))
		if startNode != nil {
			record := &models.ApprovalRecord{
				RequestID:  request.ID,
				NodeID:     startNode.ID,
				NodeName:   startNode.Label,
				ApproverID: userID,
				Action:     models.ApprovalActionStart,
				Comment:    "流程发起",
				CreatedAt:  time.Now(),
			}
			if err := tx.Create(record).Error; err != nil {
				return err
			}
		}

		nodes, _ := s.processEngine.ParseNodes(template)
		for _, node := range nodes {
			if node.Type == "approval" {
				nextNodeID, err := s.processEngine.DetermineNextNode(instance, "start", formData)
				if err == nil && nextNodeID == node.ID {
					nodeInstance := s.processEngine.FindNodeInstanceByID(instance, node.ID)
					if nodeInstance != nil {
						nodeInstance.Status = models.NodeStatusProcessing
						instanceJSON, _ = json.Marshal(instance)
						request.ProcessInstance = string(instanceJSON)
						if err := tx.Save(request).Error; err != nil {
							return err
						}
					}
					break
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return request, nil
}

func (s *ApprovalService) GetRequestByID(id uint) (*models.ApprovalRequest, error) {
	var request models.ApprovalRequest
	err := s.db.Preload("Applicant").Preload("Template").First(&request, id).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

func (s *ApprovalService) GetMyRequests(userID uint, page, pageSize int, status string) ([]models.ApprovalRequest, int64, error) {
	var requests []models.ApprovalRequest
	var total int64

	query := s.db.Model(&models.ApprovalRequest{}).Preload("Applicant").Preload("Template").Where("applicant_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&requests).Error
	if err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

func (s *ApprovalService) GetMyApprovals(userID uint, page, pageSize int, status string) ([]models.ApprovalRequest, int64, error) {
	var requests []models.ApprovalRequest
	var total int64

	query := s.db.Model(&models.ApprovalRequest{}).Preload("Applicant").Preload("Template").
		Where("status IN ?", []string{models.ApprovalStatusPending}).
		Where("EXISTS (SELECT 1 FROM approval_records ar WHERE ar.request_id = approval_requests.id AND ar.approver_id = ?)", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&requests).Error
	if err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

func (s *ApprovalService) GetApprovalRecords(requestID uint) ([]models.ApprovalRecord, error) {
	var records []models.ApprovalRecord
	err := s.db.Preload("Approver").Where("request_id = ?", requestID).Order("created_at ASC").Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (s *ApprovalService) ProcessApproval(requestID uint, approverID uint, action string, comment string) (*models.ApprovalRequest, error) {
	request, err := s.GetRequestByID(requestID)
	if err != nil {
		return nil, err
	}

	if request.Status != models.ApprovalStatusPending {
		return nil, fmt.Errorf("approval request is not pending")
	}

	var instance models.ProcessInstance
	err = json.Unmarshal([]byte(request.ProcessInstance), &instance)
	if err != nil {
		return nil, err
	}

	template, err := s.templateService.GetTemplateByID(request.TemplateID)
	if err != nil {
		return nil, err
	}

	nodes, err := s.processEngine.ParseNodes(template)
	if err != nil {
		return nil, err
	}

	currentNode := s.processEngine.FindNodeByID(nodes, instance.CurrentNodeID)
	if currentNode == nil {
		return nil, fmt.Errorf("current node not found")
	}

	nodeInstance := s.processEngine.FindNodeInstanceByID(&instance, instance.CurrentNodeID)
	if nodeInstance == nil {
		return nil, fmt.Errorf("current node instance not found")
	}

	var formData map[string]interface{}
	json.Unmarshal([]byte(request.FormData), &formData)

	if !s.processEngine.CheckApprovalPermission(approverID, currentNode, nodeInstance, request) {
		return nil, fmt.Errorf("no permission to process this approval")
	}

	var approver models.User
	s.db.First(&approver, approverID)

	updatedInstance, err := s.processEngine.ProcessApproval(
		request,
		currentNode,
		approverID,
		action,
		comment,
		nodeInstance,
		&approver,
		&instance,
		formData,
	)
	if err != nil {
		return nil, err
	}

	approvalStatus := s.processEngine.GetApprovalStatus(updatedInstance)

	err = s.db.Transaction(func(tx *gorm.DB) error {
		instanceJSON, _ := json.Marshal(updatedInstance)
		request.ProcessInstance = string(instanceJSON)
		request.CurrentNodeID = updatedInstance.CurrentNodeID
		request.Status = approvalStatus

		if err := tx.Save(request).Error; err != nil {
			return err
		}

		record := &models.ApprovalRecord{
			RequestID:  requestID,
			NodeID:     currentNode.ID,
			NodeName:   currentNode.Label,
			ApproverID: approverID,
			Action:     action,
			Comment:    comment,
			CreatedAt:  time.Now(),
		}
		if err := tx.Create(record).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return request, nil
}

func (s *ApprovalService) RevokeRequest(requestID uint, userID uint, comment string) (*models.ApprovalRequest, error) {
	request, err := s.GetRequestByID(requestID)
	if err != nil {
		return nil, err
	}

	if request.ApplicantID != userID {
		return nil, fmt.Errorf("only applicant can revoke the request")
	}

	if request.Status == models.ApprovalStatusApproved || request.Status == models.ApprovalStatusRejected {
		return nil, fmt.Errorf("cannot revoke completed request")
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		request.Status = models.ApprovalStatusRevoked
		if err := tx.Save(request).Error; err != nil {
			return err
		}

		record := &models.ApprovalRecord{
			RequestID:  requestID,
			NodeID:     request.CurrentNodeID,
			NodeName:   "撤回",
			ApproverID: userID,
			Action:     models.ApprovalActionRevoke,
			Comment:    comment,
			CreatedAt:  time.Now(),
		}
		if err := tx.Create(record).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return request, nil
}

func (s *ApprovalService) AddSign(requestID uint, approverID uint, signUserID uint, comment string) (*models.ApprovalRequest, error) {
	request, err := s.GetRequestByID(requestID)
	if err != nil {
		return nil, err
	}

	if request.Status != models.ApprovalStatusPending {
		return nil, fmt.Errorf("approval request is not pending")
	}

	var instance models.ProcessInstance
	err = json.Unmarshal([]byte(request.ProcessInstance), &instance)
	if err != nil {
		return nil, err
	}

	template, err := s.templateService.GetTemplateByID(request.TemplateID)
	if err != nil {
		return nil, err
	}

	nodes, err := s.processEngine.ParseNodes(template)
	if err != nil {
		return nil, err
	}

	currentNode := s.processEngine.FindNodeByID(nodes, instance.CurrentNodeID)
	if currentNode == nil {
		return nil, fmt.Errorf("current node not found")
	}

	if !s.processEngine.CheckApprovalPermission(approverID, currentNode, nil, request) {
		return nil, fmt.Errorf("no permission to add sign")
	}

	updatedInstance, err := s.processEngine.AddSign(&instance, instance.CurrentNodeID, signUserID, approverID, comment)
	if err != nil {
		return nil, err
	}

	var approver models.User
	s.db.First(&approver, approverID)

	var signUser models.User
	s.db.First(&signUser, signUserID)

	err = s.db.Transaction(func(tx *gorm.DB) error {
		instanceJSON, _ := json.Marshal(updatedInstance)
		request.ProcessInstance = string(instanceJSON)

		if err := tx.Save(request).Error; err != nil {
			return err
		}

		record := &models.ApprovalRecord{
			RequestID:  requestID,
			NodeID:     instance.CurrentNodeID,
			NodeName:   currentNode.Label,
			ApproverID: approverID,
			Action:     models.ApprovalActionAddSign,
			Comment:    fmt.Sprintf("%s 加签给 %s: %s", approver.Name, signUser.Name, comment),
			SignType:   "add_sign",
			CreatedAt:  time.Now(),
		}
		if err := tx.Create(record).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return request, nil
}

func (s *ApprovalService) TransferApproval(requestID uint, approverID uint, targetUserID uint, comment string) (*models.ApprovalRequest, error) {
	request, err := s.GetRequestByID(requestID)
	if err != nil {
		return nil, err
	}

	if request.Status != models.ApprovalStatusPending {
		return nil, fmt.Errorf("approval request is not pending")
	}

	var instance models.ProcessInstance
	err = json.Unmarshal([]byte(request.ProcessInstance), &instance)
	if err != nil {
		return nil, err
	}

	template, err := s.templateService.GetTemplateByID(request.TemplateID)
	if err != nil {
		return nil, err
	}

	nodes, err := s.processEngine.ParseNodes(template)
	if err != nil {
		return nil, err
	}

	currentNode := s.processEngine.FindNodeByID(nodes, instance.CurrentNodeID)
	if currentNode == nil {
		return nil, fmt.Errorf("current node not found")
	}

	if !s.processEngine.CheckApprovalPermission(approverID, currentNode, nil, request) {
		return nil, fmt.Errorf("no permission to transfer approval")
	}

	var approver models.User
	s.db.First(&approver, approverID)

	updatedInstance, err := s.processEngine.TransferApproval(request, currentNode, approverID, targetUserID, comment, &approver, &instance)
	if err != nil {
		return nil, err
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		instanceJSON, _ := json.Marshal(updatedInstance)
		request.ProcessInstance = string(instanceJSON)

		if err := tx.Save(request).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return request, nil
}

func (s *ApprovalService) getNodesFromInstance(instance *models.ProcessInstance) []models.Node {
	nodes := make([]models.Node, 0, len(instance.Nodes))
	for _, ni := range instance.Nodes {
		nodes = append(nodes, models.Node{
			ID: ni.NodeID,
		})
	}
	return nodes
}

func (s *ApprovalService) GetProcessInstance(request *models.ApprovalRequest) (*models.ProcessInstance, error) {
	var instance models.ProcessInstance
	err := json.Unmarshal([]byte(request.ProcessInstance), &instance)
	if err != nil {
		return nil, err
	}
	return &instance, nil
}

func (s *ApprovalService) GetStats(userID uint) (map[string]interface{}, error) {
	var todoCount int64
	err := s.db.Model(&models.ApprovalRequest{}).
		Where("status IN ?", []string{models.ApprovalStatusPending}).
		Where("EXISTS (SELECT 1 FROM approval_records ar WHERE ar.request_id = approval_requests.id AND ar.approver_id = ? AND ar.action IN ?)",
			userID, []string{models.ApprovalActionStart, models.ApprovalActionApprove, models.ApprovalActionAddSign}).
		Where("NOT EXISTS (SELECT 1 FROM approval_records ar WHERE ar.request_id = approval_requests.id AND ar.approver_id = ? AND ar.action IN ?)",
			userID, []string{models.ApprovalActionApprove, models.ApprovalActionReject}).
		Count(&todoCount).Error
	if err != nil {
		return nil, err
	}

	var myCount int64
	err = s.db.Model(&models.ApprovalRequest{}).Where("applicant_id = ?", userID).Count(&myCount).Error
	if err != nil {
		return nil, err
	}

	var approvedCount int64
	err = s.db.Model(&models.ApprovalRequest{}).
		Where("status = ? AND applicant_id = ?", models.ApprovalStatusApproved, userID).
		Count(&approvedCount).Error
	if err != nil {
		return nil, err
	}

	var templateCount int64
	err = s.db.Model(&models.ProcessTemplate{}).Count(&templateCount).Error
	if err != nil {
		return nil, err
	}

	var recentApprovals []models.ApprovalRequest
	err = s.db.Model(&models.ApprovalRequest{}).Preload("Applicant").Preload("Template").
		Where("status IN ?", []string{models.ApprovalStatusPending}).
		Where("EXISTS (SELECT 1 FROM approval_records ar WHERE ar.request_id = approval_requests.id AND ar.approver_id = ?)", userID).
		Order("created_at DESC").Limit(5).Find(&recentApprovals).Error
	if err != nil {
		return nil, err
	}

	var myRequests []models.ApprovalRequest
	err = s.db.Model(&models.ApprovalRequest{}).Preload("Template").
		Where("applicant_id = ?", userID).
		Order("created_at DESC").Limit(5).Find(&myRequests).Error
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"todo_count":      todoCount,
		"my_count":        myCount,
		"approved_count":  approvedCount,
		"template_count":  templateCount,
		"recent_approvals": recentApprovals,
		"my_requests":     myRequests,
	}, nil
}
