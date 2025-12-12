package models

import (
	"github.com/gofrs/uuid"
)

type User struct {
	Id   *uuid.UUID `json:"id"`
	Name string     `json:"name"`
}

type Alert struct{
	Id *uuid.UUID `json:"id"`
	Email string `json:"email"`
	IdAgenda *uuid.UUID `json:"agendaID"`
}

type Agenda struct {
	Id *uuid.UUID	`json:"id"`
	UcaId int		`json:"ucaId"`
	Name string		`json:"name"`
}