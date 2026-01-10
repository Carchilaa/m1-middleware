package agendas

import (
	"github.com/gofrs/uuid"
	"middleware/example/internal/helpers"
	"middleware/example/internal/services/agendas"
	"net/http"
)


func DeleteAgenda(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	agendaId, _ := ctx.Value("agendaId").(uuid.UUID)

	err := agendas.DeleteAgenda(agendaId)
	if err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}