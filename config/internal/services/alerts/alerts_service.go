package alert

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
	// calling repository
	alerts, err := repository.GetAlerts()
	// managing errors
	if err != nil {
		logrus.Errorf("error retrieving alerts : %s", err.Error())
		return nil, &models.ErrorGeneric{
			Message: "Something went wrong while retrieving alerts",
		}
	}

	return alerts, nil
}

func GetAlertById(id uuid.UUID) (*models.Alert, error) {
	alert, err := repository.GetAlertById(id)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, &models.ErrorNotFound{
				Message: "alert not found",
			}
		}
		logrus.Errorf("error retrieving alert %s : %s", id.String(), err.Error())
		return nil, &models.ErrorGeneric{
			Message: fmt.Sprintf("Something went wrong while retrieving alert %s", id.String()),
		}
	}

	return alert, err
}