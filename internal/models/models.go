package models

import "time"

// Usuario representa um usuario no sistema
type Usuario struct {
	ID        string    `json:"id" db:"id"`
	Nome      string    `json:"nome" db:"nome"`
	Email     string    `json:"email" db:"email"`
	SenhaHash string    `json:"-" db:"senha_hash"`
	CriadoEm  time.Time `json:"criado_em" db:"criado_em"`
}

// Sala representa uma sala de reuniao
type Sala struct {
	ID          string    `json:"id" db:"id"`
	Nome        string    `json:"nome" db:"nome"`
	Capacidade  int       `json:"capacidade" db:"capacidade"`
	Localizacao string    `json:"localizacao" db:"localizacao"`
	Status      string    `json:"status" db:"status"` // ATIVA, INATIVA
	CriadoEm    time.Time `json:"criado_em" db:"criado_em"`
}

// Reserva representa um agendamento de sala
type Reserva struct {
	ID            string    `json:"id" db:"id"`
	IDSala        string    `json:"id_sala" db:"id_sala"`
	IDUsuario     string    `json:"id_usuario" db:"id_usuario"`
	Data          time.Time `json:"data" db:"data"` 
	HorarioInicio string    `json:"horario_inicio" db:"horario_inicio"` 
	HorarioFim    string    `json:"horario_fim" db:"horario_fim"`       
	Finalidade    string    `json:"finalidade" db:"finalidade"`
	Status        string    `json:"status" db:"status"` // CONFIRMADA, CANCELADA
	CriadoEm      time.Time `json:"criado_em" db:"criado_em"`
	
	// Campos extras para listagem (joins)
	NomeSala    string `json:"nome_sala,omitempty" db:"nome_sala"`
	NomeUsuario string `json:"nome_usuario,omitempty" db:"nome_usuario"`
}

// --- DTOs (Data Transfer Objects) para Request do Swagger / Gin ---

type RegistrarRequest struct {
	Nome  string `json:"nome" binding:"required" example:"Joao Silva"`
	Email string `json:"email" binding:"required,email" example:"joao@empresa.com"`
	Senha string `json:"senha" binding:"required,min=6" example:"senha123"`
}

type LoginRequest struct {
	Email string `json:"email" binding:"required,email" example:"admin@empresa.com"`
	Senha string `json:"senha" binding:"required" example:"admin123"`
}

type AuthResponse struct {
	Token string `json:"token" example:"eyJhbG..."`
}

type SalaRequest struct {
	Nome        string `json:"nome" binding:"required" example:"Sala TI"`
	Capacidade  int    `json:"capacidade" binding:"required,gt=0" example:"10"`
	Localizacao string `json:"localizacao" binding:"required" example:"Andar 3"`
	Status      string `json:"status" binding:"omitempty" example:"ATIVA"`
}

type ReservaRequest struct {
	IDSala        string `json:"id_sala" binding:"required" example:"c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22"`
	Data          string `json:"data" binding:"required" example:"2024-05-15"` 
	HorarioInicio string `json:"horario_inicio" binding:"required" example:"09:00"` 
	HorarioFim    string `json:"horario_fim" binding:"required" example:"11:00"`    
	Finalidade    string `json:"finalidade" binding:"required" example:"Reuniao Trimestral"`
}
