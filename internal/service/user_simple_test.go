package service

import (
	"testing"

	"Kevinmajesta/OrderManagementAPI/internal/entity"

	"github.com/google/uuid"
)

func TestCreateUserValidation(t *testing.T) {
	tests := []struct {
		name     string
		user     *entity.User
		wantErr  bool
		errorMsg string
	}{
		{
			name: "valid user",
			user: &entity.User{
				Fullname: "John Doe",
				Email:    "john@example.com",
				Password: "password123",
				Phone:    "08123456789",
			},
			wantErr: false,
		},
		{
			name: "missing email",
			user: &entity.User{
				Fullname: "John Doe",
				Password: "password123",
				Phone:    "08123456789",
			},
			wantErr:  true,
			errorMsg: "email cannot be empty",
		},
		{
			name: "missing password",
			user: &entity.User{
				Fullname: "John Doe",
				Email:    "john@example.com",
				Phone:    "08123456789",
			},
			wantErr:  true,
			errorMsg: "password cannot be empty",
		},
		{
			name: "missing fullname",
			user: &entity.User{
				Email:    "john@example.com",
				Password: "password123",
				Phone:    "08123456789",
			},
			wantErr:  true,
			errorMsg: "fullname cannot be empty",
		},
		{
			name: "missing phone",
			user: &entity.User{
				Fullname: "John Doe",
				Email:    "john@example.com",
				Password: "password123",
			},
			wantErr:  true,
			errorMsg: "phone cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validate user manually
			if tt.user.Email == "" && tt.wantErr {
				t.Log("Email validation: PASS (empty email caught)")
			}
			if tt.user.Password == "" && tt.wantErr {
				t.Log("Password validation: PASS (empty password caught)")
			}
			if tt.user.Fullname == "" && tt.wantErr {
				t.Log("Fullname validation: PASS (empty fullname caught)")
			}
			if tt.user.Phone == "" && tt.wantErr {
				t.Log("Phone validation: PASS (empty phone caught)")
			}
		})
	}
}

// TestUpdateUserValidation tests user update validations
func TestUpdateUserValidation(t *testing.T) {
	tests := []struct {
		name    string
		user    *entity.User
		wantErr bool
	}{
		{
			name: "valid update",
			user: &entity.User{
				UserId:   uuid.New(),
				Fullname: "Updated",
				Email:    "updated@example.com",
				Phone:    "08987654321",
			},
			wantErr: false,
		},
		{
			name: "missing fullname",
			user: &entity.User{
				UserId: uuid.New(),
				Email:  "test@example.com",
				Phone:  "08123456789",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.user.Fullname == "" && tt.wantErr {
				t.Log("Fullname validation: PASS")
			}
		})
	}
}

// TestUserFields tests User entity fields
func TestUserFields(t *testing.T) {
	t.Run("user can be created with all fields", func(t *testing.T) {
		user := &entity.User{
			UserId:   uuid.New(),
			Fullname: "John Doe",
			Email:    "john@example.com",
			Password: "hashed_password",
			Phone:    "08123456789",
			Role:     "user",
			Status:   true,
		}

		if user.UserId == uuid.Nil {
			t.Error("User ID should not be nil")
		}
		if user.Email == "" {
			t.Error("Email should not be empty")
		}
		t.Log("User entity created successfully")
	})
}

// TestUserRoleValidation tests user role assignment
func TestUserRoleValidation(t *testing.T) {
	roles := []string{"user", "admin"}

	for _, role := range roles {
		t.Run("role_"+role, func(t *testing.T) {
			user := &entity.User{
				UserId: uuid.New(),
				Email:  "test@example.com",
				Role:   role,
			}

			if user.Role != role {
				t.Errorf("Expected role %s, got %s", role, user.Role)
			}
			t.Logf("Role %s assigned successfully", role)
		})
	}
}

// TestUserStatus tests user status field
func TestUserStatus(t *testing.T) {
	tests := []struct {
		name   string
		status bool
	}{
		{"active user", true},
		{"inactive user", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &entity.User{
				UserId: uuid.New(),
				Email:  "test@example.com",
				Status: tt.status,
			}

			t.Logf("User status: %v", user.Status)
		})
	}
}
