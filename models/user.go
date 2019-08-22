package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

//用户表模型
type User struct {
	//设置参数使用分号;分隔，设置的值如果是多个的话，使用,分隔。
	//当Filed的值为数字整形数字类型的时候，并且字段名为Id的话，默认作为主键使用。
	//如果是pk的话，就是设置为主键
	//null,数据库默认数据不能为空,如果使用null的话，代表允许为空
	//pk即是主键
	//index为摸个字段增加索引
	//unique为某个字段增加索引
	//column为字段设置名称
	//size类型默认为varchar(255)
	//auto_now每次model保存时都会对时间自动更新,auto_now_add第一次保存时才会设置时间
	//type
	//default为字段设置默认值
	Id         int64
	Username   string        `json:"username" orm:"unique;size(15);index"`
	Password   string        `orm:"size(32)"`
	Nickname   string        `json:"nickname" orm:"size(15);index"`
	Email      string        `json:"email" orm:"size(50);index"`
	Lastlogin  time.Time     `json:"lastlogin" orm:"auto_now_add;type(datetime);index"`
	Logincount int64         `orm:"index"`
	Lastip     string        `json:"lastip" orm:"size(32);index"`
	Authkey    string        `orm:"size(10)"`
	Active     int8
	Permission string        `orm:"size(100);index"`
	Avator     string        `json:"avator" orm:"size(150);default(/static/upload/default/user-default-60x60.png)"`
	Upcount    int64
	Post       []*Post       `orm:"reverse(many)"`
	Comments   []*Comments   `orm:"reverse(many)"`
}

func (m *User) TableName() string {
	return TableName("user")
}

func (m *User) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *User) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *User) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *User) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *User) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}
