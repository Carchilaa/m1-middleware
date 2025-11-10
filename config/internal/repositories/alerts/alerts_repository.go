package alerts

import (
	"github.com/gofrs/uuid"
	"middleware/example/internal/helpers"
	"middleware/example/internal/models"
)

func GetAlerts() ([]models.Alert, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	rows, err := db.Query("SELECT * FROM alerts")

	if err != nil {
		return nil, err
	}

	allAlerts := []models.Alert{}
	for rows.Next() {
		var alert models.Alert
		err = rows.Scan(&alert.Id, &alert.Email, &alert.IdAgenda)
		if err != nil {
			return nil, err
		}
		allAlerts = append(allAlerts, alert)
	}
	_ = rows.Close()

	return allAlerts, err
}

func GetAlertById(id uuid.UUID) (*models.Alert, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	row := db.QueryRow("SELECT * FROM alerts WHERE id=?", id.String())

	var alert models.Alert
	err = row.Scan(&alert.Id, &alert.Email, &alert.IdAgenda)
	if err != nil {
		return nil, err
	}
	return &alert, err
}

func PutAlert(id uuid.UUID, email string, agendaId uuid.UUID) (*models.Alert, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	query := `INSERT INTO alerts (id, email, idAgenda) VALUES (?, ?, ?)`

	_, err = db.Exec(query, id, email, agendaId)
	if err != nil {
		return nil, err
	}

	alert := models.Alert{
		Id:       &id,
		Email:    email,
		IdAgenda: &agendaId,
	}

	return &alert, nil
}




