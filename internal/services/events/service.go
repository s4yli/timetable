package events

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/models"
	repository "middleware/example/internal/repositories/events"
	"net/http"
)

func GetAllEvents() ([]models.Event, error) {
	var err error
	// calling repository
	events, err := repository.GetAllEvents()
	// managing errors
	if err != nil {
		logrus.Errorf("error retrieving events : %s", err.Error())
		return nil, &models.CustomError{
			Message: "Something went wrong",
			Code:    500,
		}
	}

	return events, nil
}

func GetEventById(id string) (*models.Event, error) {
	event, err := repository.GetEventById(id)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, &models.CustomError{
				Message: "event not found",
				Code:    http.StatusNotFound,
			}
		}
		logrus.Errorf("error retrieving events : %s", err.Error())
		return nil, &models.CustomError{
			Message: "Something went wrong",
			Code:    500,
		}
	}

	return event, err
}
