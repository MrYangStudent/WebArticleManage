package controllers

import (
	"github.com/astaxie/beego"
	"path"
	"time"
	"github.com/astaxie/beego/orm"
	"web/models"
	"math"
)

type ArticleController struct {
	beego.Controller
}

func (this *ArticleController) ShowIndex() {
	username:=this.GetSession("username")
	//获取数据
	pageIndex,err:=this.GetInt("pageIndex")
	if err!=nil{
		pageIndex=1
	}
	TypeName:=this.GetString("select")
	beego.Info(TypeName)
	o:=orm.NewOrm()
	var articles []models.Article
	var articleTypes []models.ArticleType
	qb:=o.QueryTable("Article")
	o.QueryTable("ArticleType").All(&articleTypes)
	var totalCount int64
	//数据处理
	//.RelatedSel("ArticleType")
	if TypeName==""{
		totalCount,_=qb.RelatedSel("ArticleType").Count()
	}else {
		totalCount,_=qb.RelatedSel("ArticleType").Filter("ArticleType__TypeName",TypeName).Count()
	}
	//每页显示条数
	pageSize:=2
	pagecount:=int(math.Ceil(float64(totalCount)/float64(pageSize)))
	//查询类型列表
	if TypeName==""{
		//查找所有文章
		//.RelatedSel("ArticleType")
		qb.Limit(pageSize,pageSize*(pageIndex-1)).RelatedSel("ArticleType").All(&articles)
	}else {
		qb.Limit(pageSize,pageSize*(pageIndex-1)).RelatedSel("ArticleType").Filter("ArticleType__TypeName",TypeName).All(&articles)

	}

	//渲染登陆页面
	//当前页数
	this.Data["pageIndex"]=pageIndex
	this.Data["username"]=username
	//总页数
	this.Data["pagecount"]=pagecount
	this.Data["totalCount"]=totalCount
	this.Data["articleTypes"]=articleTypes
	this.Data["TypeName"]=TypeName
	this.Data["articles"]=articles
	this.LayoutSections=make(map[string]string)
	this.LayoutSections["indexjs"]="indexjs.html"
	this.Layout="layout.html"
	this.TplName = "index.html"
}
func (this *ArticleController) ShowAddArticle() {
	username:=this.GetSession("username")
	this.Data["username"]=username
	o:=orm.NewOrm()
	var articleTypes []models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)
	TypeName:=this.GetString("select")
	this.Data["TypeName"]=TypeName
	this.Data["articleTypes"]=articleTypes
	//渲染登陆页面
	this.Layout="layout.html"
	this.TplName = "add.html"
}
func UpLoadFile(this *ArticleController,filedName string,url string)string  {
	//上传文件处理
	fil,head,err:=this.GetFile(filedName)
	defer fil.Close()
	if err!=nil{
		beego.Error("文件上传失败",err)
		this.Redirect(url,302)
		return""
	}
	//获取文件后缀
	ext:=path.Ext(head.Filename)
	//防止重名,使用时间作文件名称
	fileName:=time.Now().Format("200601021504052222")+ext
	beego.Info(fileName)

	//校验文件格式
	if ext!=".jpg"&&ext!=".png"&&ext!=".jpeg"{
		beego.Error("文件上传格式错误")
		this.Redirect(url,302)
		return ""
	}

	//文件大小校验
	if head.Size>5000000{
		beego.Error("上传文件过大")
		this.Redirect(url,302)
		return ""
	}
	//文件保存路径
	filePath:="/static/img/"+fileName
	beego.Info(filePath)
	this.SaveToFile(filedName,"./static/img/"+fileName)
	return filePath
}
func (this *ArticleController) HandleAddArticle() {
	//获取数据

	articleName:=this.GetString("articleName")
	articlecontent:=this.GetString("content")
	//调用函数完成上传文件校验
	fileName:=UpLoadFile(this,"uploadname","/article/addarticle")
	// 校验数据
	if articleName==""||articlecontent==""{
		beego.Error("数据结构不完整")
		this.Redirect("/article/addarticle",302)
		return
	}
	// 处理数据
	o:=orm.NewOrm()
	typeName:=this.GetString("select")
	var articleType models.ArticleType

	articleType.TypeName=typeName
	o.Read(&articleType,"TypeName")
	var article models.Article
	article.Acontent=articlecontent
	article.Aimg=fileName
	article.Atitle=articleName
	article.ArticleType=&articleType

	_,err:=o.Insert(&article)
	if err!=nil{
		beego.Error("文章插入失败",err)
		this.Redirect("/article/addarticle",302)
		return
	}
	// 返回数据
	this.Redirect("/article/index",302)
}
func (this *ArticleController) ShowArticleContent(){
	username:=this.GetSession("username")
	//获取数据
	id,err:=this.GetInt("id")
	// 校验数据
	if err!=nil{
		beego.Error("数据获取失败")
		this.Redirect("/article/index",302)
		return
	}
	// 处理数据
	o:=orm.NewOrm()
	var article models.Article
	//var Article models.Article
	article.Id=id
	o.Read(&article)
	article.Acount+=1
	o.Update(&article)
	o.QueryTable("Article").RelatedSel("ArticleType").Filter("Id",id).One(&article)
	m2m:=o.QueryM2M(&article,"Users")
	beego.Info(m2m)
	user:=models.User{UserName:username.(string)}
	o.Read(&user,"UserName")
	beego.Info(user)
	m2m.Add(user)
	beego.Info(article)
	var users []models.User
	o.QueryTable("User").Filter("Articles__Article__Id",id).Distinct().All(&users)

	//渲染登陆页面
	//this.Data["Article"]=Article
	this.Data["users"]=users
	this.Data["article"]=article
	this.Data["username"]=username
	this.Layout="layout.html"
	this.TplName = "content.html"
}
func (this *ArticleController) DeleteArticle(){
	//获取数据
	id,err:=this.GetInt("id")
	// 校验数据
	if err!=nil{
		beego.Error("数据获取失败")
		this.Redirect("/article/index",302)
		return
	}
	// 处理数据
	o:=orm.NewOrm()
	var article models.Article
	article.Id=id
	_,err=o.Delete(&article)
	if err!=nil{
		beego.Error("数据删除失败",err)
		this.Redirect("/article/index",302)
		return
	}
	// 返回数据
	this.Redirect("/article/index",302)
}
func (this *ArticleController) ShowEditArticle(){
	//获取数据
	id,err:=this.GetInt("id")
	// 校验数据
	if err!=nil{
		beego.Error("数据获取失败")
		this.Redirect("/article/index",302)
		return
	}
	// 处理数据
	o:=orm.NewOrm()
	var article models.Article
	article.Id=id
	o.Read(&article)
	//渲染登陆页面
	username:=this.GetSession("username")
	this.Data["username"]=username
	this.Data["article"]=article
	this.Layout="layout.html"
	this.TplName = "update.html"
}
func (this *ArticleController) HandleEditArticle(){
	//获取数据
	id,err:=this.GetInt("id")
	// 校验数据
	if err!=nil{
		beego.Error("数据获取失败",err)
		this.Redirect("/article/editArticle",302)
		return
	}
	articleName:=this.GetString("articleName")
	articlecontent:=this.GetString("content")
	//调用函数完成上传文件校验
	fileName:=UpLoadFile(this,"uploadname","/article/editArticle")
	// 校验数据
	if articleName==""||articlecontent==""{

		beego.Error("数据结构不完整")
		this.Redirect("/article/editArticle",302)
		return
	}
	// 处理数据
	o:=orm.NewOrm()
	var article models.Article
	article.Id=id
	err=o.Read(&article)
	if err!=nil{
		beego.Error("数据不存在",err)
		this.Redirect("/article/addarticle",302)
		return
	}
	article.Acontent=articlecontent
	article.Aimg=fileName
	article.Atitle=articleName
	_,err=o.Update(&article)
	if err!=nil{
		beego.Error("数据更新失败",err)
		this.Redirect("/article/editArticle",302)
		return
	}
	this.Redirect("/article/index",302)
}
func (this *ArticleController) ShowArticleType(){
	// 数据
	o:=orm.NewOrm()
	var articleType []models.ArticleType
	o.QueryTable("ArticleType").All(&articleType)
	//渲染登陆页面
	this.Data["articleType"]=articleType
	username:=this.GetSession("username")
	this.Data["username"]=username
	this.LayoutSections=make(map[string]string)
	this.LayoutSections["addTypejs"]="addTypejs.html"
	this.Layout="layout.html"
	this.TplName = "addType.html"
}
func (this *ArticleController) HandleArticleType(){
	//获取数据
	typeName:=this.GetString("typeName")
	//校验数据
	if typeName==""{
		beego.Error("数据不完整")
		this.Redirect("/article/addArticleType",302)
		return
	}
	// 处理数据
	o:=orm.NewOrm()
	var articleType models.ArticleType
	articleType.TypeName=typeName
	_,err:=o.Insert(&articleType)
	if err!=nil{
		beego.Error("插入数据错误",err)
		this.Redirect("/article/addArticleType",302)
		return
	}
	// 返回数据
	this.Redirect("/article/addArticleType",302)
}
func (this *ArticleController) DeleteType(){
	//获取数据
	id,err:=this.GetInt("id")
	// 校验数据
	if err!=nil{
		beego.Error("数据获取失败")
		this.Redirect("/article/index",302)
		return
	}
	// 处理数据
	o:=orm.NewOrm()
	var articleType models.ArticleType
	articleType.Id=id
	o.Delete(&articleType)
	this.Redirect("/article/addArticleType",302)
}