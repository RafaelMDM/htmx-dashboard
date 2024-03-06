package main

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/rafaelmdm/htmx-dashboard/views"
)

func render(ctx echo.Context, component templ.Component) error {
	return component.Render(ctx.Request().Context(), ctx.Response())
}

func homeHandler(c echo.Context) error {
	return render(c, views.Home("Home"))
}
