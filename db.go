package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	conn *sql.DB
}

func NewDB(path string) (*DB, error) {
	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	_, err = conn.Exec("CREATE TABLE IF NOT EXISTS passwords (id INTEGER PRIMARY KEY, name TEXT NOT NULL, url TEXT, username TEXT NOT NULL, password TEXT NOT NULL)")
	if err != nil {
		return nil, err
	}
	return &DB{conn: conn}, nil
}

func (d *DB) Add(up *UserPass) error {
	_, err := d.conn.Exec(`INSERT INTO passwords (name,url,username,password) VALUES (?,?,?,?)`, up.Name, up.URL, up.Username, up.Password)
	return err
}
func (d *DB) Get(name string) (*UserPass, error) {
	row := d.conn.QueryRow("SELECT name,url,username,password FROM passwords WHERE name=?", name)
	if row.Err() != nil {
		return nil, row.Err()
	}
	var up UserPass
	err := row.Scan(&up.Name, &up.URL, &up.Username, &up.Password)
	if err != nil {
		return nil, err
	}
	return &up, nil
}
func (d *DB) Delete(name string) error {
	_, err := d.conn.Exec("DELETE FROM passwords WHERE name=?", name)
	return err
}
func (d *DB) Exists(name string) (bool, error) {
	row := d.conn.QueryRow("SELECT EXISTS(SELECT id FROM passwords WHERE name=?)", name)
	if row.Err() != nil {
		return false, row.Err()
	}
	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil

}
