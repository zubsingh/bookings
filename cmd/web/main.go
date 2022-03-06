package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/zubsingh/bookings/internal/config"
	"github.com/zubsingh/bookings/internal/driver"
	"github.com/zubsingh/bookings/internal/handlers"
	"github.com/zubsingh/bookings/internal/models"
	"github.com/zubsingh/bookings/internal/render"
	"log"
	"net/http"
	"time"
)

var app config.AppConfig

const portNumber = ":8080"

var session *scs.SessionManager

func main() {
	// what i am going to put in session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	// change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// Connect to database
	db, err := driver.ConnectSQL("host=localhost port=5431 dbname=bookings user=zubinsingh password=")
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}
	defer db.SQL.Close()

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)

	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	fmt.Println("Server Started at %d", portNumber)
	err = srv.ListenAndServe()
	log.Fatal(err)

	// _ = http.ListenAndServe(portNumber, nil)
	// fmt.Println("Server Ended at %d", portNumber)
}
