package handlers

import (
	"github.com/zubsingh/bookings/cmd/pkg/config"
	"github.com/zubsingh/bookings/cmd/pkg/models"
	"github.com/zubsingh/bookings/cmd/pkg/render"
	"net/http"
)

var Repo *Repository

type Repository struct {
	app *config.AppConfig
}

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		app: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(rw http.ResponseWriter, r *http.Request) {
	//render.RenderTemplate(rw, "home.html")
	render.RenderTemplate(rw, "home.page.html", &models.TemplateData{})
}

// About is aboout page Handler
func (m *Repository) About(rw http.ResponseWriter, r *http.Request) {
	// perform logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, Again"

	// send the data to the template
	render.RenderTemplate(rw, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}
