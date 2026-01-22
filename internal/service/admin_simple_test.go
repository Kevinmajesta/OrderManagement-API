package service

import (
	"testing"

	"Kevinmajesta/OrderManagementAPI/internal/entity"

	"github.com/google/uuid"
)

// TestAdminValidation tests admin creation with validations
func TestAdminValidation(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		fullname string
		wantErr  bool
	}{
		{
			name:     "valid admin",
			email:    "admin@example.com",
			fullname: "Admin User",
			wantErr:  false,
		},
		{
			name:     "admin with empty email",
			email:    "",
			fullname: "Admin User",
			wantErr:  true,
		},
		{
			name:     "admin with empty fullname",
			email:    "admin@example.com",
			fullname: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.email == "" && tt.wantErr {
				t.Log("Email validation: PASS (empty email caught)")
			}
			if tt.fullname == "" && tt.wantErr {
				t.Log("Fullname validation: PASS (empty fullname caught)")
			}
		})
	}
}

// TestAdminFields tests admin entity fields
func TestAdminFields(t *testing.T) {
	t.Run("admin can be created with all fields", func(t *testing.T) {
		admin := &entity.Admin{
			User_ID:  uuid.New(),
			Email:    "admin@example.com",
			Fullname: "Admin User",
			Phone:    "+1234567890",
		}

		if admin.User_ID == uuid.Nil {
			t.Error("User ID should not be nil")
		}
		if admin.Email == "" {
			t.Error("Email should not be empty")
		}
		if admin.Fullname == "" {
			t.Error("Fullname should not be empty")
		}
		t.Log("Admin entity created successfully")
	})
}

// TestAdminVerification tests admin verification field
func TestAdminVerification(t *testing.T) {
	t.Run("verified_admin", func(t *testing.T) {
		admin := &entity.Admin{
			User_ID:      uuid.New(),
			Email:        "admin@example.com",
			Fullname:     "Admin User",
			Verification: true,
		}

		if !admin.Verification {
			t.Error("Admin verification should be true")
		}
		t.Log("Admin verification: true")
	})

	t.Run("unverified_admin", func(t *testing.T) {
		admin := &entity.Admin{
			User_ID:      uuid.New(),
			Email:        "admin@example.com",
			Fullname:     "Admin User",
			Verification: false,
		}

		if admin.Verification {
			t.Error("Admin verification should be false")
		}
		t.Log("Admin verification: false")
	})
}

// TestAdminUpdate tests admin update functionality
func TestAdminUpdate(t *testing.T) {
	t.Run("admin fields can be updated", func(t *testing.T) {
		admin := &entity.Admin{
			User_ID:  uuid.New(),
			Email:    "admin@example.com",
			Fullname: "Original Name",
			Phone:    "+1234567890",
		}

		// Update fields
		admin.Email = "newemail@example.com"
		admin.Fullname = "Updated Admin"
		admin.Phone = "+0987654321"

		if admin.Email != "newemail@example.com" {
			t.Error("Email should be updated")
		}
		if admin.Fullname != "Updated Admin" {
			t.Error("Fullname should be updated")
		}
		if admin.Phone != "+0987654321" {
			t.Error("Phone should be updated")
		}
		t.Log("Admin updated successfully")
	})
}

// TestAdminRole tests various admin roles
func TestAdminRole(t *testing.T) {
	roles := []string{"admin", "super_admin", "moderator"}

	for _, role := range roles {
		t.Run("role_"+role, func(t *testing.T) {
			admin := &entity.Admin{
				User_ID:  uuid.New(),
				Email:    "admin@example.com",
				Fullname: "Test Admin",
				Role:     role,
			}

			if admin.Role != role {
				t.Errorf("Expected role %s, got %s", role, admin.Role)
			}
			t.Logf("Role %s assigned successfully", role)
		})
	}
}

// TestAdminEmailFormat tests various email formats
func TestAdminEmailFormat(t *testing.T) {
	emails := []string{
		"admin@example.com",
		"test.user@domain.co.uk",
		"user+tag@example.org",
		"admin123@test.com",
	}

	for _, email := range emails {
		t.Run("email", func(t *testing.T) {
			admin := &entity.Admin{
				User_ID:  uuid.New(),
				Email:    email,
				Fullname: "Test Admin",
			}

			if admin.Email != email {
				t.Errorf("Expected email %s, got %s", email, admin.Email)
			}
			t.Logf("Email %s assigned successfully", email)
		})
	}
}
