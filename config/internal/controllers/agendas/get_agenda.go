package agendas

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	"middleware/example/internal/helpers"
	"middleware/example/internal/services/agendas"
	"net/http"
)


func GetAgenda(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	agendaId, _ := ctx.Value("agendaId").(uuid.UUID)

	agenda, err := agendas.GetAgendaById(agendaId)
	if err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(agenda)
	_, _ = w.Write(body)
	return
}