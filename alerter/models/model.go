package models

import (
	"embed"
)


const ConfigAPIUrl = "http://localhost:8080/alerts/" 
const MailAPIUrl = "https://mail-api.edu.forestier.re/mail"
const MailToken = "TImsTooRzjVIBIpSuUEejNpCrTzAYKToYJmLVCkp"


var EmbeddedTemplates embed.FS


type AlertConfig struct {
    Email string `json:"email"`
}

type MailRequest struct {
    Recipient string `json:"recipient"`
    Subject   string `json:"subject"`
    Content   string `json:"content"`
}

type FrontMatter struct {
	Subject string `yaml:"subject"`
}


type Modification struct {
    AgendaIDs []string `json:"agenda_ids"`
    EventName string `json:"event_name"`
    Message   string `json:"message"`
    
}


type Alert struct {
    ID       int    `json:"id"`
    AgendaID string `json:"agenda_id"`
    Email    string `json:"email"`
}

type MailBody struct {
    From    string   `json:"from"`
    To      []string `json:"to"`
    Subject string   `json:"subject"`
    Body    string   `json:"body"`
}

