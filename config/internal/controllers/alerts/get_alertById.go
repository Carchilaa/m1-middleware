package alerts

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	"middleware/example/internal/helpers"
	"middleware/example/internal/services/alerts"
	"net/http"
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
	alertIdAgenda, _ := ctx.Value("IdAgenda")

	alert, err := alerts.GetAlertById(alertIdAgenda)
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