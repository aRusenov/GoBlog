package controllers

import (
	"github.com/revel/revel"
	"github.com/aRusenov/GoBlog/app/models/view"
	"strconv"
)

type App struct {
	BaseController
}

const LIMIT int = 5

func (c App) Index() revel.Result {
	page, _ := strconv.Atoi(c.Params.Query.Get("page"))

	posts := []view.PostViewModel{}
	c.Txn.Table("posts").
		Select("posts.id, posts.title, posts.content, posts.created_at, users.username").
		Order("posts.created_at desc").
		Joins("left join users on users.id = posts.created_by_id").
		Limit(LIMIT).
		Offset(page * LIMIT).
		Scan(&posts)

	var totalPostCount int
	c.Txn.Model(&Post{}).Count(&totalPostCount)
	pageCount := totalPostCount / LIMIT
	if totalPostCount % LIMIT != 0 {
		pageCount++
	}

	pages := make([]int, pageCount, pageCount)
	for i := 0; i < pageCount; i++  {
		pages[i] =  i
	}

	c.RenderArgs["posts"] = posts
	c.RenderArgs["pages"] = pages
	c.RenderArgs["currentPage"] = page
	return c.Render()
}
