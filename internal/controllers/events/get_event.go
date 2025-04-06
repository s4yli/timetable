package events

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/models"
	"middleware/example/internal/services/events"
	"net/http"
)

// Get Event
// @Tags         events
// @Summary      Récupère un événement
// @Description  Récupère un événement par son ID
// @Param        id           	path      string  true  "Évènement ID"
// @Success      200            {object}  models.Event
// @Failure      422            "Cannot parse id"
// @Failure      500            "Something went wrong"
// @Router       /events/{id} [get]
func GetEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	eventId, _ := ctx.Value("eventId").(string)

	event, err := events.GetEventById(eventId)
	if err != nil {
		logrus.Errorf("error : %s", err.Error())
		customError, isCustom := err.(*models.CustomError)
		if isCustom {
			w.WriteHeader(customError.Code)
			body, _ := json.Marshal(customError)
			_, _ = w.Write(body)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(event)
	_, _ = w.Write(body)
	return
}
