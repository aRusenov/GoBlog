package controllers

import (
	"github.com/revel/revel"
)

func init() {
	revel.OnAppStart(InitDB);
	revel.InterceptMethod((*GormController).Begin, revel.BEFORE)
	revel.InterceptMethod((*GormController).Commit, revel.AFTER)
	revel.InterceptMethod((*GormController).Rollback, revel.FINALLY)

	revel.InterceptFunc(AuthFilter, revel.BEFORE, &App{})
	revel.InterceptFunc(AuthFilter, revel.BEFORE, &Post{})
	revel.InterceptFunc(AuthFilter, revel.BEFORE, &User{})
	revel.InterceptFunc(AuthFilter, revel.BEFORE, &Comment{})
}
