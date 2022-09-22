package gorm

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

type Product struct {
	gorm.Model
	Code  string `gorm:"column(code)"`
	Price uint
}

func (p *Product) TableName() string {
	return "product_t"
}

func (p *Product) BeforeSave(tx *gorm.DB) (err error) {
	println("before save")
	return
}

func (p *Product) AfterSave(tx *gorm.DB) (err error) {
	println("after save")
	return
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	println("before create")
	return
}

func (p *Product) AfterCreate(tx *gorm.DB) (err error) {
	println("after create")
	return
}

func (p *Product) BeforeUpdate(tx *gorm.DB) (err error) {
	println("before update")
	return
}

func (p *Product) AfterUpdate(tx *gorm.DB) (err error) {
	println("after update")
	return
}

func (p *Product) BeforeDelete(tx *gorm.DB) (err error) {
	println("before update")
	return
}

func (p *Product) AfterDelete(tx *gorm.DB) (err error) {
	println("after update")
	return
}

func (p *Product) AfterFind(tx *gorm.DB) (err error) {
	println("after find")
	return
}

func TestCrud(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		panic("fail to connect database")
	}
	_ = db.AutoMigrate(&Product{})

	//db.Create(&Product{Code: "D42", Price: 100})
	//db.Create(&Product{Code: "D43", Price: 100})
	//db.Create(&Product{Code: "D44", Price: 100})

	var product Product
	db.First(&product, 2)
	t.Log(product)
	p2 := new(Product)
	db.First(p2, "code=?", "D44")
	t.Log(*p2)
	p3 := Product{}
	db.First(&p3, "id=?", 3)
	t.Log(p3)

	db.Model(&product).Update("Price", 200)
	db.Model(&p3).Updates(Product{Price: 200, Code: "F42"})

	db.Delete(&p3, 3)
}
