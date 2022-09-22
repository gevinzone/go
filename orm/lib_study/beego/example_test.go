package beego

import (
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

type User struct {
	ID   int    `orm:"column(id)"`
	Name string `orm:"column(name)"`
}

func init() {
	orm.RegisterModel(new(User))
	_ = orm.RegisterDriver("sqlite3", orm.DRSqlite)
	_ = orm.RegisterDataBase("default", "sqlite3", "beego.db")
}

func TestCrud(t *testing.T) {
	_ = orm.RunSyncdb("default", false, true)
	o := orm.NewOrm()

	user := &User{Name: "mike"}
	id, err := o.Insert(user)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id)
}
