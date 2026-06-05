package models

import "time"

type Signup struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Name      string    `json:"name" gorm:"not null"`
	Phone     string    `json:"phone" gorm:"not null"`
	Age       int       `json:"age" gorm:"not null"`
	Hukou     string    `json:"hukou" gorm:"not null"`
	School    string    `json:"school" gorm:"not null"`
	Status    string    `json:"status" gorm:"default:'pending'"`
	CreatedAt time.Time `json:"created_at"`
}

type SignupListResponse struct {
	Total int      `json:"total"`
	List  []Signup `json:"list"`
}

func CreateSignup(signup *Signup) error {
	return DB.Create(signup).Error
}

func GetSignups(page, pageSize int, keyword, status string) (*SignupListResponse, error) {
	var signups []Signup
	var total int64

	query := DB.Model(&Signup{})
	if keyword != "" {
		query = query.Where("name LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, err
	}

	err = query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&signups).Error
	if err != nil {
		return nil, err
	}

	return &SignupListResponse{
		Total: int(total),
		List:  signups,
	}, nil
}

func GetAllSignupsForExport(keyword, status string) ([]Signup, error) {
	var signups []Signup

	query := DB.Model(&Signup{})
	if keyword != "" {
		query = query.Where("name LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	err := query.Order("created_at DESC").Find(&signups).Error
	if err != nil {
		return nil, err
	}

	return signups, nil
}

func UpdateSignup(signup *Signup) error {
	return DB.Save(signup).Error
}

func GetSignupByID(id uint) (*Signup, error) {
	var signup Signup
	err := DB.First(&signup, id).Error
	return &signup, err
}

func UpdateSignupStatus(id uint, status string) error {
	return DB.Model(&Signup{}).Where("id = ?", id).Update("status", status).Error
}

type StatsResponse struct {
	Pending  int `json:"pending"`
	Approved int `json:"approved"`
	Rejected int `json:"rejected"`
}

func GetStats() (*StatsResponse, error) {
	var stats StatsResponse
	var pending, approved, rejected int64
	err := DB.Model(&Signup{}).Where("status = ?", "pending").Count(&pending).Error
	if err != nil {
		return nil, err
	}
	err = DB.Model(&Signup{}).Where("status = ?", "approved").Count(&approved).Error
	if err != nil {
		return nil, err
	}
	err = DB.Model(&Signup{}).Where("status = ?", "rejected").Count(&rejected).Error
	if err != nil {
		return nil, err
	}
	stats.Pending = int(pending)
	stats.Approved = int(approved)
	stats.Rejected = int(rejected)
	return &stats, nil
}
