package events

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/helpers"
	"middleware/example/internal/models"
	repository "middleware/example/internal/repositories/events" 
)

// EventConsumer crée la connexion au flux JetStream
func EventConsumer() (jetstream.Consumer, error) {
	// 1. Récupération du contexte JetStream depuis la connexion NATS globale
	js, _ := jetstream.New(helpers.NatsConn) 
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 2. Récupération du Stream "SCHEDULER"
	stream, err := js.Stream(ctx, "SCHEDULER")
	if err != nil {
		return nil, err
	}

	// 3. Création/Récupération du Consumer Durable
	consumer, err := stream.Consumer(ctx, "timetable_consumer")
	if err != nil {
		consumer, err = stream.CreateConsumer(ctx, jetstream.ConsumerConfig{
			Durable:       "timetable_consumer", // le nom qui permet de reprendre où on s'est arrêté
			Name:          "timetable_consumer",
			FilterSubject: "SCHEDULER.events", // on ne veut que les events
			Description:   "Consumer qui met à jour la BDD Timetable",
		})
		if err != nil {
			return nil, err
		}
		logrus.Infof("Consumer créé avec succès")
	} else {
		logrus.Infof("Consumer existant récupéré")
	}

	return consumer, nil
}

// Consume lance la boucle d'écoute
func Consume(consumer jetstream.Consumer) (err error) {


	//On recupere le contexte JetStream pour pouvoir publier des alertes
	js, _ := jetstream.New(helpers.NatsConn)
	ctx := context.Background() 
	logrus.Info("Démarrage du traitement des messages...")
    
	cc, err := consumer.Consume(func(msg jetstream.Msg) {
        // 1. On accuse réception tout de suite (ou à la fin, au choix)
		_ = msg.Ack()

        // 2. Décodage
		var incomingEvent models.Event
		err := json.Unmarshal(msg.Data(), &incomingEvent)
		if err != nil {
			logrus.Errorf("Erreur décodage JSON: %v", err)
			return // On sort, tant pis pour ce message mal formé
		}

        // 3. Logique Métier (Check BDD)
		existingEvent, err := repository.GetEventByUID(incomingEvent.Uid)

		if err != nil {
			// -> NOUVEAU COURS
			logrus.Infof("[AJOUT] Nouveau cours reçu : %s", incomingEvent.Name)
            
            if incomingEvent.Id == nil {
                uid, _ := uuid.NewV4()
                incomingEvent.Id = &uid
            }
            incomingEvent.LastUpdate = time.Now()
            
			_ = repository.CreateEvent(&incomingEvent)

		} else {
			// -> COURS EXISTANT : On cherche les modifs
            hasChanged := false
            var alerteMsg string

            // changement de salle
            if incomingEvent.Location != existingEvent.Location {
				alerteMsg = "Changement de salle: " + existingEvent.Location + " -> " + incomingEvent.Location
                logrus.Warnf("ALERTE : Changement de salle pour %s (%s -> %s)", 
                    incomingEvent.Name, existingEvent.Location, incomingEvent.Location)
                hasChanged = true
                
                // TODO: envoie d'un message NATS vers Alerter
            }
            
            // changement d'heure
            if !incomingEvent.Start.Equal(existingEvent.Start) {
				alerteMsg = "Changement d'horaire pour " + incomingEvent.Name
                 logrus.Warnf("ALERTE : Changement d'horaire pour %s", incomingEvent.Name)
                 hasChanged = true

				 // TODO: envoie d'un message NATS vers Alerter
            }

            if hasChanged {
                incomingEvent.LastUpdate = time.Now()
                _ = repository.UpdateEvent(&incomingEvent)

				alertPayload := models.AlertMessage{
					AgendaID: incomingEvent.AgendaID,
					EventName: incomingEvent.Name,
					Message: alerteMsg,
				}

				payloadBytes, _ := json.Marshal(alertPayload)

				_, errPublish := js.Publish(ctx, "ALERTS.modification", payloadBytes)
				if errPublish != nil{
					logrus.Errorf("Erreur publication alerte NATS: %s", errPublish)
				}else{
					logrus.Infof("Alerte publiee vers le NATS pour %s", incomingEvent.Name)
				}
            }
		}
	})

	<-cc.Closed()
	return err
}
