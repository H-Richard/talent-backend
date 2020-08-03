package models

import (
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// Application struct represents applications
type Application struct {
	ID             uint32    `gorm:"primary_key; auto_increment" json:"id"`
	ResumeURL      string    `gorm:"size:255;" json:"resume_url"`
	LinkedInURL    string    `gorm:"size:255;" json:"linkedin_url"`
	GitHubURL      string    `gorm:"size:255;" json:"github_url"`
	PortfolioURL   string    `gorm:"size:255;" json:"portfolio_url"`
	OtherURL       string    `gorm:"size:255;" json:"other_url"`
	AdditionalInfo string    `gorm:"type:text;" json:"additional_info"`
	Applicant      User      `json:"applicant"`
	ApplicantID    uint32    `gorm:"not null" json:"applicant_id"`
	Post           Post      `json:"post"`
	PostID         uint32    `gorm:"not null" json:"post_id"`
	AppliedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"applied_at"`
	Status         int       `gorm:"default:0" json:"status"`
}

// Pre function to prepare Application struct
func (a *Application) Pre() {
	a.ID = 0
	a.LinkedInURL = html.EscapeString(strings.TrimSpace(a.LinkedInURL))
	a.GitHubURL = html.EscapeString(strings.TrimSpace(a.GitHubURL))
	a.PortfolioURL = html.EscapeString(strings.TrimSpace(a.PortfolioURL))
	a.OtherURL = html.EscapeString(strings.TrimSpace(a.OtherURL))
	a.Applicant = User{}
	a.Post = Post{}
	a.AdditionalInfo = html.EscapeString(strings.TrimSpace(a.AdditionalInfo))
}

// Validate function validates application struct
func (a *Application) Validate() error {
	// Figure out something for a.PostID
	return nil
}

// SaveApplication function for saving applications to the database
func (a *Application) SaveApplication(db *gorm.DB) (*Application, error) {
	var err error
	err = db.Debug().Model(&Application{}).Create(&a).Error
	if err != nil {
		return &Application{}, err
	}
	if a.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", a.ApplicantID).Take(&a.Applicant).Error
		if err != nil {
			return &Application{}, err
		}
		err = db.Debug().Model(&Post{}).Where("id = ?", a.PostID).Take(&a.Post).Error
		if err != nil {
			return &Application{}, err
		}
	}
	return a, nil
}

// All returns all active actions
func (a *Application) All(db *gorm.DB) (*[]Application, error) {
	var err error
	applications := []Application{}
	err = db.Debug().Model(&Application{}).Where("status >= ?", 0).Find(&applications).Error
	if err != nil {
		return &[]Application{}, err
	}
	for i := range applications {
		err := db.Debug().Model(&User{}).Where("id = ?", applications[i].ApplicantID).Take(&applications[i].Applicant).Error
		if err != nil {
			return &[]Application{}, err
		}
	}
	return &applications, nil
}

// AllByID returns all applications by the users ID
func (a *Application) AllByID(db *gorm.DB, ID uint32) (*[]Application, error) {
	var err error
	applications := []Application{}
	err = db.Debug().Model(&Application{}).Where("applicant_id = ?", ID).Find(&applications).Error
	if err != nil {
		return &[]Application{}, err
	}
	for i := range applications {
		err := db.Debug().Model(&User{}).Where("id = ?", applications[i].ApplicantID).Take(&applications[i].Applicant).Error
		if err != nil {
			return &[]Application{}, err
		}
	}
	return &applications, nil
}

// JSON gets json representation of each application
func (a *Application) JSON() map[string]interface{} {
	return map[string]interface{}{
		"id":             a.ID,
		"resumeURL":      a.ResumeURL,
		"linkedInURL":    a.LinkedInURL,
		"gitHubURL":      a.GitHubURL,
		"portfolioURL":   a.PortfolioURL,
		"otherURL":       a.OtherURL,
		"additionalInfo": a.AdditionalInfo,
		"applicant":      a.Applicant.JSON(),
		"post":           a.Post.JSON(),
		"status":         a.Status,
		"appliedAt":      a.AppliedAt,
	}
}
