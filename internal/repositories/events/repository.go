package events

import (
	"middleware/example/internal/helpers"
	"middleware/example/internal/models"
)

func GetAllEvents() ([]models.Event, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("SELECT * FROM events")
	helpers.CloseDB(db)
	if err != nil {
		return nil, err
	}

	// parsing datas in object slice
	events := []models.Event{}
	for rows.Next() {
		var data models.Event
		err = rows.Scan(&data.Id, &data.Dtstamp, &data.Dtstart, &data.Dtend, &data.Description, &data.Location, &data.Created, &data.LastModified, &data.ResourceID)
		if err != nil {
			return nil, err
		}
		events = append(events, data)
	}
	// don't forget to close rows
	_ = rows.Close()

	return events, err
}

func GetEventById(id string) (*models.Event, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	row := db.QueryRow("SELECT * FROM events WHERE id=?", id)
	helpers.CloseDB(db)

	var event models.Event
	err = row.Scan(&event.Id, &event.Dtstamp, &event.Dtstart, &event.Dtend, &event.Description, &event.Location, &event.Created, &event.LastModified, &event.ResourceID)
	if err != nil {
		return nil, err
	}
	return &event, err
}
