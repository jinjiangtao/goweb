package storage

import (
	"database/sql"
	"fmt"
	"goim/server/model"
	"time"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("sqlite", "./data/goim.db")
	if err != nil {
		panic(fmt.Sprintf("failed to open database: %v", err))
	}

	createTables()
}

func createTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			nickname TEXT NOT NULL,
			avatar TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS friends (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			friend_id TEXT NOT NULL,
			status INTEGER DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(friend_id) REFERENCES users(id),
			UNIQUE(user_id, friend_id)
		)`,
		`CREATE TABLE IF NOT EXISTS friend_requests (
			id TEXT PRIMARY KEY,
			sender_id TEXT NOT NULL,
			receiver_id TEXT NOT NULL,
			status INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(sender_id) REFERENCES users(id),
			FOREIGN KEY(receiver_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS groups (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			avatar TEXT,
			owner_id TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(owner_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS group_members (
			id TEXT PRIMARY KEY,
			group_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			role INTEGER DEFAULT 0,
			joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(group_id) REFERENCES groups(id),
			FOREIGN KEY(user_id) REFERENCES users(id),
			UNIQUE(group_id, user_id)
		)`,
		`CREATE TABLE IF NOT EXISTS messages (
			id TEXT PRIMARY KEY,
			sender_id TEXT NOT NULL,
			receiver_id TEXT NOT NULL,
			receiver_type INTEGER NOT NULL,
			content TEXT NOT NULL,
			type INTEGER DEFAULT 0,
			status INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(sender_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS message_read (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			message_id TEXT NOT NULL,
			read_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(message_id) REFERENCES messages(id)
		)`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			panic(fmt.Sprintf("failed to create table: %v", err))
		}
	}
}

func CreateUser(user *model.User) error {
	_, err := db.Exec(`INSERT INTO users (id, username, password, nickname, avatar, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		user.ID, user.Username, user.Password, user.Nickname, user.Avatar, user.CreatedAt, user.UpdatedAt)
	return err
}

func GetUserByUsername(username string) (*model.User, error) {
	user := &model.User{}
	err := db.QueryRow(`SELECT id, username, password, nickname, avatar, created_at, updated_at 
		FROM users WHERE username = ?`, username).Scan(
		&user.ID, &user.Username, &user.Password, &user.Nickname, &user.Avatar, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func GetUserByID(id string) (*model.User, error) {
	user := &model.User{}
	err := db.QueryRow(`SELECT id, username, password, nickname, avatar, created_at, updated_at 
		FROM users WHERE id = ?`, id).Scan(
		&user.ID, &user.Username, &user.Password, &user.Nickname, &user.Avatar, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func CreateFriend(friend *model.Friend) error {
	_, err := db.Exec(`INSERT INTO friends (id, user_id, friend_id, status, created_at) 
		VALUES (?, ?, ?, ?, ?)`,
		friend.ID, friend.UserID, friend.FriendID, friend.Status, friend.CreatedAt)
	return err
}

func GetFriends(userID string) ([]*model.Friend, error) {
	rows, err := db.Query(`SELECT id, user_id, friend_id, status, created_at 
		FROM friends WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []*model.Friend
	for rows.Next() {
		f := &model.Friend{}
		err := rows.Scan(&f.ID, &f.UserID, &f.FriendID, &f.Status, &f.CreatedAt)
		if err != nil {
			return nil, err
		}
		friends = append(friends, f)
	}
	return friends, nil
}

func CreateGroup(group *model.Group) error {
	_, err := db.Exec(`INSERT INTO groups (id, name, avatar, owner_id, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?)`,
		group.ID, group.Name, group.Avatar, group.OwnerID, group.CreatedAt, group.UpdatedAt)
	return err
}

func CreateGroupMember(member *model.GroupMember) error {
	_, err := db.Exec(`INSERT INTO group_members (id, group_id, user_id, role, joined_at) 
		VALUES (?, ?, ?, ?, ?)`,
		member.ID, member.GroupID, member.UserID, member.Role, member.JoinedAt)
	return err
}

func RemoveGroupMember(groupID, userID string) error {
	_, err := db.Exec(`DELETE FROM group_members WHERE group_id = ? AND user_id = ?`, groupID, userID)
	return err
}

func GetGroups(userID string) ([]*model.Group, error) {
	rows, err := db.Query(`SELECT g.id, g.name, g.avatar, g.owner_id, g.created_at, g.updated_at 
		FROM groups g 
		JOIN group_members gm ON g.id = gm.group_id 
		WHERE gm.user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*model.Group
	for rows.Next() {
		g := &model.Group{}
		err := rows.Scan(&g.ID, &g.Name, &g.Avatar, &g.OwnerID, &g.CreatedAt, &g.UpdatedAt)
		if err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}

func GetGroupMembers(groupID string) ([]*model.GroupMember, error) {
	rows, err := db.Query(`SELECT id, group_id, user_id, role, joined_at 
		FROM group_members WHERE group_id = ?`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*model.GroupMember
	for rows.Next() {
		m := &model.GroupMember{}
		err := rows.Scan(&m.ID, &m.GroupID, &m.UserID, &m.Role, &m.JoinedAt)
		if err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}

func CreateMessage(msg *model.Message) error {
	_, err := db.Exec(`INSERT INTO messages (id, sender_id, receiver_id, receiver_type, content, type, status, created_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		msg.ID, msg.SenderID, msg.ReceiverID, msg.ReceiverType, msg.Content, msg.Type, msg.Status, msg.CreatedAt)
	return err
}

func GetMessagesByUser(userID string, limit, offset int) ([]*model.Message, error) {
	rows, err := db.Query(`SELECT id, sender_id, receiver_id, receiver_type, content, type, status, created_at 
		FROM messages 
		WHERE sender_id = ? OR (receiver_id = ? AND receiver_type = 0)
		ORDER BY created_at DESC 
		LIMIT ? OFFSET ?`, userID, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		m := &model.Message{}
		err := rows.Scan(&m.ID, &m.SenderID, &m.ReceiverID, &m.ReceiverType, &m.Content, &m.Type, &m.Status, &m.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}

func GetGroupMessages(groupID string, limit, offset int) ([]*model.Message, error) {
	rows, err := db.Query(`SELECT id, sender_id, receiver_id, receiver_type, content, type, status, created_at 
		FROM messages 
		WHERE receiver_id = ? AND receiver_type = 1
		ORDER BY created_at DESC 
		LIMIT ? OFFSET ?`, groupID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		m := &model.Message{}
		err := rows.Scan(&m.ID, &m.SenderID, &m.ReceiverID, &m.ReceiverType, &m.Content, &m.Type, &m.Status, &m.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}

func GetUnreadMessages(userID string) ([]*model.Message, error) {
	rows, err := db.Query(`SELECT m.id, m.sender_id, m.receiver_id, m.receiver_type, m.content, m.type, m.status, m.created_at 
		FROM messages m
		LEFT JOIN message_read mr ON m.id = mr.message_id AND mr.user_id = ?
		WHERE m.receiver_id = ? AND m.receiver_type = 0 AND mr.id IS NULL
		ORDER BY m.created_at DESC`, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		m := &model.Message{}
		err := rows.Scan(&m.ID, &m.SenderID, &m.ReceiverID, &m.ReceiverType, &m.Content, &m.Type, &m.Status, &m.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}

func MarkMessageRead(userID, messageID string) error {
	_, err := db.Exec(`INSERT OR REPLACE INTO message_read (id, user_id, message_id, read_at) 
		VALUES (?, ?, ?, ?)`,
		fmt.Sprintf("%s_%s", userID, messageID), userID, messageID, time.Now())
	return err
}

func CheckFriendExists(userID, friendID string) (bool, error) {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM friends WHERE user_id = ? AND friend_id = ?`, userID, friendID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func CheckGroupMember(groupID, userID string) (bool, error) {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM group_members WHERE group_id = ? AND user_id = ?`, groupID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetGroupByID(groupID string) (*model.Group, error) {
	group := &model.Group{}
	err := db.QueryRow(`SELECT id, name, avatar, owner_id, created_at, updated_at FROM groups WHERE id = ?`, groupID).Scan(
		&group.ID, &group.Name, &group.Avatar, &group.OwnerID, &group.CreatedAt, &group.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return group, err
}

func GetAllUsers(limit, offset int) ([]*model.User, int, error) {
	// 获取总数
	var total int
	err := db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 获取分页用户
	rows, err := db.Query(`SELECT id, username, password, nickname, avatar, created_at, updated_at 
		FROM users 
		ORDER BY created_at DESC 
		LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		user := &model.User{}
		err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Nickname, &user.Avatar, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, total, nil
}
