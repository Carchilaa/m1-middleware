package agendas

import (
	"encoding/json"
	"net/http"
	agenda_service "middleware/example/internal/services/agendas"
	"middleware/example/internal/helpers"
)


func GetAllAgendas(w http.ResponseWriter, _ *http.Request){

	agendas, err := agenda_service.GetAllAgendas()
	if err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(agendas)
	_, _ = w.Write(body)
	return
}