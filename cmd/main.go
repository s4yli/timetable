package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/consumers"
	"middleware/example/internal/controllers/events"
	"middleware/example/internal/helpers"
	_ "middleware/example/internal/models"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Initialisation de NATS
	if err := helpers.InitNATS(); err != nil {
		logrus.Fatalf("Échec de la connexion à NATS: %v", err)
	}
	defer helpers.NatsConn.Close()

	// Lancer le consumer NATS
	go runConsumer()

	// Lancer le serveur HTTP
	r := setupRouter()

	// Gestion propre des arrêts
	waitForShutdown()

	logrus.Info("[INFO] Web server started. Now listening on *:8081")
	if err := http.ListenAndServe(":8081", r); err != nil {
		logrus.Fatalf("Erreur du serveur HTTP: %v", err)
	}
}

func runConsumer() {
	consumer, err := consumers.EventConsumer()
	if err != nil {
		logrus.Errorf("Erreur création consumer: %v", err)
		return
	}

	if err := consumers.Consume(consumer); err != nil {
		logrus.Errorf("Erreur consommation: %v", err)
	}
}

func setupRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Route("/events", func(r chi.Router) {
		r.Get("/", events.GetEvents)
		r.Route("/{id}", func(r chi.Router) {
			r.Use(events.Ctx)
			r.Get("/", events.GetEvent)
		})
	})
	return r
}

func waitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		logrus.Info("Réception du signal d'arrêt, fermeture propre...")
		os.Exit(0)
	}()
}

func init() {
	db, err := helpers.OpenDB()
	if err != nil {
		logrus.Fatalf("error while opening database : %s", err.Error())
	}
	schemes := []string{
		`CREATE TABLE IF NOT EXISTS events (
			id VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
			dtstamp VARCHAR(255) NOT NULL,
    		dtstart VARCHAR(255) NOT NULL,
    		dtend VARCHAR(255) NOT NULL,
    		location VARCHAR(255),
    		description VARCHAR(255),
    		created VARCHAR(255) NOT NULL,
    		lastModified VARCHAR(255) NOT NULL,
    		resourceId VARCHAR(255) NOT NULL
		);`,
	}
	for _, scheme := range schemes {
		if _, err := db.Exec(scheme); err != nil {
			logrus.Fatalln("Could not generate table ! Error was : " + err.Error())
		}
	}
	helpers.CloseDB(db)
}
