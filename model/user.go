package model

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AritroSaha10/htn25-backend-takehome/util"
	"github.com/rs/zerolog/log"
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

// UserUpdate represents the fields that can be updated for a user.
type UserUpdate struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	BadgeCode string `json:"badge_code"`
}

// Bind implements the render.Binder interface for UserUpdate.
func (up *UserUpdate) Bind(r *http.Request) error {
	return nil
}

// Render implements the render.Renderer interface for User.
func (u *User) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// GetUsers gets all users from the database.
func GetUsers(db *gorm.DB) ([]User, error) {
	users := []User{}
	res := db.Preload("Scans").Find(&users)
	if res.Error != nil {
		return nil, res.Error
	}
	return users, nil
}

// GetUserByID gets a user from the database by their ID.
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

// UpdateUserByID updates a user in the database by their ID, while enforcing certain
// uniqueness constraints on the fields (ex. email, badge code).
func UpdateUserByID(db *gorm.DB, id uint, userUpdate UserUpdate) (User, error) {
	// Confirm the user exists
	user := User{}
	res := db.
		Model(&user).
		Where("id = ?", id).
		Limit(1).
		Find(&user)
	if res.RowsAffected == 0 {
		return User{}, util.ErrNotFound
	}
	if res.Error != nil {
		return User{}, res.Error
	}

	// Forbid a user from claiming another user's email
	if userUpdate.Email != "" {
		otherUser := User{}
		res := db.
			Model(&otherUser).
			Where("email = ?", userUpdate.Email).
			Limit(1).
			Find(&otherUser)
		if res.RowsAffected != 0 {
			return User{}, fmt.Errorf("%w: email already in use", util.ErrBadRequest)
		}
	}

	// Forbid a user from claiming another user's badge code
	if userUpdate.BadgeCode != "" {
		otherUser := User{}
		res := db.
			Model(&otherUser).
			Where("badge_code = ?", userUpdate.BadgeCode).
			Limit(1).
			Find(&otherUser)
		if res.RowsAffected != 0 {
			return User{}, fmt.Errorf("%w: badge code already in use", util.ErrBadRequest)
		}
	}

	// Perform all updates to the user
	result := db.
		Model(&user).
		Where("id = ?", id).
		Updates(userUpdate)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("failed to update user")
		return User{}, result.Error
	}

	// Also show the scans in the response
	user.Scans = []Scan{}
	result = db.Preload("Scans").Find(&user, id)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("failed to load scans for user")
		return User{}, result.Error
	}
	return user, nil
}
