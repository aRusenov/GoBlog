package controllers

import (
	"github.com/revel/revel"
	"crypto/sha256"
	"encoding/hex"
	"encoding/base64"
	"crypto/rand"
	"fmt"
	"time"
	"github.com/aRusenov/GoBlog/app/models"
)

func CreateAuthToken(userId uint) models.AuthToken {
	selector, _ := generateRandomString(12)
	token, _ := generateRandomString(32)
	authToken := models.AuthToken{
		UserID: userId,
		Expires: time.Now().AddDate(0, 1, 0),
		Selector: selector,
		Token: token,
	}

	return authToken
}

func AuthFilter(c *revel.Controller) revel.Result {
	selector, _ := c.Session["selector"]
	token, _ := c.Session["token"]

	hashedToken := HashToken(token)
	authToken := models.AuthToken{}
	fmt.Println(selector, " ", hashedToken)
	GormDb.First(&authToken, "selector = ? and token = ?", selector, hashedToken)

	if authToken.Token != "" {
		fmt.Println("Token is valid")
		authToken.Expires = time.Now().AddDate(0, 1, 0)
		GormDb.Save(&authToken)

		c.Args["userId"] = authToken.UserID
		c.Args["tokenId"] = authToken.Id
		c.RenderArgs["logged"] = true
	}

	return nil
}

//func AuthorizeUser(c *revel.Controller) revel.Result {
//	if _, ok := c.Args["userId"]; !ok {
//		c.Flash.Error("You need to be logged to view this page.")
//		return c.Redirect(controllers.App.Index)
//	}
//
//	return nil;
//}

func HashToken(token string) string {
	h := sha256.New()
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func generateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}