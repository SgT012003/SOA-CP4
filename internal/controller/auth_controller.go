package controller

import (
	"agendamento-salas/internal/models"
	"agendamento-salas/internal/repository"
	"agendamento-salas/internal/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	Repo *repository.UsuarioRepository
}

func NewAuthController(repo *repository.UsuarioRepository) *AuthController {
	return &AuthController{Repo: repo}
}

// Registrar godoc
// @Summary Registra novo usuário
// @Description Cria conta e hasheia a senha
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.RegistrarRequest true "Dados de registro"
// @Success 201 {object} models.Usuario
// @Failure 400 {object} map[string]string
// @Router /api/auth/registrar [post]
func (ctrl *AuthController) Registrar(c *gin.Context) {
	var req models.RegistrarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos", "detalhe": err.Error()})
		return
	}

	hash, err := utils.HashPassword(req.Senha)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar hash da senha"})
		return
	}

	u := &models.Usuario{
		Nome:      req.Nome,
		Email:     req.Email,
		SenhaHash: hash,
	}

	if err := ctrl.Repo.Create(u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar usuário, email pode já estar em uso"})
		return
	}

	c.JSON(http.StatusCreated, u)
}

// Login godoc
// @Summary Faz login do usuário
// @Description Valida email/senha e retorna token JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Credenciais"
// @Success 200 {object} models.AuthResponse
// @Failure 401 {object} map[string]string
// @Router /api/auth/login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos", "detalhe": err.Error()})
		return
	}

	u, err := ctrl.Repo.GetByEmail(req.Email)
	if err != nil || u == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não encontrado ou credenciais inválidas"})
		return
	}

	if !utils.CheckPasswordHash(req.Senha, u.SenhaHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
		return
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "supersecretkey_change_in_prod"
	}

	token, err := utils.GenerateToken(u.ID, []byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token"})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{Token: token})
}
