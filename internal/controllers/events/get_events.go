package events

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/models"
	"middleware/example/internal/services/events"
	"net/http"
)

// Get Events
// @Tags         events
// @Summary      Récupère tous les événements
// @Description  Récupère tous les événements
// @Success      200            {array}  models.Event
// @Failure      500             "Something went wrong"
// @Router       /events [get]
func GetEvents(w http.ResponseWriter, _ *http.Request) {
	// calling service
	events, err := events.GetAllEvents()
	if err != nil {
		// logging error
		logrus.Errorf("error : %s", err.Error())
		customError, isCustom := err.(*models.CustomError)
		if isCustom {
			// writing http code in header
			w.WriteHeader(customError.Code)
			// writing error message in body
			body, _ := json.Marshal(customError)
			_, _ = w.Write(body)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(events)
	_, _ = w.Write(body)
	return
}
