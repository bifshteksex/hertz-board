package database

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/bifshteksex/hertz-board/internal/config"
)

// NewNATSConnection creates a new NATS connection
func NewNATSConnection(cfg *config.NATSConfig) (*nats.Conn, error) {
	opts := []nats.Option{
		nats.MaxReconnects(cfg.MaxReconnect),
		nats.ReconnectWait(time.Duration(cfg.ReconnectWait) * time.Second),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			if err != nil {
				fmt.Printf("NATS disconnected: %v\n", err)
			}
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			fmt.Printf("NATS reconnected to %s\n", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			fmt.Println("NATS connection closed")
		}),
	}

	nc, err := nats.Connect(cfg.URL, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	return nc, nil
}

// CloseNATSConnection closes the NATS connection
func CloseNATSConnection(nc *nats.Conn) {
	if nc != nil {
		nc.Close()
	}
}
