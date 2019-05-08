package models

import (
	"time"
	"github.com/astaxie/beego/orm"
	_"github.com/go-sql-driver/mysql"
)
//创建用户结构体对象
type User struct {
	Id int `orm:"pk;auto"`
	UserName string `orm:"unique;size(20)"`
	PassWord string`orm:"size(20)"`
	Articles []*Article`orm:"reverse(many)"`
}
//创建文章对象
type Article struct {
	Id int `orm:"pk;auto"`
	Atitle string `orm:"size(20)"`
	Aimg   string `orm:"null"`
	Atime time.Time`orm:"type(datatime);auto_now_add"`
	Acount int `orm:"default(0)"`
	Acontent string`orm:"size(500)"`
	ArticleType *ArticleType`orm:"rel(fk);null;on_delete(set_null)"`
	Users []*User`orm:"rel(m2m)"`
}
//创建文章类型对象
type ArticleType struct {
	Id int`orm:"pk;auto"`
	TypeName string`orm:"size(20);unique"`
	Articles []*Article`orm:"reverse(many)"`
}
func init()  {
	//连接数据库
	orm.RegisterDataBase("default","mysql","root:123456@tcp(127.0.0.1:3306)/webobject")
	//获取结构体对象，注册
	orm.RegisterModel(new(User),new(ArticleType),new(Article))
	//生成表格
	orm.RunSyncdb("default",false,true)
}
