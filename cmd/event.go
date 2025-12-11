// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"time"

	"github.com/apex/log"
	"github.com/spf13/cobra"
	"gitlab.com/mshindle/tidbits/cmd/internal/ecom"
	"gitlab.com/mshindle/tidbits/event"
	"gitlab.com/mshindle/tidbits/event/memory"
)

// eventCmd represents the gcd command
var eventCmd = &cobra.Command{
	Use:   "event",
	Short: "example of building pub/sub systems in golang",
	Long:  `example of building pub/sub systems in golang`,
	RunE:  eventMain,
}

func init() {
	rootCmd.AddCommand(eventCmd)
}

type customerOrder struct {
	customerID string
	items      []string
	amount     float64
}

func eventMain(cmd *cobra.Command, args []string) error {
	log.Info("starting event driven example")
	bus := event.NewBus()
	store := memory.NewEventStore()

	// store all events
	bus.Subscribe(ecom.WILDCARD, func(e event.Event) error {
		return store.Save(e)
	})

	// initialize services
	svcOrder := ecom.NewOrderService(bus)
	svcEmail := ecom.NewService("email")
	svcInventory := ecom.NewService("inventory")
	svcWarehouse := ecom.NewService("warehouse")
	svcAnalytics := ecom.NewService("analytics")

	log.Info("subscribing services to events")
	bus.Subscribe(ecom.OrderPlaced, svcAnalytics.HandleEvent)
	bus.Subscribe(ecom.OrderPlaced, svcEmail.HandleEvent)
	bus.Subscribe(ecom.OrderPlaced, svcInventory.HandleEvent)
	bus.Subscribe(ecom.OrderPlaced, svcWarehouse.HandleEvent)
	log.WithField("count", bus.GetSubscriberCount(ecom.OrderPlaced)).Info("subscribers registered")

	log.Info("create orders")
	orderDetails := []customerOrder{
		{
			customerID: "customer-123",
			items:      []string{"laptop", "mouse", "keyboard"},
			amount:     1299.99,
		},
		{
			customerID: "customer-123",
			items:      []string{"Monitor", "HDMI Cable"},
			amount:     399.99,
		},
		{
			customerID: "customer-456",
			items:      []string{"Desk Chair", "Standing Desk"},
			amount:     749.99,
		},
	}
	for i, orderDetail := range orderDetails {
		entry := log.WithField("order_num", i)
		entry.Info("creating order")
		if err := svcOrder.CreateOrder(orderDetail.customerID, orderDetail.items, orderDetail.amount); err != nil {
			entry.WithError(err).Error("error creating order")
		}
		time.Sleep(100 * time.Millisecond)
	}

	log.WithField("count", store.Count()).Info("events stored")
	events, _ := store.GetByType(ecom.OrderPlaced)
	log.WithField("count", len(events)).Info("order placed events")

	recentEvents, _ := store.GetLatest(3)
	for i, ev := range recentEvents {
		fmt.Printf("   %d. [%s] %s at %s\n",
			i+1,
			ev.Type,
			ev.ID[:8]+"...",
			ev.Timestamp.Format("15:04:05"))
	}
	return nil
}
