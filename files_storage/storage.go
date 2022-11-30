package storage

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SantaUser struct {
	ID       uint `gorm:"primaryKey"`
	Game     string
	Username string
	SantaID  uint
	IsAdmin  bool
	Address  string
	Wishes   string
}

type Game struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Admin string
}

type DataBase struct {
	*gorm.DB
}

func NewDB() DataBase {
	db, err := gorm.Open(sqlite.Open("dataBase.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return DataBase{db}
}

func (db *DataBase) MigrateSanta() {
	log.Println("migrating database Santa")
	db.AutoMigrate(&SantaUser{})
}

func (db *DataBase) MigrateGame() {
	log.Println("migrating database Game")
	db.AutoMigrate(&Game{})
}

func (db *DataBase) AddNewUser(username string, game string, santaID uint, isAdmin bool, address string, wishes string) {
	log.Println("adding new user to database Santa")
	// Create
	db.Create(&SantaUser{
		Username: username,
		Game:     game,
		SantaID:  santaID,
		IsAdmin:  isAdmin,
		Address:  address,
		Wishes:   wishes,
	})
}

func (db *DataBase) AddNewGame(name string, username string) uint {
	log.Println("adding new game to database Game")
	db.Create(&Game{
		Name:  name,
		Admin: username,
	})
	var game Game
	db.First(&game, "name = ?", name)
	return game.ID
}

func (db *DataBase) AddUserToGame(game *Game, username string) {
	log.Printf("adding user [%s] to game [%s]", username, game.Name)
	db.Create(&SantaUser{
		Username: username,
		Game:     game.Name,
	})
}

func (db *DataBase) AddOrUpdateWishes(text string, username string) {
	var user SantaUser
	db.First(&user, "username = ?", username)
	db.Model(&user).Update("Wishes", text)
}
