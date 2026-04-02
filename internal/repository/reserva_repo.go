package repository

import (
	"agendamento-salas/internal/models"
	"database/sql"
	"errors"
	"time"
)

type ReservaRepository struct {
	DB *sql.DB
}

func NewReservaRepository(db *sql.DB) *ReservaRepository {
	return &ReservaRepository{DB: db}
}

// CheckConflict verifica se ha conflito de horario para uma sala especifica
func (r *ReservaRepository) CheckConflict(salaID string, data string, horarioInicio string, horarioFim string) (bool, error) {
	query := `
		SELECT COUNT(id) FROM reservas
		WHERE id_sala = $1 
		  AND data = $2 
		  AND status = 'CONFIRMADA'
		  AND (horario_inicio < $4 AND horario_fim > $3)
	`
	var count int
	err := r.DB.QueryRow(query, salaID, data, horarioInicio, horarioFim).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Create cria uma nova reserva
func (r *ReservaRepository) Create(res *models.Reserva) error {
	query := `
		INSERT INTO reservas (id_sala, id_usuario, data, horario_inicio, horario_fim, finalidade, status)
		VALUES ($1, $2, $3, $4, $5, $6, 'CONFIRMADA')
		RETURNING id, criado_em, status
	`
	
	err := r.DB.QueryRow(
		query, 
		res.IDSala, res.IDUsuario, res.Data, res.HorarioInicio, res.HorarioFim, res.Finalidade,
	).Scan(&res.ID, &res.CriadoEm, &res.Status)
	
	return err
}

// List retorna todas as reservas agrupadas com nome da sala e usuario
func (r *ReservaRepository) List() ([]models.Reserva, error) {
	query := `
		SELECT 
			r.id, r.id_sala, r.id_usuario, r.data, r.horario_inicio, r.horario_fim, r.finalidade, r.status, r.criado_em,
			s.nome as nome_sala, u.nome as nome_usuario
		FROM reservas r
		INNER JOIN salas s ON r.id_sala = s.id
		INNER JOIN usuarios u ON r.id_usuario = u.id
		ORDER BY r.data DESC, r.horario_inicio DESC
	`
	
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservas []models.Reserva
	for rows.Next() {
		var res models.Reserva
		var data time.Time
		if err := rows.Scan(
			&res.ID, &res.IDSala, &res.IDUsuario, &data, &res.HorarioInicio, &res.HorarioFim, 
			&res.Finalidade, &res.Status, &res.CriadoEm,
			&res.NomeSala, &res.NomeUsuario,
		); err != nil {
			return nil, err
		}
		// Formatar data em string (YYYY-MM-DDT00...) apenas para envio seguro ou deixar time.Time e swagger formata
		res.Data = data
		reservas = append(reservas, res)
	}
	
	return reservas, nil
}

// GetByID busca reserva
func (r *ReservaRepository) GetByID(id string) (*models.Reserva, error) {
	var res models.Reserva
	var data time.Time
	query := `
		SELECT 
			r.id, r.id_sala, r.id_usuario, r.data, r.horario_inicio, r.horario_fim, r.finalidade, r.status, r.criado_em,
			s.nome as nome_sala, u.nome as nome_usuario
		FROM reservas r
		INNER JOIN salas s ON r.id_sala = s.id
		INNER JOIN usuarios u ON r.id_usuario = u.id
		WHERE r.id = $1
	`
	
	err := r.DB.QueryRow(query, id).Scan(
		&res.ID, &res.IDSala, &res.IDUsuario, &data, &res.HorarioInicio, &res.HorarioFim, 
		&res.Finalidade, &res.Status, &res.CriadoEm,
		&res.NomeSala, &res.NomeUsuario,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	
	res.Data = data
	return &res, nil
}

// Cancel atualiza status da reserva
func (r *ReservaRepository) Cancel(id string, userID string) error {
	// Apenas quem criou pode cancelar, logica simples de autorizacao (ou admin)
	query := `UPDATE reservas SET status = 'CANCELADA' WHERE id = $1 AND id_usuario = $2`
	
	result, err := r.DB.Exec(query, id, userID)
	if err != nil {
		return err
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rows == 0 {
		return errors.New("reserva nao encontrada ou voce nao tem permissao")
	}
	
	return nil
}
