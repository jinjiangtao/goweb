package cache

import (
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	onlineUsers *cache.Cache
	unreadCount sync.Map
)

func InitCache() {
	onlineUsers = cache.New(5*time.Minute, 10*time.Minute)
}

func SetOnline(userID string) {
	onlineUsers.Set(userID, true, cache.DefaultExpiration)
}

func SetOffline(userID string) {
	onlineUsers.Delete(userID)
}

func IsOnline(userID string) bool {
	_, found := onlineUsers.Get(userID)
	return found
}

func GetOnlineUsers() []string {
	var users []string
	items := onlineUsers.Items()
	for key := range items {
		users = append(users, key)
	}
	return users
}

func IncrUnreadCount(userID, senderID string) {
	key := userID + ":" + senderID
	val, _ := unreadCount.LoadOrStore(key, 0)
	unreadCount.Store(key, val.(int)+1)
}

func GetUnreadCount(userID, senderID string) int {
	key := userID + ":" + senderID
	val, found := unreadCount.Load(key)
	if !found {
		return 0
	}
	return val.(int)
}

func ClearUnreadCount(userID, senderID string) {
	key := userID + ":" + senderID
	unreadCount.Delete(key)
}

func GetAllUnreadCounts(userID string) map[string]int {
	result := make(map[string]int)
	unreadCount.Range(func(k, v interface{}) bool {
		key := k.(string)
		if len(key) > len(userID) && key[:len(userID)] == userID && key[len(userID)] == ':' {
			senderID := key[len(userID)+1:]
			result[senderID] = v.(int)
		}
		return true
	})
	return result
}