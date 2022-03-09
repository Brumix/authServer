package repository

import (
	"authServer/errors"
	"authServer/models"
	"authServer/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"strconv"
)

type RepoService interface {
	RegisterUser(user models.User) error
	ExistUser(request *models.User) (models.User, error)
	ChangePassword(request *models.ChangePass) error
	DeleteUser(request *models.User) error
}

type RepositoryStruck struct{}

func NewRepository() RepoService {
	return &RepositoryStruck{}
}

func initConnection() *gorm.DB {

	host := os.Getenv("HOST")
	user := os.Getenv("DBUSER")
	password := os.Getenv("PASSWORD")
	dbname := os.Getenv("DBNAME")
	port, _ := strconv.Atoi(os.Getenv("DBPORT"))

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Europe/Lisbon", host, user, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	errors.ErrorRepositoryF("Error connecting with database", err)

	return db
}

func init() {
	migrations()
	seeds()
}

func migrations() {
	db := initConnection()
	defer func() {
		var dbC, _ = db.DB()
		errors.ErrorRepository("ERROR closing the database", dbC.Close())
	}()
	errors.ErrorRepository("Fail executing USER Migrations", db.AutoMigrate(&models.User{}))

	log.Debug("[REPOSITORY] AutoMigration Complete!!")
}

func seeds() {
	var db = initConnection()
	defer func() {
		var con, _ = db.DB()
		con.Close()
	}()
	err := db.Transaction(func(tx *gorm.DB) error {
		var user = models.User{
			UserName: "Bruno",
			Email:    "bruno@gmail.com",
			Password: "123",
		}
		if err := tx.Create(&user).Error; err != nil {
			errors.ErrorRepository("ERROR creating the user", err)
			return err
		}
		return nil
	})
	if err != nil {
		errors.ErrorRepository("ERROR CREATING THE SEEDS", err)
	}

}

func (repo *RepositoryStruck) RegisterUser(user models.User) error {
	var db = initConnection()
	defer func() {
		var con, _ = db.DB()
		con.Close()
	}()
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			log.Error("[REPOSITORY] ERROR creating the user: %v", err)
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (repo *RepositoryStruck) ExistUser(user *models.User) (models.User, error) {
	var db = initConnection()
	var userdb models.User
	defer func() {
		var con, _ = db.DB()
		con.Close()
	}()
	err := db.Transaction(func(tx *gorm.DB) error {
		tx.Where("email = ?", user.Email).First(&userdb)
		return nil
	})

	if err != nil {
		return models.User{}, err
	}
	if utils.ComparePasswords(userdb.Password, user.Password) {
		return userdb, nil
	}

	return models.User{}, fmt.Errorf("wrong Credentials")
}

func (repo *RepositoryStruck) ChangePassword(request *models.ChangePass) error {

	var db = initConnection()
	user, err := repo.ExistUser(&models.User{Email: request.Email, Password: request.OldPassword})
	if err != nil {
		log.Error(err)
		return err
	}
	user.Password = utils.HashAndSalt(request.NewPassword)
	db.Save(&user)

	return nil
}

func (repo *RepositoryStruck) DeleteUser(request *models.User) error {
	var db = initConnection()
	user, err := repo.ExistUser(request)
	if err != nil {
		log.Error(err)
		return err
	}
	db.Delete(&user)

	return nil
}
