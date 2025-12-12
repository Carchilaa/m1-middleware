package agendas

import (
	"github.com/gofrs/uuid"
	"encoding/json"
	"middleware/example/internal/helpers"
	"middleware/example/internal/services/agendas"
	"middleware/example/internal/models"
	"net/http"
)

// PutAgenda
// @Tags         agenda
// @Summary      Update agenda with id.
// @Description  Update agenda with id.
// @Success      200            {array}  models.Alerts
// @Failure      500             "Something went wrong while updating alert {id}"
// @Router       /agenda [put]
type UpdateAgendaRequest struct{
	UcaId int `json:"UcaId"`
	Name string `json:"Name"`
}


func PutAgenda(w http.ResponseWriter, r *http.Request) {
	var req UpdateAgendaRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	id, _ := ctx.Value("agendaId").(uuid.UUID)

	UpdatedAgenda := models.Agenda{
		Id:       &id,
		UcaId:    req.UcaId,
		Name: req.Name,
	}
	agenda, err := agendas.PutAgendas(id, UpdatedAgenda)
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