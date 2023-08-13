package gormdb

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserModel struct {
	ID       uuid.UUID `gorm:"primaryKey"`
	FullName string
	Email    string
	Created  time.Time
	Updated  *time.Time
	Deleted  *time.Time `json:"-"`
}

func (UserModel) TableName() string {
	return "users"
}

// Get UserByID
func GetUserByID(db *gorm.DB, ID uuid.UUID) (*UserModel, error) {
	user := UserModel{}
	if err := db.Where("id=?", ID).Take(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Create a new user
func CreateUser(db *gorm.DB, fullName, email string) (*UserModel, error) {
	user := UserModel{
		ID:       uuid.New(),
		FullName: fullName,
		Email:    email,
		Created:  time.Now(),
	}

	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// Update user
func UpdateUser(db *gorm.DB, ID uuid.UUID, newEmail string) (int64, error) {
	updated := time.Now()

	result := db.Debug().Where("id=?", ID).Updates(&UserModel{
		Email:   newEmail,
		Updated: &updated,
	})

	return result.RowsAffected, result.Error
}

// Delete User
func DeleteUser(db *gorm.DB, ID uuid.UUID) error {
	return db.Debug().Delete(&UserModel{}, ID).Error
}
