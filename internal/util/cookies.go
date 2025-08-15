package util

import (
	"github.com/gin-gonic/gin"
)

func SetCookie(c *gin.Context, name, value string, maxAge int, secure bool) {
	c.SetCookie(name, value, maxAge, "/", "", secure, true)
}

func ClearCookie(c *gin.Context, name string, secure bool) {
	c.SetCookie(name, "", -1, "/", "", secure, true)
}
