CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS usuarios (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nome VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    senha_hash VARCHAR(255) NOT NULL,
    criado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS salas (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nome VARCHAR(255) UNIQUE NOT NULL,
    capacidade INT NOT NULL CHECK (capacidade > 0),
    localizacao VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'ATIVA' CHECK (status IN ('ATIVA', 'INATIVA')),
    criado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS reservas (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_sala UUID NOT NULL REFERENCES salas(id),
    id_usuario UUID NOT NULL REFERENCES usuarios(id),
    data DATE NOT NULL,
    horario_inicio TIME NOT NULL,
    horario_fim TIME NOT NULL,
    finalidade TEXT,
    status VARCHAR(50) DEFAULT 'CONFIRMADA' CHECK (status IN ('CONFIRMADA', 'CANCELADA')),
    criado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valida_horario CHECK (horario_inicio < horario_fim)
);

-- Seeding
-- A senha do admin sera 'admin123', o bcrypt precisa ser gerado previamente.
-- Usaremos um bcrypt fixo para "admin123" que e: $2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy
INSERT INTO usuarios (id, nome, email, senha_hash) VALUES 
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Administrador', 'admin@empresa.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy')
ON CONFLICT DO NOTHING;

INSERT INTO salas (id, nome, capacidade, localizacao, status) VALUES 
('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'Sala Alpha - Reunioes Rapidas', 5, 'Andar 1 - Setor A', 'ATIVA'),
('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'Sala Beta - Diretoria', 15, 'Andar 2 - Direcao', 'ATIVA'),
('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'Auditorio Gama', 50, 'Andar Terreo', 'ATIVA'),
('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a55', 'Sala Delta - Manutencao', 4, 'Andar 1 - Setor B', 'INATIVA')
ON CONFLICT DO NOTHING;

-- Seed uma reserva de exemplo
INSERT INTO reservas (id_sala, id_usuario, data, horario_inicio, horario_fim, finalidade, status) VALUES
('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', CURRENT_DATE, '09:00:00', '10:00:00', 'Reuniao de Alinhamento Diaria', 'CONFIRMADA')
ON CONFLICT DO NOTHING;
