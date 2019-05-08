package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"web/models"
)
//用户结构体对象
type UserController struct {
	beego.Controller
}
//显示登陆页面
func (this *UserController) ShowLogin() {
	username:=this.Ctx.GetCookie("username")
	if username!=""{
		this.Data["username"]=username
		this.Data["Checked"]="Checked"
	}else {
		this.Data["username"]=""
		this.Data["Checked"]=""
	}
	//渲染登陆页面
	this.TplName = "login.html"
}
func (this *UserController)HandleLogin()  {
	//获取数据
	userName:=this.GetString("userName")
	password:=this.GetString("password")
	//校验数据
	if userName==""||password==""{
		beego.Error("数据结构不完整")
		this.TplName="login.html"
		return
	}


	//处理数据
	o:=orm.NewOrm()
	var user models.User
	user.UserName=userName
	err:=o.Read(&user,"UserName")
	if err!=nil{
		beego.Error("用户名不存在",err)
		this.TplName="register.html"
		return
	}
	//用户输入密码校验
	if user.PassWord!=password{
		beego.Info("密码错误")
		this.TplName="login.html"
		return
	}
	remember:=this.GetString("remember")
	if remember=="on"{
		this.Ctx.SetCookie("username",userName)
	}else{
		this.Ctx.SetCookie("username","")
	}
	this.SetSession("username",userName)
	//返回数据
	this.Redirect("/article/index",302)
}
//显示注册页面
func (this *UserController) ShowRegister() {
	//渲染登陆页面
	this.TplName = "register.html"
}
//注册业务处理
func (this *UserController)HandleRegister()  {
	//获取数据
	userName:=this.GetString("userName")
	password:=this.GetString("password")
	//校验数据
	if userName==""||password==""{
		beego.Error("数据结构不完整")
		this.TplName="register.html"
		return
	}
	//处理数据
	o:=orm.NewOrm()
	var user models.User
	user.UserName=userName
	user.PassWord=password
	_,err:=o.Insert(&user)
	if err!=nil{
		beego.Error("数据插入失败",err)
		this.TplName="register.html"
		return
	}

	//返回数据
	this.TplName="login.html"
}
func (this *UserController)Logout(){
	this.DelSession("username")
	this.Redirect("/",302)
}