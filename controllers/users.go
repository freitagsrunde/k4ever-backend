package controllers

import (
	"net/http"

	"github.com/freitagsrunde/k4ever-backend/models"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/render"
	"github.com/jinzhu/gorm"
)

type UsersResource struct {
	DB *gorm.DB
}

func (rs UsersResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", rs.List)
	r.Get("/{userID}", rs.Get)

	return r
}

func (rs UsersResource) List(w http.ResponseWriter, r *http.Request) {
	users, err := models.ListUsers(rs.DB)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, &users)
}

func (rs UsersResource) Get(w http.ResponseWriter, r *http.Request) {
	var articleID string
	if articleID = chi.URLParam(r, "articleID"); articleID == "" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	user, err := models.GetUser(articleID, rs.DB)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, &user)
}
