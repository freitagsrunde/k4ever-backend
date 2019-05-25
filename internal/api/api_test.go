package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/freitagsrunde/k4ever-backend/internal/k4ever"
	"github.com/freitagsrunde/k4ever-backend/internal/models"
	"github.com/freitagsrunde/k4ever-backend/internal/test"
	"github.com/gin-gonic/gin"
)

type tokenStruct struct {
	Expire string `json:"expire"`
	Token  string `json:"token"`
}

func NewApiTest() (app *gin.Engine, router *gin.RouterGroup, conf k4ever.Config) {
	conf = test.NewConfig()
	conf.MigrateDB()
	gin.SetMode(gin.TestMode)
	app = gin.New()

	CreateAuthMiddleware(conf)
	app.POST("/api/v1/login/", AuthMiddleware.LoginHandler)
	router = app.Group("/api/v1")
	router.Use(AuthMiddleware.MiddlewareFunc())

	// Create test user
	testUser := UserTest()
	k4ever.CreateUser(&testUser, conf)

	return app, router, conf
}

func PerformRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)

	// Create token
	token := getToken(r)
	req.Header.Set("Authorization", "Bearer "+token)

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func getToken(r http.Handler) string {
	requestBody, err := json.Marshal(map[string]string{
		"password": "test",
		"name":     "test",
	})

	if err != nil {
		log.Fatal("Could not marshal request body")
		return ""
	}

	req, err := http.NewRequest("POST", "/api/v1/login/", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatal(err.Error())
	}

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	result := tokenStruct{}
	err = json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		log.Fatal("Could not get token")
		return ""
	}
	return result.Token
}

func PerformRequestWithBody(r http.Handler, method, path string, body []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func UserTest() (testUser models.User) {
	testUser.UserName = "test"
	testUser.Password = "test"
	testUser.DisplayName = "Test"

	return testUser
}
