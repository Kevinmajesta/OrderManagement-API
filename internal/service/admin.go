package service

import (
	"errors"
	"time"

	"Kevinmajesta/OrderManagementAPI/internal/entity"
	"Kevinmajesta/OrderManagementAPI/internal/repository"
	"Kevinmajesta/OrderManagementAPI/pkg/email"
	"Kevinmajesta/OrderManagementAPI/pkg/encrypt"
	"Kevinmajesta/OrderManagementAPI/pkg/token"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AdminService interface {
	LoginAdmin(email string, password string) (string, error)
	FindAllUser() ([]entity.User, error)
	CreateAdmin(admin *entity.Admin) (*entity.Admin, error)
	UpdateAdmin(admin *entity.Admin) (*entity.Admin, error)
	DeleteAdmin(admin uuid.UUID) (bool, error)
	EmailExists(email string) bool
	CheckUserExists(id uuid.UUID) (bool, error)
}

type adminService struct {
	adminRepository repository.AdminRepository
	tokenUseCase    token.TokenUseCase
	encryptTool     encrypt.EncryptTool
	emailSender     email.EmailSenderService 
}

func NewAdminService(adminRepository repository.AdminRepository, tokenUseCase token.TokenUseCase,
	encryptTool encrypt.EncryptTool, emailSender email.EmailSenderService ) *adminService {

	return &adminService{
		adminRepository: adminRepository,
		tokenUseCase:    tokenUseCase,
		encryptTool:     encryptTool,
		emailSender:     emailSender,
	}
}

func (s *adminService) LoginAdmin(email string, password string) (string, error) {
	admin, err := s.adminRepository.FindAdminByEmail(email)
	if err != nil {
		return "", errors.New("wrong input email/password")
	}
	if admin.Role != "admin" {
		return "", errors.New("you dont have access")
	}
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
	if err != nil {
		return "", errors.New("wrong input email/password")
	}

	expiredTime := time.Now().Local().Add(24 * time.Hour)

	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic(err)
	}

	admin.Phone, _ = s.encryptTool.Decrypt(admin.Phone)
	expiredTimeInJakarta := expiredTime.In(location)

	claims := token.JwtCustomClaims{
		ID:    admin.User_ID.String(),
		Email: admin.Email,
		Role:  "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Depublic",
			ExpiresAt: jwt.NewNumericDate(expiredTimeInJakarta),
		},
	}

	jwtToken, err := s.tokenUseCase.GenerateAccessToken(claims)
	if err != nil {
		return "", errors.New("there is an error in the system")
	}

	admin.JwtToken = jwtToken
	admin.JwtTokenExpiresAt = expiredTime

	if err := s.adminRepository.UpdateAdminJwtToken(admin.User_ID, jwtToken, expiredTime); err != nil {
		return "", errors.New("failed to update user token info")
	}

	if admin.JwtToken != jwtToken {
		return "", errors.New("JWT token mismatch")
	}
	return jwtToken, nil
}

func (s *adminService) FindAllUser() ([]entity.User, error) {
	admin, err := s.adminRepository.FindAllUser()
	if err != nil {
		return nil, err
	}

	formattedAdmin := make([]entity.User, 0)
	for _, v := range admin {
		v.Phone, _ = s.encryptTool.Decrypt(v.Phone)
		formattedAdmin = append(formattedAdmin, v)
	}

	return formattedAdmin, nil
}

func (s *adminService) CreateAdmin(admin *entity.Admin) (*entity.Admin, error) {
	if admin.Email == "" {
		return nil, errors.New("email cannot be empty")
	}
	if admin.Password == "" {
		return nil, errors.New("password cannot be empty")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	admin.Password = string(hashedPassword)

	newAdmin, err := s.adminRepository.CreateAdmin(admin)
	if err != nil {
		return nil, err
	}

	err = s.emailSender.SendWelcomeEmail(newAdmin.Email, newAdmin.Fullname, "")
	if err != nil {
		return nil, err
	}

	resetCode := generateResetCode()
	err = s.emailSender.SendVerificationEmail(newAdmin.Email, newAdmin.Fullname, resetCode)
	if err != nil {
		return nil, err
	}

	err = s.adminRepository.SaveVerifCode(newAdmin.User_ID, resetCode)
	if err != nil {
		return nil, err
	}

	return newAdmin, nil
}

func (s *adminService) CheckUserExists(id uuid.UUID) (bool, error) {
	return s.adminRepository.CheckUserExists(id)
}

func (s *adminService) UpdateAdmin(admin *entity.Admin) (*entity.Admin, error) {
	if admin.Email == "" {
		return nil, errors.New("email cannot be empty")
	}
	if admin.Password == "" {
		return nil, errors.New("password cannot be empty")
	}
	if admin.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		admin.Password = string(hashedPassword)
	}
	if admin.Phone != "" {
		admin.Phone, _ = s.encryptTool.Encrypt(admin.Phone)
	}

	updatedAdmin, err := s.adminRepository.UpdateAdmin(admin)
	if err != nil {
		return nil, err
	}

	return updatedAdmin, nil
}

func (s *adminService) DeleteAdmin(user_Id uuid.UUID) (bool, error) {
	user, err := s.adminRepository.FindAdminByID(user_Id)
	if err != nil {
		return false, err
	}

	return s.adminRepository.DeleteAdmin(user)
}

func (s *adminService) EmailExists(email string) bool {
	_, err := s.adminRepository.FindAdminByEmail(email)
	return err == nil
}

func generateResetCode() string {
	return time.Now().Format("150405") // optional: simple reset code
}
