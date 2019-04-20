package model

import (
	"golang.org/x/crypto/bcrypt"
	"time"

	"github.com/jinzhu/gorm"
	// gorm postgres dialect
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Model : gorm.Model definition
type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// User : Model with injected fields `ID`, `CreatedAt`, `UpdatedAt`
type User struct {
	gorm.Model
	Email     string `gorm:"type:varchar(100);not null"`
	Password  string `gorm:"not null"`
	FirstName string `gorm:"type:varchar(50);not null"`
	LastName  string `gorm:"type:varchar(50);not null"`
	Bio       string
	Avatar    string
}

// HashPassword : hashing the password
func (user *User) HashPassword() {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return
	}

	user.Password = string(hash)
}

// ComparePassword : compare the password
func (user *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return false
	}

	return true
}
