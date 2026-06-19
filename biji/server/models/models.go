package models

import (
	"time"
)

type Note struct {
	ID          string    `gorm:"primaryKey;size:36" json:"id"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Content     string    `gorm:"type:text" json:"content"`
	PlainText   string    `gorm:"type:text" json:"plainText"`
	IsEncrypted bool      `gorm:"default:false" json:"isEncrypted"`
	Password    string    `gorm:"size:255" json:"-"`
	Salt        string    `gorm:"size:32" json:"-"`
	IsPinned    bool      `gorm:"default:false;index" json:"isPinned"`
	IsArchived  bool      `gorm:"default:false;index" json:"isArchived"`
	IsDeleted   bool      `gorm:"default:false;index" json:"isDeleted"`
	Tags        []Tag     `gorm:"many2many:note_tags;constraint:OnDelete:CASCADE" json:"tags"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   time.Time `json:"deletedAt,omitempty"`
}

type Tag struct {
	ID        string    `gorm:"primaryKey;size:36" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Color     string    `gorm:"size:20;default:'#409EFF'" json:"color"`
	ParentID  *string   `gorm:"size:36;index" json:"parentId"`
	Sort      int       `gorm:"default:0;index" json:"sort"`
	GroupID   *string   `gorm:"size:36;index" json:"groupId"`
	Children  []Tag     `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Notes     []Note    `gorm:"many2many:note_tags" json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type TagGroup struct {
	ID        string    `gorm:"primaryKey;size:36" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Sort      int       `gorm:"default:0;index" json:"sort"`
	Tags      []Tag     `gorm:"foreignKey:GroupID" json:"tags,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type NoteTag struct {
	NoteID string `gorm:"primaryKey;size:36" json:"noteId"`
	TagID  string `gorm:"primaryKey;size:36" json:"tagId"`
}
