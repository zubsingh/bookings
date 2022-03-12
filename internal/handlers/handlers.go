package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
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
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(rw, errors.New("cannot get reservation from session"))
		return
	}

	roomName, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		helpers.ServerError(rw, errors.New("cannot able to fetch room Name"))
		return
	}
	res.Room = roomName

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(rw, r, "make-reservation.page.html", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// PostReservation handles the posting of a reservation page
func (m *Repository) PostReservation(rw http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}
	//2021-01-01 -- 01/02 03:04:05PM '06-0700
	layout := "2006-01-02"

	sd := r.Form.Get("start_date")
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		fmt.Println("Parsing Error startDate")
		return
	}

	ed := r.Form.Get("end_date")
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		fmt.Println(" Parsing Error endDate")
		//helpers.ServerError(rw, err)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		//helpers.ServerError(rw, err)
		fmt.Println(" strconv Atoi Error")
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

	form.Required("first_name", "last_name", "email", "start_date", "end_date")
	form.IsEmail("email")

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

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		fmt.Println("PostAvailability Parsing Error startDate")
		return
	}

	endDate, err := time.Parse(layout, end)
	if err != nil {
		fmt.Println(" PostAvailability Parsing Error endDate")
		//helpers.ServerError(rw, err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		fmt.Println(" Parsing Error SearchAvailabilityForAllRooms")
		return
	}

	for _, i := range rooms {
		fmt.Println(" printing log ", i.ID, i.RoomName)
	}
	if len(rooms) == 0 {
		// no availability
		m.App.Session.Put(r.Context(), "error", "No Availability")
		http.Redirect(rw, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}
	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(rw, r, "choose-room.page.html", &models.TemplateData{
		Data: data,
	})

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
	render.Template(rw, r, "contact.page.html", &models.TemplateData{})
}

func (m *Repository) ReservationSummary(rw http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		log.Println("cannot get the session")
	}
	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(rw, r, "reservation-summary.page.html", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

func (m *Repository) ChooseRoom(rw http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Fatal("error at chooseRoom")
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		log.Fatal("error at choose get room session")
		return
	}

	res.RoomID = roomID
	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(rw, r, "/make-reservation", http.StatusSeeOther)
}
