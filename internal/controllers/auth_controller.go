package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hugaojanuario/sentinel/internal/models"
	"github.com/hugaojanuario/sentinel/internal/services"
	"github.com/hugaojanuario/sentinel/internal/utils"
)

type AuthController struct {
	s *services.Service
}

func NewAuthController(s *services.Service) *AuthController {
	return &AuthController{s: s}
}

func (a *AuthController) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "body invalid."})
		return
	}

	token, err := a.s.Login(req)
	if err != nil {
		utils.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, token)
}
