package helpers

import (
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	natsURL = "nats://localhost:4222"
)

var NatsConn *nats.Conn

func InitNATS() error {
	nc, err := nats.Connect(natsURL,
		nats.MaxReconnects(10),
		nats.ReconnectWait(2*time.Second),
		nats.DisconnectErrHandler(func(conn *nats.Conn, err error) {
			logrus.Warnf("Déconnecté de NATS: %v", err)
		}),
		nats.ReconnectHandler(func(conn *nats.Conn) {
			logrus.Info("Reconnecté à NATS")
		}))

	if err != nil {
		return err
	}

	NatsConn = nc
	return nil
}
