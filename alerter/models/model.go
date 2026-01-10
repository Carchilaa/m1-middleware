package models

import (
	"github.com/gofrs/uuid"
	"time"
)

// Remplace par l'URL de ton API Config (localhost ou nom du service docker)
const ConfigAPIUrl = "http://localhost:8080/alerts/{idAgenda}" 
const MailAPIUrl = "https://mail.edu.forestier.re/api/send"
const MailToken = "" //TODO : mettre le token 

//go:embed templates
var embeddedTemplates embed.FS

// Structure pour la réponse de l'API Config
type AlertConfig struct {
    Email string `json:"email"`
    // Ajoute d'autres champs si ton API Config renvoie plus d'infos
}

// Structure pour l'envoi de mail (Forestier API)
type MailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

type FrontMatter struct {
	Subject string `yaml:"subject"`
}


// Modification représente le message reçu via NATS (depuis le Timetable Consumer)
type Modification struct {
    AgendaID  string `json:"agenda_id"`
    EventName string `json:"event_name"`
    Message   string `json:"message"`
    // Ajoute d'autres champs si nécessaire
}

// Alert représente une alerte reçue depuis ton API Config
type Alert struct {
    ID       int    `json:"id"`
    AgendaID string `json:"agenda_id"`
    Email    string `json:"email"`
}

// MailBody représente le corps de la requête vers l'API mail de Forestier
type MailBody struct {
    From    string   `json:"from"`
    To      []string `json:"to"`
    Subject string   `json:"subject"`
    Body    string   `json:"body"`
}

// FrontMatter sert à parser l'en-tête du template
type FrontMatter struct {
    Subject string `yaml:"subject"`
}

