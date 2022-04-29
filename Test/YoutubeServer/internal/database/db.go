package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"sync"
)

//Для ускорения поиска по юрл во время выполнения запроса создаем индекс по этому столбцу
const schemaSQL = `
 CREATE TABLE IF NOT EXISTS previews (
	 url VARCHAR(32) UNIQUE ON CONFLICT IGNORE,
	file BLOB
 );
CREATE INDEX IF NOT EXISTS url_exists ON previews(url);
`

const insertSQL = `
INSERT INTO previews (
 url, file
) VALUES (
  ?, ?
)`

const selectSQL = `
SELECT * FROM previews 
WHERE url = ?
`

type DB struct {
	sql  *sql.DB
	stmt *sql.Stmt
}

//Инициализация базы данных
func NewDB(dbFile string) (*DB, error) {
	sqlDB, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	if _, err = sqlDB.Exec(schemaSQL); err != nil {
		return nil, err
	}

	stmt, err := sqlDB.Prepare(insertSQL)
	if err != nil {
		return nil, err
	}

	db := DB{
		sql:  sqlDB,
		stmt: stmt,
	}
	return &db, nil
}

type VideoPreview struct {
	Url  string
	File []byte
}

//Добавление новых элементов в базу данных
func (db *DB) Add(preview VideoPreview, m *sync.Mutex) error {
	m.Lock()
	defer m.Unlock()
	tx, err := db.sql.Begin()
	if err != nil {
		return err
	}

	if _, err = tx.Stmt(db.stmt).Exec(preview.Url, preview.File); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

//Поиск по базе данных чтобы найти нужный нам юрл и отдать превью из бд
func (db *DB) Get(url string) (out VideoPreview, err error) {
	tx, err := db.sql.Begin()
	if err != nil {
		return VideoPreview{}, err
	}

	result, err := tx.Query(selectSQL, url)
	if err != nil {
		tx.Rollback()
		return VideoPreview{}, err
	}
	var (
		urlsql  string
		preview []byte
	)
	for result.Next() {
		if err := result.Scan(&urlsql, &preview); err != nil {
			return VideoPreview{}, err
		}
	}
	return VideoPreview{Url: urlsql, File: preview}, tx.Commit()
}

//Завершение работы с бд
func (db *DB) Close() error {
	defer func() {
		db.stmt.Close()
		db.sql.Close()
	}()
	return nil
}
