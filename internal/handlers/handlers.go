package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/zubsingh/bookings/internal/config"
	"github.com/zubsingh/bookings/internal/models"
	"github.com/zubsingh/bookings/internal/render"
	"log"
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
	render.RenderTemplate(rw, r, "home.page.html", &models.TemplateData{})
}

// About is about page Handler
func (m *Repository) About(rw http.ResponseWriter, r *http.Request) {
	// perform logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, Again"

	// send the data to the template
	render.RenderTemplate(rw, r, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Reservation renders the make a reservation page and display form
func (m *Repository) Reservation(rw http.ResponseWriter, r *http.Request) {
	//render.RenderTemplate(rw, "home.html")
	render.RenderTemplate(rw, r, "make-reservation.page.html", &models.TemplateData{})
}

// Generals renders the room page
func (m *Repository) Generals(rw http.ResponseWriter, r *http.Request) {
	//render.RenderTemplate(rw, "home.html")
	render.RenderTemplate(rw, r, "generals.page.html", &models.TemplateData{})
}

// Majors renders the room page
func (m *Repository) Majors(rw http.ResponseWriter, r *http.Request) {
	//render.RenderTemplate(rw, "home.html")
	render.RenderTemplate(rw, r, "majors.page.html", &models.TemplateData{})
}

// Availability renders the room page
func (m *Repository) Availability(rw http.ResponseWriter, r *http.Request) {
	//render.RenderTemplate(rw, "home.html")
	render.RenderTemplate(rw, r, "search-availability.page.html", &models.TemplateData{})
}

// PostAvailability renders the post request
func (m *Repository) PostAvailability(rw http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start_date")
	end := r.Form.Get("end_date")
	rw.Write([]byte(fmt.Sprintf("Posted to search %s %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJson handles request for availability and send JSON response
func (m *Repository) AvailabilityJson(rw http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		log.Println(err)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(out)
}

// Contact renders page
func (m *Repository) Contact(rw http.ResponseWriter, r *http.Request) {
	//render.RenderTemplate(rw, "home.html")
	render.RenderTemplate(rw, r, "contact.page.html", &models.TemplateData{})
}