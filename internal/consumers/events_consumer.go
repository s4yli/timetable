package consumers

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/helpers"
	"middleware/example/internal/models"
	"time"
)

const alerterSubject = "TIMETABLE.ALERTER"

// EventConsumer initialise le consumer NATS pour le stream TIMETABLE
func EventConsumer() (jetstream.Consumer, error) {
	js, err := jetstream.New(helpers.NatsConn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Vérification que le stream existe
	stream, err := js.Stream(ctx, "TIMETABLE")
	if err != nil {
		return nil, err
	}

	// Configuration du consumer
	consumerConfig := jetstream.ConsumerConfig{
		Durable:       "TIMETABLE_CONSUMER",
		AckPolicy:     jetstream.AckExplicitPolicy,
		AckWait:       30 * time.Second,
		MaxDeliver:    5,
		FilterSubject: "TIMETABLE.EVENTS",
	}

	// Création ou récupération du consumer
	consumer, err := stream.CreateOrUpdateConsumer(ctx, consumerConfig)
	if err != nil {
		return nil, err
	}

	logrus.Info("Consumer NATS TIMETABLE initialisé avec succès")
	return consumer, nil
}

// Consume démarre la consommation des messages du consumer TIMETABLE
func Consume(consumer jetstream.Consumer) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consCtx, err := consumer.Consume(func(msg jetstream.Msg) {
		processMessage(msg)
	})
	if err != nil {
		return err
	}
	defer consCtx.Stop()

	<-ctx.Done()
	return nil
}

// processMessage traite un message reçu sur TIMETABLE.EVENTS, met à jour la BDD,
// puis publie un message d'alerte si l'événement est nouveau ou a été modifié.
func processMessage(msg jetstream.Msg) {
	var event models.Event
	if err := json.Unmarshal(msg.Data(), &event); err != nil {
		logrus.Errorf("Erreur de décodage JSON: %v", err)
		_ = msg.Nak()
		return
	}

	db, err := helpers.OpenDB()
	if err != nil {
		logrus.Errorf("Erreur connexion DB: %v", err)
		_ = msg.Nak()
		return
	}
	defer helpers.CloseDB(db)

	// Vérifier si l'événement existe déjà
	var existingEvent models.Event
	err = db.QueryRow("SELECT id, dtstart, dtend, location, description FROM events WHERE id = ?", event.Id).
		Scan(&existingEvent.Id, &existingEvent.Dtstart, &existingEvent.Dtend, &existingEvent.Location, &existingEvent.Description)

	// Gestion des cas (nouvel événement ou mise à jour)
	if err == sql.ErrNoRows {
		// Nouvel événement
		if err := insertEvent(db, event); err != nil {
			logrus.Errorf("Erreur insertion événement: %v", err)
			_ = msg.Nak()
			return
		}
		logrus.Infof("Nouvel événement ajouté: %s", event.Id)

		// Publier l'alerte
		if err := publishAlert(event); err != nil {
			logrus.Errorf("Erreur lors de la publication de l'alerte: %v", err)
		}
	} else if err != nil {
		logrus.Errorf("Erreur vérification événement: %v", err)
		_ = msg.Nak()
		return
	} else {
		// Mise à jour de l'événement si des modifications sont détectées
		if hasChanges(existingEvent, event) {
			if err := updateEvent(db, event); err != nil {
				logrus.Errorf("Erreur mise à jour événement: %v", err)
				_ = msg.Nak()
				return
			}
			logrus.Infof("Événement mis à jour: %s", event.Id)

			// Publier l'alerte
			if err := publishAlert(event); err != nil {
				logrus.Errorf("Erreur lors de la publication de l'alerte: %v", err)
			}
		}
	}

	if err := msg.Ack(); err != nil {
		logrus.Errorf("Erreur ACK: %v", err)
	}
}

// publishAlert publie un message sur TIMETABLE.ALERTER contenant l'événement
func publishAlert(event models.Event) error {
	alertMsg, err := json.Marshal(event)
	if err != nil {
		return err
	}
	// Publication sur le sujet TIMETABLE.ALERTER via la connexion NATS
	return helpers.NatsConn.Publish(alerterSubject, alertMsg)
}

func insertEvent(db *sql.DB, event models.Event) error {
	_, err := db.Exec(`INSERT INTO events 
		(id, dtstamp, dtstart, dtend, location, description, created, lastModified, resourceId) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		event.Id, event.Dtstamp, event.Dtstart, event.Dtend,
		event.Location, event.Description, event.Created, event.LastModified, event.ResourceID)
	return err
}

func updateEvent(db *sql.DB, event models.Event) error {
	_, err := db.Exec(`UPDATE events SET
		dtstamp = ?,
		dtstart = ?,
		dtend = ?,
		location = ?,
		description = ?,
		lastModified = ?
		WHERE id = ?`,
		event.Dtstamp, event.Dtstart, event.Dtend,
		event.Location, event.Description, event.LastModified,
		event.Id)
	return err
}

func hasChanges(old, new models.Event) bool {
	return old.Dtstart != new.Dtstart ||
		old.Dtend != new.Dtend ||
		old.Location != new.Location ||
		old.Description != new.Description
}
