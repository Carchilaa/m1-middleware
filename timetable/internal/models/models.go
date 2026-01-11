package models

import (
	"github.com/gofrs/uuid"
	"time"
)

type Event struct {
	Id          *uuid.UUID  `json:"id"`
	AgendaIds   []uuid.UUID `json:"agendaIds"`
	Uid         string      `json:"uid"`
	Description string      `json:"description"`
	Name        string      `json:"name"`
	Start       time.Time   `json:"start"`
	End         time.Time   `json:"end"`
	Location    string      `json:"location"`
	LastUpdate  time.Time   `json:"lastUpdate"`
	IncomingAgendaID string `json:"agenda_id"`
}

type EventAgenda struct {
	EventId  uuid.UUID `json:"eventId"`
	AgendaId uuid.UUID `json:"agendaId"`
}

type AlertMessage struct {
	AgendaIds []uuid.UUID `json:"agenda_ids"`
	EventName string      `json:"event_name"`
	Message   string      `json:"message"`
}
