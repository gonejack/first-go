package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Product struct {
	gorm.Model
	Code  string `gorm:"primary_key"`
	Price uint
}

func main() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	//db.AutoMigrate(&Product{})

	//// Create
	db.Debug().Create(&Product{Code: "L1212", Price: 1000})

	// Read
	var product Product
	//db.First(&product, 1) // find product with id 1
	//db.First(&product, "code = ?", "L1212") // find product with code l1212
	//
	//// Update - update product's price to 2000
	//db.Model(&product).Update("Price", 400)

	loop := func(db *gorm.DB) *gorm.DB {
		return db.Where("code in (?)", []string{"L1212", "L1213"})
	}
	result := db.Debug().Scopes(loop, loop).First(&product)

	db.Update()

	spew.Dump(result.Error)

	// Delete - delete product
	//db.Delete(&product)
}
