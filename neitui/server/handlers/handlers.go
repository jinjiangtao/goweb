package handlers

import (
	"database/sql"
	"log"
	"neitui/models"
	"neitui/utils"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err := models.DB.QueryRow(
		"SELECT id, username, password, role, real_name, status FROM users WHERE username = ?",
		req.Username,
	).Scan(&user.ID, &user.Username, &user.Password, &user.Role, &user.RealName, &user.Status)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if user.Status == "disabled" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account disabled"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"role":      user.Role,
			"real_name": user.RealName,
		},
	})
}

func CreateJob(c *gin.Context) {
	var req struct {
		Title       string `json:"title" binding:"required"`
		Requirement string `json:"requirement" binding:"required"`
		SalaryRange string `json:"salary_range"`
		Location    string `json:"location"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetInt("user_id")
	result, err := models.DB.Exec(
		"INSERT INTO jobs (title, requirement, salary_range, location, created_by) VALUES (?, ?, ?, ?, ?)",
		req.Title, req.Requirement, req.SalaryRange, req.Location, userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job"})
		return
	}

	id, _ := result.LastInsertId()
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func GetMyJobs(c *gin.Context) {
	userID := c.GetInt("user_id")
	rows, err := models.DB.Query(
		"SELECT id, title, requirement, salary_range, location, status, created_at FROM jobs WHERE created_by = ? ORDER BY created_at DESC",
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get jobs"})
		return
	}
	defer rows.Close()

	var jobs []models.Job
	for rows.Next() {
		var job models.Job
		rows.Scan(&job.ID, &job.Title, &job.Requirement, &job.SalaryRange, &job.Location, &job.Status, &job.CreatedAt)
		jobs = append(jobs, job)
	}

	c.JSON(http.StatusOK, jobs)
}

func UpdateJobStatus(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetInt("user_id")

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := models.DB.Exec(
		"UPDATE jobs SET status = ? WHERE id = ? AND created_by = ?",
		req.Status, id, userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update job status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully"})
}

func CreateReferral(c *gin.Context) {
	jobID, _ := strconv.Atoi(c.PostForm("job_id"))
	candidateName := c.PostForm("candidate_name")
	candidatePhone := c.PostForm("candidate_phone")

	if jobID == 0 || candidateName == "" || candidatePhone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	file, err := c.FormFile("resume")
	var resumePath string
	if err == nil {
		filename := strconv.FormatInt(time.Now().Unix(), 10) + "_" + file.Filename
		resumePath = "/uploads/resumes/" + filename
		c.SaveUploadedFile(file, "./uploads/resumes/"+filename)
	}

	userID := c.GetInt("user_id")
	result, err := models.DB.Exec(
		"INSERT INTO referrals (job_id, candidate_name, candidate_phone, resume_path, employee_id) VALUES (?, ?, ?, ?, ?)",
		jobID, candidateName, candidatePhone, resumePath, userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create referral"})
		return
	}

	id, _ := result.LastInsertId()
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func GetMyReferrals(c *gin.Context) {
	userID := c.GetInt("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	offset := (page - 1) * pageSize

	query := `
		SELECT r.id, r.job_id, j.title, r.candidate_name, r.candidate_phone, 
		       r.resume_path, r.status, r.hr_remark, r.created_at, r.evaluation_score
		FROM referrals r 
		JOIN jobs j ON r.job_id = j.id 
		WHERE r.employee_id = ?
	`
	args := []interface{}{userID}

	if status != "" {
		query += " AND r.status = ?"
		args = append(args, status)
	}

	query += " ORDER BY r.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	rows, err := models.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get referrals: " + err.Error()})
		return
	}
	defer rows.Close()

	var referrals []models.Referral
	for rows.Next() {
		var r models.Referral
		var evaluationScore sql.NullInt32
		err := rows.Scan(&r.ID, &r.JobID, &r.JobTitle, &r.CandidateName, &r.CandidatePhone,
			&r.ResumePath, &r.Status, &r.HRRemark, &r.CreatedAt, &evaluationScore)
		if err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}
		if evaluationScore.Valid {
			r.EvaluationScore = int(evaluationScore.Int32)
		}
		referrals = append(referrals, r)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing rows: " + err.Error()})
		return
	}

	countQuery := "SELECT COUNT(*) FROM referrals WHERE employee_id = ?"
	countArgs := []interface{}{userID}
	if status != "" {
		countQuery += " AND status = ?"
		countArgs = append(countArgs, status)
	}
	var total int
	models.DB.QueryRow(countQuery, countArgs...).Scan(&total)

	c.JSON(http.StatusOK, gin.H{
		"data":  referrals,
		"total": total,
		"page":  page,
	})
}

func GetMyReferralStats(c *gin.Context) {
	userID := c.GetInt("user_id")

	var total int
	models.DB.QueryRow("SELECT COUNT(*) FROM referrals WHERE employee_id = ?", userID).Scan(&total)

	statuses := []string{"screening", "interviewing", "offer", "hired", "rejected"}
	stats := make(map[string]int)
	for _, s := range statuses {
		var count int
		models.DB.QueryRow("SELECT COUNT(*) FROM referrals WHERE employee_id = ? AND status = ?", userID, s).Scan(&count)
		stats[s] = count
	}

	c.JSON(http.StatusOK, gin.H{
		"total":  total,
		"status": stats,
	})
}

func GetAllReferrals(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	jobID := c.Query("job_id")
	employeeID := c.Query("employee_id")
	search := c.Query("search")

	offset := (page - 1) * pageSize

	query := `
		SELECT r.id, r.job_id, j.title, r.candidate_name, r.candidate_phone, 
		       r.resume_path, r.status, r.employee_id, u.real_name, 
		       r.created_at, r.updated_at, r.evaluation_pros, r.evaluation_cons, r.evaluation_score, r.evaluation_time
		FROM referrals r 
		JOIN jobs j ON r.job_id = j.id 
		JOIN users u ON r.employee_id = u.id 
		WHERE 1=1
	`
	args := []interface{}{}

	if status != "" {
		query += " AND r.status = ?"
		args = append(args, status)
	}
	if jobID != "" {
		query += " AND r.job_id = ?"
		args = append(args, jobID)
	}
	if employeeID != "" {
		query += " AND r.employee_id = ?"
		args = append(args, employeeID)
	}
	if search != "" {
		query += " AND (r.candidate_name LIKE ? OR r.candidate_phone LIKE ?)"
		args = append(args, "%"+search+"%", "%"+search+"%")
	}

	query += " ORDER BY r.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	rows, err := models.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get referrals: " + err.Error()})
		return
	}
	defer rows.Close()

	var referrals []models.Referral
	for rows.Next() {
		var r models.Referral
		var evaluationPros sql.NullString
		var evaluationCons sql.NullString
		var evaluationScore sql.NullInt32
		var evaluationTime sql.NullTime
		err := rows.Scan(&r.ID, &r.JobID, &r.JobTitle, &r.CandidateName, &r.CandidatePhone,
			&r.ResumePath, &r.Status, &r.EmployeeID, &r.EmployeeName, &r.CreatedAt, &r.UpdatedAt,
			&evaluationPros, &evaluationCons, &evaluationScore, &evaluationTime)
		if err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}
		if evaluationPros.Valid {
			r.EvaluationPros = evaluationPros.String
		}
		if evaluationCons.Valid {
			r.EvaluationCons = evaluationCons.String
		}
		if evaluationScore.Valid {
			r.EvaluationScore = int(evaluationScore.Int32)
		}
		if evaluationTime.Valid {
			r.EvaluationTime = &evaluationTime.Time
		}
		referrals = append(referrals, r)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing rows: " + err.Error()})
		return
	}

	countQuery := "SELECT COUNT(*) FROM referrals r WHERE 1=1"
	countArgs := []interface{}{}
	if status != "" {
		countQuery += " AND r.status = ?"
		countArgs = append(countArgs, status)
	}
	if jobID != "" {
		countQuery += " AND r.job_id = ?"
		countArgs = append(countArgs, jobID)
	}
	if employeeID != "" {
		countQuery += " AND r.employee_id = ?"
		countArgs = append(countArgs, employeeID)
	}
	if search != "" {
		countQuery += " AND (r.candidate_name LIKE ? OR r.candidate_phone LIKE ?)"
		countArgs = append(countArgs, "%"+search+"%", "%"+search+"%")
	}
	var total int
	models.DB.QueryRow(countQuery, countArgs...).Scan(&total)

	c.JSON(http.StatusOK, gin.H{
		"data":  referrals,
		"total": total,
		"page":  page,
	})
}

func UpdateReferralStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status       string `json:"status" binding:"required"`
		RejectReason string `json:"reject_reason"`
		HRRemark     string `json:"hr_remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := models.DB.Exec(
		"UPDATE referrals SET status = ?, reject_reason = ?, hr_remark = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		req.Status, req.RejectReason, req.HRRemark, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update referral"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully"})
}

func GetStats(c *gin.Context) {
	var totalReferrals int
	models.DB.QueryRow("SELECT COUNT(*) FROM referrals").Scan(&totalReferrals)

	var hired int
	models.DB.QueryRow("SELECT COUNT(*) FROM referrals WHERE status = 'hired'").Scan(&hired)

	var interviewed int
	models.DB.QueryRow("SELECT COUNT(*) FROM referrals WHERE status IN ('interviewing', 'offer', 'hired')").Scan(&interviewed)
	var totalScreened int
	models.DB.QueryRow("SELECT COUNT(*) FROM referrals WHERE status != 'screening'").Scan(&totalScreened)
	interviewRate := 0.0
	if totalScreened > 0 {
		interviewRate = float64(interviewed) / float64(totalScreened) * 100
	}

	var averageScore float64
	var evaluatedCount int
	models.DB.QueryRow("SELECT AVG(evaluation_score), COUNT(*) FROM referrals WHERE evaluation_score IS NOT NULL").Scan(&averageScore, &evaluatedCount)

	rows, err := models.DB.Query(`
		SELECT u.real_name, COUNT(*) as count 
		FROM referrals r 
		JOIN users u ON r.employee_id = u.id 
		WHERE r.status = 'hired' 
		GROUP BY u.id, u.real_name 
		ORDER BY count DESC
	`)
	type TopEmployee struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}
	var topEmployees []TopEmployee
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var t TopEmployee
			rows.Scan(&t.Name, &t.Count)
			topEmployees = append(topEmployees, t)
		}
		if err = rows.Err(); err != nil {
			log.Printf("Top employees rows error: %v", err)
		}
	} else {
		log.Printf("Top employees query error: %v", err)
	}

	rows, err = models.DB.Query(`
		SELECT u.real_name, AVG(r.evaluation_score) as avg_score 
		FROM referrals r 
		JOIN users u ON r.employee_id = u.id 
		WHERE r.evaluation_score IS NOT NULL 
		GROUP BY u.id, u.real_name 
		ORDER BY avg_score DESC
	`)
	type EmployeeScore struct {
		Name     string  `json:"name"`
		AvgScore float64 `json:"avg_score"`
	}
	var employeeScoreRanking []EmployeeScore
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var e EmployeeScore
			rows.Scan(&e.Name, &e.AvgScore)
			employeeScoreRanking = append(employeeScoreRanking, e)
		}
		if err = rows.Err(); err != nil {
			log.Printf("Employee score ranking rows error: %v", err)
		}
	} else {
		log.Printf("Employee score ranking query error: %v", err)
	}

	thirtyDaysAgo := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	rows, err = models.DB.Query(`
		SELECT DATE(created_at) as date, COUNT(*) as count 
		FROM referrals 
		WHERE created_at >= ? 
		GROUP BY DATE(created_at) 
		ORDER BY date
	`, thirtyDaysAgo)
	type DailyTrend struct {
		Date  string `json:"date"`
		Count int    `json:"count"`
	}
	var trends []DailyTrend
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var t DailyTrend
			rows.Scan(&t.Date, &t.Count)
			trends = append(trends, t)
		}
		if err = rows.Err(); err != nil {
			log.Printf("Daily trend rows error: %v", err)
		}
	} else {
		log.Printf("Daily trend query error: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"total_referrals":        totalReferrals,
		"hired":                  hired,
		"interview_rate":         interviewRate,
		"average_score":          averageScore,
		"top_employees":          topEmployees,
		"employee_score_ranking": employeeScoreRanking,
		"thirty_days_trend":      trends,
	})
}

