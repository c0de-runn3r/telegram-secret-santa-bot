package storage

import (
	"log"
	"math/rand"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SantaUser struct {
	ID       int `gorm:"primaryKey"`
	Username string
	SantaID  int
	Wishes   string
	Game     string
	ChatID   int
}

type Game struct {
	ID     int
	Name   string
	Admin  string
	Rolled bool
}

type DataBase struct {
	*gorm.DB
}

type UpdateWishesList struct {
	Wishes []*WishUpdateInfo
}

type WishUpdateInfo struct {
	ID       int
	Username string
	Wish     string
}

var ListOfWishUpdates = UpdateWishesList{}

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

func (db *DataBase) AddNewGame(name string, username string, chatID int) int {
	log.Println("adding new game to database Game")
	rand.Seed(time.Now().UnixNano())
	id := rand.Intn(9_999_999-1_000_000) + 1_000_000
	db.Create(&Game{
		ID:     id,
		Name:   name,
		Admin:  username,
		Rolled: false,
	})
	var game Game
	db.First(&game, "name = ?", name)
	db.Create(&SantaUser{
		Username: username,
		Game:     name,
		SantaID:  game.ID,
		ChatID:   chatID,
	})
	return game.ID
}

func (db *DataBase) AddUserToGame(game *Game, username string, chatID int) {
	log.Printf("adding user [%s] to game [%s]", username, game.Name)
	db.Create(&SantaUser{
		Username: username,
		Game:     game.Name,
		SantaID:  game.ID,
		ChatID:   chatID,
	})
}

func (db *DataBase) AddOrUpdateWishes(username string) {
	for _, wish := range ListOfWishUpdates.Wishes {
		if wish.Username == username {
			var user SantaUser
			db.First(&user, "username = ? AND santa_id = ?", username, wish.ID)
			db.Model(&user).Update("Wishes", wish.Wish)
		}
	}

}

func (db *DataBase) QueryAllPlayers(gameID int) ([]SantaUser, error) {
	var users []SantaUser
	db.Table("santa_users").Where("santa_id = ?", gameID).Find(&users)
	return users, nil
}

func (db *DataBase) QueryAdmin(gameID int) (string, error) {
	var game Game
	db.Table("games").Where("id = ?", gameID).First(&game)
	return game.Admin, nil
}

func (db *DataBase) DeleteUserFromGame(username string, gameID int) {
	var user SantaUser
	db.Table("santa_users").Where("username = ? AND santa_id = ?", username, gameID).Delete(&user)
}

func (db *DataBase) DeleteGameAndAllUsers(gameID int) {
	db.Table("santa_users").Where("santa_id = ?", gameID).Delete(&SantaUser{})
	var game Game
	db.Table("games").Where("id = ?", gameID).Delete(&game)
}
