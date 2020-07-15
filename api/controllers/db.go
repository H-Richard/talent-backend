package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver

	"github.com/H-Richard/talent/api/models"
)

// Server struct
type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

// Initialize connects to db and initializes routes
func (server *Server) Initialize(DBUser, DBPassword, DBPort, DBHost, DBName string) {
	var err error
	DBUrl := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DBHost, DBPort, DBUser, DBName, DBPassword)
	server.DB, err = gorm.Open("postgres", DBUrl)
	if err != nil {
		fmt.Println("Cannot connet to postgres")
		log.Fatal("Err:", err)
	} else {
		fmt.Printf("Connected to postgres on %s:%s", DBHost, DBPort)
	}

	//Migrations
	server.DB.Debug().AutoMigrate(&models.User{})

	server.Router = mux.NewRouter()

	server.initializeRoutes()
}

// Run function starts the server
func (server *Server) Run(addr string) {
	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}