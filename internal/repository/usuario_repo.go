package repository

import (
	"agendamento-salas/internal/models"
	"database/sql"
	"errors"
)

type UsuarioRepository struct {
	DB *sql.DB
}

func NewUsuarioRepository(db *sql.DB) *UsuarioRepository {
	return &UsuarioRepository{DB: db}
}

// Create cria um novo usuario
func (r *UsuarioRepository) Create(u *models.Usuario) error {
	query := `INSERT INTO usuarios (nome, email, senha_hash) VALUES ($1, $2, $3) RETURNING id, criado_em`
	
	err := r.DB.QueryRow(query, u.Nome, u.Email, u.SenhaHash).Scan(&u.ID, &u.CriadoEm)
	return err
}

// GetByEmail busca usuario
func (r *UsuarioRepository) GetByEmail(email string) (*models.Usuario, error) {
	var u models.Usuario
	query := `SELECT id, nome, email, senha_hash, criado_em FROM usuarios WHERE email = $1`
	
	err := r.DB.QueryRow(query, email).Scan(&u.ID, &u.Nome, &u.Email, &u.SenhaHash, &u.CriadoEm)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found
		}
		return nil, err
	}
	
	return &u, nil
}
