package main

import (
	"log"
	"os"

	"agendamento-salas/internal/controller"
	"agendamento-salas/internal/middleware"
	"agendamento-salas/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	// Swagger imports
	_ "agendamento-salas/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title API de Agendamento de Salas
// @version 1.0
// @description Sistema corporativo de agendamento de salas.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@empresa.com

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Carrega variaveis
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: Falha ao carregar .env, dependendo var env do sistema")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	name := os.Getenv("DB_NAME")

	db, err := repository.ConnectDB(host, port, user, pass, name)
	if err != nil {
		log.Fatalf("Nao foi possivel conectar ao banco: %v", err)
	}
	defer db.Close()

	// Repositories
	usuarioRepo := repository.NewUsuarioRepository(db)
	salaRepo := repository.NewSalaRepository(db)
	reservaRepo := repository.NewReservaRepository(db)

	// Controllers
	authCtrl := controller.NewAuthController(usuarioRepo)
	salaCtrl := controller.NewSalaController(salaRepo)
	reservaCtrl := controller.NewReservaController(reservaRepo, salaRepo)
	viewCtrl := controller.NewViewController()

	router := gin.Default()

	// Front-end files setup
	router.Static("/public", "./public")
	router.LoadHTMLGlob("templates/*")

	// View Routes
	router.GET("/", viewCtrl.RedirectLogin)
	router.GET("/login", viewCtrl.Login)
	router.GET("/salas", viewCtrl.Salas)
	router.GET("/reservas", viewCtrl.Reservas)

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API Routes
	api := router.Group("/api")
	{
		api.POST("/auth/registrar", authCtrl.Registrar)
		api.POST("/auth/login", authCtrl.Login)

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Salas
			protected.POST("/salas", salaCtrl.Cadastrar)
			protected.GET("/salas", salaCtrl.Listar)
			protected.GET("/salas/:id", salaCtrl.Buscar)

			// Reservas
			protected.POST("/reservas", reservaCtrl.Criar)
			protected.GET("/reservas", reservaCtrl.Listar)
			protected.PATCH("/reservas/:id/cancelar", reservaCtrl.Cancelar)
		}
	}

	appPort := os.Getenv("PORT")
	if appPort == "" {
		appPort = "8080"
	}
	
	log.Printf("Servidor subindo na porta %s", appPort)
	if err := router.Run(":" + appPort); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
