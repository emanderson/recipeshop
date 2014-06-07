package main

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type Ingredient struct {
	Id int64
	Name string
}

type Recipe struct {
	Id int64
	Title string
	Source string
	Steps string
}

type RecipeIngredient struct {
	RecipeId int64
	IngredientId int64
	Amount string
}

func dbmap(dbPath string) *gorp.DbMap {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalln("sql.Open failed: ", err)
	}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	return dbmap
}

func initDb(dbPath string) {
	dbmap := dbmap(dbPath)
	defer dbmap.Db.Close()

	dbmap.AddTable(Ingredient{}).SetKeys(true, "Id")
	dbmap.AddTable(Recipe{}).SetKeys(true, "Id")
	dbmap.AddTable(RecipeIngredient{}).SetUniqueTogether("RecipeId", "IngredientId")
	err := dbmap.CreateTablesIfNotExists()
	if err != nil {
		log.Fatalln("dbmap.CreateTablesIfNotExists failed: ", err)
	}
}


