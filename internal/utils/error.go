package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iki-rumondor/go-speech/internal/domain/structs/response"
)

func HandleError(c *gin.Context, err error) {
	res := err.(*response.Error)

	switch res.Code {
	case 400:
		c.AbortWithStatusJSON(http.StatusBadRequest, response.ERROR_RES(res.Message))
		return
	case 401:
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.ERROR_RES(res.Message))
		return
	case 404:
		c.AbortWithStatusJSON(http.StatusNotFound, response.ERROR_RES(res.Message))
		return
	case 406:
		c.AbortWithStatusJSON(http.StatusNotAcceptable, response.ERROR_RES(res.Message))
		return
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ERROR_RES(res.Message))
		return
	}
}

func IsErrorType(err error) bool {
	_, ok := err.(*response.Error)
	return ok
}
