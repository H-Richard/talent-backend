package seeds

import (
	"log"
	
	"github.com/jinzhu/gorm"
	"github.com/H-Richard/talent/api/models"
)

var users = []models.User{
	models.User{
		Email: "executive@gmail.com",
		Password: "password",
		FirstName: "Executive",
		LastName: "Doe",
		Executive: true,
	},
	models.User{
		Email: "applicant@gmail.com",
		Password: "password",
		FirstName: "Applicant",
		LastName: "Doe",
	},
}


// Seed function seeds the database
func Seed(db *gorm.DB) {
	err := db.Debug().DropTableIfExists(&models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	for i := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot insert to table: %v", err)
		}
	}
}