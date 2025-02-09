package model

import (
	"net/http"
	"time"

	"github.com/AritroSaha10/htn25-backend-takehome/util"
	"gorm.io/gorm"
)

// GORM model for a user. We aren't using gorm.Model so we can
// add json tags to the fields it provides.
// TODO: Add UUID field to User model.
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `json:"name"`
	Email     string         `json:"email" gorm:"unique;not null"`
	Phone     string         `json:"phone" gorm:"not null"`
	BadgeCode string         `json:"badge_code"`
	Scans     []Scan         `json:"scans"`
}

// Render implements the render.Renderer interface for User.
func (u *User) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func GetUsers(db *gorm.DB) ([]User, error) {
	users := []User{}
	res := db.Preload("Scans").Find(&users)
	if res.Error != nil {
		return nil, res.Error
	}
	return users, nil
}

func GetUserByID(db *gorm.DB, id uint) (User, error) {
	user := User{}
	result := db.Preload("Scans", "user_id = ?", id).Limit(1).Find(&user, id)
	if result.RowsAffected == 0 {
		return User{}, util.ErrNotFound
	}
	if result.Error != nil {
		return User{}, result.Error
	}
	return user, nil
}
