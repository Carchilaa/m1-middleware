package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/controllers/agendas"
	"middleware/example/internal/controllers/alerts"
	"middleware/example/internal/helpers"
	_ "middleware/example/internal/models"
	"net/http"
)

func main() {
	r := chi.NewRouter()	
	r.Route("/agendas", func(r chi.Router) {
		r.Post("/", agendas.CreateAgendaHandler) // POST /agendas // http://localhost:8080/agendas/
		r.Get("/", agendas.GetAllAgendas) // Get All Agendas
		r.Route("/{id}", func(r chi.Router){
			r.Use(agendas.Context)
			r.Get("/", agendas.GetAgenda)
			r.Delete("/", agendas.DeleteAgenda)
			r.Put("/", agendas.PutAgenda) //Update alert with id
		})
	})

	r.Route("/alerts", func(r chi.Router) {
		r.Get("/", alerts.GetAlerts) //Get all alerts
		r.Post("/", alerts.PostAlert) //Create a new alert
		r.Route("/{idAgenda}", func(r chi.Router) {
			r.Use(alerts.Context)
			r.Get("/", alerts.GetAlertById) //Get an alert with id
			r.Delete("/", alerts.DeleteAlert) //Delete an alert with id
			r.Put("/", alerts.PutAlert) //Update alert with id
		})
	})

	logrus.Info("[INFO] Web server started. Now listening on *:8080")
	logrus.Fatalln(http.ListenAndServe(":8080", r))
}

func init() {
	db, err := helpers.OpenDB()
	if err != nil {
		logrus.Fatalf("error while opening database : %s", err.Error())
	}
	schemes := []string{
		`CREATE TABLE IF NOT EXISTS alerts (
			id VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
			email VARCHAR(255) NOT NULL,
			idAgenda VARCHAR(255) NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS agenda (
			id VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
			ucaId VARCHAR(255) NOT NULL,
			name VARCHAR(255) NOT NULL
		);`,
	}
	for _, scheme := range schemes {
		if _, err := db.Exec(scheme); err != nil {
			logrus.Fatalln("Could not generate table ! Error was : " + err.Error())
		}
	}
	helpers.CloseDB(db)
}