func GetAllJobs(c *gin.Context) {
	createdBy := c.Query("created_by")

	query := `
		SELECT j.id, j.title, j.requirement, j.salary_range, j.location, j.status, j.created_by, u.real_name, j.created_at 
		FROM jobs j 
		JOIN users u ON j.created_by = u.id 
		WHERE 1=1
	`
	args := []interface{}{}
	if createdBy != "" {
		query += " AND j.created_by = ?"
		args = append(args, createdBy)
	}
	query += " ORDER BY j.created_at DESC"

	rows, err := models.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get jobs"})
		return
	}
	defer rows.Close()

	type JobWithCreator struct {
		models.Job
		CreatedByName string `json:"created_by_name"`
	}
	var jobs []JobWithCreator
	for rows.Next() {
		var j JobWithCreator
		rows.Scan(&j.ID, &j.Title, &j.Requirement, &j.SalaryRange, &j.Location, &j.Status, &j.CreatedBy, &j.CreatedByName, &j.CreatedAt)
		jobs = append(jobs, j)
	}

	c.JSON(http.StatusOK, jobs)
}

func GetUsers(c *gin.Context) {
	rows, err := models.DB.Query("SELECT id, username, role, real_name, status, created_at FROM users ORDER BY created_at DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		rows.Scan(&u.ID, &u.Username, &u.Role, &u.RealName, &u.Status, &u.CreatedAt)
		users = append(users, u)
	}

	c.JSON(http.StatusOK, users)
}

func CreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		RealName string `json:"real_name" binding:"required"`
		Role     string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	result, err := models.DB.Exec(
		"INSERT INTO users (username, password, role, real_name, status) VALUES (?, ?, ?, ?, ?)",
		req.Username, string(hashedPassword), req.Role, req.RealName, "enabled",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	id, _ := result.LastInsertId()
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func UpdateUserStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := models.DB.Exec("UPDATE users SET status = ? WHERE id = ?", req.Status, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully"})
}

func ResetPassword(c *gin.Context) {
	id := c.Param("id")
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	_, err := models.DB.Exec("UPDATE users SET password = ? WHERE id = ?", string(hashedPassword), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset to 123456"})
}

func ExportReferrals(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Export feature coming soon"})
}

func GetPublicJobs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := c.Query("search")
	salaryRange := c.Query("salary_range")
	location := c.Query("location")

	offset := (page - 1) * pageSize

	query := `
		SELECT id, title, requirement, salary_range, location, status, created_at, favorite_count 
		FROM jobs 
		WHERE status = 'active'
	`
	args := []interface{}{}

	if search != "" {
		query += " AND title LIKE ?"
		args = append(args, "%"+search+"%")
	}
	if salaryRange != "" {
		query += " AND salary_range = ?"
		args = append(args, salaryRange)
	}
	if location != "" {
		query += " AND location = ?"
		args = append(args, location)
	}

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	rows, err := models.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get jobs"})
		return
	}
	defer rows.Close()

	var jobs []models.Job
	for rows.Next() {
		var job models.Job
		rows.Scan(&job.ID, &job.Title, &job.Requirement, &job.SalaryRange, &job.Location, &job.Status, &job.CreatedAt, &job.FavoriteCount)
		jobs = append(jobs, job)
	}

	countQuery := "SELECT COUNT(*) FROM jobs WHERE status = 'active'"
	countArgs := []interface{}{}
	if search != "" {
		countQuery += " AND title LIKE ?"
		countArgs = append(countArgs, "%"+search+"%")
	}
	if salaryRange != "" {
		countQuery += " AND salary_range = ?"
		countArgs = append(countArgs, salaryRange)
	}
	if location != "" {
		countQuery += " AND location = ?"
		countArgs = append(countArgs, location)
	}
	var total int
	models.DB.QueryRow(countQuery, countArgs...).Scan(&total)

	c.JSON(http.StatusOK, gin.H{
		"data":  jobs,
		"total": total,
		"page":  page,
	})
}

