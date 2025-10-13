package events

import (
	"database/sql"
	"github.com/gofrs/uuid"
	"middleware/example/internal/helpers"
	"middleware/example/internal/models"
)

func GetAllEvents() ([]models.Event, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("SELECT id, uid, description, name, start, end, location, lastUpdate FROM events")
	helpers.CloseDB(db)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// parsing datas in object slice
	events := []models.Event{}
	for rows.Next() {
		var ev models.Event
		err = rows.Scan(
			&ev.Id,
			&ev.Uid,
			&ev.Description,
			&ev.Name,
			&ev.Start,
			&ev.End,
			&ev.Location,
			&ev.LastUpdate,
		)
		if err != nil {
			return nil, err
		}

		agendaRows, err := db.Query("SELECT agendaId FROM eventAgendas WHERE eventId = ?", ev.Id.String())
		if err != nil {
			return nil, err
		}
		var agendaIds []uuid.UUID
		for agendaRows.Next() {
			var agendaId uuid.UUID
			err = agendaRows.Scan(&agendaId)
			if err != nil {
				return nil, err
			}
			agendaIds = append(agendaIds, agendaId)
		}
		_ = agendaRows.Close()

		ev.AgendaIds = agendaIds
		events = append(events, ev)
	}
	// don't forget to close rows
	_ = rows.Close()

	return events, err
}

func GetEventById(id uuid.UUID) (*models.Event, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	row := db.QueryRow("SELECT id, uid, description, name, start, end, location, lastUpdate FROM events WHERE id=?", id.String())
	helpers.CloseDB(db)

	var ev models.Event
	err = row.Scan(
		&ev.Id,
		&ev.Uid,
		&ev.Description,
		&ev.Name,
		&ev.Start,
		&ev.End,
		&ev.Location,
		&ev.LastUpdate,
	)
	if err != nil {
		return nil, err
	}

	db2, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	agendaRows, err := db2.Query("SELECT agendaId FROM eventAgendas WHERE eventId = ?", id.String())
	helpers.CloseDB(db2)
	if err != nil {
		return nil, err
	}
	defer agendaRows.Close()

	var agendaIds []uuid.UUID
	for agendaRows.Next() {
		var agendaId uuid.UUID
		err = agendaRows.Scan(&agendaId)
		if err != nil {
			return nil, err
		}
		agendaIds = append(agendaIds, agendaId)
	}
	ev.AgendaIds = agendaIds

	return &ev, nil
}
