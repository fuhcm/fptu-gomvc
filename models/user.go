package models

import (
	"errors"
	"time"

	"github.com/gosu-team/fptu-api/config"
)

// User ...
type User struct {
	ID        int        `json:"id" gorm:"primary_key"`
	CreatedAt *time.Time `json:"created_at, omitempty"`
	UpdatedAt *time.Time `json:"updated_at, omitempty"`
	DeletedAt *time.Time `json:"deleted_at, omitempty" sql:"index"`

	Email    string `json:"email" gorm:"not null; type:varchar(250); unique_index"`
	Password string `json:"password" gorm:"not null; type:varchar(250)"`
	Admin    string `json:"admin" gorm: "not null; type:boolean"`
	Nickname string `json:"nickname" gorm: "type: varchar(250);"`
}

// TableName set User's table name to be `users`
func (User) TableName() string {
	return "users"
}

// FetchAll ...
func (u *User) FetchAll() []User {
	db := config.GetDatabaseConnection()

	var users []User
	db.Find(&users)

	return users
}

// FetchByID ...
func (u *User) FetchByID() error {
	db := config.GetDatabaseConnection()

	if err := db.Where("id = ?", u.ID).Find(&u).Error; err != nil {
		return errors.New("Could not find the user")
	}

	return nil
}

// FetchByEmail ...
func (u *User) FetchByEmail() error {
	db := config.GetDatabaseConnection()

	if err := db.Where("email = ?", u.Email).Find(&u).Error; err != nil {
		return errors.New("Could not find the user")
	}

	return nil
}

// FetchEmailByID ...
func (u *User) FetchEmailByID(userID int) string {
	db := config.GetDatabaseConnection()

	db.Where("id = ?", userID).Take(&u)

	return u.Email
}

// FetchNicknameByID ...
func (u *User) FetchNicknameByID(userID int) string {
	db := config.GetDatabaseConnection()

	db.Where("id = ?", userID).Take(&u)

	return u.Nickname
}

// Create ...
func (u *User) Create() error {
	db := config.GetDatabaseConnection()

	// Validate record
	if !db.NewRecord(u) { // => returns `true` as primary key is blank
		return errors.New("New records can not have primary key id")
	}

	if err := db.Create(&u).Error; err != nil {
		return errors.New("Could not create user")
	}

	return nil
}

// Save ...
func (u *User) Save() error {
	db := config.GetDatabaseConnection()

	if db.NewRecord(u) {
		if err := db.Create(&u).Error; err != nil {
			return errors.New("Could not create user")
		}
	} else {
		if err := db.Save(&u).Error; err != nil {
			return errors.New("Could not update user")
		}
	}

	return nil
}

// Delete ...
func (u *User) Delete() error {
	db := config.GetDatabaseConnection()

	if err := db.Delete(&u).Error; err != nil {
		return errors.New("Could not find the user")
	}

	return nil
}
