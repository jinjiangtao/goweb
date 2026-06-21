package models

import (
	"encoding/json"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	_ "modernc.org/sqlite"
)

var DB *gorm.DB

func InitDB(dbPath string) error {
	var err error
	DB, err = gorm.Open(sqlite.Dialector{
		DriverName: "sqlite",
		DSN:        dbPath,
	}, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	err = DB.AutoMigrate(&Board{}, &Element{}, &History{}, &OperationLog{})
	if err != nil {
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}

func (e *Element) Encode() error {
	pointsJSON, err := json.Marshal(e.PointsArr)
	if err != nil {
		return err
	}
	e.Points = string(pointsJSON)

	styleJSON, err := json.Marshal(e.StyleObj)
	if err != nil {
		return err
	}
	e.Style = string(styleJSON)

	return nil
}

func (e *Element) Decode() error {
	if e.Points != "" {
		err := json.Unmarshal([]byte(e.Points), &e.PointsArr)
		if err != nil {
			return err
		}
	}

	if e.Style != "" {
		err := json.Unmarshal([]byte(e.Style), &e.StyleObj)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Element) ToJSON() ([]byte, error) {
	type Alias Element
	return json.Marshal(&struct {
		*Alias
		Points interface{} `json:"points"`
		Style  interface{} `json:"style"`
	}{
		Alias:  (*Alias)(e),
		Points: e.PointsArr,
		Style:  e.StyleObj,
	})
}

func EnsureDataDir() {
	dataDir := "./data"
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		os.MkdirAll(dataDir, 0755)
	}
}
