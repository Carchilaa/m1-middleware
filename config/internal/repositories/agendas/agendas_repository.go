package agendas

import (
	"github.com/gofrs/uuid"
	"middleware/example/internal/helpers"
	"middleware/example/internal/models"
)

func PostAgenda(id uuid.UUID, ucaId int, name string) (*models.Agenda, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	_, err = db.Exec("INSERT INTO agenda (id, ucaId, name) VALUES ($1, $2, $3);", id, ucaId, name)
	if err != nil {
		return nil, err
	}

	// On retourne ce qu'on vient d'insérer 
	agenda := &models.Agenda{
		Id:    &id,
		UcaId: ucaId,
		Name:  name,
	}

	return agenda, nil
}

func GetAllAgendas()([]models.Agenda, error){
	db, err := helpers.OpenDB()
	if err != nil{
		return nil, err
	}
	rows, err := db.Query("SELECT * FROM agenda")
	defer helpers.CloseDB(db)
	if err != nil {
		return nil, err
	}

	agendas := []models.Agenda{}
	for rows.Next() {
		var data models.Agenda
		err = rows.Scan(&data.Id, &data.UcaId, &data.Name)
		if err != nil {
			return nil, err
		}
		agendas = append(agendas, data)
	}
	_ = rows.Close()

	return agendas, err

}