package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Admin struct {
	User_ID  uuid.UUID `json:"user_id"`
	Fullname string    `json:"fullname"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Role     string    `json:"role"`
	Phone    string    `json:"phone"`
	Auditable
	Verification      bool      `json:"verification"`
	JwtToken          string    `json:"jwt_token,omitempty"`
	JwtTokenExpiresAt time.Time `json:"jwt_token_expires_at,omitempty"`
}

func (u *Admin) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Role == "" {
		u.Role = "admin"
	}
	if u.Fullname == "" {
		u.Fullname = "admin"
	}
	return
}

func NewAdmin(fullname, email, password, role, phone string, verification bool) *Admin {
	return &Admin{
		User_ID:      uuid.New(),
		Fullname:     fullname,
		Email:        email,
		Password:     password,
		Role:         role,
		Phone:        phone,
		Auditable:    NewAuditable(),
		Verification: false,
	}
}

func UpdateAdmin(admin_id uuid.UUID, fullname, email, password, role, phone string, verification bool) *Admin {
	return &Admin{
		User_ID:      admin_id,
		Fullname:     fullname,
		Email:        email,
		Password:     password,
		Role:         role,
		Phone:        phone,
		Auditable:    UpdateAuditable(),
		Verification: false,
	}
}

func (Admin) TableName() string {
	return "users"
}
