package routers

import (
	"web/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	beego.InsertFilter("/article/*",beego.BeforeExec,Filter)
    //beego.Router("/", &controllers.MainController{})
    //用户登陆请求
	beego.Router("/", &controllers.UserController{},"get:ShowLogin;post:HandleLogin")
	//用户注册页面
	beego.Router("/register", &controllers.UserController{},"get:ShowRegister;post:HandleRegister")
	//文章显示页面
	beego.Router("/article/index",&controllers.ArticleController{},"get,post:ShowIndex")
	//文章添加操作
	beego.Router("/article/addarticle",&controllers.ArticleController{},"get:ShowAddArticle;post:HandleAddArticle")

	beego.Router("/article/articleContent",&controllers.ArticleController{},"get:ShowArticleContent")

	beego.Router("/article/deleteArticle",&controllers.ArticleController{},"get:DeleteArticle")

	beego.Router("/article/editArticle",&controllers.ArticleController{},"get:ShowEditArticle;post:HandleEditArticle")

	beego.Router("/article/addArticleType",&controllers.ArticleController{},"get:ShowArticleType;post:HandleArticleType")

	beego.Router("/article/DeleteType",&controllers.ArticleController{},"get:DeleteType")

	beego.Router("/article/logout",&controllers.UserController{},"get:Logout")

}
func Filter(ctx *context.Context)  {
	username:=ctx.Input.Session("username")
	if username==nil{
		ctx.Redirect(302,"/")
		return
	}
}