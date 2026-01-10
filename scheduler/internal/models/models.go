package models

import (
	"github.com/gofrs/uuid"
	"time"
)


type Agenda struct {
	Id *uuid.UUID	`json:"id"`
	UcaId int		`json:"ucaId"`
	Name string		`json:"name"`
}

type Course struct{
	AgendaName string
	Summary string
	Start time.Time
}