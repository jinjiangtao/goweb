package models

import (
	"database/sql"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	RealName  string    `json:"real_name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type Job struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Requirement  string    `json:"requirement"`
	SalaryRange  string    `json:"salary_range"`
	Location     string    `json:"location"`
	Status       string    `json:"status"`
	CreatedBy    int       `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
}

type Referral struct {
	ID             int       `json:"id"`
	JobID          int       `json:"job_id"`
	JobTitle       string    `json:"job_title,omitempty"`
	CandidateName  string    `json:"candidate_name"`
	CandidatePhone string    `json:"candidate_phone"`
	ResumePath     string    `json:"resume_path"`
	Status         string    `json:"status"`
	RejectReason   string    `json:"reject_reason,omitempty"`
	EmployeeID     int       `json:"employee_id"`
	EmployeeName   string    `json:"employee_name,omitempty"`
	HRRemark       string    `json:"hr_remark,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite", "./neitui.db")
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	if err = createTables(); err != nil {
		return err
	}

	if err = initDefaultUsers(); err != nil {
		return err
	}

	return nil
}

func createTables() error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			role TEXT NOT NULL,
			real_name TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'enabled',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS jobs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			requirement TEXT NOT NULL,
			salary_range TEXT,
			location TEXT,
			status TEXT NOT NULL DEFAULT 'active',
			created_by INTEGER NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS referrals (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			job_id INTEGER NOT NULL,
			candidate_name TEXT NOT NULL,
			candidate_phone TEXT NOT NULL,
			resume_path TEXT,
			status TEXT NOT NULL DEFAULT 'screening',
			reject_reason TEXT,
			employee_id INTEGER NOT NULL,
			hr_remark TEXT,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, table := range tables {
		_, err := DB.Exec(table)
		if err != nil {
			return err
		}
	}

	return nil
}

func initDefaultUsers() error {
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", "admin").Scan(&count)
	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		_, err := DB.Exec(
			"INSERT INTO users (username, password, role, real_name, status) VALUES (?, ?, ?, ?, ?)",
			"admin", string(hashedPassword), "admin", "管理员", "enabled",
		)
		if err != nil {
			log.Println("Failed to create admin user:", err)
		}
	}

	DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", "hr").Scan(&count)
	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		_, err := DB.Exec(
			"INSERT INTO users (username, password, role, real_name, status) VALUES (?, ?, ?, ?, ?)",
			"hr", string(hashedPassword), "hr", "HR", "enabled",
		)
		if err != nil {
			log.Println("Failed to create hr user:", err)
		}
	}

	return nil
}
