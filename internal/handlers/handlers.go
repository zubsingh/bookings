package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/zubsingh/bookings/internal/config"
	"github.com/zubsingh/bookings/internal/driver"
	"github.com/zubsingh/bookings/internal/forms"
	"github.com/zubsingh/bookings/internal/helpers"
	"github.com/zubsingh/bookings/internal/models"
	"github.com/zubsingh/bookings/internal/render"
	"github.com/zubsingh/bookings/internal/repository"
	"github.com/zubsingh/bookings/internal/repository/dbrepo"
	"log"
	"net/http"
	"strconv"
	"time"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(rw http.ResponseWriter, r *http.Request) {
	//render.RenderTemplate(rw, "home.html")
	render.Template(rw, r, "home.page.html", &models.TemplateData{})
}

// About is about page Handler
func (m *Repository) About(rw http.ResponseWriter, r *http.Request) {
	// perform logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, Again"

	// send the data to the template
	render.Template(rw, r, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Reservation renders the make a reservation page and display form
func (m *Repository) Reservation(rw http.ResponseWriter, r *http.Request) {
	//render.RenderTemplate(rw, "home.html")
	//render.RenderTemplate(rw, r, "make-reservation.page.html", &models.TemplateData{})
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.Template(rw, r, "make-reservation.page.html", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservation handles the posting of a reservation page
func (m *Repository) PostReservation(rw http.ResponseWriter, r *http.Request) {
	//fmt.Println("abc")
	//err := r.ParseForm()
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//
	//sd := r.Form.Get("start_date")
	//ed := r.Form.Get("end_date")
	//
	////2021-01-01 -- 01/02 03:04:05PM '06-0700
	//layout := "2006-01-02"
	//startDate, err := time.Parse(layout, sd)
	//endDate, err := time.Parse(layout, ed)
	//if err != nil {
	//	helpers.ServerError(rw, err)
	//	return
	//}
	//
	//roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	//if err != nil {
	//	helpers.ServerError(rw, err)
	//	return
	//}
	//
	//reservation := models.Reservation{
	//	FirstName: r.Form.Get("first_name"),
	//	LastName:  r.Form.Get("last_name"),
	//	Phone:     r.Form.Get("phone"),
	//	Email:     r.Form.Get("email"),
	//	StartDate: startDate,
	//	EndDate:   endDate,
	//	RoomID:    roomID,
	//}
	//
	//form := forms.New(r.PostForm)
	//
	//form.Has("first_name", r)
	//form.Required("first_name", "last_name", "email")
	//form.EmailLength("email")
	//
	//if !form.Valid() {
	//	data := make(map[string]interface{})
	//	data["reservation"] = reservation
	//
	//	render.Template(rw, r, "make-reservation.page.html", &models.TemplateData{
	//		Form: form,
	//		Data: data,
	//	})
	//}
	//
	//newReservationID, err := m.DB.InsertReservation(reservation)
	//if err != nil {
	//	helpers.ServerError(rw, err)
	//	return
	//}
	//
	//restriction := models.RoomRestriction{
	//	StartDate:      startDate,
	//	EndDate:        endDate,
	//	RoomID:         roomID,
	//	ReservationsID: newReservationID,
	//	RestrictionsID: 1,
	//}
	//
	//err = m.DB.InsertRoomRestriction(restriction)
	//if err != nil {
	//	helpers.ServerError(rw, err)
	//	return
	//}
	//
	//m.app.Session.Put(r.Context(), "reservation", reservation)
	//http.Redirect(rw, r, "/reservation-summary", http.StatusSeeOther)

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	//2021-01-01 -- 01/02 03:04:05PM '06-0700
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		fmt.Println(" Parsing Error")
		//helpers.ServerError(rw, err)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		//helpers.ServerError(rw, err)
		fmt.Println(" strconv.Atoi Error")
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.IsEmail("email")
	//form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.Template(rw, r, "make-reservation.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		//helpers.ServerError(rw, err)
		fmt.Println("InsertReservation Error")
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomID:        roomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		//helpers.ServerError(rw, err)
		fmt.Println("InsertRoomRestriction Error")
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(rw, r, "/reservation-summary", http.StatusSeeOther)
}

// Generals renders the room page
func (m *Repository) Generals(rw http.ResponseWriter, r *http.Request) {
	//render.RenderTemplate(rw, "home.html")
	render.Template(rw, r, "generals.page.html", &models.TemplateData{})
}

// Majors renders the room page
func (m *Repository) Majors(rw http.ResponseWriter, r *http.Request) {
	//render.RenderTemplate(rw, "home.html")
	render.Template(rw, r, "majors.page.html", &models.TemplateData{})
}

// Availability renders the room page
func (m *Repository) Availability(rw http.ResponseWriter, r *http.Request) {
	//render.RenderTemplate(rw, "home.html")
	render.Template(rw, r, "search-availability.page.html", &models.TemplateData{})
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
	render.Template(rw, r, "contact.page.html", &models.TemplateData{})
}

func (m *Repository) ReservationSummary(rw http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		log.Println("cannot get the session")
	}
	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(rw, r, "reservation-summary.page.html", &models.TemplateData{
		Data: data,
	})
}