func GetPublicJobDetail(c *gin.Context) {
	id := c.Param("id")

	var job models.Job
	err := models.DB.QueryRow(
		"SELECT id, title, requirement, salary_range, location, status, created_at, favorite_count FROM jobs WHERE id = ?",
		id,
	).Scan(&job.ID, &job.Title, &job.Requirement, &job.SalaryRange, &job.Location, &job.Status, &job.CreatedAt, &job.FavoriteCount)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	c.JSON(http.StatusOK, job)
}

func SubmitApplication(c *gin.Context) {
	jobID, _ := strconv.Atoi(c.PostForm("job_id"))
	name := c.PostForm("name")
	phone := c.PostForm("phone")

	if jobID == 0 || name == "" || phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	var count int
	models.DB.QueryRow("SELECT COUNT(*) FROM submissions WHERE job_id = ? AND phone = ?", jobID, phone).Scan(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You have already applied for this job"})
		return
	}

	file, err := c.FormFile("resume_file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Resume file is required"})
		return
	}

	uploadsPath := "./uploads/resumes/h5"
	os.MkdirAll(uploadsPath, 0755)

	filename := strconv.FormatInt(time.Now().Unix(), 10) + "_" + file.Filename
	resumePath := "/uploads/resumes/h5/" + filename
	err = c.SaveUploadedFile(file, uploadsPath+"/"+filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save resume"})
		return
	}

	result, err := models.DB.Exec(
		"INSERT INTO submissions (job_id, name, phone, resume_path, status) VALUES (?, ?, ?, ?, 'screening')",
		jobID, name, phone, resumePath,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit application"})
		return
	}

	id, _ := result.LastInsertId()
	c.JSON(http.StatusOK, gin.H{"id": id, "message": "Application submitted successfully"})
}

