package service

import (
	"errors"
	"goim/server/model"
	"goim/server/storage"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserExists      = errors.New("user already exists")
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidToken    = errors.New("invalid token")
	ErrMemberExists    = errors.New("member already exists")
)

func AddFriend(userID, friendID string) error {
	if userID == friendID {
		return errors.New("cannot add yourself as friend")
	}

	friend, err := storage.GetUserByID(friendID)
	if err != nil {
		return err
	}
	if friend == nil {
		return ErrUserNotFound
	}

	exists, err := storage.CheckFriendExists(userID, friendID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("friend already exists")
	}

	f1 := &model.Friend{
		ID:        uuid.New().String(),
		UserID:    userID,
		FriendID:  friendID,
		Status:    1,
		CreatedAt: time.Now(),
	}

	f2 := &model.Friend{
		ID:        uuid.New().String(),
		UserID:    friendID,
		FriendID:  userID,
		Status:    1,
		CreatedAt: time.Now(),
	}

	err = storage.CreateFriend(f1)
	if err != nil {
		return err
	}

	return storage.CreateFriend(f2)
}

func GetFriends(userID string) ([]*model.User, error) {
	friendRelations, err := storage.GetFriends(userID)
	if err != nil {
		return nil, err
	}

	var friends []*model.User
	for _, fr := range friendRelations {
		user, err := storage.GetUserByID(fr.FriendID)
		if err != nil {
			return nil, err
		}
		if user != nil {
			friends = append(friends, user)
		}
	}

	return friends, nil
}

func GetUserByID(userID string) (*model.User, error) {
	return storage.GetUserByID(userID)
}

func CheckFriendship(userID, friendID string) (bool, error) {
	return storage.CheckFriendExists(userID, friendID)
}

func GetAllUsers(limit, offset int) ([]*model.User, int, error) {
	return storage.GetAllUsers(limit, offset)
}
