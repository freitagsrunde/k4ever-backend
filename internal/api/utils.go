package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/freitagsrunde/k4ever-backend/internal/k4ever"
	"github.com/freitagsrunde/k4ever-backend/internal/utils"
	"github.com/gin-gonic/gin"
)

func setImage(c *gin.Context, object fmt.Stringer, config k4ever.Config) string {
	//var product models.Product
	//if err := config.DB().First(&product).Where("id = ?", c.Param("id")).Error; err != nil {
	//	c.JSON(http.StatusNotFound, gin.H{"error": "no such product"})
	//	return
	//}

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

	path, err := utils.UploadFile(bytes, fmt.Sprintf("%v/%s", object, object.String()), fileHeader.Filename, config)
	if err != nil {
		if err.Error() == "file already exists" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "file already exists"})
			return ""
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while saving file"})
		return ""
	}

	return config.HttpServerHost() + ":" + strconv.Itoa(config.HttpServerPort()) + path

	//k4ever.UpdateProduct(&product, config)
	//	c.JSON(http.StatusOK, product)

}
