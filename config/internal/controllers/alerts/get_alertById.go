package alerts

import (
	"encoding/json"
	"middleware/example/internal/helpers"
	"net/http"
	agenda_service "middleware/example/internal/services/alerts"
	"github.com/gofrs/uuid"
)

// GetUser
// @Tags         alerts
// @Summary      Get an alert by id.
// @Description  Get an alert by id.
// @Param        idAgenda          	path      string  true
// @Success      200            {object}  models.Alert
// @Failure      422            "Cannot parse id"
// @Failure      500            "Something went wrong"
// @Router       /alerts/{idAgenda} [get]
func GetAlertById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	alertIdAgenda, _ := ctx.Value("IdAgenda").(uuid.UUID)

	alert, err := agenda_service.GetAlertById(alertIdAgenda)
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