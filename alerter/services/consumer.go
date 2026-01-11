package services

import (
	
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/sirupsen/logrus"

	"alerter/models"
)

func StartAlerterConsumer(nc *nats.Conn) {
	js, _ := jetstream.New(nc)
	ctx := context.Background()

	// 1. Création du Stream "ALERTS" s'il n'existe pas
	_, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     "ALERTS",
		Subjects: []string{"ALERTS.>"},
	})
	if err != nil {
		logrus.Warnf("Stream ALERTS probablement existant ou erreur: %v", err)
	}

	// 2. Création du Consumer
	consumer, err := js.CreateOrUpdateConsumer(ctx, "ALERTS", jetstream.ConsumerConfig{
		Durable:       "alerter_consumer",
		Name:          "alerter_consumer",
		FilterSubject: "ALERTS.modification",
	})
	if err != nil {
		logrus.Fatal("Impossible de créer le consumer Alerter:", err)
	}

	logrus.Info("Alerter Consumer prêt. En attente de messages...")

	// 3. Boucle de consommation
	cc, _ := consumer.Consume(func(msg jetstream.Msg) {
        msg.Ack()

        var alertMsg models.Modification
        if err := json.Unmarshal(msg.Data(), &alertMsg); err != nil {
            logrus.Errorf("Alerter: Erreur JSON: %v", err)
            return
        }

        logrus.Infof("Alerter: Reçu modif pour Agendas : %v", alertMsg.AgendaIDs)

        if len(alertMsg.AgendaIDs) == 0 {
            logrus.Warn("Alerter: Aucun agenda associé à cet événement (liste vide).")
            return
        }

        for _, agendaID := range alertMsg.AgendaIDs {
            
            // 1. Récupérer les abonnés pour CET agenda précis
            subscribers, err := getSubscribers(agendaID)
            if err != nil {
                logrus.Errorf("Alerter: Erreur API Config pour agenda %s : %v", agendaID, err)
                continue 
            }

            if len(subscribers) == 0 {
                continue
            }

            // 2. Préparer le mail
            bodyContent, subject, err := parseTemplate("alert.txt", alertMsg)
            if err != nil {
                logrus.Errorf("Alerter: Erreur Template: %v", err)
                return
            }

            // 3. Envoyer aux abonnés de ce agenda
            for _, sub := range subscribers {
                err := sendMail(sub.Email, subject, bodyContent)
                if err != nil {
                    logrus.Errorf("Alerter: Echec envoi à %s", sub.Email)
                } else {
                    logrus.Infof("Alerter: Mail envoyé à %s (Agenda %s)", sub.Email, agendaID)
                }
            }
        }
	})

	<-cc.Closed()
}
