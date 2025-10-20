package agendas

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/models"
	agendas_repository "middleware/example/internal/repositories/agendas"
)

func PostAgenda(input models.Agenda) (*models.Agenda, error) {
	uuid, err := uuid.NewV4()
    agenda, err := agendas_repository.PostAgenda(uuid, input.UcaId, input.Name)
    if err != nil {
        logrus.Errorf("error posting agenda %s : %s", uuid.String(), err.Error())
        return nil, &models.ErrorGeneric{
            Message: fmt.Sprintf("Something went wrong while posting agenda %s", uuid.String()),
        }
    }

    return agenda, nil
}