package models

import (
	"fmt"
	"errors"
	"html"
	"strings"
	"time"

	"github.com/H-Richard/truemail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User model
type User struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	FirstName string    `gorm:"size:255;not null;" json:"firstName"`
	LastName  string    `gorm:"size:255;not null;" json:"lastName"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	Executive bool      `gorm:"default:true" json:"executive"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// HashPassword function for password encryption
func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// VerifyPassword function for password verification
func VerifyPassword(hashedPassword, password string) error {
	fmt.Printf("comparing %s with %s", hashedPassword, password)
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Pre function to prepare user struct
func (u *User) Pre() {
	u.ID = 0
	u.FirstName = html.EscapeString(strings.TrimSpace(u.FirstName))
	u.LastName = html.EscapeString(strings.TrimSpace(u.LastName))
	u.UpdatedAt = time.Now()
}

// Post function to hash user password
func (u *User) Post() error {
	preHashed := u.Password
	hashedPassword, err := HashPassword(u.Password)
	err = VerifyPassword(string(hashedPassword), preHashed)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// Helper validation function validates all user fields
func (u *User) validateAll() error {
	v:= truemail.Validator{}
	if u.FirstName == "" {
		return errors.New("FirstName musn't be empty")
	}
	if u.LastName == "" {
		return errors.New("LastName musn't be empty")
	}
	if u.Email == "" {
		return errors.New("Email musn't be empty")
	}
	if err := v.Validate(u.Email) ; err != nil {
		return errors.New("Email format is Invalid")
	}
	return nil
}

// Validate function for user validation
func (u *User) Validate(action string) error {
	v:= truemail.Validator{}
	switch strings.ToLower(action) {
	case "login":
		if u.Password == "" {
			return errors.New("Password musn't be empty")
		}
		if u.Email == "" {
			return errors.New("Email musn't be empty")
		}
		if err := v.Validate(u.Email) ; err != nil {
			return errors.New("Email format is Invalid")
		}
		return nil
	case "update":
		return u.validateAll()
	default:
		return u.validateAll()
	}
}

// SaveUser function for saving users to the database
func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	var err error
	// fmt.Printf("Prehashed = %s", u.Password)
	err = u.Post()
	// fmt.Printf("Posthashed = %s", u.Password)
	if err != nil {
		return &User{}, err
	}
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

// FindByID function finds a user by their ID
func (u *User) FindByID(db *gorm.DB, uid uint32) (*User, error) {
	var err error
	err = db.Debug().Model(User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User Not Found")
	}
	return u, err
}

// FindByEmail function finds a user by their email
func (u *User) FindByEmail(db *gorm.DB, email string) (*User, error) {
	var err error
	err = db.Debug().Model(User{}).Where("email = ?", email).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User Not Found")
	}
	return u, err
}

// All function finds all non-executive users
func (u *User) All(db *gorm.DB) (*[]User, error) {
	var err error
	var users []User = []User{}
	err = db.Debug().Where("executive = ?", true).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	return &users, err
}

// JSON function returns json representation if each user.
func (u *User) JSON() map[string]interface{} {
	return map[string]interface{} {
		"email": u.Email,
		"firstName": u.FirstName,
		"lastName": u.LastName,
		"executive": u.Executive,
		"updatedAt": u.UpdatedAt,
	}
}

