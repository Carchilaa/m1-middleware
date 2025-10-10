package alerts

import (
	"encoding/json"
	"middleware/example/internal/helpers"
	"middleware/example/internal/services/alerts"
	"net/http"
)

// GetUsers
// @Tags         alerts
// @Summary      Get all alerts.
// @Description  Get all alerts.
// @Success      200            {array}  models.Alerts
// @Failure      500             "Something went wrong"
// @Router       /alerts [get]
func GetAlerts(w http.ResponseWriter, _ *http.Request) {
	// calling service
	alerts, err := alerts.GetAlerts()
	if err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(alerts)
	_, _ = w.Write(body)
	return
}