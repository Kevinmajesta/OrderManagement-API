package test

import (
	"errors"
	"testing"


	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"Kevinmajesta/OrderManagementAPI/internal/entity"  
	"Kevinmajesta/OrderManagementAPI/internal/mocks"   
	"Kevinmajesta/OrderManagementAPI/internal/service" 

)

func TestUserService_LoginUser(t *testing.T) {
	plainPassword := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)

	userID := uuid.New()
	userEmail := "test@example.com"
	userPhoneEncrypted := "encrypted_phone_data"
	userPhoneDecrypted := "08123456789"
	expectedToken := "mocked_jwt_token_from_service"

	tests := []struct {
		name       string
		email      string
		password   string
		setupMocks func(
			*mocks.MockUserRepository, // Mengirimkan mock sebagai argumen
			*mocks.MockEncryptTool,
			*mocks.MockTokenUseCase,
			*mocks.MockEmailSenderService,
		) // Fungsi untuk mengatur ekspektasi mock
		expectedToken string
		expectedError error
	}{
		{
			name:     "Successful Login - User Role",
			email:    userEmail,
			password: plainPassword,
			setupMocks: func(
				mockUserRepo *mocks.MockUserRepository,
				mockEncryptTool *mocks.MockEncryptTool,
				mockTokenUseCase *mocks.MockTokenUseCase,
				mockEmailSender *mocks.MockEmailSenderService,
			) {
				mockUserRepo.EXPECT().FindUserByEmail(userEmail).Return(&entity.User{
					UserId:   userID,
					Email:    userEmail,
					Password: string(hashedPassword),
					Role:     "user",
					Phone:    userPhoneEncrypted,
				}, nil).Times(1)

				mockEncryptTool.EXPECT().Decrypt(userPhoneEncrypted).Return(userPhoneDecrypted, nil).Times(1)

				mockTokenUseCase.EXPECT().GenerateAccessToken(gomock.Any()).Return(expectedToken, nil).Times(1)

				mockUserRepo.EXPECT().UpdateUserJwtToken(userID, expectedToken, gomock.Any()).Return(nil).Times(1)
			},
			expectedToken: expectedToken,
			expectedError: nil,
		},
		{
			name:     "Failed Login - User Not Found",
			email:    "nonexistent@example.com",
			password: plainPassword,
			setupMocks: func(
				mockUserRepo *mocks.MockUserRepository,
				mockEncryptTool *mocks.MockEncryptTool,
				mockTokenUseCase *mocks.MockTokenUseCase,
				mockEmailSender *mocks.MockEmailSenderService,
			) {
				mockUserRepo.EXPECT().FindUserByEmail("nonexistent@example.com").Return(nil, errors.New("record not found")).Times(1)
			},
			expectedToken: "",
			expectedError: errors.New("wrong input email/password"),
		},
		{
			name:     "Failed Login - Admin Role",
			email:    userEmail,
			password: plainPassword,
			setupMocks: func(
				mockUserRepo *mocks.MockUserRepository,
				mockEncryptTool *mocks.MockEncryptTool,
				mockTokenUseCase *mocks.MockTokenUseCase,
				mockEmailSender *mocks.MockEmailSenderService,
			) {
				mockUserRepo.EXPECT().FindUserByEmail(userEmail).Return(&entity.User{
					UserId:   userID,
					Email:    userEmail,
					Password: string(hashedPassword),
					Role:     "admin",
					Phone:    userPhoneEncrypted,
				}, nil).Times(1)
			},
			expectedToken: "",
			expectedError: errors.New("you dont have access"),
		},
		{
			name:     "Failed Login - Incorrect Password",
			email:    userEmail,
			password: "wrong_password",
			setupMocks: func(
				mockUserRepo *mocks.MockUserRepository,
				mockEncryptTool *mocks.MockEncryptTool,
				mockTokenUseCase *mocks.MockTokenUseCase,
				mockEmailSender *mocks.MockEmailSenderService,
			) {
				mockUserRepo.EXPECT().FindUserByEmail(userEmail).Return(&entity.User{
					UserId:   userID,
					Email:    userEmail,
					Password: string(hashedPassword),
					Role:     "user",
					Phone:    userPhoneEncrypted,
				}, nil).Times(1)
			},
			expectedToken: "",
			expectedError: errors.New("wrong input email/password"),
		},
		{
			name:     "Failed Login - Decryption Error",
			email:    userEmail,
			password: plainPassword,
			setupMocks: func(
				mockUserRepo *mocks.MockUserRepository,
				mockEncryptTool *mocks.MockEncryptTool,
				mockTokenUseCase *mocks.MockTokenUseCase, // Perlu di sini untuk argumen
				mockEmailSender *mocks.MockEmailSenderService,
			) {
				mockUserRepo.EXPECT().FindUserByEmail(userEmail).Return(&entity.User{
					UserId:   userID,
					Email:    userEmail,
					Password: string(hashedPassword),
					Role:     "user",
					Phone:    userPhoneEncrypted,
				}, nil).Times(1)
				mockEncryptTool.EXPECT().Decrypt(userPhoneEncrypted).Return("", errors.New("decryption failed")).Times(1)
				// *** PENTING: TIDAK ADA EKSPEKTASI UNTUK mockTokenUseCase.GenerateAccessToken DI SINI ***
				// Karena LoginUser seharusnya keluar setelah dekripsi gagal.
			},
			expectedToken: "",
			expectedError: errors.New("there is an error in the system"),
		},
		{
			name:     "Failed Login - Token Generation Error",
			email:    userEmail,
			password: plainPassword,
			setupMocks: func(
				mockUserRepo *mocks.MockUserRepository,
				mockEncryptTool *mocks.MockEncryptTool,
				mockTokenUseCase *mocks.MockTokenUseCase,
				mockEmailSender *mocks.MockEmailSenderService,
			) {
				mockUserRepo.EXPECT().FindUserByEmail(userEmail).Return(&entity.User{
					UserId:   userID,
					Email:    userEmail,
					Password: string(hashedPassword),
					Role:     "user",
					Phone:    userPhoneEncrypted,
				}, nil).Times(1)
				mockEncryptTool.EXPECT().Decrypt(userPhoneEncrypted).Return(userPhoneDecrypted, nil).Times(1)
				mockTokenUseCase.EXPECT().GenerateAccessToken(gomock.Any()).Return("", errors.New("failed to generate token")).Times(1)
			},
			expectedToken: "",
			expectedError: errors.New("there is an error in the system"),
		},
		{
			name:     "Failed Login - Update User JWT Token Error",
			email:    userEmail,
			password: plainPassword,
			setupMocks: func(
				mockUserRepo *mocks.MockUserRepository,
				mockEncryptTool *mocks.MockEncryptTool,
				mockTokenUseCase *mocks.MockTokenUseCase,
				mockEmailSender *mocks.MockEmailSenderService,
			) {
				mockUserRepo.EXPECT().FindUserByEmail(userEmail).Return(&entity.User{
					UserId:   userID,
					Email:    userEmail,
					Password: string(hashedPassword),
					Role:     "user",
					Phone:    userPhoneEncrypted,
				}, nil).Times(1)
				mockEncryptTool.EXPECT().Decrypt(userPhoneEncrypted).Return(userPhoneDecrypted, nil).Times(1)
				mockTokenUseCase.EXPECT().GenerateAccessToken(gomock.Any()).Return(expectedToken, nil).Times(1)
				mockUserRepo.EXPECT().UpdateUserJwtToken(userID, expectedToken, gomock.Any()).Return(errors.New("failed to save token")).Times(1)
			},
			expectedToken: "",
			expectedError: errors.New("failed to update user token info"),
		},
	}

	// Jalankan Setiap Skenario Test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// --- Inisialisasi Gomock Controller dan Mock di dalam subtest ---
			// Ini memastikan setiap subtest memiliki lingkungan mock yang bersih.
			ctrl := gomock.NewController(t)
			defer ctrl.Finish() // Finish() akan dipanggil setelah subtest ini selesai.

			mockUserRepo := mocks.NewMockUserRepository(ctrl)
			mockEncryptTool := mocks.NewMockEncryptTool(ctrl)
			mockTokenUseCase := mocks.NewMockTokenUseCase(ctrl)
			mockEmailSender := mocks.NewMockEmailSenderService(ctrl)

			// Buat Instance UserService dengan Inject Mock
			userService := service.NewUserService(mockUserRepo, mockTokenUseCase, mockEncryptTool, mockEmailSender)

			// Atur ekspektasi mock untuk skenario saat ini.
			// Mengirimkan instance mock yang baru dibuat ke setupMocks.
			tt.setupMocks(mockUserRepo, mockEncryptTool, mockTokenUseCase, mockEmailSender)

			// Panggil fungsi LoginUser yang sebenarnya.
			token, err := userService.LoginUser(tt.email, tt.password)

			// Verifikasi Hasil (Assert)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Equal(t, tt.expectedToken, token)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedToken, token)
			}
		})
	}
}

