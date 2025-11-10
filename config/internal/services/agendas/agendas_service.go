package agendas

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/models"
	agendas_repository "middleware/example/internal/repositories/agendas"
)


func GetAllAgendas()([]models.Agenda, error) {
    var err error

    agendas, err := agendas_repository.GetAllAgendas()

    if err != nil{
        logrus.Errorf("Error retrieving agendas: %s", err.Error())
        return nil, &models.ErrorGeneric{
            Message: "Something went wrong while retrieving agendas",
        }
    }

    return agendas, nil
}

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