package alerts

import (
	"database/sql"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/models"
	repository "middleware/example/internal/repositories/alerts"
)

func GetAlerts() ([]models.Alert, error) {
	var err error
	alerts, err := repository.GetAlerts()
	if err != nil {
		logrus.Errorf("Error while retrieving alerts : %s", err.Error())
		return nil, &models.ErrorGeneric{
			Message: "Something went wrong while retrieving alerts.",
		}
	}
	return alerts, nil
}

func GetAlertById(id uuid.UUID) (*models.Alert, error) {
	alert, err := repository.GetAlertById(id)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, &models.ErrorNotFound{
				Message: fmt.Sprintf("Alert (id : %s) not found", id.String()),
			}
		}
		logrus.Errorf("Error retrieving alert %s : %s", id.String(), err.Error())
		return nil, &models.ErrorGeneric{
			Message: fmt.Sprintf("Something went wrong while retrieving alert %s", id.String()),
		}
	}
	return alert, err
}

func PutAlert(id uuid.UUID, email string, agendaId uuid.UUID) (*models.Alert, error) {
	alert, err := repository.PutAlert(id, email, agendaId)
	if err != nil {
		return nil, fmt.Errorf("Failed to create alert : %s", err)
	}
	return alert, err
}