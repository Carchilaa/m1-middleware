package alert

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
	rows, err := db.Query("SELECT * FROM alerts")
	helpers.CloseDB(db)
	if err != nil {
		return nil, err
	}

	// parsing datas in object slice
	alerts := []models.Alert{}
	for rows.Next() {
		var data models.Alert
		err = rows.Scan(&data.Id, &data.Email, &data.IdAgenda)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, data)
	}
	// don't forget to close rows
	_ = rows.Close()

	return alerts, err
}

func GetAlertById(id uuid.UUID) (*models.Alert, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	row := db.QueryRow("SELECT * FROM alerts WHERE id=?", id.String())
	helpers.CloseDB(db)

	var alert models.Alert
	err = rows.Scan(&data.Id, &data.Email, &data.IdAgenda)
	if err != nil {
		return nil, err
	}
	return &alert, err
}

//TODO : Continue to develop
// func PutAlert(id uuid.UUID, email string, agendaId uuid.UUID) (*models.Alert, error) {
// 	db, err := helpers.OpenDB()
// 	if err != nil {
// 		return nil, err
// 	}
// 	row := db.QueryRow("INSERT INTO alerts VALUES id=? email=? agendaId=?", id.String(), email.String(), agendaId.String())
// 	helpers.CloseDB(db)

// 	var alert models.Alert
// 	err = rows.Scan(&data.Id, &data.Email, &data.IdAgenda)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &alert, err
// }




