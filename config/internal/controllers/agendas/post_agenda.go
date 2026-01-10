package agendas

import (
	"encoding/json"
	"net/http"
	agenda_service "middleware/example/internal/services/agendas"
	"middleware/example/internal/models"
)


func CreateAgendaHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		http.Error(w, "Methode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var agenda models.Agenda
	if err := json.NewDecoder(r.Body).Decode(&agenda); err != nil {
        http.Error(w, "Requête invalide", http.StatusBadRequest)
        return
    }

	createdAgenda, err := agenda_service.PostAgenda(agenda)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(createdAgenda)
}
