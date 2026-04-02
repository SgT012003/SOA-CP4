# Arquitetura do Sistema: Agendamento de Salas

Este documento descreve a arquitetura geral do sistema de agendamento corporativo.

## Visão Geral

O sistema baseia-se em uma arquitetura de microserviços simplificada (Monolito Modular) encapsulada em *containers* Docker. 
A stack principal é composta por **Golang e Gin** para a API, e **PostgreSQL** para persistência, com o frontend renderizado e servido pela própria aplicação.

## Diagrama de Blocos C4-Level 2 (Container)

Abaixo o diagrama em Mermaid descrevendo as camadas da aplicação:

```mermaid
flowchart TB
 subgraph Camadas["Camadas"]
        Controllers["Controllers Auth, Salas, Reservas"]
        Repositories["Repositories Acesso a dados via database/sql"]
        Models["Models Structs e Definições de Domínio"]
  end
 subgraph subGraph1["Backend Container Aplicação Go c/ Gin"]
        API["Gin Web Framework Rotas & Middleware JWT"]
        Camadas
        Static["Recursos Estáticos & Templates"]
        Swagger["Swagger UI swaggo/gin-swagger"]
  end
 subgraph subGraph2["Ambiente Docker Compose"]
        subGraph1
        DB[("Banco de Dados PostgreSQL 16")]
  end
    Client["Navegador/Cliente HTTP HTML, Vanilla JS, TailwindCSS"] -- Chamadas REST API Autenticadas com JWT --> API
    API -- Gerencia Requisições --> Controllers
    Controllers -- Busca e Salva Dados --> Repositories
    Controllers -. Mapeamento .-> Models
    Repositories -. Mapeamento .-> Models
    API -. Renderiza views .- Static
    API -. Exibe documentação .- Swagger
    Repositories -- Queries SQL Nativas --> DB
    Client -- "Acessa Front-end" --> Static
    Client -- Lê Documentação --> Swagger

     Controllers:::container
     Repositories:::container
     Models:::container
     API:::container
     Static:::container
     Swagger:::container
     DB:::db
     Client:::client
    classDef container fill:#115e59,stroke:#0f766e,stroke-width:2px,color:#fff
    classDef db fill:#336791,stroke:#2b5578,stroke-width:2px,color:#fff
    classDef client fill:#be185d,stroke:#9d174d,stroke-width:2px,color:#fff
    style Client fill:#00C853,stroke:none
```

## Descrição dos Componentes

1. **Frontend (Cliente):** 
   - SPA híbrido: HTML servido do Go, mas com consumo dinâmico via Fetch API.
   - CSS minificado usando a CLI isolada do TailwindCSS v3 acoplada no pipeline de compilação multi-stage do Docker.
   - Todo controle e redirecionamento de estado acontecem via Vanilla JSON/JS e localStorage tokens.

2. **Backend (App Go):**
   - **Gin Framework:** Garante excelente performance de roteamento HTTP.
   - **Middleware:** Responsável pelo desacoplamento e checagem de Token (JWT).
   - **Camada Controller:** Gerencia HTTP Request/Response unicamente, transferindo dados.
   - **Camada Repository:** Interage com o PostgreSQL através de consultas SQL nativas (`database/sql`) substituindo soluções pesadas como ORMs.

3. **Database (PostgreSQL 16):**
   - Roda encapsulado na mesma virtual network do docker. 
   - A inicialização do DDL e carga primária (Seeding) corre de um script de inicialização `init.sql` na subida primária do driver DB.

4. **Documentação Integrada (Swaggo):**
   - Permite o acesso pelo frontend e geração de artefatos de spec (`swagger.json`/`yaml`) através de comentários nativos no pacote principal.
