package handlers

import (
	"neitui/models"
	"neitui/utils"
	"net/http"
	"strconv"
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
		       r.resume_path, r.status, r.hr_remark, r.created_at 
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get referrals"})
		return
	}
	defer rows.Close()

	var referrals []models.Referral
	for rows.Next() {
		var r models.Referral
		rows.Scan(&r.ID, &r.JobID, &r.JobTitle, &r.CandidateName, &r.CandidatePhone,
			&r.ResumePath, &r.Status, &r.HRRemark, &r.CreatedAt)
		referrals = append(referrals, r)
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
		       r.created_at, r.updated_at 
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get referrals"})
		return
	}
	defer rows.Close()

	var referrals []models.Referral
	for rows.Next() {
		var r models.Referral
		rows.Scan(&r.ID, &r.JobID, &r.JobTitle, &r.CandidateName, &r.CandidatePhone,
			&r.ResumePath, &r.Status, &r.EmployeeID, &r.EmployeeName, &r.CreatedAt, &r.UpdatedAt)
		referrals = append(referrals, r)
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

	rows, _ := models.DB.Query(`
		SELECT u.real_name, COUNT(*) as count 
		FROM referrals r 
		JOIN users u ON r.employee_id = u.id 
		WHERE r.status = 'hired' 
		GROUP BY u.id, u.real_name 
		ORDER BY count DESC
	`)
	defer rows.Close()
	type TopEmployee struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}
	var topEmployees []TopEmployee
	for rows.Next() {
		var t TopEmployee
		rows.Scan(&t.Name, &t.Count)
		topEmployees = append(topEmployees, t)
	}

	thirtyDaysAgo := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	rows, _ = models.DB.Query(`
		SELECT DATE(created_at) as date, COUNT(*) as count 
		FROM referrals 
		WHERE created_at >= ? 
		GROUP BY DATE(created_at) 
		ORDER BY date
	`, thirtyDaysAgo)
	defer rows.Close()
	type DailyTrend struct {
		Date  string `json:"date"`
		Count int    `json:"count"`
	}
	var trends []DailyTrend
	for rows.Next() {
		var t DailyTrend
		rows.Scan(&t.Date, &t.Count)
		trends = append(trends, t)
	}

	c.JSON(http.StatusOK, gin.H{
		"total_referrals":   totalReferrals,
		"hired":             hired,
		"interview_rate":    interviewRate,
		"top_employees":     topEmployees,
		"thirty_days_trend": trends,
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
