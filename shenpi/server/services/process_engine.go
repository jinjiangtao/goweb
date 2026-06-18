package services

import (
	"encoding/json"
	"fmt"
	"shenpi/models"
	"strconv"
	"strings"
	"time"
)

type ProcessEngine struct{}

func NewProcessEngine() *ProcessEngine {
	return &ProcessEngine{}
}

func (e *ProcessEngine) ParseNodes(template *models.ProcessTemplate) ([]models.Node, error) {
	var nodes []models.Node
	if template.Nodes == "" {
		return nodes, nil
	}
	err := json.Unmarshal([]byte(template.Nodes), &nodes)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func (e *ProcessEngine) ParseEdges(template *models.ProcessTemplate) ([]models.Edge, error) {
	var edges []models.Edge
	if template.Edges == "" {
		return edges, nil
	}
	err := json.Unmarshal([]byte(template.Edges), &edges)
	if err != nil {
		return nil, err
	}
	return edges, nil
}

func (e *ProcessEngine) ParseFormConfig(template *models.ProcessTemplate) ([]models.FormField, error) {
	var fields []models.FormField
	if template.FormConfig == "" {
		return fields, nil
	}
	err := json.Unmarshal([]byte(template.FormConfig), &fields)
	if err != nil {
		return nil, err
	}
	return fields, nil
}

func (e *ProcessEngine) InitProcessInstance(template *models.ProcessTemplate) (*models.ProcessInstance, error) {
	nodes, err := e.ParseNodes(template)
	if err != nil {
		return nil, err
	}

	edges, err := e.ParseEdges(template)
	if err != nil {
		return nil, err
	}

	nodeInstances := make([]models.NodeInstance, 0, len(nodes))
	for _, node := range nodes {
		nodeInstances = append(nodeInstances, models.NodeInstance{
			NodeID:   node.ID,
			Status:   models.NodeStatusPending,
			Approved: make([]uint, 0),
			Rejected: make([]uint, 0),
			CCUsers:  make([]uint, 0),
		})
	}

	startNode := e.FindStartNode(nodes)
	instance := &models.ProcessInstance{
		Nodes:         nodeInstances,
		Edges:         edges,
		CurrentNodeID: startNode.ID,
		History:       make([]models.HistoryRecord, 0),
	}

	return instance, nil
}

func (e *ProcessEngine) FindStartNode(nodes []models.Node) *models.Node {
	for _, node := range nodes {
		if node.Type == "start" {
			return &node
		}
	}
	return nil
}

func (e *ProcessEngine) FindEndNode(nodes []models.Node) *models.Node {
	for _, node := range nodes {
		if node.Type == "end" {
			return &node
		}
	}
	return nil
}

func (e *ProcessEngine) FindNodeByID(nodes []models.Node, nodeID string) *models.Node {
	for _, node := range nodes {
		if node.ID == nodeID {
			return &node
		}
	}
	return nil
}

func (e *ProcessEngine) FindNodeInstanceByID(instance *models.ProcessInstance, nodeID string) *models.NodeInstance {
	for i := range instance.Nodes {
		if instance.Nodes[i].NodeID == nodeID {
			return &instance.Nodes[i]
		}
	}
	return nil
}

func (e *ProcessEngine) GetOutgoingEdges(edges []models.Edge, nodeID string) []models.Edge {
	result := make([]models.Edge, 0)
	for _, edge := range edges {
		if edge.Source == nodeID {
			result = append(result, edge)
		}
	}
	return result
}

func (e *ProcessEngine) EvaluateCondition(condition string, formData map[string]interface{}) bool {
	if condition == "" {
		return true
	}

	parts := strings.Split(condition, " ")
	if len(parts) < 3 {
		return true
	}

	field := strings.TrimSpace(parts[0])
	operator := strings.TrimSpace(parts[1])
	value := strings.TrimSpace(parts[2])

	fieldValue, exists := formData[field]
	if !exists {
		return false
	}

	fieldValueStr := fmt.Sprintf("%v", fieldValue)

	switch operator {
	case "==", "=":
		return fieldValueStr == value
	case "!=":
		return fieldValueStr != value
	case ">":
		num1, _ := strconv.ParseFloat(fieldValueStr, 64)
		num2, _ := strconv.ParseFloat(value, 64)
		return num1 > num2
	case ">=":
		num1, _ := strconv.ParseFloat(fieldValueStr, 64)
		num2, _ := strconv.ParseFloat(value, 64)
		return num1 >= num2
	case "<":
		num1, _ := strconv.ParseFloat(fieldValueStr, 64)
		num2, _ := strconv.ParseFloat(value, 64)
		return num1 < num2
	case "<=":
		num1, _ := strconv.ParseFloat(fieldValueStr, 64)
		num2, _ := strconv.ParseFloat(value, 64)
		return num1 <= num2
	case "contains":
		return strings.Contains(fieldValueStr, value)
	default:
		return true
	}
}

func (e *ProcessEngine) DetermineNextNode(instance *models.ProcessInstance, currentNodeID string, formData map[string]interface{}) (string, error) {
	edges := e.GetOutgoingEdges(instance.Edges, currentNodeID)

	for _, edge := range edges {
		if e.EvaluateCondition(edge.Condition, formData) {
			return edge.Target, nil
		}
	}

	return "", fmt.Errorf("no valid next node found")
}

func (e *ProcessEngine) CheckApprovalPermission(userID uint, node *models.Node, instance *models.NodeInstance, request *models.ApprovalRequest) bool {
	if node == nil {
		return false
	}

	switch node.ApproverType {
	case "specified":
		for _, approverID := range node.Approvers {
			if approverID == userID {
				return true
			}
		}
		return false
	case "applicant":
		return request.ApplicantID == userID
	case "all":
		return true
	default:
		for _, approverID := range node.Approvers {
			if approverID == userID {
				return true
			}
		}
		return false
	}
}

func (e *ProcessEngine) CheckNodeComplete(node *models.Node, nodeInstance *models.NodeInstance) bool {
	if node == nil {
		return true
	}

	if node.Type == "start" || node.Type == "end" {
		return true
	}

	if len(node.Approvers) == 0 {
		return true
	}
	return len(nodeInstance.Approved) >= len(node.Approvers)
}

func (e *ProcessEngine) ProcessApproval(request *models.ApprovalRequest, node *models.Node, approverID uint, action string, comment string, nodeInstance *models.NodeInstance, currentUser *models.User, instance *models.ProcessInstance, formData map[string]interface{}) (*models.ProcessInstance, error) {
	nodeInstance.Status = models.NodeStatusProcessing

	historyRecord := models.HistoryRecord{
		NodeID:     node.ID,
		NodeName:   node.Label,
		Action:     action,
		ApproverID: approverID,
		Approver:   currentUser.Name,
		Comment:    comment,
		Timestamp:  time.Now(),
	}

	if action == models.ApprovalActionApprove {
		nodeInstance.Approved = append(nodeInstance.Approved, approverID)
	} else if action == models.ApprovalActionReject {
		nodeInstance.Rejected = append(nodeInstance.Rejected, approverID)
		nodeInstance.Status = models.NodeStatusRejected
		instance.History = append(instance.History, historyRecord)
		return instance, nil
	}

	if e.CheckNodeComplete(node, nodeInstance) {
		nodeInstance.Status = models.NodeStatusApproved

		if node.Type == "end" {
			instance.History = append(instance.History, historyRecord)
			return instance, nil
		}

		nextNodeID, err := e.DetermineNextNode(instance, node.ID, formData)
		if err != nil {
			return nil, err
		}

		instance.CurrentNodeID = nextNodeID
		nextNodeInstance := e.FindNodeInstanceByID(instance, nextNodeID)
		if nextNodeInstance != nil {
			nextNodeInstance.Status = models.NodeStatusProcessing
		}
	}

	instance.History = append(instance.History, historyRecord)
	return instance, nil
}

func (e *ProcessEngine) AddSign(instance *models.ProcessInstance, nodeID string, signUserID uint, currentUserID uint, comment string) (*models.ProcessInstance, error) {
	nodeInstance := e.FindNodeInstanceByID(instance, nodeID)
	if nodeInstance == nil {
		return nil, fmt.Errorf("node not found")
	}

	nodeInstance.CCUsers = append(nodeInstance.CCUsers, signUserID)

	return instance, nil
}

func (e *ProcessEngine) TransferApproval(request *models.ApprovalRequest, node *models.Node, approverID uint, targetUserID uint, comment string, currentUser *models.User, instance *models.ProcessInstance) (*models.ProcessInstance, error) {
	nodeInstance := e.FindNodeInstanceByID(instance, node.ID)
	if nodeInstance == nil {
		return nil, fmt.Errorf("node not found")
	}

	node.Approvers = append(node.Approvers, targetUserID)

	historyRecord := models.HistoryRecord{
		NodeID:     node.ID,
		NodeName:   node.Label,
		Action:     models.ApprovalActionTransfer,
		ApproverID: approverID,
		Approver:   currentUser.Name,
		Comment:    comment + fmt.Sprintf(" 转交给用户ID: %d", targetUserID),
		Timestamp:  time.Now(),
	}

	instance.History = append(instance.History, historyRecord)

	return instance, nil
}

func (e *ProcessEngine) GetApprovalStatus(instance *models.ProcessInstance) string {
	endNodeInstance := e.FindNodeInstanceByID(instance, "end")
	if endNodeInstance != nil && endNodeInstance.Status == models.NodeStatusApproved {
		return models.ApprovalStatusApproved
	}

	for _, nodeInstance := range instance.Nodes {
		if nodeInstance.Status == models.NodeStatusRejected {
			return models.ApprovalStatusRejected
		}
	}

	return models.ApprovalStatusPending
}
