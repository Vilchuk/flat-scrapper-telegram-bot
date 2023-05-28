package db

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"main.go/models"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type Database struct {
	connection *sql.DB
}

func NewDatabase(connectionString string) (*Database, error) {
	connection, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка подключения к базе данных")
	}

	err = createDatabaseIfNotExists(connection, "PolandFlatsScrapper")
	if err != nil {
		return nil, errors.Wrap(err, "ошибка создания базы данных, если она не существует")
	}

	err = createTableIfNotExists(connection)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка создания таблицы, если она не существует")
	}

	fmt.Println("Подключение к базе данных установлено")
	return &Database{connection: connection}, nil
}

func createDatabaseIfNotExists(db *sql.DB, dbName string) error {
	query := fmt.Sprintf("CREATE DATABASE %s", dbName)
	_, err := db.Exec(query)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code.Name() == "duplicate_database" {
			return nil // База данных уже существует, продолжаем выполнение
		}
		return errors.Wrap(err, "ошибка создания базы данных, если она не существует")
	}

	return nil
}

func createTableIfNotExists(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS Flats (
			Id SERIAL PRIMARY KEY,
			FoundOnSite TEXT,
			DateCreated DATE,
			Title TEXT,
			Href TEXT
		);
	`
	_, err := db.Exec(query)
	if err != nil {
		return errors.Wrap(err, "ошибка создания таблицы Flats, если она не существует")
	}

	return nil
}

func (db *Database) isFlatExists(flat models.Flat) (bool, error) {
	var id int
	query := `SELECT Id FROM Flats WHERE FoundOnSite = $1 AND Href = $2`
	err := db.connection.QueryRow(query, flat.FoundOnSite, flat.Href).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return false, errors.Wrap(err, "ошибка проверки наличия объявления")
	}

	return id != 0, nil
}

func (db *Database) insertFlatIntoDb(flat models.Flat) error {
	query := "INSERT INTO Flats (FoundOnSite, DateCreated, Title, Href) VALUES ($1, $2, $3, $4)"
	_, err := db.connection.Exec(query, flat.FoundOnSite, flat.DateCreated, flat.Title, flat.Href)
	if err != nil {
		return errors.Wrap(err, "ошибка вставки объявления в базу данных")
	}

	return nil
}

func (db *Database) IsFlatNew(flat models.Flat) (bool, error) {
	isFlatAlreadyExists, err := db.isFlatExists(flat)
	if err != nil {
		return false, errors.Wrap(err, "ошибка проверки, является ли объявление новым")
	}

	if !isFlatAlreadyExists {
		err := db.insertFlatIntoDb(flat)
		if err != nil {
			return false, errors.Wrap(err, "ошибка вставки нового объявления в базу данных")
		}
	}

	return !isFlatAlreadyExists, nil
}

func (db *Database) Close() error {
	return db.connection.Close()
}
