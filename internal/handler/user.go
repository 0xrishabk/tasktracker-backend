package handler

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ryszhio/tasktracker/internal/model"
	"github.com/ryszhio/tasktracker/internal/service"
	"github.com/ryszhio/tasktracker/internal/util"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req model.RequestCreateUser
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.userService.CreateUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var secret bool
	env := os.Getenv("ENVIRONMENT")
	if env != "prod" {
		secret = false
	} else {
		secret = true
	}

	util.SetCookie(c, "access_token", res.AccessToken, 60*60*24, secret)

	c.JSON(http.StatusCreated, res)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req model.RequestLoginUser
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.userService.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var secret bool
	env := os.Getenv("ENVIRONMENT")
	if env != "prod" {
		secret = false
	} else {
		secret = true
	}

	util.SetCookie(c, "access_token", res.AccessToken, 60*60*24, secret)

	c.JSON(http.StatusCreated, res)
}

func (h *UserHandler) Logout(c *gin.Context) {
	var secret bool
	env := os.Getenv("ENVIRONMENT")
	if env != "prod" {
		secret = false
	} else {
		secret = true
	}

	util.ClearCookie(c, "", secret)
}

func (h *UserHandler) UpdateUsername(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//user, err := h.userService.UpdateUsername(c, )
	/*
		var secret bool
		env := os.Getenv("ENVIRONMENT")
		if env != "prod" {
			secret = false
		} else {
			secret = true
		}
	*/
	//util.SetCookie(c, res.Username, res.AccessToken, 60*60*24, secret)
}

func (h *UserHandler) Delete(c *gin.Context) {
	uid := c.Param("id")

	err := h.userService.DeleteUser(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var secret bool
	env := os.Getenv("ENVIRONMENT")
	if env != "prod" {
		secret = false
	} else {
		secret = true
	}

	util.SetCookie(c, "", "", -1, secret)

	c.Status(http.StatusOK)
}
