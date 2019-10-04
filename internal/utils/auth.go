package utils

import (
	"errors"
	"net/http"
	"strconv"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/freitagsrunde/k4ever-backend/internal/models"
	"github.com/gin-gonic/gin"
)

func CheckRole(level int, c *gin.Context) bool {
	if CheckRoleWithoutAbort(level, c) {
		return true
	}
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	return false
}

func CheckRoleWithoutAbort(level int, c *gin.Context) bool {
	claims := jwt.ExtractClaims(c)
	role := int(claims["role"].(float64))
	if role >= level {
		return true
	}
	return false
}

func CheckIfUserAccess(user string, level int, c *gin.Context) bool {
	if !CheckRoleWithoutAbort(level, c) {
		return false
	}
	if user != "" && user != jwt.ExtractClaims(c)["name"].(string) {
		return false
	}
	return true
}

func ParseDefaultParams(c *gin.Context) (models.DefaultParams, error) {
	var err error
	params := models.DefaultParams{}
	params.SortBy = c.DefaultQuery("sort_by", "name")
	params.Order = c.DefaultQuery("order", "asc")
	offset := c.Query("offset")
	if offset != "" {
		params.Offset, err = strconv.Atoi(offset)
	}
	if err != nil {
		return models.DefaultParams{}, errors.New("offset is not a number")
	}
	limit := c.Query("limit")
	if limit != "" {
		params.Limit, err = strconv.Atoi(limit)
	}
	if err != nil {
		return models.DefaultParams{}, errors.New("limit is not a number")
	}
	return params, nil
}
