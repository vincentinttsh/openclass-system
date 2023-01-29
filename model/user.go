package model

import (
	"time"

	"gorm.io/gorm"
)

// User user model
type User struct {
	BaseModel
	Account        string       `gorm:"not null;unique;index"`
	Email          string       `gorm:"not null;index"`
	Password       string       `gorm:"not null"`
	Name           string       `gorm:"not null"`
	Locale         string       `gorm:"not null;default:'zh-TW'"`
	Admin          bool         `gorm:"not null;index;default:false"`
	SuperAdmin     bool         `gorm:"not null;index;default:false"`
	Department     string       `gorm:"not null;check:department in ('sh','jh', '')"`
	Subject        string       `gorm:"not null;check:subject in ('chinese','english','math','science','social','other', '')"`
	OrganizationID SQLBasePK    `gorm:"not null;index"`
	Organization   Organization `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

// GoogleOauth google oauth map to user model
type GoogleOauth struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	ID        string         `gorm:"primaryKey;not null"`
	UserID    SQLBasePK      `gorm:"not null;index"`
	User      User           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

// Organization organization model (usually a school)
type Organization struct {
	BaseModel
	Level   string `gorm:"index;not null;size:255"` // "city", "province", "country"
	Abbr    string `gorm:"index;not null;size:4"`
	Name    string `gorm:"not null;size:255"`
	Address string `gorm:"not null;size:255"`
}

// GetUserByGoogleID get user by google id
func GetUserByGoogleID(id *string) (User, error) {
	var googleOauth2User GoogleOauth
	googleOauth2User.ID = *id

	result := db.Joins("User").Preload("User.Organization").First(&googleOauth2User)

	return googleOauth2User.User, result.Error
}

// GetUserByID get user by id
func GetUserByID(id *SQLBasePK, user *User) error {
	return db.First(user, id).Error
}

// CreateUserFromGoogle create user from google oauth
func CreateUserFromGoogle(googleOauth2User *GoogleOauth) error {
	return db.Create(googleOauth2User).Error
}

// UpdateUser update user
func UpdateUser(user *User, value *User) error {
	return db.Model(user).Updates(value).Error
}

// AfterCreate set user as super admin and admin when the user is first created
func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	if u.ID == 1 {
		err = tx.Model(u).Updates(User{
			Admin:      true,
			SuperAdmin: true,
		}).Error
	}

	return
}
