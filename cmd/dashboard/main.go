package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rafaelmdm/htmx-dashboard/pkg/libsql"
)

func init_database(conn *libsql.Connection) error {
	exists, err := conn.TableExists("test")
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	_, err = conn.Execute("CREATE TABLE test (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		return err
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

	err_ch := make(chan error, 1)
	go func() {
		err_ch <- init_database(conn)
	}()

	e.GET("/", homeHandler)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := e.Start(":3000"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Shutting down the server")
		}
	}()

	err = <-err_ch
	if err != nil {
		e.Logger.Fatal(err)
	}

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
