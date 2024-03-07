package libsql

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/tursodatabase/go-libsql"
)

type Connection struct {
	Dir string
	DB  *sql.DB
}

func Connect(name string) (conn *Connection, err error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}

	dir := cacheDir + "/" + name
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return nil, err
		}
	}

	db, err := sql.Open("libsql", fmt.Sprintf("file:%v/%v.db", dir, name))
	if err != nil {
		return nil, err
	}

	return &Connection{
		Dir: dir,
		DB:  db,
	}, nil
}

func (c *Connection) Disconnect() error {
	err := c.DB.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection) TableExists(name string) (bool, error) {
	query := fmt.Sprintf("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='%v'", name)
	rows, err := c.DB.Query(query)
	if err != nil {
		return false, err
	}

	var count int
	rows.Next()
	err = rows.Scan(&count)
	if err != nil {
		return false, err
	}

	if rows.Err() != nil {
		return false, rows.Err()
	}
	return count > 0, nil
}

func (c *Connection) Execute(command string) (sql.Result, error) {
	res, err := c.DB.Exec(command)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Connection) Query(query string) (*sql.Rows, error) {
	rows, err := c.DB.Query(query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
