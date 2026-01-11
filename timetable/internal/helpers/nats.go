package helpers

import (
    "log"

    "github.com/nats-io/nats.go"
)

var NatsConn *nats.Conn
var js nats.JetStreamContext

func InitNats() {
    var err error

	// 1. Connexion au serveur
    NatsConn, err = nats.Connect(nats.DefaultURL) // "nats://127.0.0.1:4222"
    if err != nil {
        log.Fatal("Erreur connexion NATS:", err)
    }

	// 2. Contexte JetStream
    js, err = NatsConn.JetStream()
	if err != nil {
        log.Fatal("Erreur JetStream:", err)
    }

    // 3. Création du Stream (Le tuyau général)
    // Si le stream existe déjà, ça ne plantera pas (sauf config différente)
    _, err = js.AddStream(&nats.StreamConfig{
        Name:     "SCHEDULER",
        Subjects: []string{"SCHEDULER.>"}, // tous les sujets commençant par SCHEDULER
    })
    if err != nil {
        log.Printf("Info Stream (peut-être déjà existant): %v", err)
    }

    _, err = js.AddStream(&nats.StreamConfig{
        Name:     "ALERTS",
        Subjects: []string{"ALERTS.>"},
    })
    if err != nil {
        log.Printf("Info Stream ALERTS: %v", err) // Log juste l'info si existe déjà
    } else {
        log.Println("Stream 'ALERTS' initialisé avec succès.")
    }
}

// Fonction pour écouter (Consumer)
func SubscribeToEventsUpdated(handler func(data []byte)) {
    // On s'abonne au sujet spécifique "SCHEDULER.events"
    _, err := js.Subscribe("SCHEDULER.events", func(m *nats.Msg) {
        m.Ack() 
        handler(m.Data)
    })

    if err != nil {
        log.Printf("Erreur lors de la souscription : %v", err)
    } else {
        log.Println("Listening on SCHEDULER.events...")
    }
}
