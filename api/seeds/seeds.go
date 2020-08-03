package seeds

import (
	"log"

	"github.com/H-Richard/talent/api/models"
	"github.com/jinzhu/gorm"
)

var users = []models.User{
	models.User{
		Email:     "executive@gmail.com",
		Password:  "password",
		FirstName: "Executive",
		LastName:  "Doe",
		Executive: true,
	},
	models.User{
		Email:     "applicant@gmail.com",
		Password:  "password",
		FirstName: "Applicant",
		LastName:  "Doe",
	},
}

var posts = []models.Post{
	models.Post{
		Title:        "Software Developer",
		Description:  "Write Software",
		Requirements: []string{"Python", "TypeScript"},
		Desirements:  []string{"Flask", "Ember"},
	},
	models.Post{
		Title:        "VP of Marketing",
		Description:  "Market Things",
		Requirements: []string{"Instagram", "Facebook"},
		Desirements:  []string{"Skills", "Youtube"},
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
}
