package main

import (
	"flag"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/render"

	"github.com/freitagsrunde/k4ever-backend/controllers"
	"github.com/freitagsrunde/k4ever-backend/db"
)

func main() {
	flag.Parse()

	r := chi.NewRouter()
	db := db.Init()

	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		render.Render(w, r, &Test{Name: "test"})
	})

	r.Mount("/users", controllers.UsersResource{DB: db}.Routes())

	http.ListenAndServe(":8080", r)
}

type Test struct {
	Name string `json:"name"`
}

func (t *Test) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
