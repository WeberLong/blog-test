package controllers

import (
	"strings"
	"fmt"
	"myblog/models/class"
	"strconv"
)
import . "fmt"

type ArticleController struct {
	BaseController
	ret RET
}

func (c *ArticleController) Archive() {
	
	errmsg := ""

	a := class.Article{}
	if len(c.GetString("tag")) > 0 {
		tag := class.Tag{Name: c.GetString("tag")}.Get()
		if tag == nil {
			errmsg += fmt.Sprintf("Tag[%s] is not exist.\n", c.GetString("tag"))
		} else {
			a.Tags = []*class.Tag{tag}
		}
	}

	if len(c.GetString("author")) > 0 {
		author := class.User{Id: c.GetString("author")}.Get()
		if author == nil {
			errmsg += fmt.Sprintf("User[%s] is not exist.\n", c.GetString("author"))
		} else {
			a.Author = author
		}
	}

	if len(errmsg) == 0 {
		rets := a.Gets()
		c.Data["articles"] = rets
	}

	c.Data["err"] = errmsg

	c.TplName = "article/archive.html"

}
	

func (c *ArticleController) PageNew() {
	c.CheckLogin()
	c.TplName = "article/new.html"
}

func (c *ArticleController) Get() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	a := &class.Article{Id: id}
	a.ReadDB()
	a.Author.ReadDB()
	r := &class.Reply{Article: a}
	a.Replys = r.Gets()

	c.Data["article"] = a
	c.Data["replyTree"] = a.GetReplyTree()

	c.TplName = "article/article.html"
}

func (c *ArticleController) PageEdit() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	a := &class.Article{Id: id}
	a.ReadDB()
	a.Author.ReadDB()
	c.Data["article"] = a
	c.TplName = "article/edit.html"
}

func (c *ArticleController) Edit() {
	c.CheckLogin()
	u := c.GetSession("user").(class.User)

	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	a := &class.Article{Id: id}
	a.ReadDB()

	if u.Id != a.Author.Id {
		c.DoLogout()
	}

	strs := strings.Split(c.GetString("tag"), ",")
	tags := []*class.Tag{}
	for _, v := range strs {
		tags = append(tags, class.Tag{Name: strings.TrimSpace(v)}.GetOrNew())
	}

	a.Title = c.GetString("title")
	a.Content = c.GetString("content")
	a.Tags = tags

	a.Update()

	c.ret.Ok = true
	c.Data["json"] = c.ret
	c.ServeJSON()

}

func (c *ArticleController) New() {
	c.CheckLogin()

	u := c.GetSession("user").(class.User)

	defer func()  {
        c.Data["json"] = c.ret
        c.ServeJSON()
    }()

	a := &class.Article{
		Title:   c.GetString("title"),
		Content: c.GetString("content"),
		Author:  &u,
	}

	fmt.Println("title:" + a.Title)
	fmt.Println("content:" + a.Content)

	if len(a.Title) < 1 {
		c.ret.Ok = false
		c.ret.Content = "标题不能为空"
		return
	}

	if len(a.Content) < 1 {
		c.ret.Ok = false
		c.ret.Content = "内容不能为空"
		return
	}

	n, err := a.Create()

	if err == nil {
		c.ret.Ok = true
		c.ret.Content = n
		c.Data["json"] = c.ret
		c.ServeJSON()
		return
	}

	c.ret.Content = Sprint(err)
}

func (c *ArticleController) Del() {
	c.CheckLogin()
	u := c.GetSession("user").(class.User)

	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	a := &class.Article{Id: id}
	a.ReadDB()

	if u.Id != a.Author.Id {
		c.DoLogout()
	}

	a.Defunct = true
	a.Update()

	c.Redirect("/user/"+a.Author.Id, 302)
}
