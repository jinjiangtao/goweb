package service

import (
	"errors"
	"goim/server/model"
	"goim/server/storage"
	"log"
	"time"

	"github.com/google/uuid"
)

var (
	ErrOwnerExists       = errors.New("owner already exists")
	ErrOwnerMemberExists = errors.New("member already joined")
	ErrOwnerNotFound     = errors.New("owner not found")
)

func CreateOwner(name, description, avatar string) (*model.Owner, error) {
	owner := &model.Owner{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Avatar:      avatar,
		CreatedAt:   time.Now(),
	}

	err := storage.CreateOwner(owner)
	if err != nil {
		return nil, err
	}

	return owner, nil
}

func GetOwnerByID(ownerID string) (*model.Owner, error) {
	return storage.GetOwnerByID(ownerID)
}

func GetAllOwners() ([]*model.Owner, error) {
	return storage.GetAllOwners()
}

func DeleteOwner(ownerID string) error {
	owner, err := storage.GetOwnerByID(ownerID)
	if err != nil {
		return err
	}
	if owner == nil {
		return ErrOwnerNotFound
	}

	err = storage.RemoveAllOwnerMembers(ownerID)
	if err != nil {
		return err
	}

	return storage.DeleteOwner(ownerID)
}

func JoinOwner(ownerID, userID string) error {
	log.Printf("CheckOwnerMember: ownerID=%s, userID=%s", ownerID, userID)
	exists, err := storage.CheckOwnerMember(ownerID, userID)
	log.Printf("CheckOwnerMember result: exists=%v, err=%v", exists, err)
	if err != nil {
		return err
	}
	if exists {
		return ErrOwnerMemberExists
	}

	member := &model.OwnerMember{
		ID:       uuid.New().String(),
		OwnerID:  ownerID,
		UserID:   userID,
		JoinedAt: time.Now(),
	}

	return storage.CreateOwnerMember(member)
}

func LeaveOwner(ownerID, userID string) error {
	return storage.RemoveOwnerMember(ownerID, userID)
}

func GetOwnerMembers(ownerID string) ([]*model.OwnerMember, error) {
	return storage.GetOwnerMembers(ownerID)
}

func GetOwnersByUserID(userID string) ([]*model.Owner, error) {
	return storage.GetOwnersByUserID(userID)
}
