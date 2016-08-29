package controllers

import (
	"github.com/revel/revel"
	"fmt"
	"github.com/aRusenov/GoBlog/app/models/binding"
	"github.com/aRusenov/GoBlog/app/models"
	"github.com/aRusenov/GoBlog/app/models/view"
	"time"
	"github.com/microcosm-cc/bluemonday"
)

type Post struct {
	BaseController
}

func (c Post) Add() revel.Result {
	if user := c.User(); user == nil || user.Role != models.ADMIN {
		return c.Redirect(App.Index)
	}

	return c.Render()
}

func (c Post) New() revel.Result {
	if user := c.User(); user == nil || user.Role != models.ADMIN {
		return c.Redirect(App.Index)
	}

	var postData binding.PostNew
	c.Params.Bind(&postData, "post")
	c.Validation.Length(len(postData.Title), 5)
	c.Validation.Length(len(postData.Description), 10)
	c.Validation.Length(len(postData.Content), 30)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		return c.Redirect(Post.Add)
	}

	sanitizer := bluemonday.UGCPolicy()
	sanitizedContent := sanitizer.Sanitize(postData.Content)
	post := models.Post {
		Title: postData.Title,
		Content: sanitizedContent,
		Description: postData.Description,
		CreatedByID: c.Args["userId"].(uint),
		CreatedAt: time.Now(),
	}

	c.Txn.Save(&post)
	fmt.Println(post.Id)

	return c.Redirect("/post/%d/view", post.Id)
}

func (c Post) View(id int) revel.Result {
	fmt.Println(id)

	posts := []view.PostViewModel{}
	c.Txn.Table("posts").
		Select("posts.id, posts.title, posts.content, posts.created_at, users.username").
		Where("posts.id = ?", id).
		Joins("left join users on users.id = posts.created_by_id").
		Scan(&posts)

	if len(posts) == 0 {
		return c.NotFound("Post does not exist")
	}

	c.RenderArgs["post"] = posts[0]
	return c.Render()
}

func (c Post) Edit(id int) revel.Result {
	return c.Render()
}

func (c Post) Delete(id int) revel.Result {
	return c.Render()
}
