package main

import (
	"fmt"
	"github.com/Unknwon/goconfig"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
	// "os"
	"strconv"
)

type App struct {
	Id   int `pk:"auto"`
	Name string
}

type Record struct {
	Id int
}

type DataH struct {
	Apps           []App
	AppIndex       int
	CurrentAppName string
	DateInfo       []string
	DateIndex      int
}

func (this *App) Get() (*[]App, error) {
	var apps []App
	_, err := orm.NewOrm().QueryTable("app").All(&apps)
	if err != nil {
		return nil, err
	}
	return &apps, nil
}

func init() {
	// 注册驱动
	// orm.RegisterDriver("mysql", orm.DR_MySQL)
	// 注册默认数据库
	// 我的mysql的root用户密码为tom，打算把数据表建立在名为test数据库里
	// 备注：此处第一个参数必须设置为“default”（因为我现在只有一个数据库），否则编译报错说：必须有一个注册DB的别名为 default
	// orm.RegisterDataBase("default", "mysql", "root:tom@/test?charset=utf8")
	orm.Debug = true
	orm.RegisterDriver("sqlite", orm.DRSqlite)
	orm.RegisterDataBase("default", "sqlite3", "data.db")
	// 需要在 init 中注册定义的 model
	orm.RegisterModel(new(App))
	orm.RunSyncdb("default", false, false)
}

// GOOS=linux GOARCH=amd64 go build -o coolgo_linux github.com/freelifer/coolgo/*.go
// GOOS=windows GOARCH=386 go build -o github.com/freelifer/log/log_win github.com/freelifer/log/*.go
func main() {
	cfg, err := goconfig.LoadConfigFile("config.ini")
	if err != nil {
		log.Println("读取配置文件失败[conf.ini]")
		return
	}
	str, _ := cfg.GetValue(goconfig.DEFAULT_SECTION, "default_key")
	fmt.Printf("%s\n", str)

	http.HandleFunc("/", Handler)
	http.ListenAndServe(":8080", nil)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	apps, e := new(App).Get()
	if e != nil {
		checkError(e)
		return
	}

	info := []string{"全部应用", "coolgo", "pythone", "java"}
	dateInfo := []string{"今天", "昨天", "7天", "本月", "上月"}

	//t := template.New("Person template")
	//t, err := t.Parse(templ)
	r.ParseForm()
	var appId, dateIndex int
	var err error
	if len(r.Form["a"]) > 0 {
		appId, err = strconv.Atoi(r.Form["a"][0])
		checkError(err)
	} else {
		appId = 0
	}

	var currentAppName string
	for _, app := range *apps {
		if app.Id == appId {
			currentAppName = app.Name
		}
	}
	if currentAppName == "" {
		currentAppName = "全部应用"
	}
	if len(r.Form["d"]) > 0 {
		dateIndex, err = strconv.Atoi(r.Form["d"][0])
		checkError(err)
	} else {
		dateIndex = 0
	}
	if dateIndex < 0 || dateIndex > len(info) {

	}
	t, err := template.ParseFiles("tmpl.html")
	checkError(err)
	dataH := DataH{Apps: *apps, AppIndex: appId, CurrentAppName: currentAppName, DateInfo: dateInfo, DateIndex: dateIndex}
	err = t.Execute(w, dataH)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		// os.Exit(1)
	}
}
