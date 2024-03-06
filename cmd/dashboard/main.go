package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rafaelmdm/htmx-dashboard/pkg/libsql"
)

func init_database(conn *libsql.Connection) error {
	_, err := conn.Execute("CREATE TABLE test (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT)")
	if err != nil {
		return err
	}

	for i := 0; i < 10; i++ {
		_, err = conn.Execute(fmt.Sprintf("INSERT INTO test (id, name) VALUES (%d, 'test-%d')", i, i))
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/assets", "assets")

	conn, err := libsql.Connect("test")
	if err != nil {
		e.Logger.Fatal(err)
	}

	err = init_database(conn)
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.GET("/", homeHandler)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := e.Start(":3000"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Shutting down the server")
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := conn.Disconnect(); err != nil {
		e.Logger.Fatal(err)
	}
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
