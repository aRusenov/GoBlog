package controllers

import (
	"database/sql"
	"fmt"
	"github.com/aRusenov/GoBlog/app/models"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
)

type GormController struct {
	*revel.Controller
	Txn *gorm.DB
}

var GormDb *gorm.DB

func InitDB() {
	fmt.Println("Init called")
	var err error
	GormDb, err = gorm.Open("sqlite3", "new-test.db")
	if err != nil {
		panic(err)
	}

	err = GormDb.DB().Ping()
	if err != nil {
		panic(err)
	}

	GormDb.AutoMigrate(true)

	//gormDb.DropTableIfExists(&models.User{}, &models.Post{}, &models.Comment{}, &models.AuthToken{})
	GormDb.CreateTable(&models.User{}, &models.Post{}, &models.Comment{}, &models.AuthToken{})
	var count int
	if GormDb.Model(models.User{}).Count(&count); count == 0 {
		seed()
	}
}

func seed() {
	fmt.Println("Seeding db data")
	bCryptHash, _ := bcrypt.GenerateFromPassword(
		[]byte("admin"), bcrypt.DefaultCost)

	user := models.User{
		Name:           "admin",
		Username:       "admin",
		HashedPassword: bCryptHash,
		Role: 		models.ADMIN,
	}

	GormDb.Save(&user)
}

func (c *GormController) logged() *models.User {
	if username, ok := c.Session["name"]; ok {
		return &models.User{ Username: username }
	}

	return nil
}

// This method fills the c.Txn before each transaction
func (c *GormController) Begin() revel.Result {
	txn := GormDb.Begin()
	if txn.Error != nil {
		panic(txn.Error)
	}
	c.Txn = txn
	return nil
}

// This method clears the c.Txn after each transaction
func (c *GormController) Commit() revel.Result {
	if c.Txn == nil {
		return nil
	}
	c.Txn.Commit()
	if err := c.Txn.Error; err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}

// This method clears the c.Txn after each transaction, too
func (c *GormController) Rollback() revel.Result {
	if c.Txn == nil {
		return nil
	}
	c.Txn.Rollback()
	if err := c.Txn.Error; err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}
