package repository

import (
	"encoding/json"
	"errors"
	"time"

	"Kevinmajesta/OrderManagementAPI/internal/entity"
	"Kevinmajesta/OrderManagementAPI/pkg/cache"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdminRepository interface {
	FindAdminByEmail(email string) (*entity.Admin, error)
	FindAdminByID(user_id uuid.UUID) (*entity.Admin, error)
	FindByRole(role string, users *[]entity.User) error
	FindAllUser() ([]entity.User, error)
	CreateAdmin(admin *entity.Admin) (*entity.Admin, error)
	UpdateAdmin(admin *entity.Admin) (*entity.Admin, error)
	DeleteAdmin(admin *entity.Admin) (bool, error)
	SaveVerifCode(userID uuid.UUID, resetCode string) error
	UpdateAdminJwtToken(userID uuid.UUID, token string, expiresAt time.Time) error
	CheckUserExists(id uuid.UUID) (bool, error)
}

type adminRepository struct {
	db        *gorm.DB
	cacheable cache.Cacheable
}

func NewAdminRepository(db *gorm.DB, cacheable cache.Cacheable) *adminRepository {
	return &adminRepository{db: db, cacheable: cacheable}
}

func (r *adminRepository) FindAdminByID(user_id uuid.UUID) (*entity.Admin, error) {
	admin := new(entity.Admin)
	if err := r.db.Where("user_id = ?", user_id).Take(admin).Error; err != nil {
		return admin, err
	}
	return admin, nil
}

func (r *adminRepository) FindAdminByEmail(email string) (*entity.Admin, error) {
	admin := new(entity.Admin)
	if err := r.db.Where("email = ?", email).Take(admin).Error; err != nil {
		return admin, err
	}
	return admin, nil
}

func (r *adminRepository) FindByRole(role string, users *[]entity.User) error {
	return r.db.Where("role = ?", role).Find(users).Error
}

func (r *adminRepository) FindAllUser() ([]entity.User, error) {
	users := make([]entity.User, 0)

	key := "FindAllUsers"

	data, _ := r.cacheable.Get(key)
	if data == "" {
		if err := r.db.Find(&users).Error; err != nil {
			return users, err
		}
		marshalledUsers, _ := json.Marshal(users)
		err := r.cacheable.Set(key, marshalledUsers, 5*time.Minute)
		if err != nil {
			return users, err
		}
	} else {
		// Data ditemukan di Redis, unmarshal data ke users
		err := json.Unmarshal([]byte(data), &users)
		if err != nil {
			return users, err
		}
	}
	return users, nil
}

func (r *adminRepository) CreateAdmin(admin *entity.Admin) (*entity.Admin, error) {
	if err := r.db.Create(&admin).Error; err != nil {
		return admin, err
	}
	return admin, nil
}

func (r *adminRepository) CheckUserExists(id uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.Model(&entity.Admin{}).Where("user_id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *adminRepository) UpdateAdmin(admin *entity.Admin) (*entity.Admin, error) {
	// Use map to store fields to be updated.
	fields := make(map[string]interface{})

	// Update fields only if they are not empty.
	if admin.Email != "" {
		fields["email"] = admin.Email
	}
	if admin.Password != "" {
		fields["password"] = admin.Password
	}
	if admin.Role != "" {
		fields["role"] = admin.Role
	}

	// Update the database in one query.
	if err := r.db.Model(admin).Where("user_id = ?", admin.User_ID).Updates(fields).Error; err != nil {
		return admin, err
	}

	return admin, nil
}

func (r *adminRepository) DeleteAdmin(admin *entity.Admin) (bool, error) {
	if err := r.db.Delete(&entity.Admin{}, admin.User_ID).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (r *adminRepository) SaveVerifCode(user_ID uuid.UUID, resetCode string) error {
	return r.db.Model(&entity.User{}).Where("user_id = ?", user_ID).Updates(map[string]interface{}{
		"verification_code": resetCode,
	}).Error
}

func (u *adminRepository) UpdateAdminJwtToken(userID uuid.UUID, token string, expiresAt time.Time) error {
	result := u.db.Model(&entity.User{}).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Updates(map[string]interface{}{
			"jwt_token":            token,
			"jwt_token_expires_at": expiresAt,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found or already deleted")
	}

	return nil
}
