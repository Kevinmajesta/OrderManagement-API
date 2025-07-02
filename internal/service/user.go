package service

import (
	"errors"
	"time"

	"Kevinmajesta/OrderManagementAPI/internal/entity"
	"Kevinmajesta/OrderManagementAPI/internal/repository"
	"Kevinmajesta/OrderManagementAPI/pkg/email"
	"Kevinmajesta/OrderManagementAPI/pkg/encrypt"
	"Kevinmajesta/OrderManagementAPI/pkg/token"
	"Kevinmajesta/OrderManagementAPI/worker"


	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	LoginUser(email string, password string) (string, error)
	CreateUser(user *entity.User) (*entity.User, error)
	UpdateUser(user *entity.User) (*entity.User, error)
	DeleteUser(user_id uuid.UUID) (bool, error)
	RequestPasswordReset(email string) error
	ResetPassword(resetCode string, newPassword string) error
	EmailExists(email string) bool
	GetUserProfileByID(userID string) (*entity.User, error)
	VerifUser(resetCode string) error
	CheckUserExists(id uuid.UUID) (bool, error)
}

type userService struct {
	userRepository repository.UserRepository
	tokenUseCase   token.TokenUseCase
	encryptTool    encrypt.EncryptTool
	emailSender    *email.EmailSender
}

var InternalError = "internal server error"

func NewUserService(userRepository repository.UserRepository, tokenUseCase token.TokenUseCase,
	encryptTool encrypt.EncryptTool, emailSender *email.EmailSender) *userService {

	return &userService{
		userRepository: userRepository,
		tokenUseCase:   tokenUseCase,
		encryptTool:    encryptTool,
		emailSender:    emailSender,
	}
}

func (s *userService) LoginUser(email string, password string) (string, error) {
	user, err := s.userRepository.FindUserByEmail(email)
	if err != nil {
		return "", errors.New("wrong input email/password")
	}
	if user.Role != "user" {
		return "", errors.New("you dont have access")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("wrong input email/password")
	}

	expiredTime := time.Now().Local().Add(24 * time.Hour)

	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic(err)
	}

	user.Phone, _ = s.encryptTool.Decrypt(user.Phone)
	expiredTimeInJakarta := expiredTime.In(location)

	claims := token.JwtCustomClaims{
		ID:    user.UserId.String(),
		Email: user.Email,
		Role:  "user",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Depublic",
			ExpiresAt: jwt.NewNumericDate(expiredTimeInJakarta),
		},
	}

	jwtToken, err := s.tokenUseCase.GenerateAccessToken(claims)
	if err != nil {
		return "", errors.New("there is an error in the system")
	}

	user.JwtToken = jwtToken
	user.JwtTokenExpiresAt = expiredTime

	if err := s.userRepository.UpdateUserJwtToken(user.UserId, jwtToken, expiredTime); err != nil {
		return "", errors.New("failed to update user token info")
	}

	if user.JwtToken != jwtToken {
		return "", errors.New("JWT token mismatch")
	}

	return jwtToken, nil
}

func (s *userService) CreateUser(user *entity.User) (*entity.User, error) {
	if user.Email == "" {
		return nil, errors.New("email cannot be empty")
	}
	if user.Password == "" {
		return nil, errors.New("password cannot be empty")
	}
	if user.Fullname == "" {
		return nil, errors.New("fullname cannot be empty")
	}
	if user.Phone == "" {
		return nil, errors.New("phone cannot be empty")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	newUser, err := s.userRepository.CreateUser(user)
	if err != nil {
		return nil, err
	}

	// Kirim job welcome email
	worker.EmailQueue <- worker.EmailJob{
		Type: "welcome",
		To:   newUser.Email,
		Name: newUser.Fullname,
	}

	// Kirim job verifikasi
	resetCode := generateResetCode()

	worker.EmailQueue <- worker.EmailJob{
		Type:      "verification",
		To:        newUser.Email,
		Name:      newUser.Fullname,
		ResetCode: resetCode,
	}

	err = s.userRepository.SaveVerifCode(user.UserId, resetCode)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *userService) CheckUserExists(id uuid.UUID) (bool, error) {
	return s.userRepository.CheckUserExists(id)
}

func (s *userService) UpdateUser(user *entity.User) (*entity.User, error) {
	if user.Email == "" {
		return nil, errors.New("email cannot be empty")
	}
	if user.Password == "" {
		return nil, errors.New("password cannot be empty")
	}
	if user.Fullname == "" {
		return nil, errors.New("fullname cannot be empty")
	}
	if user.Phone == "" {
		return nil, errors.New("phone cannot be empty")
	}

	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}

	updatedUser, err := s.userRepository.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (s *userService) DeleteUser(user_Id uuid.UUID) (bool, error) {
	user, err := s.userRepository.FindUserByID(user_Id)
	if err != nil {
		return false, err
	}

	return s.userRepository.DeleteUser(user)
}

func (s *userService) RequestPasswordReset(email string) error {
	user, err := s.userRepository.FindUserByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	resetCode := generateResetCode()
	expiresAt := time.Now().Add(1 * time.Hour)

	err = s.userRepository.SaveResetCode(user.UserId, resetCode, expiresAt)
	if err != nil {
		return errors.New("failed to save reset code")
	}

	return s.emailSender.SendResetPasswordEmail(user.Email, user.Fullname, resetCode)
}

func (s *userService) ResetPassword(resetCode string, newPassword string) error {
	user, err := s.userRepository.FindUserByResetCode(resetCode)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("invalid reset code")
	}
	if newPassword == "" {
		return errors.New("password cannot be empty")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	_, err = s.userRepository.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) EmailExists(email string) bool {
	_, err := s.userRepository.FindUserByEmail(email)
	return err == nil
}

func (s *userService) GetUserProfileByID(userID string) (*entity.User, error) {
	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	return s.userRepository.FindUserByID(userIDUUID)
}

func (s *userService) VerifUser(verifCode string) error {
	user, err := s.userRepository.FindUserByVerifCode(verifCode)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("invalid verification code")
	}

	user.Verification = true
	_, err = s.userRepository.UpdateUser(user)
	return err
}
