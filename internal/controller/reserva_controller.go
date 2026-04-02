package controller

import (
	"time"
	"agendamento-salas/internal/models"
	"agendamento-salas/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReservaController struct {
	RepoReserva *repository.ReservaRepository
	RepoSala    *repository.SalaRepository
}

func NewReservaController(rr *repository.ReservaRepository, sr *repository.SalaRepository) *ReservaController {
	return &ReservaController{RepoReserva: rr, RepoSala: sr}
}

// Criar godoc
// @Summary Cria nova reserva
// @Description Verifica disponibilidade antes de criar
// @Tags Reservas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.ReservaRequest true "Dados reserva"
// @Success 201 {object} models.Reserva
// @Router /api/reservas [post]
func (ctrl *ReservaController) Criar(c *gin.Context) {
	userID := c.GetString("user_id")

	var req models.ReservaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetros inválidos", "detalhe": err.Error()})
		return
	}

	// 1. Verifica se sala existe e está ativa
	sala, err := ctrl.RepoSala.GetByID(req.IDSala)
	if err != nil || sala == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sala não encontrada"})
		return
	}
	if sala.Status != "ATIVA" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sala inativa ou indisponível no momento"})
		return
	}

	// 2. Verifica Conflito de horário
	conflito, err := ctrl.RepoReserva.CheckConflict(req.IDSala, req.Data, req.HorarioInicio, req.HorarioFim)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar disponibilidade"})
		return
	}
	if conflito {
		c.JSON(http.StatusConflict, gin.H{"error": "Horário já reservado para esta sala nesta data"})
		return
	}

	// 3. Cria e formata a Data (YYYY-MM-DD -> time.Time)
	dataParsed, errTime := time.Parse("2006-01-02", req.Data)
	if errTime != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de data invalido, experado YYYY-MM-DD"})
		return
	}

	r := &models.Reserva{
		IDSala:        req.IDSala,
		IDUsuario:     userID,
		Data:          dataParsed,          
		HorarioInicio: req.HorarioInicio, 
		HorarioFim:    req.HorarioFim,
		Finalidade:    req.Finalidade,
	}

	if err := ctrl.RepoReserva.Create(r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao registrar reserva no banco", "detalhe": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, r)
}

// Listar godoc
// @Summary Lista todas reservas
// @Description Inclui informações de Usuários e Salas atreladas
// @Tags Reservas
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Reserva
// @Router /api/reservas [get]
func (ctrl *ReservaController) Listar(c *gin.Context) {
	reservas, err := ctrl.RepoReserva.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao listar reservas"})
		return
	}
	if reservas == nil {
		reservas = []models.Reserva{}
	}
	c.JSON(http.StatusOK, reservas)
}

// Cancelar godoc
// @Summary Cancela reserva (Apenas proprietario)
// @Tags Reservas
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID Reserva"
// @Success 200 {object} map[string]string
// @Router /api/reservas/{id}/cancelar [patch]
func (ctrl *ReservaController) Cancelar(c *gin.Context) {
	reservaID := c.Param("id")
	userID := c.GetString("user_id")

	err := ctrl.RepoReserva.Cancel(reservaID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reserva cancelada com sucesso"})
}
