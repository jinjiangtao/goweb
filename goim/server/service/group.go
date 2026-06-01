package service

import (
	"goim/server/model"
	"goim/server/storage"
	"time"

	"github.com/google/uuid"
)

func CreateGroup(name, ownerID string) (*model.Group, error) {
	group := &model.Group{
		ID:        uuid.New().String(),
		Name:      name,
		Avatar:    "",
		OwnerID:   ownerID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := storage.CreateGroup(group)
	if err != nil {
		return nil, err
	}

	member := &model.GroupMember{
		ID:       uuid.New().String(),
		GroupID:  group.ID,
		UserID:   ownerID,
		Role:     1,
		JoinedAt: time.Now(),
	}

	err = storage.CreateGroupMember(member)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func AddGroupMember(groupID, userID string) error {
	exists, err := storage.CheckGroupMember(groupID, userID)
	if err != nil {
		return err
	}
	if exists {
		return ErrMemberExists
	}

	member := &model.GroupMember{
		ID:       uuid.New().String(),
		GroupID:  groupID,
		UserID:   userID,
		Role:     0,
		JoinedAt: time.Now(),
	}

	return storage.CreateGroupMember(member)
}

func RemoveGroupMember(groupID, userID string) error {
	return storage.RemoveGroupMember(groupID, userID)
}

func GetGroups(userID string) ([]*model.Group, error) {
	return storage.GetGroups(userID)
}

func GetGroupMembers(groupID string) ([]*model.GroupMember, error) {
	return storage.GetGroupMembers(groupID)
}

func GetGroupByID(groupID string) (*model.Group, error) {
	return storage.GetGroupByID(groupID)
}