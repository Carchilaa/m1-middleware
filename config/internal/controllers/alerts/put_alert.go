package alerts

import (
	"github.com/gofrs/uuid"
	"encoding/json"
	"middleware/example/internal/helpers"
	"middleware/example/internal/services/alerts"
	"net/http"
)

// PutAlert
// @Tags         alerts
// @Summary      Create new alert.
// @Description  Create new alert.
// @Success      200            {array}  models.Alerts
// @Failure      500             "Something went wrong while creating alert {id}"
// @Router       /alerts [put]
type PutAlertRequest struct{
	Email string `json:"email"`
	AgendaId uuid.UUID `json:"agendaID"`
}


func PutAlert(w http.ResponseWriter, r *http.Request) {
	var req PutAlertRequest

	//Vérifier que le body n'est pas nul
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	//Génère nouvel id
	id := uuid.Must(uuid.NewV4())

	alert, err := alerts.PutAlert(id, req.Email, req.AgendaId)
	if err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(alert)
	_, _ = w.Write(body)
	return
}