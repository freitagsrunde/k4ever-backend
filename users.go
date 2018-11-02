package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/render"
)

type User struct {
	Name string `json:"name"`
}

func (u *User) Bind(r *http.Request) error {
	return nil
}

func (u *User) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type usersResource struct{}

func (rs usersResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", rs.List)

	return r
}

func (rs usersResource) List(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, &User{Name: "name"})
}
