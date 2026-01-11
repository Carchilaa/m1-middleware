package services

import (
	
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/sirupsen/logrus"

	"alerter/models"
)

// StartAlerterConsumer lance le processus
func StartAlerterConsumer(nc *nats.Conn) {
	js, _ := jetstream.New(nc)
	ctx := context.Background()

	// 1. Création du Stream "ALERTS" s'il n'existe pas
	_, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     "ALERTS",
		Subjects: []string{"ALERTS.>"},
	})
	if err != nil {
        // On ignore l'erreur si le stream existe déjà, sinon on log
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

		// A. Décoder le message
		var alertMsg models.Modification
		if err := json.Unmarshal(msg.Data(), &alertMsg); err != nil {
			logrus.Errorf("Alerter: Erreur JSON: %v", err)
			return
		}

		logrus.Infof("Alerter: Reçu modif pour Agenda %s", alertMsg.AgendaID)

		// B. Récupérer les abonnés via l'API Config
		subscribers, err := getSubscribers(alertMsg.AgendaID)
		if err != nil {
			logrus.Errorf("Alerter: Erreur API Config: %v", err)
			return
		}

		if len(subscribers) == 0 {
			logrus.Infof("Aucun abonné pour l'agenda %s", alertMsg.AgendaID)
			return
		}

		// C. Préparer le contenu du mail (Template)
        // Assure-toi d'avoir "templates/notification.html" dans ton dossier
		bodyContent, subject, err := parseTemplate("alert.html", alertMsg)
		if err != nil {
			logrus.Errorf("Alerter: Erreur Template: %v", err)
			return
		}

		// D. Envoyer un mail à chaque abonné
		for _, sub := range subscribers {
			err := sendMail(sub.Email, subject, bodyContent)
			if err != nil {
				logrus.Errorf("Alerter: Echec envoi mail à %s : %v", sub.Email, err)
			} else {
				logrus.Infof("Alerter: Mail envoyé à %s", sub.Email)
			}
		}
	})

	<-cc.Closed()
}
