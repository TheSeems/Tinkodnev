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
	query := "SELECT * FROM `Members` WHERE SecondName LIKE '%" + data + "%' OR FirstName LIKE '%" + data + "%' OR Patronymic LIKE '%" + data + "%'" +
		"OR concat(FirstName, ' ', SecondName) LIKE '%" + data + "%' OR concat(SecondName, ' ', FirstName) LIKE '%" + data + "%' OR concat(FirstName, ' ', Patronymic) LIKE '%" + data + "%'" +
		"OR concat(SecondName, ' ', FirstName, ' ', Patronymic) LIKE '%" + data + "%' " +
		"OR concat(FirstName, ' ', Patronymic, ' ', SecondName) LIKE '%" + data + "%' LIMIT " + strconv.Itoa(limit)
	err := db.Connection.Select(&items, query)
	return items, err
}

func (db *MySQLMemDB) Get(id uint64) (engine.Member, error) {
	var item engine.Member
	err := db.Connection.Get(&item, "SELECT * FROM `Members` WHERE Id = ?", id)
	return item, err
}

func (db *MySQLMemDB) Add(member engine.Member) (uint64, error) {
	res, err := db.Connection.Exec("INSERT INTO `Members` (Id, FirstName, SecondName, Patronymic, Photo, Status) VALUES (DEFAULT, ?, ?, ?, ?, ?)",
		member.FirstName, member.SecondName, member.Patronymic, member.Photo, member.Status)
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
		"CREATE TABLE IF NOT EXISTS `Members` (Id BIGINT PRIMARY KEY AUTO_INCREMENT, FirstName VARCHAR(16) NOT NULL, SecondName VARCHAR(16) NOT NULL, Patronymic VARCHAR(16), Photo VARCHAR(100), Status INTEGER(1))")
	if err != nil {
		fmt.Println("[MemDB/MySQL] WARNING: Cannot execute init statement")
	}
}
