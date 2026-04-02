# Prompt para Agente de IA: Sistema de Agendamento de Salas

> **Atue como um Engenheiro de Software Full-Stack Sênior e Arquiteto de Soluções.**
> 
> Seu objetivo é desenvolver do zero um Sistema de Agendamento de Salas de Reunião corporativo, pronto para produção. O projeto deve ser modular, com separação clara de responsabilidades (Clean Architecture ou MVC adaptado para Go).

## ⚙️ Stack Tecnológica Obrigatória
* **Back-end:** Golang com o framework GinGonic.
* **Banco de Dados:** PostgreSQL 16.
* **Front-end:** HTML5 puro, TailwindCSS (via CDN ou build simples) e Vanilla JavaScript (ES5).
* **Documentação:** Swagger (usando `swaggo/swag` ou similar para auto-gerar a partir de comentários no Go).
* **Infraestrutura:** Docker e Docker Compose.

## 📦 Estrutura de Infraestrutura (Docker)
Crie um arquivo `docker-compose.yml` utilizando 2 containers:
1. **`backend`**: Container Go (Alpine) que expõe a API REST e também serve os arquivos estáticos do Front-end (HTML/CSS/JS) em uma rota específica (ex: `/public` ou `/`).
2. **`database`**: Container PostgreSQL 16.
*(Gere também os respectivos `Dockerfile` e arquivos `.env.example`)*

## 🗄️ Modelagem de Dados e Regras de Negócio
Crie as migrações/scripts SQL de inicialização para as seguintes tabelas:

### 1. Entidade: Usuario
* `id` (UUID ou Serial, PK)
* `nome` (String, Not Null)
* `email` (String, Unique, Not Null)
* `senha_hash` (String, Not Null) - **Obrigatório o uso de bcrypt antes de salvar**

### 2. Entidade: Sala
* `id` (UUID ou Serial, PK)
* `nome` (String, Unique, Not Null)
* `capacidade` (Integer, Not Null)
* `localizacao` (String, Not Null)
* `status` (Enum/String: 'ATIVA' ou 'INATIVA', Default 'ATIVA')

### 3. Entidade: Reserva
* `id` (UUID ou Serial, PK)
* `id_sala` (FK para Sala, Not Null)
* `id_usuario` (FK para Usuario, Not Null) - **Garante a rastreabilidade de quem fez a reserva**
* `data` (Date, Not Null)
* `horario_inicio` (Time, Not Null)
* `horario_fim` (Time, Not Null)
* `finalidade` (Text)
* `status` (Enum/String: 'CONFIRMADA', 'CANCELADA')
> **Regra:** O back-end deve validar se a sala está 'ATIVA' antes de permitir a reserva e se não há conflito de horários para a mesma sala na mesma data.

## 🔌 API REST e Endpoints

**Autenticação e Segurança (Login Interno com JWT):**
* Utilize `golang.org/x/crypto/bcrypt` para hash de senhas e `github.com/golang-jwt/jwt/v5` para geração de tokens.
* O payload do JWT deve conter o `id` do usuário para garantir a rastreabilidade.
* Implemente um middleware no Gin que proteja os endpoints de Salas e Reservas, exigindo o JWT via Bearer Token no header `Authorization`.

**Endpoints Públicos (Auth):**
* `POST /api/auth/registrar` - Recebe nome, email e senha. Hashea a senha e salva o usuário.
* `POST /api/auth/login` - Recebe email e senha. Valida o hash e retorna o token JWT.

**Endpoints Protegidos (Salas):**
* `POST /api/salas` - Cadastrar sala
* `GET /api/salas` - Listar salas
* `GET /api/salas/:id` - Buscar sala por ID

**Endpoints Protegidos (Reservas):**
* `POST /api/reservas` - Criar reserva (Extrair o `id_usuario` direto do JWT da requisição. Validar conflitos de horário).
* `GET /api/reservas` - Listar reservas (Retornar os dados da sala e do usuário atrelados).
* `GET /api/reservas/:id` - Buscar reserva por ID
* `PATCH /api/reservas/:id/cancelar` - Cancelar reserva (Muda o status)

## 🎨 Front-end (UI/UX)
Desenvolva uma Single Page Application (SPA) ou páginas HTML simples consumindo a API via Fetch API (ES5).
* **Tema:** Claro (Light theme).
* **Estilo Visual:** Utilize cores vibrantes em botões, cards e cabeçalhos (ex: Rosa, Azul ciano). Use gradientes sutis do Tailwind (ex: `bg-gradient-to-r from-blue-400 to-pink-500`) onde for pertinente.
* **Componentes:** Crie telas para Login/Registro, listagem/cadastro de salas e listagem/criação de reservas. O token JWT retornado no login deve ser armazenado localmente para ser enviado nas próximas requisições.

## 🚀 Entregáveis Esperados (Divida sua resposta em partes)
1. Scripts SQL de criação do banco.
2. Código fonte do Back-end em Go (Main, Rotas, Controllers, Models, Middleware).
3. Código do Front-end (HTML e JS ES5).
4. Arquivos do Docker (`Dockerfile` e `docker-compose.yml`).
5. Instruções claras de como rodar o projeto e acessar o Swagger.

---

## 🏗️ Diagrama Lógico da Arquitetura do Sistema

Para facilitar o entendimento estrutural, o sistema funcionará com a seguinte arquitetura de comunicação:

### 1. Camada de Apresentação (Browser/Front-end)
* Construída em **HTML5**, estilizada com utilitários do **TailwindCSS** e com a lógica de requisição feita em **JS (ES5)**.
* Realiza o login, guarda o JWT em armazenamento local (ex: localStorage) e faz chamadas HTTP enviando o token no cabeçalho.

### 2. Camada de Aplicação (Docker Container 1 - Back-end Go)
* O serviço expõe a porta (ex: `8080`).
* Atua em duas frentes:
    1. **File Server:** Entrega os arquivos estáticos do front-end (`.html`, `.js`, `.css`) para o navegador do usuário.
    2. **API REST (GinGonic):** Intercepta as requisições com o **Auth Middleware** (Bearer token), valida as regras de negócio (bcrypt, validação de horários) e documenta tudo na rota `/swagger/index.html`.

### 3. Camada de Persistência (Docker Container 2 - Banco de Dados)
* Serviço do **PostgreSQL 16**, rodando em uma rede interna do Docker (invisível para o mundo exterior, mas acessível pelo container do Go).
* Armazena de forma relacional os registros de Usuários, Salas e Reservas.