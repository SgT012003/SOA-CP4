package repository

import (
	"agendamento-salas/internal/models"
	"database/sql"
	"errors"
)

type SalaRepository struct {
	DB *sql.DB
}

func NewSalaRepository(db *sql.DB) *SalaRepository {
	return &SalaRepository{DB: db}
}

// Create cadastra sala
func (r *SalaRepository) Create(s *models.Sala) error {
	query := `INSERT INTO salas (nome, capacidade, localizacao, status) VALUES ($1, $2, $3, $4) RETURNING id, criado_em`
	status := s.Status
	if status == "" {
		status = "ATIVA"
	}

	err := r.DB.QueryRow(query, s.Nome, s.Capacidade, s.Localizacao, status).Scan(&s.ID, &s.CriadoEm)
	return err
}

// GetByID busca sala
func (r *SalaRepository) GetByID(id string) (*models.Sala, error) {
	var s models.Sala
	query := `SELECT id, nome, capacidade, localizacao, status, criado_em FROM salas WHERE id = $1`
	
	err := r.DB.QueryRow(query, id).Scan(&s.ID, &s.Nome, &s.Capacidade, &s.Localizacao, &s.Status, &s.CriadoEm)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

// List todas as salas
func (r *SalaRepository) List() ([]models.Sala, error) {
	query := `SELECT id, nome, capacidade, localizacao, status, criado_em FROM salas ORDER BY nome ASC`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var salas []models.Sala
	for rows.Next() {
		var s models.Sala
		if err := rows.Scan(&s.ID, &s.Nome, &s.Capacidade, &s.Localizacao, &s.Status, &s.CriadoEm); err != nil {
			return nil, err
		}
		salas = append(salas, s)
	}
	return salas, nil
}
