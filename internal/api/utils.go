package api

import (
	"fmt"
	"net/http"
	"path"
	"strconv"

	"github.com/freitagsrunde/k4ever-backend/internal/k4ever"
	"github.com/freitagsrunde/k4ever-backend/internal/utils"
	"github.com/gin-gonic/gin"
)

// Set image will just handle one image per object for now but could be easily extended to support multiple files per object
func setImage(c *gin.Context, object fmt.Stringer, config k4ever.Config) string {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not get file from key file"})
		return ""
	}

	f, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not open file"})
		return ""
	}

	bytes := utils.StreamToByte(f)

	// Delete the old file if the is one
	utils.DeleteFiles(fmt.Sprintf("%T/%s/%s", object, object.String(), object.String()), config)

	path, err := utils.UploadFile(bytes, fmt.Sprintf("%T/%s/%s%s", object, object.String(), object.String(), path.Ext(fileHeader.Filename)), config)
	if err != nil {
		if err.Error() == "file already exists" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "file already exists"})
			return ""
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while saving file: " + err.Error()})
		return ""
	}

	return config.HttpServerHost() + ":" + strconv.Itoa(config.HttpServerPort()) + path
}
