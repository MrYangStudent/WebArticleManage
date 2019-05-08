package main

import (
	_ "web/routers"
	"github.com/astaxie/beego"
	_"web/models"
)

func main() {
	beego.AddFuncMap("Pre",PrePage)
	beego.AddFuncMap("Next",NextPage)
	beego.AddFuncMap("Addval",Add)
	beego.Run()
}
func PrePage(a int) int {
	if a==1{
		return a
	}
	return a-1
}
func NextPage(total,b int) int {
	if b==total{
		return b
	}
	return b+1
}
func Add(b int) int {
	return b+1
}