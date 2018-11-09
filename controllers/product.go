package controllers

import (
	"net/http"

	"github.com/freitagsrunde/k4ever-backend/models"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/render"
)

type ProductHandler struct {
	PR *models.ProductResource
}

func (ph ProductHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", ph.List)
	r.Get("/{productID}", ph.Get)
	r.Post("/", ph.Create)

	return r
}

func (ph ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	products, err := ph.PR.ListProducts()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, &products)
}

func (ph ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	var productID string
	if productID = chi.URLParam(r, "productID"); productID == "" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	product, err := ph.PR.GetProduct(productID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, &product)
}

func (ph ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	product := &models.Product{}
	if err := render.Bind(r, product); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err := ph.PR.CreateProduct(product); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, &product)
}
