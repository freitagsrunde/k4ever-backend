package controllers

import (
	"net/http"

	"github.com/freitagsrunde/k4ever-backend/models"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/render"
	"github.com/jinzhu/gorm"
)

type UserHandler struct {
	DB *gorm.DB
	UR *models.UserResource
}

func (rs UserHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", rs.List)
	r.Get("/{userID}", rs.Get)

	return r
}

func (rs UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := rs.UR.ListUsers(rs.DB)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, &users)
}

func (rs UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	var userID string
	if userID = chi.URLParam(r, "userID"); userID == "" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	user, err := models.GetUser(userID, rs.DB)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, &user)
}
