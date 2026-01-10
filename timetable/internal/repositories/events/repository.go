package events

import (
	"github.com/gofrs/uuid"
	"middleware/example/internal/helpers"
	"middleware/example/internal/models"
)

// API REST fait uniquement les Read

func GetAllEvents() ([]models.Event, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	rows, err := db.Query("SELECT id, uid, description, name, start, end, location, lastUpdate FROM events")
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

// Consumer NATS fait les Create, Update et Read (avec le UID)

func GetEventByUID(uid string) (*models.Event, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	row := db.QueryRow("SELECT id, uid, description, name, start, end, location, lastUpdate FROM events WHERE uid=?", uid)

	var ev models.Event
	err = row.Scan(&ev.Id, &ev.Uid, &ev.Description, &ev.Name, &ev.Start, &ev.End, &ev.Location, &ev.LastUpdate)
	if err != nil {
		return nil, err // retourne une erreur si non trouvé, c'est ce qu'on veut
	}
	return &ev, nil
}

func CreateEvent(ev *models.Event) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	statement, _ := db.Prepare("INSERT INTO events (id, uid, description, name, start, end, location, lastUpdate) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	_, err = statement.Exec(ev.Id.String(), ev.Uid, ev.Description, ev.Name, ev.Start, ev.End, ev.Location, ev.LastUpdate)
	
    return err
}

func UpdateEvent(ev *models.Event) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)
	statement, _ := db.Prepare("UPDATE events SET description=?, name=?, start=?, end=?, location=?, lastUpdate=? WHERE uid=?")
	_, err = statement.Exec(ev.Description, ev.Name, ev.Start, ev.End, ev.Location, ev.LastUpdate, ev.Uid)
	return err
}
