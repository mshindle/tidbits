package ecom

import (
	"fmt"
	"time"

	"github.com/apex/log"
	"github.com/google/uuid"
	"gitlab.com/mshindle/tidbits/event"
)

const OrderPlaced = "order.placed"
const WILDCARD = "*"

// OrderPlacedPayload represents the data for an order placed event
type OrderPlacedPayload struct {
	OrderID    string
	CustomerID string
	Amount     float64
	Items      []string
}

type EcomService struct {
	name string
}

// NewService creates a new service
func NewService(name string) *EcomService {
	return &EcomService{name}
}

// HandleEvent tracks analytics when an order is placed
func (s *EcomService) HandleEvent(event event.Event) error {
	switch payload := event.Payload.(type) {
	case OrderPlacedPayload:
		log.WithFields(log.Fields{
			"service":   s.name,
			"customer":  payload.CustomerID,
			"amount":    fmt.Sprintf("%.2f", payload.Amount),
			"num_items": len(payload.Items),
			"order_id":  payload.OrderID[:8] + "...",
		}).Info("processing order placed event")
	default:
		return fmt.Errorf("unknow payload type for %s service, event: %s", s.name, event.Type)
	}

	return nil
}

type OrderService struct {
	bus *event.Bus
}

func NewOrderService(bus *event.Bus) *OrderService {
	return &OrderService{bus: bus}
}

// CreateOrder creates a new order and publishes an event
func (os *OrderService) CreateOrder(customerID string, items []string, amount float64) error {
	orderID := uuid.New().String()

	ev := event.Event{
		ID:        uuid.New().String(),
		Type:      OrderPlaced,
		Timestamp: time.Now(),
		Payload: OrderPlacedPayload{
			OrderID:    orderID,
			CustomerID: customerID,
			Amount:     amount,
			Items:      items,
		},
	}

	return os.bus.Publish(ev)
}
