package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

// Post struct represents job posts
type Post struct {
	ID uint32 `gorm:"primary_key;auto_increment" json:"id"`
	Title string `gorm:"size:255; not null" json:"title"`
	Description string `gorm:"type:text; not null" json:"description"`
	Requirements pq.StringArray `gorm:"type:text[]" json:"requirements"`
	Desirements pq.StringArray `gorm:"type:text[]" json:"desirements"`
	Active string `gorm:"default:true" json:"active"`
	Author User `json:"author"`
	AuthorID uint32 `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	ExpiresAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"expires_at"`
}

// Pre function to prepare post struct
func (p *Post) Pre() {
	p.ID = 0
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Description = html.EscapeString(strings.TrimSpace(p.Description))
	p.Author = User{}
	p.Desirements = []string{}
	p.Requirements = []string{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	p.ExpiresAt = time.Now()
}

// Validate function validates post struct
func (p *Post) Validate() error {
	if p.Title == "" {
		return errors.New("Title musn't be empty")
	}
	if p.Description == "" {
		return errors.New("Description musn't be empty")
	}
	return nil;
}

// SavePost for saving posts to the database
func (p *Post) SavePost(db *gorm.DB) (*Post, error) {
	var err error
	err = db.Debug().Model(&Post{}).Create(&p).Error
	if err != nil {
		return &Post{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Post{}, err
		}
	}
	return p, nil
}

// FindByID finds posts by their ID
func (p *Post) FindByID(db *gorm.DB, pid uint32) (*Post, error) {
	var err error
	err = db.Debug().Model(&Post{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Post{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Post{}, err
		}
	}
	return p, nil
}

// FindAllActivePosts finds all active job posts
func (p *Post) FindAllActivePosts(db *gorm.DB) (*[]Post, error) {
	var err error
	posts := []Post{}
	err = db.Debug().Model(&Post{}).Where("active = ?", true).Find(&posts).Error
	if err != nil {
		return &[]Post{}, err
	}
	for i := range posts {
		err := db.Debug().Model(&User{}).Where("id = ?", posts[i].AuthorID).Take(&posts[i].Author).Error
		if err != nil {
			return &[]Post{}, err
		}
	}
	return &posts, nil
}

// UpdatePost lets users update their post
func (p *Post) UpdatePost(db *gorm.DB) (*Post, error) {
	var err error
	err = db.Debug().Model(&Post{}).Where("id = ?", p.ID).Updates(
		Post{
			Title: p.Title,
			Description: p.Description,
			Desirements: p.Desirements,
			Requirements: p.Requirements,
			UpdatedAt: time.Now(),
			ExpiresAt: p.ExpiresAt,
		},
	).Error
	if err != nil {
		return &Post{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Post{}, err
		}
	}
	return p, nil
}

// DeletePost deletes a users post
func (p *Post) DeletePost(db *gorm.DB, pid, uid uint32) (int64, error) {
	db = db.Debug().Model(&Post{}).Where("id = ? and author_id = ?", pid, uid).Take(&Post{}).Delete(&Post{})
	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Post not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// JSON gets json representation of each post
func (p *Post) JSON () map[string]interface{} {
	return map[string]interface{}{
		"id": p.ID,
		"title": p.Title,
		"description": p.Description,
		"requirements": p.Requirements,
		"desirements": p.Desirements,
		"active": p.Active,
		"createdAt": p.CreatedAt,
		"updatedAt": p.UpdatedAt,
		"expiresAt": p.ExpiresAt,
		"author": p.Author.JSON(),
	}
}