func GetMySubmissions(c *gin.Context) {
	phone := c.Query("phone")
	if phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number is required"})
		return
	}

	query := `
		SELECT s.id, s.job_id, j.title, s.name, s.phone, s.resume_path, s.status, s.created_at 
		FROM submissions s 
		JOIN jobs j ON s.job_id = j.id 
		WHERE s.phone = ?
		ORDER BY s.created_at DESC
	`

	rows, err := models.DB.Query(query, phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get submissions"})
		return
	}
	defer rows.Close()

	type SubmissionWithJobTitle struct {
		models.Submission
	}
	var submissions []SubmissionWithJobTitle
	for rows.Next() {
		var s SubmissionWithJobTitle
		rows.Scan(&s.ID, &s.JobID, &s.JobTitle, &s.Name, &s.Phone, &s.ResumePath, &s.Status, &s.CreatedAt)
		submissions = append(submissions, s)
	}

	c.JSON(http.StatusOK, submissions)
}

func GetAllSubmissions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	jobID := c.Query("job_id")
	search := c.Query("search")

	offset := (page - 1) * pageSize

	query := `
		SELECT s.id, s.job_id, j.title, s.name, s.phone, s.resume_path, s.status, 
		       s.converted, s.converted_at, s.created_at 
		FROM submissions s 
		JOIN jobs j ON s.job_id = j.id 
		WHERE 1=1
	`
	args := []interface{}{}

	if status != "" {
		query += " AND s.status = ?"
		args = append(args, status)
	}
	if jobID != "" {
		query += " AND s.job_id = ?"
		args = append(args, jobID)
	}
	if search != "" {
		query += " AND (s.name LIKE ? OR s.phone LIKE ?)"
		args = append(args, "%"+search+"%", "%"+search+"%")
	}

	query += " ORDER BY s.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	rows, err := models.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get submissions"})
		return
	}
	defer rows.Close()

	type SubmissionWithJobTitle struct {
		models.Submission
	}
	var submissions []SubmissionWithJobTitle
	for rows.Next() {
		var s SubmissionWithJobTitle
		rows.Scan(&s.ID, &s.JobID, &s.JobTitle, &s.Name, &s.Phone, &s.ResumePath, &s.Status,
			&s.Converted, &s.ConvertedAt, &s.CreatedAt)
		submissions = append(submissions, s)
	}

	countQuery := "SELECT COUNT(*) FROM submissions s WHERE 1=1"
	countArgs := []interface{}{}
	if status != "" {
		countQuery += " AND s.status = ?"
		countArgs = append(countArgs, status)
	}
	if jobID != "" {
		countQuery += " AND s.job_id = ?"
		countArgs = append(countArgs, jobID)
	}
	if search != "" {
		countQuery += " AND (s.name LIKE ? OR s.phone LIKE ?)"
		countArgs = append(countArgs, "%"+search+"%", "%"+search+"%")
	}
	var total int
	models.DB.QueryRow(countQuery, countArgs...).Scan(&total)

	c.JSON(http.StatusOK, gin.H{
		"data":  submissions,
		"total": total,
		"page":  page,
	})
}

func UpdateSubmissionStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := models.DB.Exec(
		"UPDATE submissions SET status = ? WHERE id = ?",
		req.Status, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update submission"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully"})
}

