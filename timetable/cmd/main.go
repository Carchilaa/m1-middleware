package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/controllers/events"
	"middleware/example/internal/helpers"
	_ "middleware/example/internal/models"
	"net/http"
)

func main() {
	// ⚠️ IMPORTANT : Initialiser NATS AVANT toute utilisation !
    helpers.InitNats()

    // TEST : publier un message
    err := helpers.PublishEventUpdated(map[string]string{
        "message": "Hello depuis Timetable",
    })
    if err != nil {
        logrus.Println("Erreur Publish:", err)
    }

    // TEST : souscription
    helpers.SubscribeToEventsUpdated(func(data []byte) {
        logrus.Println("Message reçu:", string(data))
    })

    // Ici tu lanceras ensuite :
    // runMyConsumer()
    // runMyServer()

    select {} // empêcher l'arrêt du programme

	r := chi.NewRouter()

	r.Route("/events", func(r chi.Router) { // route /events
		r.Get("/", events.GetEvents)            // GET /events
		r.Route("/{id}", func(r chi.Router) { // route /events/{id}
			r.Use(events.Context)      // Use Context method to get event ID
			r.Get("/", events.GetEvent) // GET /events/{id}
		})
	})

	logrus.Info("[INFO] Web server started. Now listening on *:8081")
	logrus.Fatalln(http.ListenAndServe(":8081", r))
}

func init() {
	db, err := helpers.OpenDB()
	if err != nil {
		logrus.Fatalf("error while opening database : %s", err.Error())
	}
	schemes := []string{
		`CREATE TABLE IF NOT EXISTS events (
			id VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
			uid VARCHAR(255) NOT NULL,
			description TEXT,
			name VARCHAR(255) NOT NULL,
			start DATETIME NOT NULL,
			end DATETIME NOT NULL,
			location VARCHAR(255),
			lastUpdate DATETIME NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS eventAgendas (
			eventId VARCHAR(255) NOT NULL,
			agendaId VARCHAR(255) NOT NULL,
			FOREIGN KEY (eventId) REFERENCES events(id) ON DELETE CASCADE,
			FOREIGN KEY (agendaId) REFERENCES agendas(id) ON DELETE CASCADE
		);`,
	}
	for _, scheme := range schemes {
		if _, err := db.Exec(scheme); err != nil {
			logrus.Fatalln("Could not generate table ! Error was : " + err.Error())
		}
	}
	helpers.CloseDB(db)
}
