package blog

import (
	"github.com/astaxie/beego"
	"blog-master/models"
	"strconv"
	"strings"
	"blog-master/controllers/ipfilter"
)

type baseController struct {
	beego.Controller
	options           map[string]string
	right             string
	page              int
	pagesize          int
	clientip          string
	allowconn         bool
	allowconnmsg      string
}

func (this *baseController) Prepare() {//这个函数是用于用户扩展用的，这个函数会在下面定义的这些Method方法之前执行，用户可以重写这个函数实现类似用户验证之类。
	this.clientip = this.getClientIp()
	this.allowconn = true
	this.allowconn, this.allowconnmsg = ipfilter.ConnFilterCtx().OnConnected(this.clientip)//对ip进行一定的过滤
	if !this.allowconn {
		//超过3次异常访问，返回500
		this.Controller.Abort("500")
	}
	this.Data["IsLogin"] = this.IsLogin()
	this.options = models.GetOptions()
	this.right = "right.html"
	this.Data["options"] = this.options
	this.Data["latestblog"] = models.GetLatestBlog()//最新文章
	this.Data["hotblog"] = models.GetHotBlog()//点击排行
	this.Data["newcomments"] = models.GetNewComments()//最新评论
	this.Data["links"] = models.GetLinks()//友情链接
	//下面的是二维码链接
	this.Data["hidejs"] = `<!--[if lt IE 9]>
	<script src="/static/js/html5shiv.min.js"></script>
	<![endif]-->`

	var (
		pagesize int
		err      error
		page     int
	)

	if page, err = strconv.Atoi(this.Ctx.Input.Param(":page")); err != nil || page < 1 {
		page = 1
	}

	if pagesize, err = strconv.Atoi(this.getOption("pagesize")); err != nil || pagesize < 1 {
		pagesize = 10
	}
	this.page = page
	this.pagesize = pagesize
}

func (this *baseController) display(tpl string) {
	theme := "double"
	if v, ok := this.options["theme"]; ok && v != "" {
		theme = v
	}

	this.Layout = theme + "/layout.html"
	this.Data["root"] = "/" + beego.BConfig.WebConfig.ViewsPath + "/" + theme + "/"
	this.TplName = theme + "/" + tpl + ".html"

	this.LayoutSections = make(map[string]string)
	this.LayoutSections["head"] = theme + "/head.html"

	if tpl == "index" {
		this.LayoutSections["banner"] = theme + "/banner.html"
	}
	if this.right != "" {
		this.LayoutSections["right"] = theme + "/" + this.right
	}
	this.LayoutSections["foot"] = theme + "/foot.html"
}

func (this *baseController) getOption(name string) string {
	if v, ok := this.options[name]; ok {
		return v
	} else {
		return ""
	}
}

func (this *baseController) setHeadMetas(params ...string) {
	title_buf := make([]string, 0, 3)
	if len(params) == 0 && this.getOption("sitename") != "" {
		title_buf = append(title_buf, this.getOption("sitename"))
	}
	if len(params) > 0 {
		title_buf = append(title_buf, params[0])
	}
	title_buf = append(title_buf, this.getOption("subtitle"))
	this.Data["title"] = strings.Join(title_buf, " - ")

	if len(params) > 1 {
		this.Data["keywords"] = params[1]
	} else {
		this.Data["keywords"] = this.getOption("keywords")
	}

	if len(params) > 2 {
		this.Data["description"] = params[2]
	} else {
		this.Data["description"] = this.getOption("description")
	}
}

func (this *baseController) IsLogin() bool {//这个函数用于验证登录
	arr := strings.Split(this.Ctx.GetCookie("auth"), "|")
	if len(arr) == 2 {
		idstr, password := arr[0], arr[1]
		userid, _ := strconv.ParseInt(idstr, 10, 0)
		if userid > 0 {
			var user models.User
			user.Id = userid
			if user.Read() == nil && password == models.Md5([]byte(this.getClientIp()+"|"+user.Password)) {
				return true
			}
		}
	}
	return false
}

//获取用户IP地址
func (this *baseController) getClientIp() string {
	s := this.Ctx.Request.Header.Get("X-Real-IP")
	if s == "" {
		s = strings.Split(this.Ctx.Request.RemoteAddr, ":")[0]
	}
	return s
}