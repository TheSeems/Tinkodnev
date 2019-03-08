package main

import (
	"tinkodnev/api"
	"tinkodnev/engine"
	"tinkodnev/storages"
	"tinkodnev/utils"
)

func loadDb() {
	dbConfig := utils.MustParseJsonConfig("db.json")
	db := storages.MySQLMemDB{}
	db.Init(dbConfig["user"].(string) + ":" + dbConfig["pass"].(string) +
		"@(" + dbConfig["host"].(string) + ")/" + dbConfig["db"].(string))
	engine.Database = &db
}

func loadApi() {
	engine.Router.HandleFunc("/api/get", api.GetMemberMethod)
	engine.Router.HandleFunc("/api/search", api.SearchMemberMethod)
}

func main() {
	loadDb()
	engine.Load()

	loadApi()
	engine.Start()
}
