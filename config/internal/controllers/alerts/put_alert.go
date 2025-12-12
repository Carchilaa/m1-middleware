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
// @Summary      Update alert with id.
// @Description  Update alert with id.
// @Success      200            {array}  models.Alerts
// @Failure      500             "Something went wrong while updating alert {id}"
// @Router       /alerts [put]
type UpdateAlertRequest struct{
	Email string `json:"email"`
	AgendaId uuid.UUID `json:"agendaID"`
}


func PutAlert(w http.ResponseWriter, r *http.Request) {
	var req UpdateAlertRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	id, _ := ctx.Value("alertId").(uuid.UUID)

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