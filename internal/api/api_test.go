package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/freitagsrunde/k4ever-backend/internal/k4ever"
	"github.com/freitagsrunde/k4ever-backend/internal/test"
	"github.com/gin-gonic/gin"
)

func NewApiTest() (app *gin.Engine, router *gin.RouterGroup, conf k4ever.Config) {
	conf = test.NewConfig()
	conf.MigrateDB()
	gin.SetMode(gin.TestMode)
	app = gin.New()

	router = app.Group("/api/v1")

	return app, router, conf
}

func PerformRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func PerformRequestWithBody(r http.Handler, method, path string, body []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
