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
	"time"
)

type RepoService interface {
	RegisterUser(user models.User) error
	ExistUser(request *models.User) (models.User, error)
	ChangePassword(request *models.ChangePass) error
	DeleteUser(request *models.User) error
}

var db = initConnection()

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
	dbGorm, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	errors.ErrorRepositoryF("Error connecting with database", err)
	sqlDB, _ := dbGorm.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)
	return dbGorm
}

func init() {
	migrations()
	seeds()
}

func migrations() {
	errors.ErrorRepository("Fail executing USER Migrations", db.AutoMigrate(&models.User{}))

	log.Debug("[REPOSITORY] AutoMigration Complete!!")
}

func seeds() {
	err := db.Transaction(func(tx *gorm.DB) error {
		var user = models.User{
			UserName: "Bruno",
			Email:    "bruno@gmail.com",
			Password: utils.HashAndSalt("123"),
			Role:     "ADMIN",
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

	var userdb models.User

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
	user, err := repo.ExistUser(request)
	if err != nil {
		log.Error(err)
		return err
	}
	db.Delete(&user)

	return nil
}
