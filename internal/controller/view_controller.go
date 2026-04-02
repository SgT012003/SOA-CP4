package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ViewController struct{}

func NewViewController() *ViewController {
	return &ViewController{}
}

func (v *ViewController) RedirectLogin(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/login")
}

func (v *ViewController) Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func (v *ViewController) Salas(c *gin.Context) {
	c.HTML(http.StatusOK, "salas.html", gin.H{"Title": "Salas de Reunião"})
}

func (v *ViewController) Reservas(c *gin.Context) {
	c.HTML(http.StatusOK, "reservas.html", gin.H{"Title": "Reservas"})
}
