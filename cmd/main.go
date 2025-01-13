package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplates() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type Count struct {
	Count int
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Renderer = NewTemplates()

	count := Count{Count: 0}

	e.Static("/css", "css")

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "login", count)
	})

	e.POST("/login", func(c echo.Context) error {
		username := c.FormValue("username")
		password := c.FormValue("password")

		// Simple credential check - in a real app, you'd want to check against a database
		if username == "admin" && password == "admin" {
			// Use HX-Redirect header for HTMX to handle the redirect
			c.Response().Header().Set("HX-Redirect", "/dashboard")
			return c.NoContent(200)
		}

		// If credentials are wrong, redirect back to login
		c.Response().Header().Set("HX-Redirect", "/")
		return c.NoContent(200)
	})

	e.GET("/dashboard", func(c echo.Context) error {
		return c.Render(200, "dashboard", count)
	})

	e.POST("/count", func(c echo.Context) error {
		count.Count++
		return c.Render(200, "count", count)
	})

	e.POST("/upload", func(c echo.Context) error {
		provider := c.FormValue("provider")
		file, err := c.FormFile("file-upload")
		if err != nil {
			return c.HTML(400, "<pre>Error getting file: "+err.Error()+"</pre>")
		}

		src, err := file.Open()
		if err != nil {
			return c.HTML(400, "<pre>Error opening file: "+err.Error()+"</pre>")
		}
		defer src.Close()

		content, err := io.ReadAll(src)
		if err != nil {
			return c.HTML(400, "<pre>Error reading file: "+err.Error()+"</pre>")
		}

		return c.HTML(200, "<pre>Uploaded to "+provider+"\n\n"+string(content)+"</pre>")
	})

	e.Logger.Fatal(e.Start(":42069"))

}
