# Estagio 1: Tailwind CSS builder
FROM node:20-alpine AS css-builder
WORKDIR /app
RUN npm install tailwindcss@3
COPY tailwind.config.js input.css ./
COPY templates/ ./templates/
COPY public/js/ ./public/js/
RUN mkdir -p public/css
RUN npx tailwindcss -i input.css -o public/css/style.css --minify

# Estagio 2: Backend Go builder
FROM golang:alpine AS go-builder
WORKDIR /app
RUN go install github.com/swaggo/swag/cmd/swag@latest
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Gera pasta docs p/ o swagger ser embutido (main.go importa _ "agendamento-salas/docs")
RUN swag init -g cmd/api/main.go

RUN go build -ldflags="-w -s" -o main ./cmd/api/main.go

# Estagio 3: Servidor Prod
FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates tzdata

# Copia binario e assets
COPY --from=go-builder /app/main .
COPY --from=go-builder /app/.env.example .env
# Copia o folder docs gerado (talvez o binario nao reclame de não ter o /docs fisico dependendo se os assets staticos sao servidos pelo go embed, o swaggo usa Go code imports `docs/docs.go`, o q embute no proprio binario! Entao a pasta nao precisa ser movida).

# Copia front templates e public
COPY --from=go-builder /app/templates/ ./templates/
COPY --from=go-builder /app/public/js/ ./public/js/
# Pega o CSS do builder do Node
COPY --from=css-builder /app/public/css/ ./public/css/

EXPOSE 8080
CMD ["./main"]
