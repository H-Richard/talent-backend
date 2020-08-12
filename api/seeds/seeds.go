package seeds

import (
	"log"

	"github.com/H-Richard/talent/api/models"
	"github.com/jinzhu/gorm"
)

var users = []models.User{
	{
		Email:     "executive@gmail.com",
		Password:  "password",
		FirstName: "Executive",
		LastName:  "Doe",
		Executive: true,
	},
	{
		Email:     "applicant1@gmail.com",
		Password:  "password",
		FirstName: "Applicant1",
		LastName:  "Doe",
	},
	{
		Email:     "applicant2@gmail.com",
		Password:  "password",
		FirstName: "Applicant2",
		LastName:  "Doe",
	},
}

var posts = []models.Post{
	{
		Title:        "Software Developer",
		Description:  "Write Software",
		Requirements: []string{"Python", "TypeScript"},
		Desirements:  []string{"Flask", "Ember"},
	},
	{
		Title:        "VP of Marketing",
		Description:  "Market Things",
		Requirements: []string{"Instagram", "Facebook"},
		Desirements:  []string{"Skills", "Youtube"},
	},
}

var applications = []models.Application{
	{
		ResumeURL:      "google.ca",
		LinkedInURL:    "google.ca",
		GitHubURL:      "google.ca",
		PortfolioURL:   "google.ca",
		AdditionalInfo: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
		OtherURL:       "youtube.com",
	},
	{
		ResumeURL:      "google.ca",
		LinkedInURL:    "google.ca",
		GitHubURL:      "google.ca",
		PortfolioURL:   "google.ca",
		AdditionalInfo: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
		OtherURL:       "youtube.com",
	},
}

// Migrate runs migrations
func Migrate(db *gorm.DB) {
	err := db.Debug().DropTableIfExists(&models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}
	err = db.Debug().DropTableIfExists(&models.Post{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.Post{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}
	err = db.Debug().DropTableIfExists(&models.Application{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.Application{}).AddUniqueIndex("one_per_post", "applicant_id", "post_id").Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}
}

// Seed function seeds the database
func Seed(db *gorm.DB) {
	Migrate(db)
	for i := range users {
		_, err := users[i].SaveUser(db)
		if err != nil {
			log.Fatalf("cannot insert to table: %v", err)
		}
	}
	for i := range posts {
		posts[i].AuthorID = users[0].ID
		err := db.Model(&models.Post{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("cannot insert to table: %v", err)
		}
	}
	for i := range applications {
		applications[i].ApplicantID = users[1+i].ID
		applications[i].PostID = posts[i].ID
		err := db.Model(&models.Application{}).Create(&applications[i]).Error
		if err != nil {
			log.Fatalf("cannot insert to table: %v", err)
		}
	}
}
