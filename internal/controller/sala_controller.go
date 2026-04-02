package controller

import (
	"agendamento-salas/internal/models"
	"agendamento-salas/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SalaController struct {
	Repo *repository.SalaRepository
}

func NewSalaController(repo *repository.SalaRepository) *SalaController {
	return &SalaController{Repo: repo}
}

// Cadastrar godoc
// @Summary Cadastra uma nova sala
// @Description Cria sala nova, default status ATIVA
// @Tags Salas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.SalaRequest true "Dados da sala"
// @Success 201 {object} models.Sala
// @Failure 400 {object} map[string]string
// @Router /api/salas [post]
func (ctrl *SalaController) Cadastrar(c *gin.Context) {
	var req models.SalaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos", "detalhe": err.Error()})
		return
	}

	s := &models.Sala{
		Nome:        req.Nome,
		Capacidade:  req.Capacidade,
		Localizacao: req.Localizacao,
		Status:      req.Status,
	}

	if err := ctrl.Repo.Create(s); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar sala. Nome deve ser unico."})
		return
	}

	c.JSON(http.StatusCreated, s)
}

// Listar godoc
// @Summary Lista todas as salas
// @Description Busca todas as salas, independente do status
// @Tags Salas
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Sala
// @Router /api/salas [get]
func (ctrl *SalaController) Listar(c *gin.Context) {
	salas, err := ctrl.Repo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao listar salas"})
		return
	}
	// Garante que retorne um array vazio invés de null se não tiver nada
	if salas == nil {
		salas = []models.Sala{}
	}
	c.JSON(http.StatusOK, salas)
}

// Buscar godoc
// @Summary Busca sala por ID
// @Tags Salas
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID da Sala"
// @Success 200 {object} models.Sala
// @Router /api/salas/{id} [get]
func (ctrl *SalaController) Buscar(c *gin.Context) {
	id := c.Param("id")
	sala, err := ctrl.Repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar sala"})
		return
	}
	if sala == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sala não encontrada"})
		return
	}
	c.JSON(http.StatusOK, sala)
}
