package controllers

import (
	"github.com/aRusenov/GoBlog/app/models"
	"fmt"
)

type BaseController struct {
	GormController
}

func (c *BaseController) Authenticated() bool {
	if _, ok := c.Args["userId"]; ok {
		return true
	}

	return false
}

func (c *BaseController) User() *models.User {
	if userId, ok := c.Args["userId"]; ok {
		var user models.User
		c.Txn.Where("id = ?", userId).First(&user)
		fmt.Println(user)
		if user.Username == "" {
			return nil
		}

		return &user
	}

	return nil
}
