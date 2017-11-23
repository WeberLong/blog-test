package main

import (
	"encoding/gob"
	"github.com/astaxie/beego"
	_ "myblog/routers"
	_ "myblog/models"
	"strings"
	"myblog/models/class"
)

func init() {
	
	gob.Register(class.User{})

	beego.AddFuncMap("split", SplitHobby)

}

func main() {

	beego.Run()
}

/*	Template Function	*/

func SplitHobby(s string, sep string) []string {
	return strings.Split(s, sep)
}