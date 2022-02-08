package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/zubsingh/bookings/cmd/pkg/config"
	"github.com/zubsingh/bookings/cmd/pkg/handlers"
	"github.com/zubsingh/bookings/cmd/pkg/render"
	"log"
	"net/http"
	"time"
)

var app config.AppConfig

const portNumber = ":8080"

var session *scs.SessionManager

func main() {

	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.SetConfig(&app)

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
