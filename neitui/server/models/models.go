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
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Requirement   string    `json:"requirement"`
	SalaryRange   string    `json:"salary_range"`
	Location      string    `json:"location"`
	Status        string    `json:"status"`
	CreatedBy     int       `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
	FavoriteCount int       `json:"favorite_count"`
}

type Referral struct {
	ID              int        `json:"id"`
	JobID           int        `json:"job_id"`
	JobTitle        string     `json:"job_title,omitempty"`
	CandidateName   string     `json:"candidate_name"`
	CandidatePhone  string     `json:"candidate_phone"`
	ResumePath      string     `json:"resume_path"`
	Status          string     `json:"status"`
	RejectReason    string     `json:"reject_reason,omitempty"`
	EmployeeID      int        `json:"employee_id"`
	EmployeeName    string     `json:"employee_name,omitempty"`
	HRRemark        string     `json:"hr_remark,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	EvaluationPros  string     `json:"evaluation_pros,omitempty"`
	EvaluationCons  string     `json:"evaluation_cons,omitempty"`
	EvaluationScore int        `json:"evaluation_score,omitempty"`
	EvaluationTime  *time.Time `json:"evaluation_time,omitempty"`
}

type Submission struct {
	ID          int        `json:"id"`
	JobID       int        `json:"job_id"`
	JobTitle    string     `json:"job_title,omitempty"`
	Name        string     `json:"name"`
	Phone       string     `json:"phone"`
	ResumePath  string     `json:"resume_path"`
	Status      string     `json:"status"`
	Converted   bool       `json:"converted"`
	ConvertedAt *time.Time `json:"converted_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
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

func columnExists(tableName, columnName string) bool {
	rows, err := DB.Query("PRAGMA table_info(" + tableName + ")")
	if err != nil {
		return false
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name, typeCol string
		var notNull, dfltValue, pk sql.NullString
		if err := rows.Scan(&cid, &name, &typeCol, &notNull, &dfltValue, &pk); err != nil {
			continue
		}
		if name == columnName {
			return true
		}
	}
	return false
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
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			favorite_count INTEGER NOT NULL DEFAULT 0
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
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			evaluation_pros TEXT,
			evaluation_cons TEXT,
			evaluation_score INTEGER,
			evaluation_time DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS submissions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			job_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			phone TEXT NOT NULL,
			resume_path TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'screening',
			converted BOOLEAN NOT NULL DEFAULT 0,
			converted_at DATETIME,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, table := range tables {
		_, err := DB.Exec(table)
		if err != nil {
			return err
		}
	}

	// 迁移现有表，添加新字段（SQLite 不支持 IF NOT EXISTS，所以我们手动检查）
	// 检查并添加 jobs 表的 favorite_count 列
	if !columnExists("jobs", "favorite_count") {
		_, _ = DB.Exec(`ALTER TABLE jobs ADD COLUMN favorite_count INTEGER NOT NULL DEFAULT 0`)
	}

	// 检查并添加 referrals 表的评价相关列
	if !columnExists("referrals", "evaluation_pros") {
		_, _ = DB.Exec(`ALTER TABLE referrals ADD COLUMN evaluation_pros TEXT`)
	}
	if !columnExists("referrals", "evaluation_cons") {
		_, _ = DB.Exec(`ALTER TABLE referrals ADD COLUMN evaluation_cons TEXT`)
	}
	if !columnExists("referrals", "evaluation_score") {
		_, _ = DB.Exec(`ALTER TABLE referrals ADD COLUMN evaluation_score INTEGER`)
	}
	if !columnExists("referrals", "evaluation_time") {
		_, _ = DB.Exec(`ALTER TABLE referrals ADD COLUMN evaluation_time DATETIME`)
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

	// 确保 yg 用户存在并且密码正确，不管之前是否已经存在
	DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", "yg").Scan(&count)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if count == 0 {
		_, err := DB.Exec(
			"INSERT INTO users (username, password, role, real_name, status) VALUES (?, ?, ?, ?, ?)",
			"yg", string(hashedPassword), "employee", "员工", "enabled",
		)
		if err != nil {
			log.Println("Failed to create yg user:", err)
		}
	} else {
		// 如果用户已存在，重置密码以确保能登录
		_, err := DB.Exec(
			"UPDATE users SET password = ? WHERE username = ?",
			string(hashedPassword), "yg",
		)
		if err != nil {
			log.Println("Failed to reset yg user password:", err)
		}
	}

	// 添加默认测试职位
	DB.QueryRow("SELECT COUNT(*) FROM jobs").Scan(&count)
	if count == 0 {
		_, err := DB.Exec(
			"INSERT INTO jobs (title, requirement, salary_range, location, status, created_by) VALUES (?, ?, ?, ?, ?, ?)",
			"前端开发工程师", "1. 熟悉Vue3、React等前端框架\n2. 熟悉HTML、CSS、JavaScript\n3. 有良好的编码习惯", "15k-25k", "北京", "active", 1,
		)
		if err != nil {
			log.Println("Failed to create test job 1:", err)
		}

		_, err = DB.Exec(
			"INSERT INTO jobs (title, requirement, salary_range, location, status, created_by) VALUES (?, ?, ?, ?, ?, ?)",
			"后端开发工程师", "1. 熟悉Go或Java语言\n2. 熟悉数据库操作\n3. 有后端开发经验", "18k-30k", "上海", "active", 1,
		)
		if err != nil {
			log.Println("Failed to create test job 2:", err)
		}

		_, err = DB.Exec(
			"INSERT INTO jobs (title, requirement, salary_range, location, status, created_by) VALUES (?, ?, ?, ?, ?, ?)",
			"产品经理", "1. 熟悉产品设计流程\n2. 有互联网产品经验\n3. 良好的沟通能力", "20k-35k", "深圳", "active", 1,
		)
		if err != nil {
			log.Println("Failed to create test job 3:", err)
		}
	}

	return nil
}
