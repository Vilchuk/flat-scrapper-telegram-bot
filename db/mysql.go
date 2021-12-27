package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func mySqlConnection(dbName string) *sql.DB {
	//conn := fmt.Sprintf("root:root@unix(/Applications/MAMP/tmp/mysql/mysql.sock)/%s", dbName)
	conn := fmt.Sprintf("root:@tcp(localhost:3306)/%s", dbName)
	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err)
	}
	return db
}

func isHashExists(db *sql.DB, siteName string, hash string) bool {
	var hashRow int
	db.QueryRow("SELECT Id FROM Flats WHERE SiteName = ? AND Hash = ?", siteName, hash).Scan(&hashRow)

	if hashRow != 0 {
		return true
	}

	return false
}

func insertHash(db *sql.DB, siteName string, flatUrl, hash string) {
	query := "INSERT INTO Flats (SiteName, FlatUrl, Hash) VALUES(?,?,?)"
	insert, err := db.Prepare(query)
	defer insert.Close()

	if err != nil {
		println(err)
	}

	_, err = insert.Exec(siteName, flatUrl, hash)
	if err != nil {
		println(err)
	}
}

func IsFlatNew(siteName string, flatUrl string, hash string) bool {
	database := mySqlConnection("go-flat-scrapper")

	isFlatHasAlreadyViewed := isHashExists(database, siteName, hash)
	if isFlatHasAlreadyViewed == false {
		insertHash(database, siteName, flatUrl, hash)
	}

	defer database.Close()
	return !isFlatHasAlreadyViewed
}
