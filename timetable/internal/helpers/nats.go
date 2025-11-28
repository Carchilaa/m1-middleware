package helpers

import (
    "encoding/json"
    "errors"
    "log"

    "github.com/nats-io/nats.go"
)

var NatsConn *nats.Conn
var JSContext nats.JetStreamContext

// Init NATS + JetStream + Stream
func InitNats() {
    var err error

    // Connexion NATS
    NatsConn, err = nats.Connect(nats.DefaultURL)
    if err != nil {
        log.Fatal("Erreur connexion NATS:", err)
    }

    // JetStream
    JSContext, err = NatsConn.JetStream()
    if err != nil {
        log.Fatal("Erreur JetStream:", err)
    }

    // Création du stream EVENTS s'il n'existe pas
    _, err = JSContext.AddStream(&nats.StreamConfig{
        Name:     "EVENTS",
        Subjects: []string{"EVENTS.*"},
    })
    if err != nil {
        // Si déjà existant => OK
        if err != nats.ErrStreamNameAlreadyInUse {
            log.Fatal("Erreur création stream:", err)
        }
    }

    log.Println("NATS initialisé ✔️")
}

// Publication d'un message EVÈNEMENT MODIFIÉ
func PublishEventUpdated(payload interface{}) error {

    if JSContext == nil {
        return errors.New("JetStream non initialisé")
    }

    data, err := json.Marshal(payload)
    if err != nil {
        return err
    }

    ack, err := JSContext.Publish("EVENTS.updated", data)
    if err != nil {
        return err
    }

    log.Printf("Message publié dans EVENTS.updated (seq=%d)\n", ack.Sequence)
    return nil
}

// Exemple de souscription
func SubscribeToEventsUpdated(handler func(data []byte)) error {

    if NatsConn == nil {
        return errors.New("Connexion NATS non initialisée")
    }

    // Simple subscriber
    _, err := NatsConn.Subscribe("EVENTS.updated", func(msg *nats.Msg) {
        handler(msg.Data)
    })

    if err != nil {
        return err
    }

    log.Println("Souscription sur EVENTS.updated OK")
    return nil
}
