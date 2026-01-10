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

type Course struct {
    Uid string      `json:"uid"` 
    Name string     `json:"name"` 
    Start time.Time `json:"start"`
    End time.Time   `json:"end"` 
    Location string `json:"location"` 
}
