package storages

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"strconv"
	"tinkodnev/engine"
)

type MySQLMemDB struct {
	Connection sqlx.DB
}

func (db *MySQLMemDB) Search(data string, limit int) ([]engine.Member, error) {
	var items []engine.Member
	err := db.Connection.Select(&items, "SELECT * FROM `Members` WHERE SecondName LIKE '%"+data+"%' OR FirstName LIKE '%"+data+"%' LIMIT "+strconv.Itoa(limit))
	return items, err
}

func (db *MySQLMemDB) Get(id uint64) (engine.Member, error) {
	var item engine.Member
	err := db.Connection.Get(&item, "SELECT * FROM `Members` WHERE Id = ?", id)
	return item, err
}

func (db *MySQLMemDB) Add(member engine.Member) (uint64, error) {
	res, err := db.Connection.Exec("INSERT INTO `Members` (Id, FirstName, SecondName, Status) VALUES (DEFAULT, ?, ?, ?)", member.FirstName, member.SecondName, member.Status)
	lastId, _ := res.LastInsertId()
	if err != nil {
		return 0, err
	} else {
		return uint64(lastId), nil
	}
}

func (db *MySQLMemDB) Init(data string) {
	fmt.Println("[MemDB/MySQL] Connecting to database")
	var (
		err error
		con *sqlx.DB
	)

	con, err = sqlx.Connect("mysql", data)
	if err == nil {
		db.Connection = *con
		fmt.Println("[MemDB/MySQL] Successfully connected to database")
		db.Connection.MapperFunc(func(s string) string { return s })
	} else {
		panic("[MemDB/MySQL] Error connecting to database: " + err.Error())
	}

	db.Connection.SetMaxIdleConns(0)
	_, _ = db.Connection.Exec(
		"CREATE TABLE IF NOT EXISTS `Members` (Id BIGINT PRIMARY KEY AUTO_INCREMENT, FirstName VARCHAR(100) NOT NULL, SecondName VARCHAR(500), Status INTEGER(10))")
	if err != nil {
		fmt.Println("[MemDB/MySQL] WARNING: Cannot execute init statement")
	}
}