func ConvertSubmissionToReferral(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		EmployeeID int `json:"employee_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var submission struct {
		JobID      int
		Name       string
		Phone      string
		ResumePath string
	}
	err := models.DB.QueryRow(
		"SELECT job_id, name, phone, resume_path FROM submissions WHERE id = ?",
		id,
	).Scan(&submission.JobID, &submission.Name, &submission.Phone, &submission.ResumePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
		return
	}

	result, err := models.DB.Exec(
		`INSERT INTO referrals (job_id, candidate_name, candidate_phone, resume_path, employee_id, status) 
		VALUES (?, ?, ?, ?, ?, 'screening')`,
		submission.JobID, submission.Name, submission.Phone, submission.ResumePath, req.EmployeeID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create referral"})
		return
	}

	_, err = models.DB.Exec(
		"UPDATE submissions SET converted = 1, converted_at = CURRENT_TIMESTAMP WHERE id = ?",
		id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update submission"})
		return
	}

	referralID, _ := result.LastInsertId()
	c.JSON(http.StatusOK, gin.H{"referral_id": referralID, "message": "Converted successfully"})
}

func UpdateFavorite(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Action string `json:"action" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var query string
	if req.Action == "add" {
		query = "UPDATE jobs SET favorite_count = favorite_count + 1 WHERE id = ?"
	} else if req.Action == "remove" {
		query = "UPDATE jobs SET favorite_count = favorite_count - 1 WHERE id = ? AND favorite_count > 0"
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action"})
		return
	}

	_, err := models.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update favorite count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully"})
}

func GetFavoritesByIDs(c *gin.Context) {
	var req struct {
		IDs []int `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.IDs) == 0 {
		c.JSON(http.StatusOK, []models.Job{})
		return
	}

	placeholders := make([]string, len(req.IDs))
	args := make([]interface{}, len(req.IDs))
	for i, id := range req.IDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := `
		SELECT id, title, requirement, salary_range, location, status, created_at, favorite_count 
		FROM jobs 
		WHERE id IN (` + strings.Join(placeholders, ",") + `)
		ORDER BY created_at DESC
	`

	rows, err := models.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get favorite jobs: " + err.Error()})
		return
	}
	defer rows.Close()

	var jobs []models.Job
	for rows.Next() {
		var job models.Job
		err := rows.Scan(&job.ID, &job.Title, &job.Requirement, &job.SalaryRange, &job.Location, &job.Status, &job.CreatedAt, &job.FavoriteCount)
		if err != nil {
			log.Printf("Scan error in favorites: %v", err)
			continue
		}
		jobs = append(jobs, job)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Rows error in favorites: %v", err)
	}

	c.JSON(http.StatusOK, jobs)
}

func UpdateEvaluation(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		EvaluationPros  string `json:"evaluation_pros"`
		EvaluationCons  string `json:"evaluation_cons"`
		EvaluationScore int    `json:"evaluation_score" binding:"min=1,max=5"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := models.DB.Exec(
		"UPDATE referrals SET evaluation_pros = ?, evaluation_cons = ?, evaluation_score = ?, evaluation_time = CURRENT_TIMESTAMP WHERE id = ?",
		req.EvaluationPros, req.EvaluationCons, req.EvaluationScore, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update evaluation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully"})
}

func GetEvaluation(c *gin.Context) {
	id := c.Param("id")

	var r models.Referral
	var evaluationPros sql.NullString
	var evaluationCons sql.NullString
	var evaluationScore sql.NullInt32
	var evaluationTime sql.NullTime
	
	err := models.DB.QueryRow(
		"SELECT evaluation_pros, evaluation_cons, evaluation_score, evaluation_time FROM referrals WHERE id = ?",
		id,
	).Scan(&evaluationPros, &evaluationCons, &evaluationScore, &evaluationTime)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Referral not found: " + err.Error()})
		return
	}

	if evaluationPros.Valid {
		r.EvaluationPros = evaluationPros.String
	}
	if evaluationCons.Valid {
		r.EvaluationCons = evaluationCons.String
	}
	if evaluationScore.Valid {
		r.EvaluationScore = int(evaluationScore.Int32)
	}
	if evaluationTime.Valid {
		r.EvaluationTime = &evaluationTime.Time
	}

	c.JSON(http.StatusOK, r)
}
