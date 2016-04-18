package main

import (
	"fmt"
	"html/template"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"

	"labix.org/v2/mgo"
)

func unescape(x string) interface{} {
	return template.HTML(x)
}

func main() {
	fmt.Println("Listening on port :3000")

	mongoSession, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	db := mongoSession.DB("blog")

	m := martini.Classic()

	unescapeFuncMap := template.FuncMap{"unescape": unescape}

	m.Map(db)

	m.Use(session.Middleware)

	m.Use(render.Renderer(render.Options{
		Directory:  "templates",                         // Specify what path to load the templates from.
		Layout:     "layout",                            // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".tmpl", ".html"},          // Specify extensions to load for templates.
		Funcs:      []template.FuncMap{unescapeFuncMap}, // Specify helper function maps for templates to access.
		Charset:    "UTF-8",                             // Sets encoding for json and html content-types. Default is "UTF-8".
		IndentJSON: true,                                // Output human readable JSON
	}))

	staticOptions := martini.StaticOptions{Prefix: "assets"}
	m.Use(martini.Static("assets", staticOptions))

	m.Get("/", routes.IndexHandler)
	m.Get("/user", routes.GetLoginHandler)
	m.Get("/new", routes.WriteHandler)
	m.Get("/edit/:id", routes.EditHandler)
	m.Get("/delete/:id", routes.DeleteHandler)
	m.Post("/post", routes.SavePostHandler)

	m.Run()
}
