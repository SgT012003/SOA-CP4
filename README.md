# Agendamento de Salas (AgendamentoPRO)

Sistema moderno e responsivo para gerenciamento corporativo de salas de reunião, desenhado para prover de forma simples o cadastro, visualização e reserva de espaços físicos.

---

### Grupo

|Nome|RM|
|:-:|:-:|
|Diogo Julio|553837|
|Jonata Rafael|552939|
|Matheus Zottis|94119|
|Victor Didoff|552965|
|Vinicius Silva|553240|

---

## Índice
1. [Sobre o Projeto](#-sobre-o-projeto)
2. [Funcionalidades Principais](#-funcionalidades-principais)
3. [Stack Tecnológica](#%EF%B8%8F-stack-tecnológica)
4. [Estrutura de Arquivos](#-estrutura-de-arquivos-file-tree)
5. [Como Rodar com Docker (Recomendado)](#-como-rodar-com-docker-recomendado)
6. [Como Rodar Localmente (Desenvolvimento)](#-como-rodar-localmente-desenvolvimento)
7. [Exemplos de Uso (cURL)](#-exemplos-de-uso-curl)
8. [Documentação da API (Swagger)](#-documentação-da-api-swagger)

---

## Sobre o Projeto
O **AgendamentoPRO** é um projeto modular desenvolvido em **Golang** (utilizando Gin) com **PostgreSQL**. Ele serve sua própria interface (HTML5 e Vanilla JS estilizada com TailwindCSS), dispensando a necessidade de frameworks front-end pesados.

## Funcionalidades Principais
- **Autenticação:** Login e Registro seguro via JWT e hashing de senha (bcrypt).
- **Gestão de Salas:** Listagem dinâmica das capacidades, localizações e status (Ativa/Inativa) das salas.
- **Reservas:** Sistema prático para agendamento de janelas de tempo atrelando usuário logado à sala requerida.
- **Visualização:** Uma Single Page Application fluida que processa as requisições API de forma silenciosa via `fetch`.

## Stack Tecnológica
- **Back-end:** Go 1.22+, Gin Web Framework
- **Banco de Dados:** PostgreSQL 16 (Interação via DB nativa `database/sql`)
- **Front-end:** HTML, Vanilla JS (ES6+)
- **Estilização:** Tailwind CSS v3 (Compilação via Docker/Node)
- **Infraestrutura:** Docker e Docker Compose multi-stage build.

---

## Estrutura de Arquivos (File Tree)

Esta é a estrutura modular simplificada:

```text
CP4 - SOA/
├── cmd/
│   └── api/
│       └── main.go                 # Entrypoint da Aplicação Golang
├── docs/                           # Documentação Swagger e Diagramas (C4/Mermaid)
├── internal/
│   ├── controller/                 # Camada de Requisição/Resposta HTTP e Roteamento
│   ├── middleware/                 # Interceptadores (Log, Autenticação JWT)
│   ├── models/                     # Entidades, Structs Visuais e Domínio
│   └── repository/                 # Camada de Banco de Dados Postgre (SQL Nativo)
├── public/                         # Arquivos Estáticos compilados (JS e CSS)
├── scripts/
│   └── init.sql                    # Script contendo DDL e Seeding do Banco
├── templates/                      # Componentes UI (Arquivos HTML processados pelo Gin)
├── docker-compose.yml              # Cluster Infraestrutura
├── Dockerfile                      # Etapas de Build Customizado (Node TailWind + Go Alpine)
├── go.mod / go.sum                 # Gerenciamento de Dependencias Go
└── tailwind.config.js              # Paramêtros do Tailwindcss
```

---

## Como Rodar com Docker (Recomendado)

A infraestrutura foi desenhada para subir com zero configuração nativa, compilando Front-end e Back-end e criando os contêineres e o banco de dados magicamente.

1. Clone o repositório ou navegue até a pasta do projeto:
   ```bash
   cd "SOA-CP4"
   ```

2. Suba o projeto de forma desvinculada (`-d`) e reconstruindo as imagens (`--build`):
   ```bash
   docker compose up -d --build
   ```

3. Acesse no seu navegador:
   - **Sistema Web:** [http://localhost:8080/](http://localhost:8080/)
   - **Banco de Dados:** Irá rodar internamente na porta `5432` automaticamente. O script inicial (`init.sql`) irá providenciar as tabelas e usuários básicos, caso necessite.

---

## Como Rodar Localmente (Desenvolvimento)

Se você preferir instalar localmente para programar e debugar (Host direto):

### 1. Requisitos:
- Instale [Go (1.22+)](https://go.dev/dl/)
- Instale [Node.js](https://nodejs.org/) (Para rodar o TailwindCSS)
- Um banco de dados **PostgreSQL** rodando localmente (Porta padrão: `5432`)

### 2. Configurações:
Crie um arquivo `.env` na raiz espelhado no `.env.example`:
```ini
DB_HOST=localhost
DB_PORT=5432
DB_USER=seu_usuario
DB_PASS=sua_senha
DB_NAME=agendamento_db
JWT_SECRET=super_senha_ultra_secreta
PORT=8080
```

### 3. Compilando o CSS Analiticamente (Tailwind):
Para rodar a CLI do tailwind localmente:
```bash
npm install tailwindcss@3
npx tailwindcss -i input.css -o public/css/style.css --watch
```

### 4. Compilando as Documentações da API (Swagger):
Sempre que você modificar o código das rotas ou baixar localmente pela primeira vez, é necessário gerar/atualizar o Swagger:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/api/main.go
```

### 5. Executando o Back-end (Go):
Em outro terminal paralelo, inicie o backend:
```bash
go mod tidy
go run ./cmd/api/main.go
```
O servidor estará acessível de idêntica maneira em [http://localhost:8080/](http://localhost:8080/).

---

## Exemplos de Uso (cURL)

Você pode testar manualmente o funcionamento da API sem depender da GUI usando requisições como estas:

### 1. Criar um Registro (Usuário)
```bash
curl -X POST http://localhost:8080/api/auth/registrar \
-H "Content-Type: application/json" \
-d '{"nome":"Carlos Silva", "email":"carlos@empresa.com", "senha":"minhasenha123"}'
```

### 2. Fazer Login e adquirir um Token JWT
```bash
curl -X POST http://localhost:8080/api/auth/login \
-H "Content-Type: application/json" \
-d '{"email":"carlos@empresa.com", "senha":"minhasenha123"}'
```
*Guarde o `"token"` retornado na resposta para os próximos passos.*

### 3. Criar uma Nova Sala (Requer JWT Header)
```bash
curl -X POST http://localhost:8080/api/salas \
-H "Authorization: Bearer <SEU_TOKEN_AQUI>" \
-H "Content-Type: application/json" \
-d '{"nome":"Sala Ômega", "capacidade":15, "localizacao":"4º Andar", "status":"ATIVA"}'
```

### 4. Consultar todas as Salas Abertas
```bash
curl -X GET http://localhost:8080/api/salas \
-H "Authorization: Bearer <SEU_TOKEN_AQUI>"
```

---

## Documentação da API (Swagger)

A aplicação conta com a integração nativa ao Swagger (Auto-gerada via comentários de pacotes em Go).

Você pode checar visualmente todos os *schemas*, testar as rotas, e validar formatos diretamente via browser acessando:

**Acesso Automático Dashboard Swagger:**
👉 [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
