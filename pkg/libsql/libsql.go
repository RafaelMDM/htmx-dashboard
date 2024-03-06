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
	dir, err := os.MkdirTemp("", "libsql-*")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			os.RemoveAll(dir)
		}
	}()

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
	err := os.RemoveAll(c.Dir)
	if err != nil {
		return err
	}
	err = c.DB.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection) Execute(query string) (*sql.Result, error) {
	res, err := c.DB.Exec(query)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
