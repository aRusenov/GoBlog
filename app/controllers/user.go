package controllers

import (
	"github.com/revel/revel"
	"github.com/aRusenov/GoBlog/app/models/binding"
	"fmt"
	"github.com/aRusenov/GoBlog/app/models"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	BaseController
}

const MSG_ALREADY_LOGGED = "Already logged in. Please sign out."

func (c User) SignIn() revel.Result {
	return c.Render()
}

func (c User) Login() revel.Result {
	if _, ok := c.Args["userId"]; ok {
		c.Flash.Error(MSG_ALREADY_LOGGED)
		return c.Redirect(App.Index);
	}

	var loginData binding.UserLogin
	c.Params.Bind(&loginData, "user")
	fmt.Println(loginData)

	user := models.User{}
	c.Txn.Where("username = ? ", loginData.Username).First(&user)
	err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(loginData.Password))
	if err != nil {
		c.Flash.Error("Invalid username or password.")
		return c.Redirect(User.SignIn);
	}

	c.setAuthToken(user.ID)
	c.Flash.Success("Welcome back, " + user.Username + "!")

	return c.Redirect(App.Index)
}
func (c User) Register() revel.Result {
	if _, ok := c.Args["userId"]; ok {
		c.Flash.Error(MSG_ALREADY_LOGGED)
		return c.Redirect(App.Index);
	}

	var userData binding.UserRegister
	c.Params.Bind(&userData, "user")
	fmt.Println(userData)

	c.Validation.Required(userData.Username)
	c.Validation.Range(len(userData.Username), 5, 20)
	c.Validation.Required(userData.Password)
	c.Validation.Min(len(userData.Password), 3)
	c.Validation.Required(userData.ConfirmPassword)
	if c.Validation.HasErrors() {
		return nil;
	}

	if userData.Password != userData.ConfirmPassword {
		return nil;
	}

	duplicateUser := models.User{}
	c.Txn.First(&duplicateUser, "username = ?", userData.Username)
	if (duplicateUser.Username != "") {
		return nil;
	}

	passwordHash, _ := bcrypt.GenerateFromPassword(
		[]byte(userData.Password), bcrypt.DefaultCost)

	var user = models.User {
		Username: userData.Username,
		HashedPassword: passwordHash,
	}
	c.Txn.Save(&user)
	fmt.Println("Saved user ", user)

	c.setAuthToken(user.ID)
	c.Flash.Success("Welcome to our site!")

	return c.Redirect(App.Index)
}

func (c User) Logout() revel.Result {
	if _, ok := c.Args["userId"]; !ok {
		c.Flash.Error("You need to be logged in to logout.")
		return c.Redirect(App.Index)
	}

	c.Txn.Where("id = ?", c.Args["tokenId"]).Delete(&models.AuthToken{})
	delete(c.Session, "selector")
	delete(c.Session, "token")

	c.Flash.Success("Logout successful.")
	return c.Redirect(App.Index)
}

func (c User) setAuthToken(userId uint) {
	authToken := CreateAuthToken(userId)
	c.Session["selector"] = authToken.Selector
	c.Session["token"] = authToken.Token

	authToken.Token = HashToken(authToken.Token)
	c.Txn.Save(&authToken)
}
