package controllers

import (
	"net/http"
	"strconv"

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
	r.Post("/", ph.Create)
	r.Get("/{productID}", ph.Get)
	r.Put("/{productID}", ph.Update)

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
	productID, err := strconv.ParseUint(chi.URLParam(r, "productID"), 10, 0)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	product, err := ph.PR.GetProduct(uint(productID))
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

func (ph ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	productID, err := strconv.ParseUint(chi.URLParam(r, "productID"), 10, 0)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	product := &models.Product{}
	if err := render.Bind(r, product); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	(*product).ID = uint(productID)
	err = ph.PR.UpdateProduct(product)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, product)
}
