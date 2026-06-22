package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type StyleConfig struct {
	Size        int    `json:"size"`
	Foreground  string `json:"foreground"`
	Background  string `json:"background"`
	CornerStyle string `json:"cornerStyle"`
	EccLevel    string `json:"eccLevel"`
	Margin      int    `json:"margin"`
	Logo        string `json:"logo"`
}

func (s *StyleConfig) Scan(src interface{}) error {
	var b []byte
	switch v := src.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	case nil:
		return nil
	default:
		return errors.New("style_config: unsupported scan type")
	}
	if len(b) == 0 {
		return nil
	}
	return json.Unmarshal(b, s)
}

func (s StyleConfig) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s StyleConfig) JSON() string {
	b, _ := json.Marshal(s)
	return string(b)
}

func DefaultStyle() StyleConfig {
	return StyleConfig{
		Size:        512,
		Foreground:  "#0E0E11",
		Background:  "#FFFFFF",
		CornerStyle: "square",
		EccLevel:    "M",
		Margin:      4,
		Logo:        "",
	}
}

func (s StyleConfig) WithDefaults() StyleConfig {
	d := DefaultStyle()
	if s.Size == 0 {
		s.Size = d.Size
	}
	if s.Foreground == "" {
		s.Foreground = d.Foreground
	}
	if s.Background == "" {
		s.Background = d.Background
	}
	if s.CornerStyle == "" {
		s.CornerStyle = d.CornerStyle
	}
	if s.EccLevel == "" {
		s.EccLevel = d.EccLevel
	}
	if s.Margin == 0 {
		s.Margin = d.Margin
	}
	return s
}

type Record struct {
	ID          int64       `json:"id" db:"id"`
	Content     string      `json:"content" db:"content"`
	ContentType string      `json:"contentType" db:"content_type"`
	Source      string      `json:"source" db:"source"`
	StyleConfig StyleConfig `json:"styleConfig" db:"style_config"`
	FilePath    string      `json:"filePath" db:"file_path"`
	Status      string      `json:"status" db:"status"`
	TemplateID  *int64      `json:"templateId,omitempty" db:"template_id"`
	BatchID     string      `json:"batchId,omitempty" db:"batch_id"`
	Remark      string      `json:"remark,omitempty" db:"remark"`
	CreatedAt   string      `json:"createdAt" db:"created_at"`
	UpdatedAt   string      `json:"updatedAt" db:"updated_at"`
}

type Template struct {
	ID          int64       `json:"id" db:"id"`
	Name        string      `json:"name" db:"name"`
	StyleConfig StyleConfig `json:"styleConfig" db:"style_config"`
	CreatedAt   string      `json:"createdAt" db:"created_at"`
	UpdatedAt   string      `json:"updatedAt" db:"updated_at"`
}

type ParseLog struct {
	ID        int64  `json:"id" db:"id"`
	Source    string `json:"source" db:"source"`
	Result    string `json:"result" db:"result"`
	Type      string `json:"type" db:"type"`
	CreatedAt string `json:"createdAt" db:"created_at"`
}

type ParseResult struct {
	Raw     string `json:"raw"`
	Type    string `json:"type"`
	Content string `json:"content"`
}

type Stats struct {
	Total     int64    `json:"total"`
	Today     int64    `json:"today"`
	Templates int64    `json:"templates"`
	Invalid   int64    `json:"invalid"`
	Recent    []Record `json:"recent"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func OK(data interface{}) Response {
	return Response{Code: 0, Message: "ok", Data: data}
}

func Err(msg string) Response {
	return Response{Code: 1, Message: msg, Data: nil}
}
