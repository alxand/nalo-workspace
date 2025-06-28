package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserRole represents user roles
type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleUser    UserRole = "user"
	RoleManager UserRole = "manager"
)

// User represents a user in the system
type User struct {
	ID        int64      `gorm:"primaryKey" json:"id"`
	Email     string     `gorm:"uniqueIndex;not null" json:"email" validate:"required,email"`
	Username  string     `gorm:"uniqueIndex;not null" json:"username" validate:"required,min=3,max=50"`
	Password  string     `gorm:"not null" json:"-" validate:"required,min=8"` // "-" means don't include in JSON
	FirstName string     `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string     `json:"last_name" validate:"required,min=2,max=50"`
	Role      UserRole   `gorm:"default:'user'" json:"role" validate:"required,oneof=admin user manager"`
	IsActive  bool       `gorm:"default:true" json:"is_active"`
	CountryID *int64     `json:"country_id"` // Optional - user may not have a country
	CompanyID *int64     `json:"company_id"` // Optional - user may not have a company
	LastLogin *time.Time `json:"last_login,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

	// Relationships
	DailyTasks []DailyTask `gorm:"foreignKey:UserID" json:"daily_tasks,omitempty"`
	Country    *Country    `gorm:"foreignKey:CountryID" json:"country,omitempty"`
	Company    *Company    `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

// BeforeCreate is a GORM hook that runs before creating a user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Hash password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// BeforeUpdate is a GORM hook that runs before updating a user
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// Only hash password if it has changed
	if tx.Statement.Changed("Password") {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// CheckPassword compares the provided password with the stored hash
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// HasRole checks if the user has a specific role
func (u *User) HasRole(role UserRole) bool {
	return u.Role == role
}

// HasAnyRole checks if the user has any of the specified roles
func (u *User) HasAnyRole(roles ...UserRole) bool {
	for _, role := range roles {
		if u.Role == role {
			return true
		}
	}
	return false
}

// IsAdmin checks if the user is an admin
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// FullName returns the user's full name
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}
